package alibaba

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/bssopenapi"
	"github.com/galaxy-future/BridgX/pkg/cloud"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

func (p *AlibabaCloud) GetOrders(req cloud.GetOrdersRequest) (cloud.GetOrdersResponse, error) {
	request := bssopenapi.CreateQueryOrdersRequest()
	request.Scheme = "https"
	request.CreateTimeStart = req.StartTime.Format("2006-01-02T15:04:05Z")
	request.CreateTimeEnd = req.EndTime.Format("2006-01-02T15:04:05Z")
	request.PageNum = requests.NewInteger(req.PageNum)
	request.PageSize = requests.NewInteger(req.PageSize)
	response, err := p.bssClient.QueryOrders(request)
	if err != nil {
		return cloud.GetOrdersResponse{}, err
	}
	if !response.Success {
		return cloud.GetOrdersResponse{}, errors.New(response.Message)
	}
	orderNum := len(response.Data.OrderList.Order)
	if orderNum == 0 {
		return cloud.GetOrdersResponse{}, nil
	}

	orders := make([]cloud.Order, 0, orderNum*_subOrderNumPerMain)
	detailReq := bssopenapi.CreateGetOrderDetailRequest()
	detailReq.Scheme = "https"
	for _, row := range response.Data.OrderList.Order {
		detailReq.OrderId = row.OrderId
		detailRsp, err := p.bssClient.GetOrderDetail(detailReq)
		if err != nil {
			return cloud.GetOrdersResponse{}, err
		}
		if !detailRsp.Success {
			return cloud.GetOrdersResponse{}, errors.New(detailRsp.Message)
		}
		if len(detailRsp.Data.OrderList.Order) == 0 {
			continue
		}

		for _, subOrder := range detailRsp.Data.OrderList.Order {
			orderTime, _ := time.Parse("2006-01-02T15:04:05Z", subOrder.CreateTime)
			var usageStartTime, usageEndTime time.Time
			if subOrder.UsageStartTime == "" {
				usageStartTime = orderTime
			} else {
				usageStartTime, _ = time.Parse("2006-01-02T15:04:05Z", subOrder.UsageStartTime)
			}
			if subOrder.UsageEndTime == "" {
				usageEndTime = orderTime
			} else {
				usageEndTime, _ = time.Parse("2006-01-02T15:04:05Z", subOrder.UsageEndTime)
			}
			if _orderChargeType[subOrder.SubscriptionType] == cloud.PostPaid && usageEndTime.Sub(usageStartTime).Hours() > 24*365*20 {
				usageEndTime, _ = time.Parse("2006-01-02 15:04:05", "2038-01-01 00:00:00")
			}

			orders = append(orders, cloud.Order{
				OrderId:        subOrder.SubOrderId,
				OrderTime:      orderTime,
				Product:        subOrder.ProductCode,
				Quantity:       cast.ToInt32(subOrder.Quantity),
				UsageStartTime: usageStartTime,
				UsageEndTime:   usageEndTime,
				RegionId:       subOrder.Region,
				ChargeType:     _orderChargeType[subOrder.SubscriptionType],
				PayStatus:      _payStatus[subOrder.PaymentStatus],
				Currency:       subOrder.Currency,
				Cost:           cast.ToFloat32(subOrder.PretaxAmount),
				Extend: map[string]interface{}{
					"main_order_id": subOrder.OrderId,
					"order_type":    subOrder.OrderType,
				},
			})
		}
	}
	return cloud.GetOrdersResponse{Orders: orders}, nil
}
