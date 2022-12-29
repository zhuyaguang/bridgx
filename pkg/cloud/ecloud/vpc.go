package ecloud

import (
	"errors"
	"github.com/galaxy-future/BridgX/pkg/cloud"
)

func (p *ECloud) CreateVPC(req cloud.CreateVpcRequest) (cloud.CreateVpcResponse, error) {
	// TODO implement me
	return cloud.CreateVpcResponse{}, errors.New("implement me")
}

func (p *ECloud) GetVPC(req cloud.GetVpcRequest) (cloud.GetVpcResponse, error) {
	// TODO implement me
	return cloud.GetVpcResponse{}, errors.New("implement me")
}

func (p *ECloud) CreateSwitch(req cloud.CreateSwitchRequest) (cloud.CreateSwitchResponse, error) {
	// TODO implement me
	return cloud.CreateSwitchResponse{}, errors.New("implement me")
}

func (p *ECloud) GetSwitch(req cloud.GetSwitchRequest) (cloud.GetSwitchResponse, error) {
	// TODO implement me
	return cloud.GetSwitchResponse{}, errors.New("implement me")
}

func (p *ECloud) DescribeVpcs(req cloud.DescribeVpcsRequest) (cloud.DescribeVpcsResponse, error) {
	// TODO implement me
	return cloud.DescribeVpcsResponse{}, errors.New("implement me")
}

func (p *ECloud) DescribeSwitches(req cloud.DescribeSwitchesRequest) (cloud.DescribeSwitchesResponse, error) {
	// TODO implement me
	return cloud.DescribeSwitchesResponse{}, errors.New("implement me")
}
