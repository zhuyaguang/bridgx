package huawei

import (
	"github.com/galaxy-future/BridgX/pkg/cloud"
)

func (p *AwsCloud) BatchCreate(m cloud.Params, num int) ([]string, error) {

	return nil, nil
}

func (p *AwsCloud) GetInstances(ids []string) (instances []cloud.Instance, err error) {

	return nil, nil
}

func (p *AwsCloud) GetInstancesByTags(regionId string, tags []cloud.Tag) (instances []cloud.Instance, err error) {
	return nil, nil
}

func (p *AwsCloud) GetInstancesByCluster(regionId, clusterName string) (instances []cloud.Instance, err error) {
	return p.GetInstancesByTags(regionId, []cloud.Tag{{
		Key:   cloud.ClusterName,
		Value: clusterName,
	}})
}

// BatchDelete 华为云限制一次最多操作_maxNumEcsPerOperation台
func (p *AwsCloud) BatchDelete(ids []string, regionId string) error {

	return nil
}

func (p *AwsCloud) StartInstances(ids []string) error {

	return nil
}

func (p *AwsCloud) StopInstances(ids []string) error {

	return nil
}

// GetZones 华为云无ZoneId字段用ZoneName填充
func (p *AwsCloud) GetZones(req cloud.GetZonesRequest) (cloud.GetZonesResponse, error) {
	return cloud.GetZonesResponse{}, nil
}

func (p *AwsCloud) DescribeAvailableResource(req cloud.DescribeAvailableResourceRequest) (cloud.DescribeAvailableResourceResponse, error) {
	return cloud.DescribeAvailableResourceResponse{}, nil
}

//DescribeInstanceTypes NovaShowFlavor 华为云还没实现
func (p *AwsCloud) DescribeInstanceTypes(req cloud.DescribeInstanceTypesRequest) (cloud.DescribeInstanceTypesResponse, error) {
	return cloud.DescribeInstanceTypesResponse{}, nil
}
