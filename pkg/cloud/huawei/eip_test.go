package huawei

import (
	"testing"

	"github.com/galaxy-future/BridgX/internal/logs"
	"github.com/galaxy-future/BridgX/pkg/cloud"
	"github.com/stretchr/testify/assert"
)

var (
	_AK       = ""
	_SK       = ""
	_regionId = "cn-north-1"
)

func init() {
	logs.Init()
}

func TestHuaweiCloud_AllocateEip(t *testing.T) {

	type args struct {
		req cloud.AllocateEipRequest
	}
	tests := []struct {
		name    string
		args    args
		wantIds []string
		wantErr bool
	}{
		{
			name: "PayByTraffic",
			args: args{req: cloud.AllocateEipRequest{
				Name: "按需计算-流量",
				Charge: &cloud.Charge{
					ChargeType: "PayByTraffic",
				},
				Bandwidth: 5,
				Num:       1,
				RegionId:  "cn-north-4",
			}},
		},
		{
			name: "PayByBandwidth",
			args: args{req: cloud.AllocateEipRequest{
				Name: "按需计算-宽带",
				Charge: &cloud.Charge{
					ChargeType: "PayByTraffic",
				},
				Bandwidth: 2,
				Num:       1,
				RegionId:  "cn-north-4",
			}},
		},
		{
			name: "PrePaid",
			args: args{req: cloud.AllocateEipRequest{
				Name: "包周期性",
				Charge: &cloud.Charge{
					ChargeType: "PrePaid",
					PeriodUnit: "Month",
				},
				Bandwidth: 2,
				Num:       1,
				RegionId:  "cn-north-4",
			}},
		},
	}
	p, err := New(_AK, _SK, _regionId)
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIds, err := p.AllocateEip(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("AllocateEip() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("AllocateEip() gotIds = %v", gotIds)
		})
	}
}

func TestHuaweiCloud_GetEips(t *testing.T) {
	p, err := New(_AK, _SK, _regionId)
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}
	response, err := p.GetEips([]string{"117.78.40.234"}, "cn-north-4")
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}
	if assert.Equal(t, 1, len(response)) {
		for _, eip := range response {
			assert.NotEmpty(t, eip.Name)
			assert.NotEmpty(t, eip.Id)
			assert.NotEmpty(t, eip.Ip)
			t.Log(eip.Id)
		}
	}
}

func TestHuaweiCloud_DescribeEip(t *testing.T) {
	p, err := New(_AK, _SK, _regionId)
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}
	response, err := p.DescribeEip(cloud.DescribeEipRequest{
		PageSize: 10,
		RegionId: "cn-north-4",
	})
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}
	t.Log(response.List)
	assert.Equal(t, response.TotalCount, len(response.List))
}

func TestHuaweiCloud_AssociateEip(t *testing.T) {
	p, err := New(_AK, _SK, _regionId)
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}
	err = p.AssociateEip("de5828e3-df17-472c-a9b8-766664ace268", "de5828e3-df17-472c-a9b8-766664ace268", "")
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}
}

func TestHuaweiCloud_DisassociateEip(t *testing.T) {
	p, err := New(_AK, _SK, _regionId)
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}
	err = p.DisassociateEip("de5828e3-df17-472c-a9b8-766664ace268")
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}
}

func TestHuaweiCloud_ReleaseEip(t *testing.T) {
	p, err := New(_AK, _SK, _regionId)
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}
	err = p.ReleaseEip([]string{"8d5b672c-9ab3-4f78-95c2-0fa2d53d4ece"})
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}
}
