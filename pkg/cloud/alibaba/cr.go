package alibaba

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/galaxy-future/BridgX/pkg/cloud"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cast"
)

func (p *AlibabaCloud) ContainerInstanceList(region string, pageNumber, pageSize int) ([]cloud.RegistryInstance, int, error) {
	request := getEnterpriseContainerRequest(region, "POST")
	request.ApiName = "ListInstance"
	request.QueryParams["InstanceStatus"] = "RUNNING"
	request.QueryParams["PageNo"] = cast.ToString(pageNumber)
	request.QueryParams["PageSize"] = cast.ToString(pageSize)
	apiRes, err := p.sdkClient.ProcessCommonRequest(request)
	instances := make([]cloud.RegistryInstance, 0)
	if err != nil {
		return instances, 0, err
	}
	res := cloud.AcrInstanceListResponse{}
	err = jsoniter.Unmarshal(apiRes.GetHttpContentBytes(), &res)
	if err != nil {
		return instances, 0, err
	}
	return res.Instances, res.TotalCount, nil
}

func (p *AlibabaCloud) EnterpriseNamespaceList(region, instanceId string, pageNumber, pageSize int) ([]cloud.Namespace, int, error) {
	request := getEnterpriseContainerRequest(region, "POST")
	request.ApiName = "ListNamespace"
	request.QueryParams["InstanceId"] = instanceId
	request.QueryParams["PageNo"] = cast.ToString(pageNumber)
	request.QueryParams["PageSize"] = cast.ToString(pageSize)
	request.QueryParams["NamespaceStatus"] = "NORMAL"
	apiRes, err := p.sdkClient.ProcessCommonRequest(request)
	namespaces := make([]cloud.Namespace, 0)
	if err != nil {
		return namespaces, 0, err
	}
	res := cloud.EnterpriseNamespaceListResponse{}
	err = jsoniter.Unmarshal(apiRes.GetHttpContentBytes(), &res)
	if err != nil {
		return namespaces, 0, err
	}
	for _, namespace := range res.Namespaces {
		namespaces = append(namespaces, cloud.Namespace{
			Name: namespace.NamespaceName,
		})
	}
	return namespaces, res.TotalCount, nil
}

func (p *AlibabaCloud) EnterpriseRepositoryList(region, instanceId, namespace string, pageNumber, pageSize int) ([]cloud.Repository, int, error) {
	request := getEnterpriseContainerRequest(region, "POST")
	request.ApiName = "ListRepository"
	request.QueryParams["InstanceId"] = instanceId
	request.QueryParams["RepoStatus"] = "NORMAL"
	request.QueryParams["NamespaceName"] = namespace
	request.QueryParams["PageNo"] = cast.ToString(pageNumber)
	request.QueryParams["PageSize"] = cast.ToString(pageSize)
	apiRes, err := p.sdkClient.ProcessCommonRequest(request)
	repositories := make([]cloud.Repository, 0)
	if err != nil {
		return repositories, 0, err
	}
	res := cloud.EnterpriseRepositoryListResponse{}
	err = jsoniter.Unmarshal(apiRes.GetHttpContentBytes(), &res)
	if err != nil {
		return repositories, 0, err
	}
	for _, r := range res.Repositories {
		repositories = append(repositories, cloud.Repository{
			Name: r.RepoName,
			ID:   r.RepoId,
		})
	}
	return repositories, res.TotalCount, nil
}

func (p *AlibabaCloud) EnterpriseImageList(region, instanceId, repoId, namespace, repoName string, pageNumber, pageSize int) ([]cloud.DockerArtifact, int, error) {
	request := getEnterpriseContainerRequest(region, "POST")
	request.ApiName = "ListRepoTag"
	request.QueryParams["InstanceId"] = instanceId
	request.QueryParams["RepoId"] = repoId
	request.QueryParams["PageNo"] = cast.ToString(pageNumber)
	request.QueryParams["PageSize"] = cast.ToString(pageSize)
	apiRes, err := p.sdkClient.ProcessCommonRequest(request)
	images := make([]cloud.DockerArtifact, 0)
	if err != nil {
		return images, 0, err
	}
	res := cloud.EnterpriseImageListResponse{}
	err = jsoniter.Unmarshal(apiRes.GetHttpContentBytes(), &res)
	if err != nil {
		return images, 0, err
	}
	for _, tag := range res.Images {
		images = append(images, cloud.DockerArtifact{
			Name: fmt.Sprintf("/%s/%s:%s", namespace, repoName, tag.Tag),
		})
	}
	return images, res.TotalCount, nil
}

func (p *AlibabaCloud) PersonalNamespaceList(region string) ([]cloud.Namespace, error) {
	request := getPersonalContainerRequest(region, "GET")
	request.PathPattern = "/namespace"
	request.QueryParams["Status"] = "NORMAL"
	apiRes, err := p.sdkClient.ProcessCommonRequest(request)
	namespaces := make([]cloud.Namespace, 0)
	if err != nil {
		return namespaces, err
	}
	res := cloud.PersonalNamespaceListResponse{}
	err = jsoniter.Unmarshal(apiRes.GetHttpContentBytes(), &res)
	if err != nil {
		return namespaces, err
	}
	for _, namespace := range res.Data.Namespaces {
		namespaces = append(namespaces, cloud.Namespace{
			Name: namespace.Namespace,
		})
	}
	return namespaces, nil
}

func (p *AlibabaCloud) PersonalRepositoryList(region, namespace string, pageNumber, pageSize int) ([]cloud.Repository, int, error) {
	request := getPersonalContainerRequest(region, "GET")
	request.PathPattern = fmt.Sprintf("/repos/%s", namespace)
	request.QueryParams["Status"] = "NORMAL"
	request.QueryParams["Page"] = cast.ToString(pageNumber)
	request.QueryParams["PageSize"] = cast.ToString(pageSize)
	apiRes, err := p.sdkClient.ProcessCommonRequest(request)
	repos := make([]cloud.Repository, 0)
	if err != nil {
		return repos, 0, err
	}
	res := cloud.PersonalRepositoryListResponse{}
	err = jsoniter.Unmarshal(apiRes.GetHttpContentBytes(), &res)
	if err != nil {
		return repos, 0, err
	}
	for _, r := range res.Data.Repos {
		repos = append(repos, cloud.Repository{
			Name: r.RepoName,
		})
	}
	return repos, res.Data.Total, nil
}

func (p *AlibabaCloud) PersonalImageList(region, repoNamespace, repoName string, pageNumber, pageSize int) ([]cloud.DockerArtifact, int, error) {
	request := getPersonalContainerRequest(region, "GET")
	request.PathPattern = fmt.Sprintf("/repos/%s/%s/tags", repoNamespace, repoName)
	request.QueryParams["Page"] = cast.ToString(pageNumber)
	request.QueryParams["PageSize"] = cast.ToString(pageSize)
	apiRes, err := p.sdkClient.ProcessCommonRequest(request)
	resp := make([]cloud.DockerArtifact, 0)
	if err != nil {
		return resp, 0, err
	}
	tags := cloud.TagsResponse{}
	err = jsoniter.Unmarshal(apiRes.GetHttpContentBytes(), &tags)
	if err != nil {
		return resp, 0, err
	}
	for _, tag := range tags.Data.Tags {
		resp = append(resp, cloud.DockerArtifact{
			Name: fmt.Sprintf("/%s/%s:%s", repoNamespace, repoName, tag.Tag),
		})
	}
	return resp, tags.Data.Total, err
}

func getEnterpriseContainerRequest(region, method string) *requests.CommonRequest {
	request := getCommonContainerRequest(region, method)
	request.Version = "2018-12-01"
	return request
}

func getPersonalContainerRequest(region, method string) *requests.CommonRequest {
	request := getCommonContainerRequest(region, method)
	request.Headers["Content-Type"] = "application/json"
	request.Version = "2016-06-07"
	return request
}

func getCommonContainerRequest(region, method string) *requests.CommonRequest {
	request := requests.NewCommonRequest()
	request.Method = method
	request.Scheme = "https"
	request.Domain = fmt.Sprintf("cr.%s.aliyuncs.com", region)
	return request
}
