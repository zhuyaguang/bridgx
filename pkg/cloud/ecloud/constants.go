package ecloud

import "github.com/galaxy-future/BridgX/pkg/cloud"

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
