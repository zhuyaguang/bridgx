package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/galaxy-future/BridgX/internal/logs"
	"github.com/galaxy-future/BridgX/pkg/cloud"
)

type AwsCloud struct {
	ec2Client *ec2.EC2
}

func New(ak, sk, regionId string) (*AwsCloud, error) {
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(ak, sk, ""),
		Region:      aws.String(regionId),
	})
	if err != nil {
		return nil, err
	}
	svc := ec2.New(sess)
	return &AwsCloud{svc}, nil
}

func (AwsCloud) ProviderType() string {
	return cloud.AwsCloud
}

// GetRegions miss LocalName
func (p *AwsCloud) GetRegions() (cloud.GetRegionsResponse, error) {
	input := &ec2.DescribeRegionsInput{}
	result, err := p.ec2Client.DescribeRegions(input)
	if err != nil {
		logs.Logger.Errorf("GetRegions AwsCloud failed.err: [%v]", err)
		return cloud.GetRegionsResponse{}, nil
	}
	var regions = make([]cloud.Region, 0, len(result.Regions))
	for _, region := range result.Regions {
		regions = append(regions, cloud.Region{
			RegionId: *region.RegionName,
			//LocalName: ,
		})
	}
	return cloud.GetRegionsResponse{Regions: regions}, nil
}

//DescribeImages InsType isn't use
func (p *AwsCloud) DescribeImages(req cloud.DescribeImagesRequest) (cloud.DescribeImagesResponse, error) {
	var owners = make([]*string, 0, 1)
	var filters = make([]*ec2.Filter, 0, 1)
	var images = make([]cloud.Image, 0, _pageSize)
	input := &ec2.DescribeImagesInput{
		Owners: append(owners, aws.String(_imageType[req.ImageType])),
		Filters: append(filters, &ec2.Filter{
			Name:   aws.String("state"),
			Values: []*string{aws.String("available")},
		}),
	}
	result, err := p.ec2Client.DescribeImages(input)
	if err != nil {
		logs.Logger.Errorf("DescribeImages AwsCloud failed.err: [%v] req[%v]", err, req)
		return cloud.DescribeImagesResponse{}, nil
	}
	for _, image := range result.Images {
		images = append(images, cloud.Image{
			OsType:  *image.Platform,
			OsName:  *image.Name,
			ImageId: *image.ImageId,
		})
	}
	return cloud.DescribeImagesResponse{Images: images}, nil
}

func (p *AwsCloud) payOrders(orderId string) error {

	return nil
}

//up to 50 at once
func (p *AwsCloud) listPrePaidResources(ids []string) (map[string]prePaidResources, error) {

	return nil, nil
}

func (p *AwsCloud) GetOrders(req cloud.GetOrdersRequest) (cloud.GetOrdersResponse, error) {
	return cloud.GetOrdersResponse{}, nil
}
