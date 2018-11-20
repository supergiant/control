package account

import (
	"context"
	"strconv"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/digitalocean/godo"
	"github.com/pkg/errors"
	compute "google.golang.org/api/compute/v1"

	"github.com/supergiant/supergiant/pkg/clouds"
	"github.com/supergiant/supergiant/pkg/clouds/digitaloceansdk"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/util"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
	"github.com/supergiant/supergiant/pkg/workflows/steps/gce"
)

var (
	ErrNilAccount          = errors.New("nil account")
	ErrUnsupportedProvider = errors.New("unsupported provider")
)

//Region represents
type Region struct {
	//Human readable name, e.g. New York City 1 or EU West 1 Frankfurt
	Name string `json:"name"`
	//API specific ID, e.g. t2.micro
	ID string `json:"id"`

	//API specific IDs for a node sizes/type
	AvailableSizes []string
}

type Size struct {
	RAM string `json:"ram"`
	CPU string `json:"cpu"`
}

//RegionSizes represents aggregated information about available regions/azs and node sizes/types
type RegionSizes struct {
	Provider clouds.Name            `json:"provider"`
	Regions  []*Region              `json:"regions"`
	Sizes    map[string]interface{} `json:"sizes"`
}

type ZonesGetter interface {
	GetZones(context.Context, steps.Config) ([]string, error)
}

type TypesGetter interface {
	GetTypes(context.Context, steps.Config) ([]string, error)
}

//RegionsGetter is used to find a list of available regions(availability zones, etc) with available vm types
//in a given cloud provider using given account credentials
type RegionsGetter interface {
	//GetRegions returns a slice of cloud specific regions/az's
	//if not found would return an empty slice
	GetRegions(context.Context) (*RegionSizes, error)
}

//NewRegionsGetter returns finder attached to corresponding account as it has all credentials for a cloud provider
func NewRegionsGetter(account *model.CloudAccount, config *steps.Config) (RegionsGetter, error) {
	if account == nil {
		return nil, ErrNilAccount
	}

	switch account.Provider {
	case clouds.DigitalOcean:
		return NewDOFinder(account)
	case clouds.AWS:
		// We need to provide region to AWS even if our
		// request does not specify region
		config.AWSConfig.Region = "us-west-1"
		return NewAWSFinder(account, config)
	case clouds.GCE:
		return NewGCEFinder(account, config)
	}
	return nil, ErrUnsupportedProvider
}

//NewZonesGetter returns finder attached to corresponding
// account as it has all credentials for a cloud provider
func NewZonesGetter(account *model.CloudAccount, config *steps.Config) (ZonesGetter, error) {
	if account == nil {
		return nil, ErrNilAccount
	}

	switch account.Provider {
	case clouds.AWS:
		return NewAWSFinder(account, config)
	case clouds.GCE:
		return NewGCEFinder(account, config)
	}
	return nil, ErrUnsupportedProvider
}

//NewTypesGetter returns finder attached to corresponding
// account as it has all credentials for a cloud provider
func NewTypesGetter(account *model.CloudAccount, config *steps.Config) (TypesGetter, error) {
	if account == nil {
		return nil, ErrNilAccount
	}

	switch account.Provider {
	case clouds.AWS:
		return NewAWSFinder(account, config)
	case clouds.GCE:
		return NewGCEFinder(account, config)
	}
	return nil, ErrUnsupportedProvider
}

type digitalOceanRegionFinder struct {
	sdk         *digitaloceansdk.SDK
	getServices func() (godo.SizesService, godo.RegionsService)
}

func NewDOFinder(acc *model.CloudAccount) (*digitalOceanRegionFinder, error) {
	sdk, err := digitaloceansdk.NewFromAccount(acc)
	if err != nil {
		return nil, err
	}
	return &digitalOceanRegionFinder{
		getServices: func() (godo.SizesService, godo.RegionsService) {
			client := sdk.GetClient()
			return client.Sizes, client.Regions
		},
	}, nil
}

func (rf *digitalOceanRegionFinder) GetRegions(ctx context.Context) (*RegionSizes, error) {
	sizeService, regionService := rf.getServices()

	var wg sync.WaitGroup
	var sizes []godo.Size
	var sizeErr error

	var doRegions []godo.Region
	var regionErr error

	wg.Add(2)
	go func() {
		defer wg.Done()
		sizes, _, sizeErr = sizeService.List(ctx, nil)
	}()
	go func() {
		defer wg.Done()
		doRegions, _, regionErr = regionService.List(ctx, nil)
	}()
	//assignment will work fine because of the memory barrier
	wg.Wait()

	if sizeErr != nil {
		return nil, sizeErr
	}
	if regionErr != nil {
		return nil, regionErr
	}

	nodeSizes := make(map[string]interface{})
	regions := make([]*Region, 0, len(doRegions))

	for _, s := range sizes {
		convertSize(s, nodeSizes)
	}

	for _, r := range doRegions {
		regions = append(regions, convertRegion(r))
	}

	rs := &RegionSizes{
		Provider: clouds.DigitalOcean,
		Regions:  regions,
		Sizes:    nodeSizes,
	}

	return rs, nil
}

func convertSize(s godo.Size, nodeSizes map[string]interface{}) {
	ns := Size{
		RAM: strconv.Itoa(s.Memory),
		CPU: strconv.Itoa(s.Vcpus),
	}
	nodeSizes[s.Slug] = ns
}

func convertRegion(r godo.Region) *Region {
	region := &Region{
		ID:             r.Slug,
		Name:           r.Name,
		AvailableSizes: r.Sizes,
	}

	return region
}

type AWSFinder struct {
	defaultClient *ec2.EC2
}

func NewAWSFinder(acc *model.CloudAccount, config *steps.Config) (*AWSFinder, error) {
	if acc.Provider != clouds.AWS {
		return nil, ErrUnsupportedProvider
	}

	err := util.FillCloudAccountCredentials(context.Background(), acc, config)

	if err != nil {
		return nil, errors.Wrap(err, "aws new finder")
	}

	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String(config.AWSConfig.Region),
			Credentials: credentials.NewStaticCredentials(
				config.AWSConfig.KeyID, config.AWSConfig.Secret,
				""),
		},
	})

	if err != nil {
		return nil, errors.Wrap(err, "aws authentication: ")
	}

	client := ec2.New(sess)

	return &AWSFinder{
		defaultClient: client,
	}, nil
}

func (af *AWSFinder) GetRegions(ctx context.Context) (*RegionSizes, error) {
	regionsOut, err := af.defaultClient.DescribeRegionsWithContext(ctx, &ec2.DescribeRegionsInput{})

	if err != nil {
		return nil, errors.Wrap(err, "failed to read aws regions")
	}

	regions := make([]*Region, 0)
	for _, r := range regionsOut.Regions {
		regions = append(regions, &Region{
			ID:   *r.RegionName,
			Name: *r.RegionName,
		})
	}

	rs := &RegionSizes{
		Provider: clouds.AWS,
		Regions:  regions,
	}

	return rs, nil
}

func (af *AWSFinder) GetZones(ctx context.Context, config steps.Config) ([]string, error) {
	azsOut, err := af.defaultClient.DescribeAvailabilityZonesWithContext(ctx, &ec2.DescribeAvailabilityZonesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("region-name"),
				Values: []*string{
					aws.String(config.AWSConfig.Region),
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	zones := make([]string, 0)
	for _, az := range azsOut.AvailabilityZones {
		zones = append(zones, *az.ZoneName)
	}

	return zones, nil
}

func (af *AWSFinder) GetTypes(ctx context.Context, config steps.Config) ([]string, error) {
	out, err := af.defaultClient.DescribeReservedInstancesOfferingsWithContext(ctx, &ec2.DescribeReservedInstancesOfferingsInput{})

	if err != nil {
		return nil, errors.Wrap(err, "failed to read aws types")
	}

	instances := make([]string, 0)
	for _, of := range out.ReservedInstancesOfferings {
		instances = append(instances, *of.InstanceType)
	}

	return instances, nil
}

type GCEResourceFinder struct {
	client *compute.Service
	config steps.Config

	listRegions      func(*compute.Service, string) (*compute.RegionList, error)
	getRegion        func(*compute.Service, string, string) (*compute.Region, error)
	listMachineTypes func(*compute.Service, string, string) (*compute.MachineTypeList, error)
}

func NewGCEFinder(acc *model.CloudAccount, config *steps.Config) (*GCEResourceFinder, error) {
	if acc.Provider != clouds.GCE {
		return nil, ErrUnsupportedProvider
	}

	err := util.FillCloudAccountCredentials(context.Background(),
		acc, config)

	if err != nil {
		return nil, errors.Wrap(err, "create gce finder")
	}

	client, err := gce.GetClient(context.Background(),
		config.GCEConfig.ClientEmail, config.GCEConfig.PrivateKey,
		config.GCEConfig.TokenURI)

	if err != nil {
		return nil, err
	}

	return &GCEResourceFinder{
		client: client,
		config: *config,
		listRegions: func(client *compute.Service, projectID string) (*compute.RegionList, error) {
			return client.Regions.List(projectID).Do()
		},
		getRegion: func(client *compute.Service, projectID, regionID string) (*compute.Region, error) {
			return client.Regions.Get(projectID, regionID).Do()
		},
		listMachineTypes: func(client *compute.Service, projectID, availabilityZone string) (*compute.MachineTypeList, error) {
			return client.MachineTypes.List(projectID, availabilityZone).Do()
		},
	}, nil
}

func (g *GCEResourceFinder) GetRegions(ctx context.Context) (*RegionSizes, error) {
	regionsOutput, err := g.listRegions(g.client, g.config.GCEConfig.ProjectID)

	if err != nil {
		return nil, errors.Wrap(err, "gce find regions")
	}

	regions := make([]*Region, 0)
	for _, r := range regionsOutput.Items {
		regions = append(regions, &Region{
			ID:   r.Name,
			Name: r.Name,
		})
	}

	rs := &RegionSizes{
		Provider: clouds.GCE,
		Regions:  regions,
	}

	return rs, nil
}

func (g *GCEResourceFinder) GetZones(ctx context.Context, config steps.Config) ([]string, error) {
	regionOutput, err := g.getRegion(g.client, config.GCEConfig.ProjectID,
		config.GCEConfig.Region)

	if err != nil {
		return nil, errors.Wrap(err, "gce get availability zones")
	}

	zones := make([]string, 0, len(regionOutput.Zones))

	for _, zoneLink := range regionOutput.Zones {
		splitted := strings.Split(zoneLink, "/")
		zones = append(zones, splitted[len(splitted)-1])
	}
	return zones, nil
}

func (g *GCEResourceFinder) GetTypes(ctx context.Context, config steps.Config) ([]string, error) {
	machineOutput, err := g.listMachineTypes(g.client, config.GCEConfig.ProjectID,
		config.GCEConfig.AvailabilityZone)

	if err != nil {
		return nil, errors.Wrap(err, "gce get machine types")
	}

	machineTypes := make([]string, 0)
	for _, machineType := range machineOutput.Items {
		machineTypes = append(machineTypes, machineType.Name)
	}

	return machineTypes, nil
}
