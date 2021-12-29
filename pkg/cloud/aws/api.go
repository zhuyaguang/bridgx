package huawei

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/galaxy-future/BridgX/pkg/cloud"
)

type AwsCloud struct {
	ec2 *ec2.EC2
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

// GetRegions 暂时返回中文名字
func (p *AwsCloud) GetRegions() (cloud.GetRegionsResponse, error) {

	return cloud.GetRegionsResponse{}, nil
}

func (p *AwsCloud) DescribeImages(req cloud.DescribeImagesRequest) (cloud.DescribeImagesResponse, error) {

	return cloud.DescribeImagesResponse{}, nil
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
