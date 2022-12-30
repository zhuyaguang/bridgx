package ecloud

import (
	"github.com/galaxy-future/BridgX/pkg/cloud"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	_AK       = ""
	_SK       = ""
	_regionId = "CIDC-RP-29"
)

func TestEcloud_AllocateEip(t *testing.T) {
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
					//Period:     1,
					//PeriodUnit: "year",
				},
				Bandwidth: 1,
				Num:       1,
				RegionId:  "CIDC-RP-29",
			}},
		},
		//{
		//	name: "PayByBandwidth",
		//	args: args{req: cloud.AllocateEipRequest{
		//		Name: "按需计算-宽带",
		//		Charge: &cloud.Charge{
		//			ChargeType: "PayByBandwidth",
		//			Period:     10,
		//			PeriodUnit: "month",
		//		},
		//		Bandwidth: 2,
		//		Num:       1,
		//		RegionId:  "CIDC-RP-02",
		//	}},
		//},
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

func TestEcloud_ReleaseEip(t *testing.T) {
	p, err := New(_AK, _SK, _regionId)
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}
	err = p.ReleaseEip([]string{"8929eb44-f481-4cf8-adab-ea480267079c"})
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}
}

func TestEcloud_GetEips(t *testing.T) {
	p, err := New(_AK, _SK, _regionId)
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}
	response, err := p.GetEips([]string{"117.78.40.234"}, "CIDC-RP-29")
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}
	if assert.Equal(t, 3, len(response)) {
		for _, eip := range response {
			assert.NotEmpty(t, eip.Id)
			assert.NotEmpty(t, eip.Ip)
			t.Log(eip.Id, eip.Ip)
		}
	}
}

func TestEcloud_AssociateEip(t *testing.T) {
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

func TestEcloud_DisassociateEip(t *testing.T) {
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

func TestEcloud_DescribeEip(t *testing.T) {
	p, err := New(_AK, _SK, _regionId)
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}
	response, err := p.DescribeEip(cloud.DescribeEipRequest{
		PageNum:  1,
		PageSize: 10,
		RegionId: "CIDC-RP-29",
	})
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}
	t.Log(response.List)
	assert.Equal(t, response.TotalCount, len(response.List))
}
