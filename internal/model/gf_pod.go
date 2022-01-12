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

func DeletePodByPodNameFromDB(podName string) error {
	if err := clients.WriteDBCli.Where("pod_name = ?", podName).Delete(&gf_cluster.Pod{}).Error; err != nil {
		logErr("DeletePodByPodNameFromDB from write db", err)
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

func UpdatePodByPodNameFromDB(pod *gf_cluster.Pod) error {
	if err := clients.WriteDBCli.Where("pod_name = ? and instance_group_name = ?", pod.PodName, pod.InstanceGroupName).Updates(map[string]interface{}{"node_name": pod.NodeName, "node_ip": pod.NodeIp, "pod_ip": pod.PodIP, "allocated_cpu_cores": pod.AllocatedCpuCores, "allocated_memory_gi": pod.AllocatedMemoryGi, "allocated_disk_gi": pod.AllocatedDiskGi, "instance_group_id": pod.InstanceGroupId, "instance_group_name": pod.InstanceGroupName, "running_time": pod.RunningTime, "status": pod.Status, "start_time": pod.StartTime, "created_user_id": pod.CreatedUserId}).Error; err != nil {
		logErr("UpdatePodByPodNameFromDB from write db", err)
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

func ListPodByCreatedUserFromDB(podIp string, nodeIp string, instanceGroupName string, createdUserId int64, pageNumber int, pageSize int) ([]*gf_cluster.PodSummary, int, error) {
	clients := clients.ReadDBCli.Model(gf_cluster.Pod{}).Where("created_user_id = ?", createdUserId)
	if podIp != "" {
		clients.Where("pod_ip like ?", "%"+podIp+"%")
	}

	if nodeIp != "" {
		clients.Where("node_ip like ?", "%"+nodeIp+"%")
	}

	if instanceGroupName != "" {
		clients.Where("instance_group_name like ?", "%"+instanceGroupName+"%")
	}

	var total int64
	if err := clients.Count(&total).Error; err != nil {
		logErr("ListPodByCreatedUserFromDB from read db", err)
		return nil, 0, err
	}
	var podList []*gf_cluster.Pod
	if err := clients.Order("id desc").Offset((pageNumber - 1) * pageSize).Limit(pageSize).Find(&podList).Error; err != nil {
		logErr("ListPodByCreatedUserFromDB from read db", err)
		return nil, 0, err
	}
	var result []*gf_cluster.PodSummary
	for _, pod := range podList {
		podSummary := &gf_cluster.PodSummary{
			NodeName:          pod.NodeName,
			NodeIp:            pod.NodeIp,
			PodName:           pod.PodName,
			PodIP:             pod.PodIP,
			AllocatedCpuCores: pod.AllocatedCpuCores,
			AllocatedMemoryGi: pod.AllocatedMemoryGi,
			AllocatedDiskGi:   pod.AllocatedDiskGi,
			GroupName:         pod.InstanceGroupName,
			RunningTime:       pod.RunningTime,
			Status:            pod.Status,
			GroupId:           pod.InstanceGroupId,
			StartTime:         pod.StartTime,
		}
		result = append(result, podSummary)
	}
	return result, int(total), nil
}

func ListPodByClusterIdFromDB(podIp string, nodeIp string, clusterId int64, pageNumber int, pageSize int) ([]*gf_cluster.PodSummary, int, error) {
	clients := clients.ReadDBCli.Model(gf_cluster.Pod{}).Where("kubernetes_infos.id = ?", clusterId)
	if podIp != "" {
		clients.Where("pod.pod_ip like ?", "%"+podIp+"%")
	}
	if nodeIp != "" {
		clients.Where("pod.node_ip like ?", "%"+nodeIp+"%")
	}

	var total int64
	if err := clients.Count(&total).Error; err != nil {
		logErr("ListPodByClusterIdFromDB from read db", err)
		return nil, 0, err
	}
	var podList []*gf_cluster.Pod
	if err := clients.Order("pod.id desc").Offset((pageNumber - 1) * pageSize).Limit(pageSize).Select("pod.pod_name,pod.pod_ip,pod.node_name,pod.node_ip,pod.allocated_cpu_cores,pod.allocated_memory_gi,pod.allocated_disk_gi,pod.running_time,pod.status,pod.instance_group_id,pod.start_time,instance_groups.name AS instance_group_name").Joins("LEFT JOIN instance_groups ON pod.instance_group_id = instance_groups.id LEFT JOIN kubernetes_infos ON instance_groups.kubernetes_id = kubernetes_infos.id").Find(&podList).Error; err != nil {
		logErr("ListPodByClusterIdFromDB from read db", err)
		return nil, 0, err
	}
	var result []*gf_cluster.PodSummary
	for _, pod := range podList {
		podSummary := &gf_cluster.PodSummary{
			NodeName:          pod.NodeName,
			NodeIp:            pod.NodeIp,
			PodName:           pod.PodName,
			PodIP:             pod.PodIP,
			AllocatedCpuCores: pod.AllocatedCpuCores,
			AllocatedMemoryGi: pod.AllocatedMemoryGi,
			AllocatedDiskGi:   pod.AllocatedDiskGi,
			GroupName:         pod.InstanceGroupName,
			RunningTime:       pod.RunningTime,
			Status:            pod.Status,
			GroupId:           pod.InstanceGroupId,
			StartTime:         pod.StartTime,
		}
		result = append(result, podSummary)
	}
	return result, int(total), nil
}
