package ecloud

import (
	"github.com/galaxy-future/BridgX/pkg/cloud"
	"gitlab.ecloud.com/ecloud/ecloudsdkecs/model"
)

const (
	_payByBandwidth = "bandwidthCharge"
	_payByTraffic   = "trafficCharge"
	_ipTypeMobile   = "MOBILE"
	_productType    = "ip"
	_pageSize       = 100
	_State_OK       = "OK"
)

var _eipChargeType = map[string]string{
	cloud.BandwidthPayByTraffic: _payByBandwidth,
	cloud.BandwidthPayByFix:     _payByTraffic,
}

var _vmType = map[string]model.VmCreateBodyVmTypeEnum{
	"memImprove":         model.VmCreateBodyVmTypeEnumMemimprove,
	"common":             model.VmCreateBodyVmTypeEnumCommon,
	"gpu":                model.VmCreateBodyVmTypeEnumGpu,
	"commonIntroductory": model.VmCreateBodyVmTypeEnumCommonintroductory,
	"commonNetImprove":   model.VmCreateBodyVmTypeEnumCommonnetimprove,
	"compute":            model.VmCreateBodyVmTypeEnumCompute,
	"computeNetImprove":  model.VmCreateBodyVmTypeEnumComputenetimprove,
	"memNetImprove":      model.VmCreateBodyVmTypeEnumMemnetimprove,
	"localStorage":       model.VmCreateBodyVmTypeEnumLocalstorage,
	"xlargeMemory":       model.VmCreateBodyVmTypeEnumXlargememory,
	"highFrequency":      model.VmCreateBodyVmTypeEnumHighfrequency,
	"vgpu":               model.VmCreateBodyVmTypeEnumVgpu,
	"fpga":               model.VmCreateBodyVmTypeEnumFpga,
	"highIO":             model.VmCreateBodyVmTypeEnumHighio,
	"exclusive":          model.VmCreateBodyVmTypeEnumExclusive,
}
