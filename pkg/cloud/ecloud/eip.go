package ecloud

import "github.com/galaxy-future/BridgX/pkg/cloud"

func (p *ECloud) AllocateEip(req cloud.AllocateEipRequest) (ids []string, err error) {
	// TODO implement me
	panic("implement me")
}

func (p *ECloud) GetEips(ids []string, regionId string) (map[string]cloud.Eip, error) {
	// TODO implement me
	panic("implement me")
}

func (p *ECloud) ReleaseEip(ids []string) (err error) {
	// TODO implement me
	panic("implement me")
}

func (p *ECloud) AssociateEip(id, instanceId, vpcId string) error {
	// TODO implement me
	panic("implement me")
}

func (p *ECloud) DisassociateEip(id string) error {
	// TODO implement me
	panic("implement me")
}

func (p *ECloud) DescribeEip(req cloud.DescribeEipRequest) (cloud.DescribeEipResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (p *ECloud) ConvertPublicIpToEip(req cloud.ConvertPublicIpToEipRequest) error {
	// TODO implement me
	panic("implement me")
}
