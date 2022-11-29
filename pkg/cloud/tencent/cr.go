package tencent

import (
	"github.com/galaxy-future/BridgX/internal/logs"
	"github.com/galaxy-future/BridgX/pkg/cloud"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	tcr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcr/v20190924"
	"strings"
)

//此处与阿里云对齐，各名称定义如下示例
//namespce : galaxy-future
//repository name : bridgX
//image : galaxy-future/bridgX:v0.01
//需检查各函数中的相关实现

func (p *TencentCloud) PersonalImageList(region, repoNamespace, repoName string, pageNumber, pageSize int) ([]cloud.DockerArtifact, int, error) {

	var DockerArtifactList []cloud.DockerArtifact

	request := tcr.NewDescribeImagePersonalRequest()
	ComRepoName := repoNamespace + "/" + repoName
	request.RepoName = &ComRepoName
	request.Offset, request.Limit = SizeNumConvert(pageNumber, pageSize)
	// 返回的resp是一个DescribeImagePersonalResponse的实例，与请求对象对应
	response, err := p.tcrClient.DescribeImagePersonal(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		logs.Logger.Errorf("An API error has returned: %s", err)
		return DockerArtifactList, 0, err
	}
	if err != nil {
		return DockerArtifactList, 0, err
	}

	for _, d := range response.Response.Data.TagInfo {
		var docker cloud.DockerArtifact
		docker.Name = repoNamespace + "/" + repoName + ":" + *d.TagName
		DockerArtifactList = append(DockerArtifactList, docker)
	}
	return DockerArtifactList, int(*response.Response.Data.TagCount), nil
}

func (p *TencentCloud) EnterpriseImageList(region, instanceId, repoId, namespace, repoName string, pageNumber, pageSize int) ([]cloud.DockerArtifact, int, error) {
	request := tcr.NewDescribeImagesRequest()
	request.RegistryId = &instanceId
	request.NamespaceName = &namespace
	request.RepositoryName = &repoName
	request.Offset, request.Limit = SizeNumConvert(pageNumber, pageSize)

	var DockerArtifactList []cloud.DockerArtifact
	// 返回的resp是一个DescribeImagesResponse的实例，与请求对象对应
	response, err := p.tcrClient.DescribeImages(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		logs.Logger.Errorf("An API error has returned: %s", err)
		return DockerArtifactList, 0, err
	}
	if err != nil {
		return DockerArtifactList, 0, err
	}

	for _, d := range response.Response.ImageInfoList {
		var docker cloud.DockerArtifact
		docker.Name = *d.Digest
		DockerArtifactList = append(DockerArtifactList, docker)
	}
	return DockerArtifactList, int(*response.Response.TotalCount), nil
}

func (p *TencentCloud) ContainerInstanceList(region string, pageNumber, pageSize int) ([]cloud.RegistryInstance, int, error) {

	request := tcr.NewDescribeInstancesRequest()
	request.Offset, request.Limit = SizeNumConvert(pageNumber, pageSize)
	var RegistryInstanceList []cloud.RegistryInstance

	response, err := p.tcrClient.DescribeInstances(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		logs.Logger.Errorf("An API error has returned: %s", err)
		return RegistryInstanceList, 0, err
	}
	if err != nil {
		return RegistryInstanceList, 0, err
	}

	for _, r := range response.Response.Registries {
		var registryInstance cloud.RegistryInstance
		registryInstance.InstanceId = *r.RegistryId
		registryInstance.InstanceName = *r.RegionName
		RegistryInstanceList = append(RegistryInstanceList, registryInstance)
	}

	return RegistryInstanceList, int(*response.Response.TotalCount), nil
}

func (p *TencentCloud) EnterpriseNamespaceList(region, instanceId string, pageNumber, pageSize int) ([]cloud.Namespace, int, error) {
	request := tcr.NewDescribeNamespacesRequest()
	var NamespaceList []cloud.Namespace
	request.RegistryId = common.StringPtr(instanceId)
	request.Offset, request.Limit = SizeNumConvert(pageNumber, pageSize)

	response, err := p.tcrClient.DescribeNamespaces(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		logs.Logger.Errorf("An API error has returned: %s", err)
		return NamespaceList, 0, err
	}
	if err != nil {
		return NamespaceList, 0, err
	}
	for _, n := range response.Response.NamespaceList {
		namespace := cloud.Namespace{Name: *n.Name}
		NamespaceList = append(NamespaceList, namespace)
	}
	return NamespaceList, int(*response.Response.TotalCount), nil
}

func (p *TencentCloud) PersonalNamespaceList(region string) ([]cloud.Namespace, error) {
	var NamespaceList []cloud.Namespace
	requestROP := tcr.NewDescribeRepositoryOwnerPersonalRequest()
	response, err := p.tcrClient.DescribeRepositoryOwnerPersonal(requestROP)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		logs.Logger.Errorf("An API error has returned: %s", err)
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	for _, r := range response.Response.Data.RepoInfo {
		namespace := ""
		name := *r.RepoName
		resultArr := strings.Split(name, "/")
		if len(resultArr) != 0 {
			namespace = resultArr[0]
		}
		NamespaceList = append(NamespaceList, cloud.Namespace{Name: namespace})
	}
	//request := tcr.NewDescribeNamespacePersonalRequest()
	//request.Namespace = &namespace
	//request.Offset, request.Limit = SizeNumConvert(0, 10)
	//response, err := p.tcrClient.DescribeNamespacePersonal(request)
	//if _, ok := err.(*errors.TencentCloudSDKError); ok {
	//	logs.Logger.Errorf("An API error has returned: %s", err)
	//	return NamespaceList, err
	//}
	//if err != nil {
	//	return NamespaceList, err
	//}
	//for _, n := range response.Response.Data.NamespaceInfo {
	//	namespace := cloud.Namespace{Name: *n.Namespace}
	//	NamespaceList = append(NamespaceList, namespace)
	//}

	return NamespaceList, nil
}

func (p *TencentCloud) EnterpriseRepositoryList(region, instanceId, namespace string, pageNumber, pageSize int) ([]cloud.Repository, int, error) {
	request := tcr.NewDescribeRepositoriesRequest()
	var RepositoryList []cloud.Repository
	request.RegistryId = common.StringPtr(instanceId)
	request.Offset, request.Limit = SizeNumConvert(pageNumber, pageSize)

	response, err := p.tcrClient.DescribeRepositories(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		logs.Logger.Errorf("An API error has returned: %s", err)
		return RepositoryList, 0, err
	}
	if err != nil {
		return RepositoryList, 0, err
	}
	for _, r := range response.Response.RepositoryList {
		repository := cloud.Repository{
			Name: *r.Name,
			ID:   "",
		}
		RepositoryList = append(RepositoryList, repository)
	}
	return RepositoryList, int(*response.Response.TotalCount), nil
}

func (p *TencentCloud) PersonalRepositoryList(region, namespace string, pageNumber, pageSize int) ([]cloud.Repository, int, error) {
	request := tcr.NewDescribeRepositoryOwnerPersonalRequest()
	request.Offset, request.Limit = SizeNumConvert(pageNumber, pageSize)
	var RepositoryList []cloud.Repository
	response, err := p.tcrClient.DescribeRepositoryOwnerPersonal(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		logs.Logger.Errorf("An API error has returned: %s", err)
		return RepositoryList, 0, err
	}
	if err != nil {
		return RepositoryList, 0, err
	}
	for _, r := range response.Response.Data.RepoInfo {
		repository := cloud.Repository{
			Name: *r.RepoName,
			ID:   "",
		}
		RepositoryList = append(RepositoryList, repository)
	}
	return RepositoryList, int(*response.Response.Data.TotalCount), nil
}

func SizeNumConvert(pageNumber, pageSize int) (offset, limit *int64) {
	// offset 是偏移量，与pageNumber 有区别
	var num, size int64
	size = int64(pageSize)
	num = (int64(pageNumber) - 1) * size
	return &num, &size
}
