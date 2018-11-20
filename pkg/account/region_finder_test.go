package account

import (
	"context"
	"strconv"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"

	compute "google.golang.org/api/compute/v1"

	"github.com/supergiant/supergiant/pkg/clouds"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/sgerrors"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/aws"
)

type mockSizeService struct {
	mock.Mock
}

type mockRegionService struct {
	mock.Mock
}

func (m *mockSizeService) List(ctx context.Context, options *godo.ListOptions) ([]godo.Size,
	*godo.Response, error) {
	args := m.Called(ctx, options)
	val, ok := args.Get(0).([]godo.Size)
	if !ok {
		return nil, nil, args.Error(1)
	}
	return val, nil, args.Error(1)
}

func (m *mockRegionService) List(ctx context.Context, options *godo.ListOptions) ([]godo.Region,
	*godo.Response, error) {
	args := m.Called(ctx, options)
	val, ok := args.Get(0).([]godo.Region)
	if !ok {
		return nil, nil, args.Error(1)
	}
	return val, nil, args.Error(1)
}

func TestGetRegionFinder(t *testing.T) {
	testCases := []struct {
		account *model.CloudAccount
		err     error
	}{
		{
			account: nil,
			err:     ErrNilAccount,
		},
		{
			account: &model.CloudAccount{
				Provider: "Unknown",
			},
			err: ErrUnsupportedProvider,
		},
		{
			account: &model.CloudAccount{
				Provider: clouds.DigitalOcean,
				Credentials: map[string]string{
					"dumb": "1234",
				},
			},
			err: sgerrors.ErrInvalidCredentials,
		},
		{
			account: &model.CloudAccount{
				Provider: clouds.DigitalOcean,
				Credentials: map[string]string{
					"accessToken": "1234",
				},
			},
		},
	}

	for _, testCase := range testCases {
		rf, err := NewRegionsGetter(testCase.account, &steps.Config{})

		if err != testCase.err {
			t.Errorf("expected error %v actual %v", testCase.err, err)
		}

		if err == nil && rf == nil {
			t.Error("region finder must not be nil")
		}
	}
}

func TestFind(t *testing.T) {
	errRegion := errors.New("region")
	errSize := errors.New("sizes")

	testCases := []struct {
		regions []godo.Region
		sizes   []godo.Size

		sizeErr   error
		regionErr error

		expectedErr    error
		expectedOutput *RegionSizes
	}{
		{
			sizeErr:     errSize,
			expectedErr: errSize,
		},
		{
			regionErr:   errRegion,
			expectedErr: errRegion,
		},
		{
			sizes:   []godo.Size{},
			regions: []godo.Region{},

			regionErr: nil,
			sizeErr:   nil,

			expectedErr: nil,
			expectedOutput: &RegionSizes{
				Regions: []*Region{},
				Sizes:   map[string]interface{}{},
			},
		},
	}

	for _, testCase := range testCases {
		sizeSvc := &mockSizeService{}
		sizeSvc.On("List", mock.Anything, mock.Anything).
			Return(testCase.sizes, testCase.sizeErr)

		regionSvc := &mockRegionService{}
		regionSvc.On("List", mock.Anything, mock.Anything).
			Return(testCase.regions, testCase.regionErr)

		rf := digitalOceanRegionFinder{
			getServices: func() (godo.SizesService, godo.RegionsService) {
				return sizeSvc, regionSvc
			},
		}

		regionSizes, err := rf.GetRegions(context.Background())

		if err != testCase.expectedErr {
			t.Errorf("expected error %v actual %v", testCase.expectedErr, err)
		}

		if err == nil && regionSizes == nil {
			t.Error("output must not be nil")
		}

		if testCase.expectedErr == nil {
			if regionSizes.Provider != clouds.DigitalOcean {
				t.Errorf("Wrong cloud provider expected %s actual %s",
					clouds.DigitalOcean, regionSizes.Provider)
			}

			if len(regionSizes.Regions) != len(testCase.regions) {
				t.Errorf("wrong count of regions expected %d actual %d",
					len(testCase.regions), len(regionSizes.Regions))
			}

			if len(regionSizes.Sizes) != len(testCase.sizes) {
				t.Errorf("wrong count of sizes expected %d actual %d",
					len(testCase.sizes), len(regionSizes.Sizes))
			}
		}
	}
}

func TestConvertSize(t *testing.T) {
	memory := 16
	vcpus := 4

	size := godo.Size{
		Slug:   "test",
		Memory: memory,
		Vcpus:  vcpus,
	}

	nodeSizes := map[string]interface{}{}
	convertSize(size, nodeSizes)

	if _, ok := nodeSizes[size.Slug]; !ok {
		t.Errorf("size with slug %s not found in %v",
			size.Slug, nodeSizes)
		return
	}

	s, ok := nodeSizes[size.Slug].(Size)

	if !ok {
		t.Errorf("Wrong type of value %v expected Size", nodeSizes[size.Slug])
		return
	}

	if s.CPU != strconv.Itoa(size.Vcpus) {
		t.Errorf("wrong vcpu count expected %d actual %s", size.Vcpus, s.CPU)
	}

	if s.RAM != strconv.Itoa(size.Memory) {
		t.Errorf("wrong memory count expected %d actual %s", size.Memory, s.RAM)
	}
}

func TestConvertRegions(t *testing.T) {
	region := godo.Region{
		Slug:  "fra1",
		Name:  "Frankfurt1",
		Sizes: []string{"size-1", "size-2"},
	}

	r := convertRegion(region)

	if r.Name != region.Name {
		t.Errorf("Wrong name of region expected %s actual %s", region.Name, r.Name)
	}

	if r.ID != region.Slug {
		t.Errorf("Wrong ID of region expected %s actual %s", region.Slug, r.ID)
	}

	if len(r.AvailableSizes) != len(region.Sizes) {
		t.Errorf("Wrong count of sizes expected %d actual %d",
			len(region.Sizes), len(r.AvailableSizes))
	}
}

func TestGCEResourceFinder_GetRegions(t *testing.T) {
	testCases := []struct {
		projectID  string
		err        error
		regionList *compute.RegionList
	}{
		{
			projectID:  "test",
			err:        sgerrors.ErrNotFound,
			regionList: nil,
		},
		{
			projectID: "test",
			err:       nil,
			regionList: &compute.RegionList{
				Items: []*compute.Region{
					{
						Name: "europe-1",
					},
					{
						Name: "ap-north-2",
					},
					{
						Name: "us-west-3",
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		finder := &GCEResourceFinder{
			client: nil,
			config: steps.Config{
				GCEConfig: steps.GCEConfig{
					ProjectID: testCase.projectID,
				},
			},
			listRegions: func(client *compute.Service, projectID string) (*compute.RegionList, error) {
				if projectID != testCase.projectID {
					t.Errorf("Expected projectID %s actual %s",
						testCase.projectID, projectID)
				}

				return testCase.regionList, testCase.err
			},
		}

		regionSizes, err := finder.GetRegions(context.Background())

		if testCase.err != nil && !sgerrors.IsNotFound(err) {
			t.Errorf("Expected err %v actual %v", testCase.err, err)
		}

		if testCase.err == nil {
			if len(regionSizes.Regions) != len(testCase.regionList.Items) {
				t.Errorf("Wrong count of regions expected %d actual %d",
					len(testCase.regionList.Items), len(regionSizes.Regions))
			}
		}
	}
}

func TestGCEResourceFinder_GetZones(t *testing.T) {
	testCases := []struct {
		projectID string
		regionID  string
		err       error
		region    *compute.Region
	}{
		{
			projectID: "test",
			regionID:  "us-east33",
			err:       sgerrors.ErrNotFound,
			region:    nil,
		},
		{
			projectID: "test",
			regionID:  "us-east1",
			err:       nil,
			region: &compute.Region{
				Zones: []string{"us-east1-b", "us-east1-c", "us-east1-d"},
			},
		},
	}

	for _, testCase := range testCases {
		finder := &GCEResourceFinder{
			client: nil,
			config: steps.Config{
				GCEConfig: steps.GCEConfig{
					ProjectID: testCase.projectID,
				},
			},
			getRegion: func(client *compute.Service, projectID, regionID string) (*compute.Region, error) {
				if projectID != testCase.projectID {
					t.Errorf("Expected projectID %s actual %s",
						testCase.projectID, projectID)
				}

				if regionID != testCase.regionID {
					t.Errorf("Expected regionID %s actual %s",
						testCase.regionID, regionID)
				}

				return testCase.region, testCase.err
			},
		}

		config := steps.Config{
			GCEConfig: steps.GCEConfig{
				ProjectID: testCase.projectID,
				Region:    testCase.regionID,
			},
		}
		zones, err := finder.GetZones(context.Background(), config)

		if testCase.err != nil && !sgerrors.IsNotFound(err) {
			t.Errorf("Expected err %v actual %v", testCase.err, err)
		}

		if testCase.err == nil {
			if len(zones) != len(testCase.region.Zones) {
				t.Errorf("Wrong count of zones expected %d actual %d",
					len(testCase.region.Zones), len(zones))
			}
		}
	}
}

func TestGCEResourceFinder_GetTypes(t *testing.T) {
	testCases := []struct {
		projectID string
		zoneID    string
		err       error
		types     *compute.MachineTypeList
	}{
		{
			projectID: "test",
			zoneID:    "us-east33-a",
			err:       sgerrors.ErrNotFound,
			types:     nil,
		},
		{
			projectID: "test",
			zoneID:    "us-east1-b",
			err:       nil,
			types: &compute.MachineTypeList{
				Items: []*compute.MachineType{
					{
						Name: "n1-standard-8",
					},
					{
						Name: "n1-highmem-32",
					},
					{
						Name: "n1-highcpu-96",
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		config := steps.Config{
			GCEConfig: steps.GCEConfig{
				ProjectID:        testCase.projectID,
				AvailabilityZone: testCase.zoneID,
			},
		}

		finder := &GCEResourceFinder{
			client: nil,
			config: config,
			listMachineTypes: func(client *compute.Service, projectID, zoneID string) (*compute.MachineTypeList, error) {
				if projectID != testCase.projectID {
					t.Errorf("Expected projectID %s actual %s",
						testCase.projectID, projectID)
				}

				if zoneID != testCase.zoneID {
					t.Errorf("Expected types %s actual %s",
						testCase.zoneID, zoneID)
				}

				return testCase.types, testCase.err
			},
		}

		types, err := finder.GetTypes(context.Background(), config)

		if testCase.err != nil && !sgerrors.IsNotFound(err) {
			t.Errorf("Expected err %v actual %v", testCase.err, err)
		}

		if testCase.err == nil {
			if len(types) != len(testCase.types.Items) {
				t.Errorf("Wrong count of types expected %d actual %d",
					len(testCase.types.Items), len(types))
			}
		}
	}
}

func TestAWSFinder_GetRegions(t *testing.T) {
	testCases := []struct {
		err       error
		resp     *ec2.DescribeRegionsOutput
	}{
		{
			err: sgerrors.ErrNotFound,
			resp: nil,
		},
		{
			err: nil,
			resp: &ec2.DescribeRegionsOutput{
				Regions:[]*ec2.Region{
					{
						RegionName: aws.String("ap-northeast-1"),
					},
					{
						RegionName: aws.String("eu-west-2"),
					},
					{
						RegionName: aws.String("us-west-1"),
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		awsFinder :=  &AWSFinder{
			getRegions: func(ctx context.Context, client *ec2.EC2,
				input *ec2.DescribeRegionsInput) (*ec2.DescribeRegionsOutput, error) {
				return testCase.resp, testCase.err
			},
		}

		resp, err := awsFinder.GetRegions(context.Background())

		if testCase.err != nil && !sgerrors.IsNotFound(err) {
			t.Errorf("wrong error expected %v actual %v", testCase.err, err)
		}

		if err == nil && len(resp.Regions) != len(testCase.resp.Regions) {
			t.Errorf("Wrong count of regions expected %d actual %d",
				len(testCase.resp.Regions), len(resp.Regions))
		}
	}
}

func TestAWSFinder_GetZones(t *testing.T) {
	testCases := []struct {
		err       error
		resp     *ec2.DescribeAvailabilityZonesOutput
	}{
		{
			err: sgerrors.ErrNotFound,
			resp: nil,
		},
		{
			err: nil,
			resp: &ec2.DescribeAvailabilityZonesOutput{
				AvailabilityZones:[]*ec2.AvailabilityZone{
					{
						ZoneName: aws.String("ap-northeast1-b"),
					},
					{
						ZoneName: aws.String("eu-west2-a"),
					},
					{
						ZoneName: aws.String("us-west1-c"),
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		awsFinder :=  &AWSFinder{
			getZones: func(ctx context.Context, client *ec2.EC2,
				input *ec2.DescribeAvailabilityZonesInput) (*ec2.DescribeAvailabilityZonesOutput, error) {
				return testCase.resp, testCase.err
			},
		}

		resp, err := awsFinder.GetZones(context.Background(), steps.Config{})

		if testCase.err != nil && !sgerrors.IsNotFound(err) {
			t.Errorf("wrong error expected %v actual %v", testCase.err, err)
		}

		if err == nil && len(resp) != len(testCase.resp.AvailabilityZones) {
			t.Errorf("Wrong count of regions expected %d actual %d",
				len(testCase.resp.AvailabilityZones), len(resp))
		}
	}
}

func TestAWSFinder_GetTypes(t *testing.T) {
	testCases := []struct {
		err       error
		resp     *ec2.DescribeReservedInstancesOfferingsOutput
	}{
		{
			err: sgerrors.ErrNotFound,
			resp: nil,
		},
		{
			err: nil,
			resp: &ec2.DescribeReservedInstancesOfferingsOutput{
				ReservedInstancesOfferings:[]*ec2.ReservedInstancesOffering{
					{
						InstanceType: aws.String("t3.medium"),
					},
					{
						InstanceType: aws.String("c5.2xlarge"),
					},
					{
						InstanceType: aws.String("r5d.xlarge"),
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		awsFinder :=  &AWSFinder{
			getTypes: func(ctx context.Context, client *ec2.EC2,
				input *ec2.DescribeReservedInstancesOfferingsInput) (*ec2.DescribeReservedInstancesOfferingsOutput,
					error) {
				return testCase.resp, testCase.err
			},
		}

		resp, err := awsFinder.GetTypes(context.Background(), steps.Config{})

		if testCase.err != nil && !sgerrors.IsNotFound(err) {
			t.Errorf("wrong error expected %v actual %v", testCase.err, err)
		}

		if err == nil && len(resp) != len(testCase.resp.ReservedInstancesOfferings) {
			t.Errorf("Wrong count of regions expected %d actual %d",
				len(testCase.resp.ReservedInstancesOfferings), len(resp))
		}
	}
}
