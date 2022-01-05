package aws

import (
	"errors"
	"time"

	"github.com/galaxy-future/BridgX/pkg/cloud"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2/model"
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
	_errInstanceIdsEmpty = errors.New("at least one instance id")
)

type prePaidResources struct {
	Id            string
	Name          string
	RegionId      string
	EffectiveTime time.Time //资源生效时间
	ExpireTime    time.Time //资源过期时间
	ExpirePolicy  int       //资源到期后的扣费策略: 1:到期转按需2:到期后自动删除(从生效中直接删除)3:到期后自动续费4:到期后冻结5:到期后删除(从保留期删除)
	Status        int       //2:使用中3:已关闭(页面不展示这个状态)4:已冻结5:已过期
}

//in
var _inEcsChargeType = map[string]model.PrePaidServerExtendParamChargingMode{
	cloud.InstanceChargeTypePrePaid:  model.GetPrePaidServerExtendParamChargingModeEnum().PRE_PAID,
	cloud.InstanceChargeTypePostPaid: model.GetPrePaidServerExtendParamChargingModeEnum().POST_PAID,
}

var _ecsPeriodType = map[string]model.PrePaidServerExtendParamPeriodType{
	"Month": model.GetPrePaidServerExtendParamPeriodTypeEnum().MONTH,
	"Year":  model.GetPrePaidServerExtendParamPeriodTypeEnum().YEAR,
}

var _imageType = map[string]string{
	cloud.ImageGlobal:  "amazon",
	cloud.ImageShared:  "aws-marketplace",
	cloud.ImagePrivate: "self",
}

var _rootDiskCategory = map[string]model.PrePaidServerRootVolumeVolumetype{
	"SATA":  model.GetPrePaidServerRootVolumeVolumetypeEnum().SATA,
	"SAS":   model.GetPrePaidServerRootVolumeVolumetypeEnum().SAS,
	"SSD":   model.GetPrePaidServerRootVolumeVolumetypeEnum().SSD,
	"GPSSD": model.GetPrePaidServerRootVolumeVolumetypeEnum().GPSSD,
	"CO_P1": model.GetPrePaidServerRootVolumeVolumetypeEnum().CO_P1,
	"UH_L1": model.GetPrePaidServerRootVolumeVolumetypeEnum().UH_L1,
}

var _dataDiskCategory = map[string]model.PrePaidServerDataVolumeVolumetype{
	"SATA":  model.GetPrePaidServerDataVolumeVolumetypeEnum().SATA,
	"SAS":   model.GetPrePaidServerDataVolumeVolumetypeEnum().SAS,
	"SSD":   model.GetPrePaidServerDataVolumeVolumetypeEnum().SSD,
	"GPSSD": model.GetPrePaidServerDataVolumeVolumetypeEnum().GPSSD,
	"CO_P1": model.GetPrePaidServerDataVolumeVolumetypeEnum().CO_P1,
	"UH_L1": model.GetPrePaidServerDataVolumeVolumetypeEnum().UH_L1,
}

var _bandwidthChargeMode = map[string]string{
	cloud.BandwidthPayByTraffic: "traffic",
	cloud.BandwidthPayByFix:     "",
}

var _protocol = map[string]string{
	cloud.ProtocolIcmp:   "icmp",
	cloud.ProtocolIcmpV6: "icmpv6",
	cloud.ProtocolTcp:    "tcp",
	cloud.ProtocolUdp:    "udp",
	cloud.ProtocolAll:    "",
}

//out
var _ecsChargeType = map[string]string{
	"0": cloud.InstanceChargeTypePostPaid,
	"1": cloud.InstanceChargeTypePrePaid,
}

var _ecsStatus = map[string]string{
	"pending":       cloud.EcsBuilding,
	"running":       cloud.EcsRunning,
	"shutting-down": cloud.EcsStopped,
	"terminated":    cloud.EcsAbnormal,
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

var _osType = map[string]string{
	"\"Linux\"\n":   cloud.OsLinux,
	"\"Windows\"\n": cloud.OsWindows,
	"\"Other\"\n":   cloud.OsOther,
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
