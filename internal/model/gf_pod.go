package model

import (
	"github.com/galaxy-future/BridgX/internal/clients"
	gf_cluster "github.com/galaxy-future/BridgX/pkg/gf-cluster"
)

func CreatePodFromDB(pod *gf_cluster.Pod) error {
	if err := clients.WriteDBCli.Create(pod).Error; err != nil {
		logErr("CreatePodFromDB from write db", err)
		return err
	}
	return nil
}

func DeletePodFromDB(podId int64) error {
	if err := clients.WriteDBCli.Delete(&gf_cluster.Pod{}, podId).Error; err != nil {
		logErr("DeletePodFromDB from write db", err)
		return err
	}
	return nil
}

func UpdatePodFromDB(pod *gf_cluster.Pod) error {
	if err := clients.WriteDBCli.Save(pod).Error; err != nil {
		logErr("UpdatePodFromDB from write db", err)
		return err
	}
	return nil
}

func GetPodFromDB(podId int64) (*gf_cluster.Pod, error) {
	var pod gf_cluster.Pod
	if err := clients.ReadDBCli.Where("id = ?", podId).First(&pod).Error; err != nil {
		logErr("GetPodFromDB from read db", err)
		return nil, err
	}
	return &pod, nil
}
