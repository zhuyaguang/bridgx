package cloud

type Provider interface {
	BatchCreate(m Params, num int) (instanceIds []string, err error)
	ProviderType() string
	GetInstances(ids []string) (instances []Instance, err error)
	GetInstancesByTags(region string, tags []Tag) (instances []Instance, err error)
	GetInstancesByCluster(regionId, clusterName string) (instances []Instance, err error)
	BatchDelete(ids []string, regionId string) error
	StartInstances(ids []string) error
	StopInstances(ids []string) error
	CreateVPC(req CreateVpcRequest) (CreateVpcResponse, error)
	GetVPC(req GetVpcRequest) (GetVpcResponse, error)
	CreateSwitch(req CreateSwitchRequest) (CreateSwitchResponse, error)
	GetSwitch(req GetSwitchRequest) (GetSwitchResponse, error)
	CreateSecurityGroup(req CreateSecurityGroupRequest) (CreateSecurityGroupResponse, error)
	AddIngressSecurityGroupRule(req AddSecurityGroupRuleRequest) error
	AddEgressSecurityGroupRule(req AddSecurityGroupRuleRequest) error
	DescribeSecurityGroups(req DescribeSecurityGroupsRequest) (DescribeSecurityGroupsResponse, error)
	GetRegions() (GetRegionsResponse, error)
	GetZones(req GetZonesRequest) (GetZonesResponse, error)
	DescribeAvailableResource(req DescribeAvailableResourceRequest) (DescribeAvailableResourceResponse, error)
	DescribeInstanceTypes(req DescribeInstanceTypesRequest) (DescribeInstanceTypesResponse, error)
	DescribeImages(req DescribeImagesRequest) (DescribeImagesResponse, error)
	DescribeVpcs(req DescribeVpcsRequest) (DescribeVpcsResponse, error)
	DescribeSwitches(req DescribeSwitchesRequest) (DescribeSwitchesResponse, error)
	DescribeGroupRules(req DescribeGroupRulesRequest) (DescribeGroupRulesResponse, error)
	// order
	GetOrders(req GetOrdersRequest) (GetOrdersResponse, error)
	// key pairs
	CreateKeyPair(req CreateKeyPairRequest) (CreateKeyPairResponse, error)
	ImportKeyPair(req ImportKeyPairRequest) (ImportKeyPairResponse, error)
	DescribeKeyPairs(req DescribeKeyPairsRequest) (DescribeKeyPairsResponse, error)
	AllocateEip(req AllocateEipRequest) (ids []string, err error)
	GetEips(ids []string, regionId string) (map[string]Eip, error)
	// eip
	ReleaseEip(ids []string) (err error)
	AssociateEip(id, instanceId, vpcId string) error
	DisassociateEip(id string) error
	DescribeEip(req DescribeEipRequest) (DescribeEipResponse, error)
	ConvertPublicIpToEip(req ConvertPublicIpToEipRequest) error
	// s3
	ListObjects(endpoint, bucketName, prefix string) ([]ObjectProperties, error)
	ListBucket(endpoint string) ([]BucketProperties, error)
	GetOssDownloadUrl(string, string, string) string
	GetObjectDownloadUrl(bucketName, objectKey string) (string, error)

	// container registry
	ContainerInstanceList(region string, pageNumber, pageSize int) ([]RegistryInstance, int, error)
	EnterpriseNamespaceList(region, instanceId string, pageNumber, pageSize int) ([]Namespace, int, error)
	PersonalNamespaceList(region string) ([]Namespace, error)
	EnterpriseRepositoryList(region, instanceId, namespace string, pageNumber, pageSize int) ([]Repository, int, error)
	PersonalRepositoryList(region, namespace string, pageNumber, pageSize int) ([]Repository, int, error)
	EnterpriseImageList(region, instanceId, repoId, namespace, repoName string, pageNumber, pageSize int) ([]DockerArtifact, int, error)
	PersonalImageList(region, repoNamespace, repoName string, pageNum, pageSize int) ([]DockerArtifact, int, error)
}

type ProviderDriverFunc func(keyId ...string) (Provider, error)

var registeredPlugins = map[string]ProviderDriverFunc{}

func RegisterProviderDriver(name string, f ProviderDriverFunc) {
	registeredPlugins[name] = f
}
