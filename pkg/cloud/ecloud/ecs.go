package ecloud

import (
	"errors"
	"github.com/galaxy-future/BridgX/internal/logs"
	"github.com/galaxy-future/BridgX/pkg/cloud"
	"gitlab.ecloud.com/ecloud/ecloudsdkecs/model"
	"strings"
	"time"
)

func (p *ECloud) BatchCreate(m cloud.Params, num int) (instanceIds []string, err error) {
	// 参数不对齐 num 循环？
	request := &model.VmCreateRequest{}
	vmCreateBody := &model.VmCreateBody{}
	vmCreateBody.Region = m.Region

	if m.Charge.PeriodUnit == cloud.Year {
		vmCreateBody.BillingType = model.VmCreateBodyBillingTypeEnumYear
	} else if m.Charge.PeriodUnit == cloud.Month {
		vmCreateBody.BillingType = model.VmCreateBodyBillingTypeEnumMonth
	} else {
		vmCreateBody.BillingType = model.VmCreateBodyBillingTypeEnumHour
	}
	vmCreateBody.VmType = _vmType[m.InstanceType]

	request.VmCreateBody = vmCreateBody

	response, err := p.ecsClient.VmCreate(request)
	if err != nil {
		logs.Logger.Errorf(err.Error())
		return []string{}, err
	}
	logs.Logger.Info("%+v\n", response.Body)

	return instanceIds, nil
}

func (p *ECloud) GetInstances(ids []string) (instances []cloud.Instance, err error) {
	for _, id := range ids {
		ins, err := p.generateInstances(id)
		if err != nil {
			return []cloud.Instance{}, err
		}
		instances = append(instances, ins)
	}
	return instances, nil
}

func (p *ECloud) generateInstances(id string) (instances cloud.Instance, err error) {
	request := &model.VmGetServerDetailRequest{}
	vmGetServerDetailPath := &model.VmGetServerDetailPath{}
	VmGetServerDetailQuery := &model.VmGetServerDetailQuery{}
	vmGetServerDetailPath.ServerId = id
	b := true
	VmGetServerDetailQuery.Detail = &b
	request.VmGetServerDetailPath = vmGetServerDetailPath
	request.VmGetServerDetailQuery = VmGetServerDetailQuery

	response, err := p.ecsClient.VmGetServerDetail(request)
	if err != nil {
		return instances, err
	}
	if response.State == model.VmGetServerDetailResponseStateEnumOk {
		expireAt, err := time.Parse("2006-01-02T15:04Z", response.Body.CreatedTime)
		if err != nil {
			return instances, err
		}
		ipInner := ""
		for _, p := range *response.Body.Ports {
			ipInner = ipInner + "," + strings.Join(p.PrivateIp, ",")
		}
		ipOuter := ""
		if len(*response.Body.Ports) > 0 {
			if len((*response.Body.Ports)[0].PublicIp) > 0 {
				ipOuter = (*response.Body.Ports)[0].PublicIp[0]
			}
		}

		instances = cloud.Instance{
			Id:       response.Body.Id,
			CostWay:  "",
			Provider: cloud.ECloud,
			IpInner:  ipInner,
			IpOuter:  ipOuter,
			Network:  nil,
			ImageId:  response.Body.ImageId,
			Status:   string(*response.Body.Status),
			ExpireAt: &expireAt,
		}
		return instances, nil
	} else {
		err := errors.New(response.ErrorMessage)
		return instances, err
	}
}

func (p *ECloud) GetInstanceStatus(id string) (status string, err error) {
	ins, err := p.generateInstances(id)
	if err != nil {
		return "", err
	}
	return ins.Status, nil
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
	for _, id := range ids {
		request := &model.VmDeleteRequest{}
		vmDeletePath := &model.VmDeletePath{}
		vmDeletePath.ServerId = id
		request.VmDeletePath = vmDeletePath
		response, err := p.ecsClient.VmDelete(request)
		if err != nil {
			return err
		}

		if response.State == model.VmDeleteResponseStateEnumOk {
			logs.Logger.Info("BatchDelete %s", id)
		} else {
			err := errors.New(response.ErrorMessage)
			return err
		}

	}
	return nil
}

func (p *ECloud) StartInstances(ids []string) error {
	for _, id := range ids {
		request := &model.VmStartRequest{}
		vmStartPath := &model.VmStartPath{}
		vmStartPath.ServerId = id
		request.VmStartPath = vmStartPath
		response, err := p.ecsClient.VmStart(request)
		if err != nil {
			return err
		}
		if response.State == model.VmStartResponseStateEnumOk {
			logs.Logger.Info("StartInstances %s", id)
		} else {
			err := errors.New(response.ErrorMessage)
			return err
		}

	}
	return nil
}

func (p *ECloud) StopInstances(ids []string) error {
	for _, id := range ids {
		request := &model.VmStopRequest{}
		vmStopPath := &model.VmStopPath{}
		vmStopPath.ServerId = id
		request.VmStopPath = vmStopPath
		response, err := p.ecsClient.VmStop(request)
		if err != nil {
			return err
		}
		if response.State == model.VmStopResponseStateEnumOk {
			logs.Logger.Info("StopInstances %s", id)
		} else {
			err := errors.New(response.ErrorMessage)
			return err
		}
	}
	return nil
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
