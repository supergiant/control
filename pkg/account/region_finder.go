package account

import (
	"context"
	"strconv"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/digitalocean/godo"
	"github.com/pkg/errors"

	"github.com/supergiant/supergiant/pkg/clouds"
	"github.com/supergiant/supergiant/pkg/clouds/awssdk"
	"github.com/supergiant/supergiant/pkg/clouds/digitaloceanSDK"
	"github.com/supergiant/supergiant/pkg/model"
)

var (
	ErrUnsupportedProvider = errors.New("unsupported provider")
)

//Region represents
type Region struct {
	//Human readable name, e.g. New York City 1 or EU West 1 Frankfurt
	Name string `json:"name"`
	//API specific ID, e.g. t2.micro
	ID string `json:"id"`

	//API specific IDs for a node size/type
	AvailableSizes []string
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
	switch account.Provider {
	case clouds.DigitalOcean:
		sdk, err := digitaloceanSDK.NewFromAccount(account)
		if err != nil {
			return nil, err
		}
		return &digitalOceanRegionFinder{
			sdk: sdk,
		}, nil
	case clouds.AWS:
		sdk, err := awssdk.NewFromAccount(account)
		if err != nil {
			return nil, err
		}
		return &awsRegionFinder{
			sdk: sdk,
		}, nil
	}
	return nil, ErrUnsupportedProvider
}

type digitalOceanRegionFinder struct {
	sdk *digitaloceanSDK.SDK
}

func (rf *digitalOceanRegionFinder) Find(ctx context.Context) (*RegionSizes, error) {
	cl := rf.sdk.GetClient()
	regions := make([]*Region, 0)

	var wg sync.WaitGroup
	var sizes []godo.Size
	var sizeErr error

	var doRegions []godo.Region
	var doErr error

	wg.Add(2)
	go func() {
		defer wg.Done()
		sizes, _, sizeErr = cl.Sizes.List(ctx, nil)
	}()
	go func() {
		defer wg.Done()
		doRegions, _, doErr = cl.Regions.List(ctx, nil)
	}()
	//assignment will work fine because of the memory barrier
	wg.Wait()

	if sizeErr != nil {
		return nil, sizeErr
	}
	if doErr != nil {
		return nil, doErr
	}

	nodeSizes := make(map[string]interface{})
	for _, s := range sizes {
		ns := struct {
			RAM string `json:"ram"`
			CPU string `json:"cpu"`
		}{
			RAM: strconv.Itoa(s.Memory),
			CPU: strconv.Itoa(s.Vcpus),
		}
		nodeSizes[s.Slug] = ns
	}

	for _, r := range doRegions {
		region := &Region{
			ID:             r.Slug,
			Name:           r.Name,
			AvailableSizes: r.Sizes,
		}
		regions = append(regions, region)
	}

	rs := &RegionSizes{
		Provider: clouds.DigitalOcean,
		Regions:  regions,
		Sizes:    nodeSizes,
	}

	return rs, nil
}

type awsRegionFinder struct {
	sdk *awssdk.SDK
}

func (rf *awsRegionFinder) Find(ctx context.Context) (*RegionSizes, error) {
	regionsOut, err := rf.sdk.EC2.DescribeRegionsWithContext(ctx, &ec2.DescribeRegionsInput{})
	if err != nil {
		return nil, err
	}

	rf.sdk.EC2.DescribeReservedInstancesOfferingsWithContext(ctx,
		&ec2.DescribeReservedInstancesOfferingsInput{
			Filters: []*ec2.Filter{
				{
					Name: aws.String(""),
					Values: aws.StringSlice(
						[]string{},
					),
				},
			},
		})
	return nil, nil
}
