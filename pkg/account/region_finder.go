package account

import (
	"context"
	"strconv"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/digitalocean/godo"
	"github.com/pkg/errors"

	"github.com/supergiant/supergiant/pkg/clouds"
	"github.com/supergiant/supergiant/pkg/clouds/digitaloceansdk"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
	"github.com/supergiant/supergiant/pkg/util"
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

//RegionFinder is used to find a list of available regions(availability zones, etc) with available vm types
//in a given cloud provider using given account credentials
type RegionFinder interface {
	//Find returns a slice of cloud specific regions/az's
	//if not found would return an empty slice
	Find(context.Context) (*RegionSizes, error)
}

//GetRegionFinder returns finder attached to corresponding account as it has all credentials for a cloud provider
func GetRegionFinder(account *model.CloudAccount) (RegionFinder, error) {
	if account == nil {
		return nil, ErrNilAccount
	}

	switch account.Provider {
	case clouds.DigitalOcean:
		sdk, err := digitaloceansdk.NewFromAccount(account)
		if err != nil {
			return nil, err
		}
		return &digitalOceanRegionFinder{
			getServices: func() (godo.SizesService, godo.RegionsService) {
				client := sdk.GetClient()
				return client.Sizes, client.Regions
			},
		}, nil
	case clouds.AWS:
		return NewAWSFinder(account)
	case clouds.GCE:
		return NewGCEFinder(account)
	}
	return nil, ErrUnsupportedProvider
}

type digitalOceanRegionFinder struct {
	sdk         *digitaloceansdk.SDK
	getServices func() (godo.SizesService, godo.RegionsService)
}

func (rf *digitalOceanRegionFinder) Find(ctx context.Context) (*RegionSizes, error) {
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
	accessKey string
	secret    string
}

func (af *AWSFinder) Find(ctx context.Context) (*RegionSizes, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region:      aws.String("us-west-1"),
			Credentials: credentials.NewStaticCredentials(af.accessKey, af.secret, ""),
		},
	})

	if err != nil {
		return nil, errors.Wrap(err, "aws authentication: ")
	}

	EC2 := ec2.New(sess)

	regionsOut, err := EC2.DescribeRegionsWithContext(ctx, &ec2.DescribeRegionsInput{})

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

func (af *AWSFinder) GetAZ(ctx context.Context, region string) ([]string, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region:      aws.String(region),
			Credentials: credentials.NewStaticCredentials(af.accessKey, af.secret, ""),
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "aws authentication: ")
	}

	EC2 := ec2.New(sess)
	azsOut, err := EC2.DescribeAvailabilityZonesWithContext(ctx, &ec2.DescribeAvailabilityZonesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("region-name"),
				Values: []*string{
					aws.String(region),
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

func (af *AWSFinder) GetTypes(ctx context.Context, region, az string) ([]string, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region:      aws.String(region),
			Credentials: credentials.NewStaticCredentials(af.accessKey, af.secret, ""),
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "aws authentication: ")
	}
	EC2 := ec2.New(sess)

	out, err := EC2.DescribeReservedInstancesOfferingsWithContext(ctx, &ec2.DescribeReservedInstancesOfferingsInput{})

	if err != nil {
		return nil, errors.Wrap(err, "failed to read aws types")
	}

	instances := make([]string, 0)
	for _, of := range out.ReservedInstancesOfferings {
		instances = append(instances, *of.InstanceType)
	}

	return instances, nil
}

func NewAWSFinder(acc *model.CloudAccount) (*AWSFinder, error) {
	if acc.Provider != clouds.AWS {
		return nil, ErrUnsupportedProvider
	}

	accessKey := acc.Credentials[clouds.AWSAccessKeyID]
	secret := acc.Credentials[clouds.AWSSecretKey]

	if accessKey == "" {
		return nil, errors.New("no access key id provided for AWS account")
	}

	if secret == "" {
		return nil, errors.New("no secret key provided for AWS account")
	}

	return &AWSFinder{
		accessKey: accessKey,
		secret:    secret,
	}, nil
}

type GCEFinder struct{
	config steps.Config
}

func NewGCEFinder(acc *model.CloudAccount) (*GCEFinder, error) {
	if acc.Provider != clouds.GCE {
		return nil, ErrUnsupportedProvider
	}

	gceFinder := &GCEFinder{
		config: steps.Config{},
	}

	err := util.FillCloudAccountCredentials(context.Background(),
		acc, &gceFinder.config)

	if err != nil {
		return nil, errors.Wrap(err, "create gce finder")
	}

	return gceFinder, nil
}

func (g *GCEFinder) Find(ctx context.Context) (*RegionSizes, error) {
	client, err := gce.GetClient(ctx, g.config.GCEConfig.ClientEmail,
		g.config.GCEConfig.PrivateKey, g.config.GCEConfig.TokenURI)

	regionsOutput, err := client.Regions.List(g.config.GCEConfig.ProjectID).Do()

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
		Provider: clouds.AWS,
		Regions:  regions,
	}

	return rs, nil
}

func (g *GCEFinder) GetAZ(ctx context.Context, region string) ([]string, error) {
	client, err := gce.GetClient(ctx, g.config.GCEConfig.ClientEmail,
		g.config.GCEConfig.PrivateKey, g.config.GCEConfig.TokenURI)

	regionOutput, err := client.Regions.Get(g.config.GCEConfig.ProjectID, region).Do()

	if err != nil {
		return nil, errors.Wrap(err, "gce get availability zone")
	}

	return regionOutput.Zones, nil
}

func (g *GCEFinder) GetTypes(ctx context.Context, zone string) ([]string, error) {
	client, err := gce.GetClient(ctx, g.config.GCEConfig.ClientEmail,
		g.config.GCEConfig.PrivateKey, g.config.GCEConfig.TokenURI)

	machineOutput, err := client.MachineTypes.List(g.config.GCEConfig.ProjectID, zone).Do()

	if err != nil {
		return nil, errors.Wrap(err, "gce get machine types")
	}

	machineTypes := make([]string, 0)
	for _, machineType := range machineOutput.Items {
		machineTypes = append(machineTypes, machineType.Name)
	}

	return machineTypes, nil
}