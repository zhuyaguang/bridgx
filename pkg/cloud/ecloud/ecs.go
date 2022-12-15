package ecloud

import (
	"errors"
	"github.com/galaxy-future/BridgX/pkg/cloud"
)

func (p *ECloud) BatchCreate(m cloud.Params, num int) (instanceIds []string, err error) {
	// TODO implement me
	return []string{}, errors.New("implement me")
}

func (p *ECloud) GetInstances(ids []string) (instances []cloud.Instance, err error) {
	// TODO implement me
	return []cloud.Instance{}, errors.New("implement me")
}

func (p *ECloud) GetInstancesByTags(region string, tags []cloud.Tag) (instances []cloud.Instance, err error) {
	// TODO implement me
	return []cloud.Instance{}, errors.New("implement me")
}

func (p *ECloud) GetInstancesByCluster(regionId, clusterName string) (instances []cloud.Instance, err error) {
	// TODO implement me
	return []cloud.Instance{}, errors.New("implement me")
}

func (p *ECloud) BatchDelete(ids []string, regionId string) error {
	// TODO implement me
	return errors.New("implement me")
}

func (p *ECloud) StartInstances(ids []string) error {
	// TODO implement me
	return errors.New("implement me")
}

func (p *ECloud) StopInstances(ids []string) error {
	// TODO implement me
	return errors.New("implement me")
}

func (p *ECloud) GetRegions() (cloud.GetRegionsResponse, error) {
	// TODO implement me
	return cloud.GetRegionsResponse{}, errors.New("implement me")
}

func (p *ECloud) GetZones(req cloud.GetZonesRequest) (cloud.GetZonesResponse, error) {
	// TODO implement me
	return cloud.GetZonesResponse{}, errors.New("implement me")
}

func (p *ECloud) DescribeAvailableResource(req cloud.DescribeAvailableResourceRequest) (cloud.DescribeAvailableResourceResponse, error) {
	// TODO implement me
	return cloud.DescribeAvailableResourceResponse{}, errors.New("implement me")
}

func (p *ECloud) DescribeInstanceTypes(req cloud.DescribeInstanceTypesRequest) (cloud.DescribeInstanceTypesResponse, error) {
	// TODO implement me
	return cloud.DescribeInstanceTypesResponse{}, errors.New("implement me")
}
