package account

import (
	"context"
	"strconv"
	"strings"
	"sync"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/regions"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/availabilityzones"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
	"github.com/supergiant/control/pkg/clouds/gcesdk"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2017-09-01/skus"
	"github.com/Azure/azure-sdk-for-go/services/preview/subscription/mgmt/2018-03-01-preview/subscription"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/digitalocean/godo"
	"github.com/pkg/errors"
	gcecomputev1 "google.golang.org/api/compute/v1"

	"github.com/supergiant/control/pkg/clouds"
	"github.com/supergiant/control/pkg/clouds/digitaloceansdk"
	"github.com/supergiant/control/pkg/model"
	"github.com/supergiant/control/pkg/sgerrors"
	"github.com/supergiant/control/pkg/util"
	"github.com/supergiant/control/pkg/util/strset"
	"github.com/supergiant/control/pkg/workflows/steps"
)

var (
	ErrNilAccount          = errors.New("nil account")
	ErrUnsupportedProvider = errors.New("unsupported provider")

	VMResourceType = "virtualMachines"
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
	ID   string `json:"id"`
	Name string `json:"name"`

	RAM  string `json:"ram"`
	CPU  string `json:"cpu"`
	Disk string `json:"disk"`
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
	case clouds.Azure:
		return NewAzureFinder(account, config)
	case clouds.OpenStack:
		return NewOpenstackFinder(account, config)
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
	case clouds.OpenStack:
		return NewOpenstackFinder(account, config)
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
	case clouds.Azure:
		return NewAzureFinder(account, config)
	case clouds.OpenStack:
		return NewOpenstackFinder(account, config)
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
	machines      MachineTypes

	getZones func(ctx context.Context, client *ec2.EC2,
		input *ec2.DescribeAvailabilityZonesInput) (*ec2.DescribeAvailabilityZonesOutput, error)
}

func NewAWSFinder(acc *model.CloudAccount, config *steps.Config) (*AWSFinder, error) {
	if acc.Provider != clouds.AWS {
		return nil, ErrUnsupportedProvider
	}

	err := util.FillCloudAccountCredentials(acc, config)

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
		machines:      awsMachines,

		getZones: func(ctx context.Context, client *ec2.EC2,
			input *ec2.DescribeAvailabilityZonesInput) (*ec2.DescribeAvailabilityZonesOutput, error) {
			return client.DescribeAvailabilityZonesWithContext(ctx, input)
		},
	}, nil
}

func (af *AWSFinder) GetRegions(ctx context.Context) (*RegionSizes, error) {
	return &RegionSizes{
		Provider: clouds.AWS,
		Regions:  toRegions(af.machines.Regions()),
	}, nil
}

func (af *AWSFinder) GetZones(ctx context.Context, config steps.Config) ([]string, error) {
	azsOut, err := af.getZones(ctx, af.defaultClient, &ec2.DescribeAvailabilityZonesInput{
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
	return af.machines.RegionTypes(config.AWSConfig.Region)
}

type GCEResourceFinder struct {
	client *gcecomputev1.Service
	config steps.Config

	listRegions      func(*gcecomputev1.Service, string) (*gcecomputev1.RegionList, error)
	getRegion        func(*gcecomputev1.Service, string, string) (*gcecomputev1.Region, error)
	listMachineTypes func(*gcecomputev1.Service, string, string) (*gcecomputev1.MachineTypeList, error)
}

func NewGCEFinder(acc *model.CloudAccount, config *steps.Config) (*GCEResourceFinder, error) {
	if acc.Provider != clouds.GCE {
		return nil, ErrUnsupportedProvider
	}

	err := util.FillCloudAccountCredentials(acc, config)

	if err != nil {
		return nil, errors.Wrap(err, "create gce finder")
	}

	client, err := gcesdk.GetClient(context.Background(),
		config.GCEConfig)

	if err != nil {
		return nil, err
	}

	return &GCEResourceFinder{
		client: client,
		config: *config,
		listRegions: func(client *gcecomputev1.Service, projectID string) (*gcecomputev1.RegionList, error) {
			return client.Regions.List(projectID).Do()
		},
		getRegion: func(client *gcecomputev1.Service, projectID, regionID string) (*gcecomputev1.Region, error) {
			return client.Regions.Get(projectID, regionID).Do()
		},
		listMachineTypes: func(client *gcecomputev1.Service, projectID, availabilityZone string) (*gcecomputev1.MachineTypeList, error) {
			return client.MachineTypes.List(projectID, availabilityZone).Do()
		},
	}, nil
}

func (g *GCEResourceFinder) GetRegions(ctx context.Context) (*RegionSizes, error) {
	regionsOutput, err := g.listRegions(g.client, g.config.GCEConfig.ServiceAccount.ProjectID)

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
	regionOutput, err := g.getRegion(g.client, g.config.GCEConfig.ServiceAccount.ProjectID,
		g.config.GCEConfig.Region)

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
	machineOutput, err := g.listMachineTypes(g.client, g.config.GCEConfig.ServiceAccount.ProjectID,
		g.config.GCEConfig.AvailabilityZone)

	if err != nil {
		return nil, errors.Wrap(err, "gce get machine types")
	}

	machineTypes := make([]string, 0)
	for _, machineType := range machineOutput.Items {
		machineTypes = append(machineTypes, machineType.Name)
	}

	return machineTypes, nil
}

type SubscriptionsInterface interface {
	ListLocations(ctx context.Context, subscriptionID string) (subscription.LocationListResult, error)
}

type SKUSInterface interface {
	ListComplete(ctx context.Context) (skus.ResourceSkusResultIterator, error)
}

type AzureFinder struct {
	subscriptionID      string
	location            string
	subscriptionsClient SubscriptionsInterface
	skusClient          SKUSInterface
}

func NewAzureFinder(acc *model.CloudAccount, cfg *steps.Config) (*AzureFinder, error) {
	err := util.FillCloudAccountCredentials(acc, cfg)
	if err != nil {
		return nil, errors.Wrap(err, "retrieve cloud credentials")
	}

	token, err := auth.NewClientCredentialsConfig(cfg.AzureConfig.ClientID, cfg.AzureConfig.ClientSecret, cfg.AzureConfig.TenantID).Authorizer()
	if err != nil {
		return nil, errors.Wrap(err, "get authorization token")
	}

	sclient := subscription.NewSubscriptionsClient()
	sclient.Authorizer = token
	skusclient := skus.NewResourceSkusClient(cfg.AzureConfig.SubscriptionID)
	skusclient.Authorizer = token

	return &AzureFinder{
		subscriptionID:      cfg.AzureConfig.SubscriptionID,
		location:            acc.Credentials["region"],
		subscriptionsClient: sclient,
		skusClient:          skusclient,
	}, nil
}

func (f AzureFinder) GetRegions(ctx context.Context) (*RegionSizes, error) {
	res, err := f.subscriptionsClient.ListLocations(ctx, f.subscriptionID)
	if err != nil {
		return nil, err
	}
	if res.Value == nil {
		return nil, errors.Wrap(sgerrors.ErrNilEntity, "list locations")
	}

	sizes, err := f.getVMSizes(ctx)
	if err != nil {
		return nil, err
	}

	// TODO: cache this data
	return f.toRegionSizes(*res.Value, sizes), nil
}

func (f AzureFinder) GetZones(ctx context.Context, config steps.Config) ([]string, error) {
	// TODO: add azure zones getter: https://github.com/Azure/azure-rest-api-specs/issues/3594
	return []string{}, nil
}

func (f AzureFinder) GetTypes(ctx context.Context, cfg steps.Config) ([]string, error) {
	sizes, err := f.getVMSizes(ctx)
	if err != nil {
		return nil, err
	}
	return vmSizesNames(sizes), nil
}

func (f AzureFinder) getVMSizes(ctx context.Context) ([]skus.ResourceSku, error) {
	sizesList, err := f.skusClient.ListComplete(ctx)
	if err != nil {
		return nil, err
	}

	machineTypes := make([]skus.ResourceSku, 0)
	for {
		err = sizesList.NextWithContext(ctx)
		if err != nil {
			return nil, err
		}

		if !sizesList.NotDone() {
			// exit the loop
			break
		}

		if to.String(sizesList.Value().ResourceType) != VMResourceType {
			continue
		}

		if hasLocationRestriction(sizesList.Value()) {
			continue
		}

		// filter by location
		if sizesList.Value().Locations != nil && f.location != "" && !contains(*sizesList.Value().Locations, f.location) {
			continue
		}

		machineTypes = append(machineTypes, sizesList.Value())
	}

	return machineTypes, nil
}

func (f AzureFinder) toRegionSizes(locations []subscription.Location, machineSizes []skus.ResourceSku) *RegionSizes {
	// regions map with unique sizes
	regionSizes := make(map[string]*strset.Set)
	sizes := make(map[string]interface{})
	for _, vm := range machineSizes {
		if to.String(vm.Name) == "" || to.String(vm.Size) == "" {
			continue
		}

		sizes[to.String(vm.Name)] = vm
		for _, l := range to.StringSlice(vm.Locations) {
			l = strings.ToLower(l) // could be camelCase, for instance SouthAfricaWest
			if regionSizes[l] == nil {
				regionSizes[l] = strset.New()
			}
			regionSizes[l].Add(to.String(vm.Name))
		}
	}

	regions := make([]*Region, 0)
	for _, l := range locations {
		regions = append(regions, &Region{
			ID:             to.String(l.Name),
			Name:           to.String(l.DisplayName),
			AvailableSizes: regionSizes[to.String(l.Name)].ToSlice(),
		})
	}

	return &RegionSizes{
		Provider: clouds.Azure,
		Regions:  regions,
		Sizes:    sizes,
	}
}

func toRegions(names []string) []*Region {
	regions := make([]*Region, len(names))
	for i, r := range names {
		regions[i] = &Region{
			ID:   r,
			Name: r,
		}
	}
	return regions
}

func vmSizesNames(sizes []skus.ResourceSku) []string {
	names := strset.New()
	for _, s := range sizes {
		if to.String(s.Name) == "" || to.String(s.Size) == "" {
			continue
		}
		names.Add(to.String(s.Name))
	}
	return names.ToSlice()
}

func contains(list []string, item string) bool {
	for _, li := range list {
		if li == item {
			return true
		}
	}
	return false
}

func hasLocationRestriction(r skus.ResourceSku) bool {
	return hasRestriction(r, skus.Location)
}

func hasZoneRestriction(r skus.ResourceSku) bool {
	return hasRestriction(r, skus.Zone)
}

func hasRestriction(s skus.ResourceSku, restr skus.ResourceSkuRestrictionsType) bool {
	if s.Restrictions == nil {
		return false
	}
	for _, r := range *s.Restrictions {
		if r.Type == restr {
			return true
		}
	}
	return false
}

type OpenstackFinder struct {
	config steps.OpenStackConfig

	getClient func(config steps.OpenStackConfig) (*gophercloud.ServiceClient, error)
}

func NewOpenstackFinder(acc *model.CloudAccount, config *steps.Config) (*OpenstackFinder, error) {
	err := util.FillCloudAccountCredentials(acc, config)

	if err != nil {
		return nil, errors.Wrap(err, "openstack new finder")
	}

	return &OpenstackFinder{
		config: config.OpenStackConfig,
	}, nil
}

func (of *OpenstackFinder) GetRegions(context.Context) (*RegionSizes, error) {
	opts := gophercloud.AuthOptions{
		IdentityEndpoint: of.config.AuthURL,
		Username:         of.config.UserName,
		Password:         of.config.Password,
		TenantName:       of.config.TenantName,
		DomainID:         of.config.DomainID,
		DomainName:       of.config.DomainName,
	}

	client, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		return nil, err
	}

	computeClient, err := openstack.NewComputeV2(client, gophercloud.EndpointOpts{
		Region: of.config.Region,
	})

	if err != nil {
		return nil, errors.Wrapf(err, "get compute client for get zones")
	}

	page, err := regions.List(computeClient, regions.ListOpts{}).AllPages()

	if err != nil {
		return nil, errors.Wrapf(err, "getting regions")
	}

	regions, err := regions.ExtractRegions(page)

	if err != nil {
		return nil, errors.Wrap(err, "extract regions")
	}

	ids := make([]*Region, 0, 0)

	for _, region := range regions {
		ids = append(ids, &Region{
			ID:   region.ID,
			Name: region.ID,
		})
	}

	rs := &RegionSizes{
		Provider: clouds.GCE,
		Regions:  ids,
	}

	return rs, nil
}

func (of *OpenstackFinder) GetZones(ctx context.Context, config steps.Config) ([]string, error) {
	opts := gophercloud.AuthOptions{
		IdentityEndpoint: config.OpenStackConfig.AuthURL,
		Username:         config.OpenStackConfig.UserName,
		Password:         config.OpenStackConfig.Password,
		TenantName:       config.OpenStackConfig.TenantName,
		DomainID:         config.OpenStackConfig.DomainID,
	}

	client, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		return nil, err
	}

	computeClient, err := openstack.NewComputeV2(client, gophercloud.EndpointOpts{
		Region: of.config.Region,
	})

	if err != nil {
		return nil, errors.Wrapf(err, "get compute client for get zones")
	}

	page, err := availabilityzones.List(computeClient).AllPages()

	if err != nil {
		return nil, errors.Wrap(err, "get az pages")
	}

	azs, err := regions.ExtractRegions(page)

	if err != nil {
		return nil, errors.Wrap(err, "extract regions")
	}

	ids := make([]string, 0, 0)

	for _, az := range azs {
		ids = append(ids, az.ID)
	}

	return ids, nil
}

func (of *OpenstackFinder) GetTypes(ctx context.Context, config steps.Config) ([]string, error) {
	opts := gophercloud.AuthOptions{
		IdentityEndpoint: config.OpenStackConfig.AuthURL,
		Username:         config.OpenStackConfig.UserName,
		Password:         config.OpenStackConfig.Password,
		TenantName:       config.OpenStackConfig.TenantName,
		DomainID:         config.OpenStackConfig.DomainID,
	}

	client, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		return nil, err
	}

	computeClient, err := openstack.NewComputeV2(client, gophercloud.EndpointOpts{
		Region: of.config.Region,
	})

	page, err := flavors.ListDetail(computeClient, flavors.ListOpts{}).AllPages()

	if err != nil {
		return nil, errors.Wrap(err, "list flavor pages")
	}

	flvrs, err := flavors.ExtractFlavors(page)

	if err != nil {
		return nil, errors.Wrap(err, "extract flavors")
	}

	flavorNames := make([]string, 0, 0)

	for _, flavor := range flvrs {
		flavorNames = append(flavorNames, flavor.Name)
	}

	return flavorNames, nil
}
