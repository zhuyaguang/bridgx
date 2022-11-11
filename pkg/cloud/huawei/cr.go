package huawei

import "github.com/galaxy-future/BridgX/pkg/cloud"

func (p *HuaweiCloud) PersonalImageList(region, repoNamespace, repoName string, pageNum, pageSize int) ([]cloud.DockerArtifact, int, error) {
	return []cloud.DockerArtifact{}, 0, nil
}

func (p *HuaweiCloud) EnterpriseImageList(region, instanceId, repoId, namespace, repoName string, pageNumber, pageSize int) ([]cloud.DockerArtifact, int, error) {
	return []cloud.DockerArtifact{}, 0, nil
}

func (p *HuaweiCloud) ContainerInstanceList(region string, pageNumber, pageSize int) ([]cloud.RegistryInstance, int, error) {
	return []cloud.RegistryInstance{}, 0, nil
}

func (p *HuaweiCloud) EnterpriseNamespaceList(region, instanceId string, pageNumber, pageSize int) ([]cloud.Namespace, int, error) {
	return []cloud.Namespace{}, 0, nil
}

func (p *HuaweiCloud) PersonalNamespaceList(region string) ([]cloud.Namespace, error) {
	return []cloud.Namespace{}, nil
}

func (p *HuaweiCloud) EnterpriseRepositoryList(region, instanceId, namespace string, pageNumber, pageSize int) ([]cloud.Repository, int, error) {
	return []cloud.Repository{}, 0, nil
}

func (p *HuaweiCloud) PersonalRepositoryList(region, namespace string, pageNumber, pageSize int) ([]cloud.Repository, int, error) {
	return []cloud.Repository{}, 0, nil
}
