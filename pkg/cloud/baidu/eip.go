package baidu

import (
	"github.com/baidubce/bce-sdk-go/services/eip"
	"github.com/galaxy-future/BridgX/internal/logs"
	"github.com/galaxy-future/BridgX/pkg/cloud"
	"github.com/pkg/errors"
)

func (b *BaiduCloud) AllocateEip(req cloud.AllocateEipRequest) (ids []string, err error) {
	name := ""
	if req.Name != "" {
		name = req.Name
	}
	billing := &eip.Billing{}
	if req.Charge.ChargeType == "PayByTraffic" {
		billing.BillingMethod = "ByTraffic"
		billing.PaymentTiming = "Postpaid"
	} else if req.Charge.ChargeType == "PayByBandwidth" {
		billing.BillingMethod = "ByBandwidth"
		billing.PaymentTiming = "Postpaid"
	} else if req.Charge.ChargeType == "PrePaid" {
		billing.Reservation = &eip.Reservation{}
		billing.PaymentTiming = "Prepaid"
		billing.Reservation.ReservationTimeUnit = "Month"
		if req.Charge.PeriodUnit == cloud.Year {
			billing.Reservation.ReservationLength = req.Charge.Period * 12
		} else {
			billing.Reservation.ReservationLength = req.Charge.Period
		}
	}
	// 创建EIP
	args := &eip.CreateEipArgs{
		Name:            name,
		BandWidthInMbps: req.Bandwidth,
		Billing:         billing,
	}
	idChan := make(chan string, req.Num)
	errChan := make(chan error, req.Num)

	for i := 0; i < req.Num; i++ {
		go func() {
			res, err := b.eipClient.CreateEip(args)
			if err != nil {
				logs.Logger.Errorf("AllocateEip BaiduCloud failed.err:[%v] req:[%v]", err, req)
				errChan <- err
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
	//ids = append(ids, res.Eip)
	return ids, err
}

func (b *BaiduCloud) GetEips(ids []string, regionId string) (map[string]cloud.Eip, error) {
	res := map[string]cloud.Eip{}
	for _, id := range ids {
		args := &eip.ListEipArgs{
			Eip: id,
		}
		result, err := b.eipClient.ListEip(args)
		if err != nil {
			return nil, err
		}
		if len(result.EipList) == 0 {
			return res, nil
		}
		tmp := &cloud.Eip{
			Id:         result.EipList[0].EipId,
			Name:       result.EipList[0].Name,
			Ip:         result.EipList[0].Eip,
			InstanceId: result.EipList[0].InstanceId,
		}
		res[id] = *tmp
	}

	return res, nil
}

func (b *BaiduCloud) ReleaseEip(ids []string) (err error) {

	for _, eip := range ids {
		err = b.eipClient.DeleteEip(eip, "")
		if err != nil {
			break
		}
	}
	return
}

func (b *BaiduCloud) AssociateEip(id, instanceId, vpcId string) error {
	args := &eip.BindEipArgs{
		InstanceType: "BCC",
		InstanceId:   instanceId,
	}
	err := b.eipClient.BindEip(id, args)
	if err != nil {
		return err
	}
	return nil
}

func (b *BaiduCloud) DisassociateEip(id string) error {
	//id  == eip
	err := b.eipClient.UnBindEip(id, "")
	if err != nil {
		return err
	}
	return nil
}

func (b *BaiduCloud) DescribeEip(req cloud.DescribeEipRequest) (cloud.DescribeEipResponse, error) {
	args := &eip.ListEipArgs{
		InstanceId: req.InstanceId,
		Status:     "available",
		Marker:     req.OlderMarker,
		MaxKeys:    req.PageSize,
	}
	result, _ := b.eipClient.ListEip(args)
	list := []cloud.Eip{}
	for _, eip := range result.EipList {
		e := cloud.Eip{
			Id:         eip.Eip,
			Name:       eip.Name,
			Ip:         eip.Eip,
			InstanceId: eip.InstanceId,
		}
		list = append(list, e)
	}

	return cloud.DescribeEipResponse{
		List:       list,
		NewMarker:  result.NextMarker,
		TotalCount: 0,
	}, nil
}

func (b *BaiduCloud) ConvertPublicIpToEip(req cloud.ConvertPublicIpToEipRequest) error {
	return errors.New("not Implemented") // do not support in baidu cloud
}
