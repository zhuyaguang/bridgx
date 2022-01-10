package aws

import (
	"errors"

	"github.com/galaxy-future/BridgX/pkg/cloud"
)

const (
	_maxNumEcsPerOperation = 100
	_pageSize              = 100

	_filterNameLocation     = "location"
	_locationTypeNameRegion = "region"
	_locationTypeNameZone   = "availability-zone"
	_locationTypeNameZoneId = "availability-zone-id"
)

var (
	_errInstanceIdsEmpty    = errors.New("at least one instance id")
	_errCodeDryRunOperation = "DryRunOperation"
)

var _imageType = map[string]string{
	cloud.ImageGlobal:  "amazon",
	cloud.ImageShared:  "aws-marketplace",
	cloud.ImagePrivate: "self",
}

var _ecsStatus = map[string]string{
	"pending":       cloud.EcsBuilding,
	"running":       cloud.EcsRunning,
	"shutting-down": cloud.EcsStopped,
	"terminated":    cloud.EcsDeleted,
	"stopping":      cloud.EcsStopping,
	"stopped":       cloud.EcsStopped,
}

var _insTypeStat = map[string]string{
	"normal":    cloud.InsTypeAvailable,
	"promotion": cloud.InsTypeAvailable,
	"":          cloud.InsTypeAvailable,
	"obt":       cloud.InsTypeAvaSoon,
	"abandon":   cloud.InsTypeSellOut,
	"sellout":   cloud.InsTypeSellOut,
}

var _secGrpRuleDirection = map[bool]string{
	false: cloud.SecGroupRuleIn,
	true:  cloud.SecGroupRuleOut,
}

var _vpcStatus = map[string]string{
	"pending":   cloud.VPCStatusPending,
	"available": cloud.VPCStatusAvailable,
}

var _subnetStatus = map[string]string{
	"pending":   cloud.VPCStatusPending,
	"available": cloud.VPCStatusAvailable,
}

var _subnetIsDefault = map[bool]int{
	false: 0,
	true:  1,
}
