package tencent

import "github.com/galaxy-future/BridgX/pkg/cloud"

func (p *TencentCloud) PersonalImageList(region, repoNamespace, repoName string, pageNum, pageSize int) ([]cloud.DockerArtifact, int, error) {
	return []cloud.DockerArtifact{}, 0, nil
}

func (p *TencentCloud) EnterpriseImageList(region, instanceId, repoId, namespace, repoName string, pageNumber, pageSize int) ([]cloud.DockerArtifact, int, error) {
	return []cloud.DockerArtifact{}, 0, nil
}

func (p *TencentCloud) ContainerInstanceList(region string, pageNumber, pageSize int) ([]cloud.RegistryInstance, int, error) {
	return []cloud.RegistryInstance{}, 0, nil
}

func (p *TencentCloud) EnterpriseNamespaceList(region, instanceId string, pageNumber, pageSize int) ([]cloud.Namespace, int, error) {
	return []cloud.Namespace{}, 0, nil
}

func (p *TencentCloud) PersonalNamespaceList(region string) ([]cloud.Namespace, error) {
	return []cloud.Namespace{}, nil
}

func (p *TencentCloud) EnterpriseRepositoryList(region, instanceId, namespace string, pageNumber, pageSize int) ([]cloud.Repository, int, error) {
	return []cloud.Repository{}, 0, nil
}

func (p *TencentCloud) PersonalRepositoryList(region, namespace string, pageNumber, pageSize int) ([]cloud.Repository, int, error) {
	return []cloud.Repository{}, 0, nil
}
