package cloud

import (
	"time"
)

type Params struct {
	Provider     string
	InstanceType string
	ImageId      string
	Network      *Network
	Zone         string
	Region       string
	Disks        *Disks
	Charge       *Charge
	Password     string
	Tags         []Tag
	DryRun       bool
	KeyPairId    string
	KeyPairName  string
}

type Tag struct {
	Key   string
	Value string
}

type Instance struct {
	Id       string     `json:"id"`
	CostWay  string     `json:"cost_way"`
	Provider string     `json:"provider"`
	IpInner  string     `json:"ip_inner"`
	IpOuter  string     `json:"ip_outer"`
	Network  *Network   `json:"network"`
	ImageId  string     `json:"image_id"`
	Status   string     `json:"status"`
	ExpireAt *time.Time `json:"expire_at"`
}

type Network struct {
	VpcId                   string `json:"vpc_id"`
	SubnetId                string `json:"subnet_id"`
	SecurityGroup           string `json:"security_group"`
	InternetChargeType      string `json:"internet_charge_type"`
	InternetMaxBandwidthOut int    `json:"internet_max_bandwidth_out"`
	InternetIpType          string `json:"internet_ip_type"`
}

type Charge struct {
	Period     int    `json:"period"`
	PeriodUnit string `json:"period_unit"`
	ChargeType string `json:"charge_type"`
}

type Disks struct {
	SystemDisk DiskConf   `json:"system_disk"`
	DataDisk   []DiskConf `json:"data_disk"`
}

type DiskConf struct {
	Category         string `json:"category"`
	Size             int    `json:"size"`
	PerformanceLevel string `json:"performance_level"`
}

type CreateVpcRequest struct {
	RegionId  string
	VpcName   string
	CidrBlock string
}

type CreateVpcResponse struct {
	//RouterId       *string `json:"VRouterId,omitempty" xml:"VRouterId,omitempty"`
	//RouteTableId    *string `json:"RouteTableId,omitempty" xml:"RouteTableId,omitempty"`
	VpcId     string
	RequestId string
	//ResourceGroupId *string `json:"ResourceGroupId,omitempty" xml:"ResourceGroupId,omitempty"`
}

type GetVpcRequest struct {
	VpcId    string
	RegionId string
	VpcName  string
}

type GetVpcResponse struct {
	Vpc VPC
}

type VPC struct {
	VpcId     string
	VpcName   string
	CidrBlock string
	RegionId  string
	Status    string
	CreateAt  string
}

type CreateSwitchRequest struct {
	RegionId    string
	ZoneId      string
	CidrBlock   string
	VSwitchName string
	VpcId       string
	GatewayIp   string
}
type CreateSwitchResponse struct {
	SwitchId  string
	RequestId string
}

type CreateSecurityGroupRequest struct {
	RegionId          string
	SecurityGroupName string
	VpcId             string
	SecurityGroupType string
}

type CreateSecurityGroupResponse struct {
	SecurityGroupId string
	RequestId       string
}

type AddSecurityGroupRuleRequest struct {
	RegionId        string
	VpcId           string
	SecurityGroupId string
	IpProtocol      string
	PortFrom        int
	PortTo          int
	GroupId         string
	CidrIp          string
	PrefixListId    string
}

type DescribeSecurityGroupsRequest struct {
	VpcId    string
	RegionId string
}

type DescribeSecurityGroupsResponse struct {
	Groups []SecurityGroup
}

type SecurityGroup struct {
	SecurityGroupId   string
	SecurityGroupType string
	SecurityGroupName string
	CreateAt          string
	VpcId             string
	RegionId          string
}
type AddSecurityGroupRuleResponse struct {
}

type GetSwitchRequest struct {
	SwitchId string
}

type Switch struct {
	VpcId                   string
	SwitchId                string
	Name                    string
	IsDefault               int
	AvailableIpAddressCount int
	VStatus                 string
	CreateAt                string
	ZoneId                  string
	CidrBlock               string
	GatewayIp               string
}

type GetSwitchResponse struct {
	Switch Switch
}

type GetRegionsResponse struct {
	Regions []Region
}

type Region struct {
	RegionId  string
	LocalName string
}

type GetZonesRequest struct {
	RegionId string
}

type GetZonesResponse struct {
	Zones []Zone
}

type Zone struct {
	ZoneId    string
	LocalName string
}

type InstanceType struct {
	ChargeType  string `json:"charge_type"`
	IsGpu       bool   `json:"is_gpu"`
	Core        int    `json:"core"`
	Memory      int    `json:"memory"`
	Family      string `json:"instance_type_family"`
	InsTypeName string `json:"instance_type"`
	Status      string `json:"status"`
}

type DescribeAvailableResourceRequest struct {
	RegionId string
	ZoneId   string
}

type DescribeAvailableResourceResponse struct {
	InstanceTypes map[string][]InstanceType
}

type AvailableZone struct {
	ZoneId string
	Status string
}

type DescribeInstanceTypesRequest struct {
	TypeName []string
}
type DescribeInstanceTypesResponse struct {
	Infos []InstanceType
}
type DescribeImagesRequest struct {
	RegionId  string
	InsType   string
	ImageType string
}

type DescribeImagesResponse struct {
	Images []Image
}

type Image struct {
	Platform  string `json:"platform"`
	OsType    string `json:"os_type"`
	OsName    string `json:"os_name"`
	Size      int    `json:"size"` //GB
	ImageId   string `json:"image_id"`
	ImageName string `json:"image_name"`
}

type DescribeVpcsRequest struct {
	RegionId string
}

type DescribeVpcsResponse struct {
	Vpcs []VPC
}

type DescribeSwitchesRequest struct {
	VpcId string
}

type DescribeSwitchesResponse struct {
	Switches []Switch
}

type SecurityGroupRule struct {
	VpcId           string
	SecurityGroupId string
	PortFrom        int
	PortTo          int
	Protocol        string
	Direction       string
	GroupId         string
	CidrIp          string
	PrefixListId    string
	CreateAt        string
}

type DescribeGroupRulesRequest struct {
	RegionId        string
	SecurityGroupId string
}

type DescribeGroupRulesResponse struct {
	Rules []SecurityGroupRule
}

type GetOrdersRequest struct {
	StartTime time.Time
	EndTime   time.Time
	PageNum   int
	PageSize  int
}

type GetOrdersResponse struct {
	Orders []Order
}

type Order struct {
	OrderId        string
	OrderTime      time.Time
	Product        string
	Quantity       int32
	UsageStartTime time.Time
	UsageEndTime   time.Time
	RegionId       string
	ChargeType     string
	PayStatus      int8
	Currency       string
	Cost           float32
	Extend         map[string]interface{}
}
type CreateKeyPairRequest struct {
	RegionId    string
	KeyPairName string
}

type CreateKeyPairResponse struct {
	KeyPairId   string
	KeyPairName string
	PrivateKey  string
	PublicKey   string
}

type ImportKeyPairRequest struct {
	RegionId    string
	KeyPairName string
	PublicKey   string
}

type ImportKeyPairResponse struct {
	KeyPairId   string
	KeyPairName string
}

type DescribeKeyPairsRequest struct {
	RegionId    string
	PageNumber  int
	OlderMarker string
	PageSize    int
}

type DescribeKeyPairsResponse struct {
	TotalCount int
	KeyPairs   []KeyPair
	NewMarker  string
}

type KeyPair struct {
	KeyPairId   string
	KeyPairName string
}

type AllocateEipRequest struct {
	RegionId                string
	Name                    string
	InternetServiceProvider string
	Bandwidth               int
	Charge                  *Charge
	Num                     int
}

type DescribeEipRequest struct {
	RegionId    string
	InstanceId  string
	PageNum     int
	OlderMarker string
	PageSize    int
}
type DescribeEipResponse struct {
	List       []Eip  `json:"list"`
	NewMarker  string `json:"new_maker"`
	TotalCount int    `json:"total_count"`
}
type Eip struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Ip         string `json:"ip"`
	InstanceId string `json:"instance_id,omitempty"`
}

type ConvertPublicIpToEipRequest struct {
	RegionId   string
	InstanceId string
}

type ObjectProperties struct {
	Name string `json:"name"`
}

type BucketProperties struct {
	Name string `json:"name"`
}

type AcrInstanceListResponse struct {
	EnterpriseContainerCommon
	Instances []RegistryInstance `json:"Instances"`
}

type EnterpriseContainerCommon struct {
	TotalCount int    `json:"TotalCount"`
	PageSize   int    `json:"PageSize"`
	PageNum    int    `json:"PageNo"`
	Code       string `json:"Code"`
}

type RegistryInstance struct {
	InstanceName string `json:"InstanceName"`
	InstanceId   string `json:"InstanceId"`
}

type Namespace struct {
	Name string `json:"namespace"`
}

type EnterpriseNamespaceListResponse struct {
	EnterpriseContainerCommon
	Namespaces []EnterpriseNamespace `json:"Namespaces"`
}

type EnterpriseNamespace struct {
	NamespaceName string `json:"NamespaceName"`
	NamespaceId   string `json:"NamespaceId"`
}

type EnterpriseRepositoryListResponse struct {
	EnterpriseContainerCommon
	Repositories []EnterpriseRepository `json:"repositories"`
}

type Repository struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type EnterpriseRepository struct {
	RepoNamespaceName string `json:"RepoNamespaceName"`
	RepoName          string `json:"RepoName"`
	RepoId            string `json:"RepoId"`
}

type EnterpriseImageListResponse struct {
	EnterpriseContainerCommon
	Images []EnterpriseImage `json:"Images"`
}

type EnterpriseImage struct {
	Status string `json:"Status"`
	Tag    string `json:"Tag"`
}

type PersonalNamespaceListResponse struct {
	Data PersonalNamespaceListData `json:"data"`
}

type PersonalNamespaceListData struct {
	Namespaces []PersonalNamespace `json:"namespaces"`
}

type PersonalNamespace struct {
	Namespace string `json:"namespace"`
}

type PersonalRepositoryListResponse struct {
	Data PersonalRepositoryListData `json:"data"`
}

type PersonalRepositoryListData struct {
	Total    int            `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"pageSize"`
	Repos    []PersonalRepo `json:"repos"`
}

type PersonalRepo struct {
	RepoName string `json:"repoName"`
}

type DockerArtifact struct {
	Name string `json:"name"`
}

type TagsResponse struct {
	Data TagData `json:"data"`
}

type TagData struct {
	Total    int        `json:"total"`
	Page     int        `json:"page"`
	PageSize int        `json:"page_size"`
	Tags     []ImageTag `json:"tags"`
}

type ImageTag struct {
	Tag    string `json:"tag"`
	Status string `json:"status"`
}
