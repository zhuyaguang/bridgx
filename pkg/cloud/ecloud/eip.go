package ecloud

import (
	"fmt"
	"github.com/galaxy-future/BridgX/internal/logs"
	"github.com/galaxy-future/BridgX/pkg/cloud"
	"github.com/galaxy-future/BridgX/pkg/utils"
	"gitlab.ecloud.com/ecloud/ecloudsdkeip/model"
)

func (p *ECloud) AllocateEip(req cloud.AllocateEipRequest) (ids []string, err error) {
	if req.Charge == nil {
		return nil, fmt.Errorf("charge is nil")
	}
	allocateEipAddressRequest := &model.FloatingIpOrderCreateRequest{
		FloatingIpOrderCreateBody: &model.FloatingIpOrderCreateBody{
			IpType:           _ipTypeMobile,
			Quantity:         utils.Int32(int32(req.Num)),
			BandwidthSize:    utils.Int32(int32(req.Bandwidth)),
			ChargeModeEnum:   model.FloatingIpOrderCreateBodyChargeModeEnumEnum(_eipChargeType[req.Charge.ChargeType]),
			ChargePeriodEnum: model.FloatingIpOrderCreateBodyChargePeriodEnumEnum(req.Charge.PeriodUnit),
			Duration:         utils.Int32(int32(req.Charge.Period)),
		},
	}
	idChan := make(chan string, req.Num)
	errChan := make(chan error, req.Num)
	for i := 0; i < req.Num; i++ {
		go func(req *model.FloatingIpOrderCreateRequest) {
			rsp, err := p.eipClient.FloatingIpOrderCreate(req)
			if err != nil {
				logs.Logger.Errorf("AllocateEip Ecloud failed.err:[%v] req:[%v]", err, req)
				errChan <- err
				return
			}
			if rsp == nil || rsp.Body == nil {
				errChan <- fmt.Errorf("AllocateEip Ecloud resp nil req:[%v]", req)
				return
			}
			if rsp.State != _State_OK {
				errChan <- fmt.Errorf("AllocateEip Ecloud resp state not ok:[%v] req:[%v]", rsp.ErrorMessage, req)
				return
			}
			if rsp.Body != nil {
				idChan <- rsp.Body.OrderId
			}
		}(allocateEipAddressRequest)
	}

	for i := 0; i < req.Num; i++ {
		select {
		case err = <-errChan:
		case id := <-idChan:
			ids = append(ids, id)
		}
	}
	return ids, err
}

func (p *ECloud) GetEips(ids []string, regionId string) (map[string]cloud.Eip, error) {
	idNum := len(ids)
	eipMap := make(map[string]cloud.Eip, idNum)
	if idNum < 1 {
		return eipMap, nil
	}

	for _, id := range ids {
		getFipWithBandwidthRequest := &model.GetFipWithBandwidthRequest{
			GetFipWithBandwidthPath: &model.GetFipWithBandwidthPath{
				IpId: id,
			},
		}
		rsp, err := p.eipClient.GetFipWithBandwidth(getFipWithBandwidthRequest)
		if err != nil {
			return eipMap, err
		}
		if rsp == nil || rsp.Body == nil {
			return eipMap, fmt.Errorf("GetEips Ecloud resp nil req:[%v]", id)
		}
		if rsp.State != _State_OK {
			return eipMap, fmt.Errorf("GetEips Ecloud resp state not ok:[%v] req:[%v]", rsp.ErrorMessage, id)
		}
		if rsp.Body != nil {
			eipMap[utils.StringValue(&rsp.Body.Id)] = eip2Ecloud(rsp.Body)
		}
	}
	return eipMap, nil
}

func (p *ECloud) ReleaseEip(ids []string) (err error) {
	num := len(ids)
	if num < 1 {
		return nil
	}
	idChan := make(chan string, num)
	errChan := make(chan error, num)
	for _, id := range ids {
		go func(id string) {
			getFipWithBandwidthRequest := &model.GetFipWithBandwidthRequest{
				GetFipWithBandwidthPath: &model.GetFipWithBandwidthPath{
					IpId: id,
				},
			}
			fipWithBandwidth, err1 := p.eipClient.GetFipWithBandwidth(getFipWithBandwidthRequest)
			if err1 != nil {
				errChan <- err1
				return
			}
			if fipWithBandwidth == nil || fipWithBandwidth.Body == nil {
				errChan <- fmt.Errorf("GetFipWithBandwidth Ecloud resp nil req:[%v]", id)
				return
			}
			if fipWithBandwidth.State != _State_OK {
				errChan <- fmt.Errorf("GetFipWithBandwidth Ecloud resp state not ok:[%v] req:[%v]", fipWithBandwidth.ErrorMessage, id)
				return
			}
			resourceId := fipWithBandwidth.Body.BandwidthId
			relatedResourceId := fipWithBandwidth.Body.Id
			commonMopOrderDeleteIpRequest := &model.CommonMopOrderDeleteIpRequest{
				CommonMopOrderDeleteIpBody: &model.CommonMopOrderDeleteIpBody{
					ProductType:       _productType,
					ResourceId:        resourceId,
					RelatedResourceId: relatedResourceId,
				},
			}
			rsp, err := p.eipClient.CommonMopOrderDeleteIp(commonMopOrderDeleteIpRequest)
			if err != nil {
				errChan <- err
				return
			}
			if rsp == nil || rsp.Body == nil {
				errChan <- fmt.Errorf("CommonMopOrderDeleteIp Ecloud resp nil req:[%v]", id)
				return
			}
			if rsp.State != _State_OK {
				errChan <- fmt.Errorf("CommonMopOrderDeleteIp Ecloud resp state not ok:[%v] req:[%v]", fipWithBandwidth.ErrorMessage, id)
				return
			}
			idChan <- id
		}(id)
	}

	for i := 0; i < num; i++ {
		select {
		case err = <-errChan:
		case <-idChan:
		}
	}
	return err
}

func (p *ECloud) AssociateEip(id, instanceId, vpcId string) error {
	floatingIpBindRequest := &model.FloatingIpBindRequest{
		FloatingIpBindBody: &model.FloatingIpBindBody{
			IpId:       id,
			ResourceId: instanceId,
			PortId:     vpcId,
		},
	}
	rsp, err := p.eipClient.FloatingIpBind(floatingIpBindRequest)
	if rsp == nil {
		return fmt.Errorf("AssociateEip Ecloud resp nil id:[%v] instanceId:[%v] vpcId:[%v]", id, instanceId, vpcId)

	}
	if rsp.OpenApiReturnValue != _State_OK {
		return fmt.Errorf("AssociateEip Ecloud resp OpenApiReturnValue not ok id:[%v] instanceId:[%v] vpcId:[%v]", id, instanceId, vpcId)
	}
	if err != nil {
		return err
	}
	return err
}

func (p *ECloud) DisassociateEip(id string) error {
	floatingIpUnbindRequest := &model.FloatingIpUnbindRequest{
		FloatingIpUnbindPath: &model.FloatingIpUnbindPath{
			IpId: id,
		},
	}
	rsp, err := p.eipClient.FloatingIpUnbind(floatingIpUnbindRequest)
	if err != nil {
		return err
	}
	if rsp == nil {
		return fmt.Errorf("DisassociateEip Ecloud resp nil id:[%v]", id)

	}
	if rsp.OpenApiReturnValue != _State_OK {
		return fmt.Errorf("DisassociateEip Ecloud resp OpenApiReturnValue not ok id:[%v]", id)
	}
	return err
}

func (p *ECloud) DescribeEip(req cloud.DescribeEipRequest) (cloud.DescribeEipResponse, error) {
	if req.PageNum < 1 {
		req.PageNum = 1
	}
	if req.PageSize < 1 || req.PageSize > _pageSize {
		req.PageSize = _pageSize
	}

	listFipWithBandwidthRequest := &model.ListFipWithBandwidthRequest{
		ListFipWithBandwidthQuery: &model.ListFipWithBandwidthQuery{
			Page: utils.Int32(int32(req.PageNum)),
			Size: utils.Int32(int32(req.PageSize)),
		},
	}
	rsp, err := p.eipClient.ListFipWithBandwidth(listFipWithBandwidthRequest)
	if err != nil {
		return cloud.DescribeEipResponse{}, err
	}

	if rsp == nil || rsp.Body == nil {
		return cloud.DescribeEipResponse{}, fmt.Errorf("DescribeEip Ecloud resp nil req:[%v]", req)
	}
	if rsp.State != _State_OK {
		return cloud.DescribeEipResponse{}, fmt.Errorf("DescribeEip Ecloud resp state not ok:[%v] req:[%v]", rsp.ErrorMessage, req)
	}
	if rsp.Body != nil {
		ret := cloud.DescribeEipResponse{
			List:       []cloud.Eip{},
			TotalCount: int(utils.Int32Value(rsp.Body.Total)),
		}
		for _, v := range *rsp.Body.Content {
			ret.List = append(ret.List, eip2Cloud(&v))
		}
		return ret, nil
	}
	return cloud.DescribeEipResponse{}, nil
}

func (p *ECloud) ConvertPublicIpToEip(req cloud.ConvertPublicIpToEipRequest) error {
	// TODO implement me
	panic("implement me")
}

func eip2Cloud(eip *model.ListFipWithBandwidthResponseContent) cloud.Eip {
	return cloud.Eip{
		Id:         eip.Id,
		Name:       eip.NicName,
		Ip:         eip.Name,
		InstanceId: eip.ResourceId,
	}
}

func eip2Ecloud(eip *model.GetFipWithBandwidthResponseBody) cloud.Eip {
	return cloud.Eip{
		Id:         eip.Id,
		Name:       eip.NicName,
		Ip:         eip.Name,
		InstanceId: eip.ResourceId,
	}
}
