package alibaba

import (
	"fmt"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/bssopenapi"
)

var (
	_AK = "test_ak"
	_SK = "test_sk"
)

func TestDescribeInstanceBill(t *testing.T) {
	config := sdk.NewConfig()
	credential := credentials.NewAccessKeyCredential(_AK, _SK)
	client, err := bssopenapi.NewClientWithOptions("cn-beijing", config, credential)
	if err != nil {
		panic(err)
	}

	now := time.Now()
	duration, _ := time.ParseDuration("-24h")
	request := bssopenapi.CreateDescribeInstanceBillRequest()
	request.Scheme = "https"
	request.BillingCycle = now.Format("2006-01")
	//request.IsBillingItem = requests.NewBoolean(true)
	request.IsHideZeroCharge = requests.NewBoolean(true)
	request.Granularity = "DAILY"
	request.ProductCode = "ecs"
	request.BillingDate = now.Add(duration).Format("2006-01-02")
	request.MaxResults = requests.NewInteger(300) // aliyun max limit
	response, err := client.DescribeInstanceBill(request)
	if err != nil {
		fmt.Printf("response is %#v\n", response)
	}
	totalCount := response.Data.TotalCount
	rem := 0
	if (totalCount-300)%300 != 0 {
		rem = 1
	}
	for i := 0; i < (totalCount-300)/300+rem; i++ {
		request.NextToken = response.Data.NextToken
		nextResponse, err := client.DescribeInstanceBill(request)
		if err != nil {
			fmt.Printf("response is %#v\n", response)
		}
		response.Data.Items = append(response.Data.Items, nextResponse.Data.Items...)
	}
	fmt.Printf("response is %#v\n", len(response.Data.Items))
}

func TestQueryInstanceBill(t *testing.T) {
	config := sdk.NewConfig()
	credential := credentials.NewAccessKeyCredential("_AK", "_SK")
	client, err := bssopenapi.NewClientWithOptions("cn-beijing", config, credential)
	if err != nil {
		panic(err)
	}

	duration, _ := time.ParseDuration("-24h")
	date := time.Now().Add(duration)
	request := bssopenapi.CreateQueryInstanceBillRequest()
	request.Scheme = "https"
	request.BillingCycle = date.Format("2006-01")
	request.BillingDate = date.Format("2006-01-02")
	//request.IsBillingItem = requests.NewBoolean(true)
	request.IsHideZeroCharge = requests.NewBoolean(true)
	request.Granularity = "DAILY"
	request.ProductCode = "ecs"
	request.PageNum = requests.NewInteger(1)
	request.PageSize = requests.NewInteger(300) // aliyun max limit
	response, err := client.QueryInstanceBill(request)
	if err != nil {
		fmt.Printf("response is %#v\n", response)
	}
	totalCount := response.Data.TotalCount
	rem := 0
	if (totalCount-300)%300 != 0 {
		rem = 1
	}
	for i := 0; i < 3+rem; i++ {
		request.PageNum = requests.NewInteger(i + 2)
		nextResponse, err := client.QueryInstanceBill(request)
		if err != nil {
			fmt.Printf("response is %#v\n", response)
		}
		response.Data.Items.Item = append(response.Data.Items.Item, nextResponse.Data.Items.Item...)
	}
	fmt.Printf("response is %#v\n", len(response.Data.Items.Item))
}

func TestQueryInstanceBill2(t *testing.T) {
	config := sdk.NewConfig()
	credential := credentials.NewAccessKeyCredential("_AK", "_SK")
	client, err := bssopenapi.NewClientWithOptions("cn-beijing", config, credential)
	if err != nil {
		panic(err)
	}

	now := time.Now()
	duration, _ := time.ParseDuration("-24h")
	request := bssopenapi.CreateQueryInstanceBillRequest()
	request.Scheme = "https"
	request.BillingCycle = now.Format("2006-01")
	//request.IsBillingItem = requests.NewBoolean(true)
	request.IsHideZeroCharge = requests.NewBoolean(true)
	request.Granularity = "DAILY"
	request.BillingDate = now.Add(duration).Format("2006-01-02")
	//request.MaxResults = requests.NewInteger(300) // aliyun max limit
	response, err := client.QueryInstanceBill(request)
	if err != nil {
		fmt.Print(err.Error())
	}
	fmt.Printf("response is %#v\n", response)
}
