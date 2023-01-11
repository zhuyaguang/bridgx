package ecloud

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/galaxy-future/BridgX/internal/logs"
	"github.com/galaxy-future/BridgX/pkg/cloud"
	"gitlab.ecloud.com/ecloud/ecloudsdkecs/model"
)

func (p *ECloud) BatchCreate(m cloud.Params, num int) (instanceIds []string, err error) {
	request := &model.VmCreateRequest{}
	vmCreateBody := &model.VmCreateBody{
		BootVolume: &model.VmCreateRequestBootVolume{
			VolumeType: "",
			Size:       nil,
		}, Networks: &model.VmCreateRequestNetworks{
			NetworkId: "",
			PortId:    "",
		}}
	vmCreateBody.Region = m.Region
	Num := int32(num)
	vmCreateBody.Quantity = &Num

	if m.Charge.PeriodUnit == cloud.Year {
		vmCreateBody.BillingType = model.VmCreateBodyBillingTypeEnumYear
	} else if m.Charge.PeriodUnit == cloud.Month {
		vmCreateBody.BillingType = model.VmCreateBodyBillingTypeEnumMonth
	} else {
		vmCreateBody.BillingType = model.VmCreateBodyBillingTypeEnumHour
	}
	vmCreateBody.VmType = _vmType[m.InstanceType]
	vmCreateBody.Cpu = &Num
	vmCreateBody.Ram = &Num
	DiskNum := int32(m.Disks.DataDisk[0].Size)
	vmCreateBody.Disk = &DiskNum
	BootDiskNum := int32(m.Disks.SystemDisk.Size)
	vmCreateBody.BootVolume.Size = &BootDiskNum
	if m.Disks.SystemDisk.Category == "highPerformance" {
		vmCreateBody.BootVolume.VolumeType = model.VmCreateRequestBootVolumeVolumeTypeEnumHighperformance
	} else if m.Disks.SystemDisk.Category == "performanceOptimization" {
		vmCreateBody.BootVolume.VolumeType = model.VmCreateRequestBootVolumeVolumeTypeEnumPerformanceoptimization
	}
	vmCreateBody.ImageName = m.ImageId
	vmCreateBody.Networks.NetworkId = m.Network.VpcId
	// 云主机名 region + 时间戳
	vmCreateBody.Name = fmt.Sprintf(m.Region + "-" + strconv.FormatInt(time.Now().Unix(), 10))
	Duration := int32(m.Charge.Period)
	vmCreateBody.Duration = &Duration

	request.VmCreateBody = vmCreateBody

	response, err := p.ecsClient.VmCreate(request)
	if err != nil {
		logs.Logger.Errorf(err.Error())
		return []string{}, err
	}
	logs.Logger.Info("[BatchCreate] %+v\n", response.Body)

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
	b := true
	request := &model.VmGetServerDetailRequest{
		VmGetServerDetailPath: &model.VmGetServerDetailPath{
			ServerId: id},
		VmGetServerDetailQuery: &model.VmGetServerDetailQuery{
			Detail: &b},
	}

	response, err := p.ecsClient.VmGetServerDetail(request)
	if err != nil {
		return instances, err
	}
	if response.State != model.VmGetServerDetailResponseStateEnumOk {
		return instances, errors.New(response.ErrorMessage)
	}

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
		request := &model.VmDeleteRequest{
			VmDeletePath: &model.VmDeletePath{
				ServerId: id,
			},
		}
		response, err := p.ecsClient.VmDelete(request)
		if err != nil {
			return err
		}
		if response.State != model.VmDeleteResponseStateEnumOk {
			return errors.New(response.ErrorMessage)
		} else {
			logs.Logger.Infof("[BatchDelete] requestId: %s", response.RequestId)
		}

	}
	return nil
}

func (p *ECloud) StartInstances(ids []string) error {
	for _, id := range ids {
		request := &model.VmStartRequest{VmStartPath: &model.VmStartPath{
			ServerId: id,
		}}
		response, err := p.ecsClient.VmStart(request)
		if err != nil {
			return err
		}
		if response.State != model.VmStartResponseStateEnumOk {
			return errors.New(response.ErrorMessage)
		} else {
			logs.Logger.Infof("[StartInstances] requestId: %s", response.RequestId)
		}
	}
	return nil
}

func (p *ECloud) StopInstances(ids []string) error {
	for _, id := range ids {
		request := &model.VmStopRequest{VmStopPath: &model.VmStopPath{
			ServerId: id,
		}}
		response, err := p.ecsClient.VmStop(request)
		if err != nil {
			return err
		}
		if response.State != model.VmStopResponseStateEnumOk {
			return errors.New(response.ErrorMessage)
		} else {
			logs.Logger.Infof("[StopInstances] requestId: %s", response.RequestId)
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
