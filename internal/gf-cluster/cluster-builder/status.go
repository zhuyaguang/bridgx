package cluster_builder

import (
	"github.com/galaxy-future/BridgX/internal/clients"
	"github.com/galaxy-future/BridgX/internal/logs"
	gf_cluster "github.com/galaxy-future/BridgX/pkg/gf-cluster"
	"go.uber.org/zap"
)

func updateStatus(id int64, status string) {
	kubernetes := gf_cluster.KubernetesInfo{
		Id:     id,
		Status: status,
	}

	_ = update(kubernetes)
}

// TODO: drop this func on future
//deprecated
func updateInstallStep(id int64, step string) {
	kubernetes := gf_cluster.KubernetesInfo{
		Id:          id,
		InstallStep: step,
	}

	_ = update(kubernetes)
}

func recordStep(kubernetesId int64, ip, step string, err error) {
	var msg string
	if err != nil {
		msg = err.Error()
	} else {
		msg = "success"
	}

	installStep := gf_cluster.KubernetesInstallStep{
		KubernetesId: kubernetesId,
		HostIp:       ip,
		Operation:    step,
		Message:      msg,
	}

	connection := clients.WriteDBCli
	if connection == nil {
		return
	}

	connection.Create(&installStep)
}

func recordConfig(id int64, config string) error {
	kubernetes := gf_cluster.KubernetesInfo{
		Id:     id,
		Config: config,
	}

	return update(kubernetes)
}

func failed(id int64, message string) {
	logs.Logger.Error("创建集群失败", zap.Int64("kubernetes id ", id), zap.String("message", message))

	kubernetes := gf_cluster.KubernetesInfo{
		Id:      id,
		Status:  gf_cluster.KubernetesStatusFailed,
		Message: message,
	}
	_ = update(kubernetes)
}

func update(kubernetes gf_cluster.KubernetesInfo) error {
	connection := clients.WriteDBCli
	tx := connection.Model(&gf_cluster.KubernetesInfo{}).Where("id = ?", kubernetes.Id)
	kubernetes.Id = 0
	if err := tx.Updates(kubernetes).Error; err != nil {
		logs.Logger.Error("mysql update失败", zap.Int64("kubernetes id ", kubernetes.Id), zap.Error(err))
		return err
	}

	return nil
}
