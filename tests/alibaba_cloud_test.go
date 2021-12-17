package tests

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/bssopenapi"
	"github.com/galaxy-future/BridgX/pkg/cloud"
	"github.com/galaxy-future/BridgX/pkg/cloud/alibaba"
	"github.com/stretchr/testify/assert"
)

func getAlibabaClient() (*alibaba.AlibabaCloud, error) {
	c, err := alibaba.New("ak", "sk", "cn-beijing")
	if err != nil {
		return nil, err
	}

	return c, nil
}

func TestDescribeAvailableResource(t *testing.T) {
	cli, err := getAlibabaClient()
	if err != nil {
		t.Log(err.Error())
		return
	}

	var res interface{}
	var resStr []byte

	res, err = cli.DescribeAvailableResource(cloud.DescribeAvailableResourceRequest{RegionId: "cn-beijing", ZoneId: ""})
	if err != nil {
		t.Log(err.Error())
		return
	}
	resStr, _ = json.Marshal(res)
	t.Log(string(resStr))

	res, err = cli.DescribeInstanceTypes(cloud.DescribeInstanceTypesRequest{TypeName: []string{"ecs.g6.large"}})
	if err != nil {
		t.Log(err.Error())
		return
	}
	resStr, _ = json.Marshal(res)
	t.Log(string(resStr))

	res, err = cli.DescribeImages(cloud.DescribeImagesRequest{RegionId: "cn-beijing", InsType: "ecs.g6.large"})
	if err != nil {
		t.Log(err.Error())
		return
	}
	resStr, _ = json.Marshal(res)
	t.Log(string(resStr))
}

func TestQueryOrders(t *testing.T) {
	cli, err := getAlibabaClient()
	if err != nil {
		t.Log(err.Error())
		return
	}

	//endTime := time.Now().UTC()
	//duration, _ := time.ParseDuration("-5h")
	//startTime := endTime.Add(duration)
	startTime, _ := time.Parse("2006-01-02 15:04:05", "2021-11-19 11:40:02")
	endTime, _ := time.Parse("2006-01-02 15:04:05", "2021-11-19 11:45:02")
	pageNum := 1
	pageSize := 100
	for {
		res, err := cli.GetOrders(cloud.GetOrdersRequest{StartTime: startTime, EndTime: endTime,
			PageNum: pageNum, PageSize: pageSize})
		if err != nil {
			t.Log(err.Error())
			return
		}
		cnt := 0
		t.Log("len:", len(res.Orders))
		for _, row := range res.Orders {
			cnt += 1
			if cnt > 3 {
				t.Log("---------------")
				break
			}
			rowStr, _ := json.Marshal(row)
			t.Log(string(rowStr))
		}
		if len(res.Orders) < pageSize {
			break
		}
		pageNum += 1
	}
	t.Log(pageNum)
}

func TestGetOrderDetail(t *testing.T) {
	client, err := bssopenapi.NewClientWithAccessKey("cn-beijing", "a", "b")
	if err != nil {
		t.Log(err.Error())
		return
	}
	request := bssopenapi.CreateGetOrderDetailRequest()
	request.Scheme = "https"
	request.OrderId = "211577282350149"
	response, err := client.GetOrderDetail(request)
	if err != nil {
		t.Log(err.Error())
		return
	}

	orders, err := json.Marshal(response.Data.OrderList)
	if err != nil {
		t.Log(err.Error())
		return
	}
	t.Log(string(orders))
}

func TestInstanceExpireTimeParse(t *testing.T) {
	expireAt, err := time.Parse("2006-01-02T15:04:05Z", "2099-11-01T01:03:04Z")
	assert.Nil(t, err)
	t.Logf("expire at:%v", expireAt)
	expireAt, err = time.Parse("2006-01-02T15:04:05Z", "2099x-11-01T01:03:04Z")
	assert.NotNil(t, err)
	t.Logf("expire at:%v", expireAt)
	var tt *time.Time
	assert.Nil(t, tt)
	t.Logf("tt:%v", tt)
}
