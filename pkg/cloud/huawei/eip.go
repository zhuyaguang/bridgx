package huawei

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/galaxy-future/BridgX/internal/logs"
	"github.com/galaxy-future/BridgX/pkg/cloud"
	eip "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/eip/v2"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/eip/v2/model"
)

type allocateEipResult struct {
	Eip string
}

type allocateEip interface {
	handle(cli *eip.EipClient, req cloud.AllocateEipRequest) (allocateEipResult, error)
}

type allocateEipPayByTraffic struct {
}

// allocateEipPayByTraffic 按使用流量计费
func (eip allocateEipPayByTraffic) handle(cli *eip.EipClient, req cloud.AllocateEipRequest) (allocateEipResult, error) {
	chargeModeBandwidth := model.GetCreatePublicipBandwidthOptionChargeModeEnum().TRAFFIC
	shareType := model.GetCreatePublicipBandwidthOptionShareTypeEnum().PER
	size := int32(req.Bandwidth)
	request := &model.CreatePublicipRequest{}
	request.Body = &model.CreatePublicipRequestBody{
		Bandwidth: &model.CreatePublicipBandwidthOption{
			ChargeMode: &chargeModeBandwidth,
			Name:       &req.Name,
			ShareType:  shareType,
			Size:       &size,
		},
		Publicip: &model.CreatePublicipOption{
			Type: "5_bgp",
		},
	}
	response, err := cli.CreatePublicip(request)
	if err != nil {
		return allocateEipResult{}, err
	}
	if response.HttpStatusCode != http.StatusOK {
		return allocateEipResult{}, fmt.Errorf("httpcode %d", response.HttpStatusCode)
	}
	return allocateEipResult{
		Eip: *response.Publicip.PublicIpAddress,
	}, nil
}

type allocateEipPayByBandwidth struct {
}

// allocateEipPayByBandwidth 按固定宽带
func (eip allocateEipPayByBandwidth) handle(cli *eip.EipClient, req cloud.AllocateEipRequest) (allocateEipResult, error) {
	request := &model.CreatePublicipRequest{}
	chargeModeBandwidth := model.GetCreatePublicipBandwidthOptionChargeModeEnum().BANDWIDTH
	shareType := model.GetCreatePublicipBandwidthOptionShareTypeEnum().PER

	request.Body = &model.CreatePublicipRequestBody{
		Bandwidth: &model.CreatePublicipBandwidthOption{
			ChargeMode: &chargeModeBandwidth,
			Name:       &req.Name,
			ShareType:  shareType,
		},
		Publicip: &model.CreatePublicipOption{
			Type: "5_bgp",
		},
	}
	response, err := cli.CreatePublicip(request)
	if err != nil {
		return allocateEipResult{}, err
	}
	if response.HttpStatusCode != http.StatusOK {
		return allocateEipResult{}, fmt.Errorf("httpcode %d", response.HttpStatusCode)
	}
	return allocateEipResult{
		Eip: *response.Publicip.PublicIpAddress,
	}, nil
}

type allocateEipPrePaid struct {
}

// allocateEipPrePaid 预付费
func (eip allocateEipPrePaid) handle(cli *eip.EipClient, req cloud.AllocateEipRequest) (allocateEipResult, error) {
	autoPay := true
	autoRenew := true
	request := &model.CreatePrePaidPublicipRequest{}
	chargeModeBandwidth := model.GetCreatePublicipBandwidthOptionChargeModeEnum().BANDWIDTH
	chargeModePrePaid := model.GetCreatePrePaidPublicipExtendParamOptionChargeModeEnum().PRE_PAID
	shareType := model.GetCreatePublicipBandwidthOptionShareTypeEnum().PER
	periodType, periodNum, err := getPeriod(req.Charge)
	if err != nil {
		return allocateEipResult{}, err
	}
	request.Body = &model.CreatePrePaidPublicipRequestBody{
		Bandwidth: &model.CreatePublicipBandwidthOption{
			ChargeMode: &chargeModeBandwidth,
			Name:       &req.Name,
			ShareType:  shareType,
		},
		ExtendParam: &model.CreatePrePaidPublicipExtendParamOption{
			ChargeMode:  &chargeModePrePaid,
			PeriodType:  periodType,
			PeriodNum:   periodNum,
			IsAutoRenew: &autoRenew,
			IsAutoPay:   &autoPay,
		},
	}
	response, err := cli.CreatePrePaidPublicip(request)
	if err != nil {
		return allocateEipResult{}, err
	}
	if response.HttpStatusCode != http.StatusOK {
		return allocateEipResult{}, fmt.Errorf("httpcode %d", response.HttpStatusCode)
	}
	return allocateEipResult{
		Eip: *response.Publicip.PublicIpAddress,
	}, nil
}

func (p *HuaweiCloud) AllocateEip(req cloud.AllocateEipRequest) (ids []string, err error) {
	if req.Charge == nil {
		return nil, errors.New("invalid charge type")
	}

	var handler allocateEip
	switch req.Charge.ChargeType {
	case "PayByTraffic":
		handler = allocateEipPayByTraffic{}
	case "PayByBandwidth":
		handler = allocateEipPayByBandwidth{}
	case "PrePaid":
		handler = allocateEipPrePaid{}
	default:
		return nil, errors.New("invalid charge type")
	}

	idChan := make(chan string, req.Num)
	errChan := make(chan error, req.Num)

	for i := 0; i < req.Num; i++ {
		go func() {
			res, err1 := handler.handle(p.eipClient, req)
			if err1 != nil {
				logs.Logger.Errorf("AllocateEip HuaweiCloud failed.err:[%v] req:[%v]", err1, req)
				errChan <- err1
				return
			}

			idChan <- res.Eip
		}()
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

func getPeriod(change *cloud.Charge) (periodType *model.CreatePrePaidPublicipExtendParamOptionPeriodType, periodNum *int32, err error) {
	switch change.PeriodUnit {
	case "Month":
		periodTypeEnum := model.GetCreatePrePaidPublicipExtendParamOptionPeriodTypeEnum().MONTH
		periodType = &periodTypeEnum
	case "Year":
		periodTypeEnum := model.GetCreatePrePaidPublicipExtendParamOptionPeriodTypeEnum().YEAR
		periodType = &periodTypeEnum
	default:
		err = errors.New("invalid change PeriodUnit")
		return
	}
	if change.Period > 0 {
		periodNumInt32 := int32(change.Period)
		periodNum = &periodNumInt32
	}
	return
}

func (p *HuaweiCloud) GetEips(ids []string, regionId string) (map[string]cloud.Eip, error) {
	idNum := len(ids)
	eipMap := make(map[string]cloud.Eip, idNum)
	if idNum < 1 {
		return eipMap, nil
	}

	limit := int32(idNum)
	request := &model.ListPublicipsRequest{
		Limit:           &limit,
		PublicIpAddress: &ids,
	}
	response, err := p.eipClient.ListPublicips(request)
	if err != nil {
		return nil, err
	}
	if response.HttpStatusCode != http.StatusOK {
		return nil, fmt.Errorf("httpcode %d", response.HttpStatusCode)
	}
	for _, i := range *response.Publicips {
		eipMap[*i.Id] = eip2Cloud(i)
	}

	return eipMap, nil
}

func eip2Cloud(eip model.PublicipShowResp) cloud.Eip {
	cloudEip := cloud.Eip{
		Id: *eip.Id,
		Ip: *eip.PublicIpAddress,
	}
	if *eip.IpVersion == model.GetPublicipShowRespIpVersionEnum().E_6 {
		cloudEip.Ip = *eip.PublicIpv6Address
	}
	// TODO 这里赋值是否合适？
	if eip.Alias != nil {
		cloudEip.Name = *eip.Alias
	} else {
		cloudEip.Name = *eip.BandwidthName
	}
	if eip.Profile != nil && eip.Profile.ProductId != nil {
		cloudEip.InstanceId = *eip.Profile.ProductId
	}
	return cloudEip
}

func (p *HuaweiCloud) ReleaseEip(ids []string) (err error) {
	num := len(ids)
	if num < 1 {
		return nil
	}
	idChan := make(chan string, num)
	errChan := make(chan error, num)
	for _, id := range ids {
		err1 := p.DisassociateEip(id)
		if err1 != nil {
			errChan <- err1
			return
		}

		request := &model.DeletePublicipRequest{PublicipId: id}
		response, err2 := p.eipClient.DeletePublicip(request)
		if err1 != nil {
			errChan <- err2
			return
		}
		if response.HttpStatusCode != http.StatusOK {
			errChan <- fmt.Errorf("httpcode %d", response.HttpStatusCode)
			return
		}

		idChan <- id
	}

	for i := 0; i < num; i++ {
		select {
		case err = <-errChan:
		case <-idChan:
		}
	}

	return err
}

func (p *HuaweiCloud) AssociateEip(id, instanceId, vpcId string) error {
	request := &model.UpdatePublicipRequest{
		PublicipId: id,
		Body: &model.UpdatePublicipsRequestBody{
			Publicip: &model.UpdatePublicipOption{
				PortId: &instanceId,
			},
		},
	}
	response, err := p.eipClient.UpdatePublicip(request)
	if err != nil {
		return err
	}
	if response.HttpStatusCode != http.StatusOK {
		return fmt.Errorf("httpcode %d", response.HttpStatusCode)
	}
	return err
}

func (p *HuaweiCloud) DisassociateEip(id string) error {
	request := &model.UpdatePublicipRequest{
		PublicipId: id,
		Body: &model.UpdatePublicipsRequestBody{
			Publicip: &model.UpdatePublicipOption{},
		},
	}
	response, err := p.eipClient.UpdatePublicip(request)
	if err != nil {
		return err
	}
	if response.HttpStatusCode != http.StatusOK {
		return fmt.Errorf("httpcode %d", response.HttpStatusCode)
	}
	return err
}

func (p *HuaweiCloud) DescribeEip(req cloud.DescribeEipRequest) (cloud.DescribeEipResponse, error) {
	limit := int32(req.PageSize)
	request := &model.ListPublicipsRequest{
		Limit: &limit,
	}
	if req.OlderMarker != "" {
		request.Marker = &req.OlderMarker
	}

	response, err := p.eipClient.ListPublicips(request)
	if err != nil {
		return cloud.DescribeEipResponse{}, err
	}
	ret := cloud.DescribeEipResponse{}
	if response.HttpStatusCode != http.StatusOK {
		return ret, fmt.Errorf("httpcode %d", response.HttpStatusCode)
	}
	ret.TotalCount = len(*response.Publicips)
	for _, ip := range *response.Publicips {
		ret.List = append(ret.List, eip2Cloud(ip))
	}

	return ret, nil
}

func (p *HuaweiCloud) ConvertPublicIpToEip(req cloud.ConvertPublicIpToEipRequest) error {
	return errors.New("not Implemented")
}
