package huawei

import "github.com/galaxy-future/BridgX/pkg/cloud"

func (p *HuaweiCloud) AllocateEip(req cloud.AllocateEipRequest) (ids []string, err error) {

	return nil, nil
}

func (p *HuaweiCloud) GetEips(ids []string, regionId string) (map[string]cloud.Eip, error) {

	return map[string]cloud.Eip{}, nil
}

func (p *HuaweiCloud) ReleaseEip(ids []string) (err error) {

	return nil
}

func (p *HuaweiCloud) AssociateEip(id, instanceId, vpcId string) error {

	return nil
}

func (p *HuaweiCloud) DisassociateEip(id string) error {

	return nil
}

func (p *HuaweiCloud) DescribeEip(req cloud.DescribeEipRequest) (cloud.DescribeEipResponse, error) {
	ret := cloud.DescribeEipResponse{}
	return ret, nil
}

func (p *HuaweiCloud) ConvertPublicIpToEip(req cloud.ConvertPublicIpToEipRequest) error {
	// return errors.New("not Implemented"), if not supported in huawei cloud
	return nil
}
