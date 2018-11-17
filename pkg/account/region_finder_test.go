package account

import (
	"context"
	"strconv"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"

	"github.com/supergiant/supergiant/pkg/clouds"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/sgerrors"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
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
		rf, err := GetRegionsGetter(testCase.account, &steps.Config{})

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
