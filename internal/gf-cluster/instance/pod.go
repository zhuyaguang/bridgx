package instance

import (
	"github.com/galaxy-future/BridgX/internal/gf-cluster/cluster"
	"github.com/galaxy-future/BridgX/internal/logs"
	"github.com/galaxy-future/BridgX/internal/model"
	gf_cluster "github.com/galaxy-future/BridgX/pkg/gf-cluster"
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"strconv"
	"time"
)

func createPod(instanceGroup *gf_cluster.InstanceGroup, pod *v1.Pod) error {
	cpuUsed, memoryUsed, storageUsed := getPodResourceInfo(pod)
	runningTime := getPodRunningTime(pod)
	podInfo := &gf_cluster.Pod{
		NodeName:          pod.Spec.NodeName,
		NodeIp:            pod.Status.HostIP,
		PodName:           pod.GetName(),
		PodIP:             pod.Status.PodIP,
		AllocatedCpuCores: cpuUsed,
		AllocatedMemoryGi: memoryUsed,
		AllocatedDiskGi:   storageUsed,
		InstanceGroupName: instanceGroup.Name,
		RunningTime:       runningTime,
		Status:            string(pod.Status.Phase),
		InstanceGroupId:   instanceGroup.Id,
		StartTime:         pod.Status.StartTime.Time.Unix(),
	}
	err := model.CreatePodFromDB(podInfo)
	if err != nil {
		return err
	}
	return nil
}

func deletePod(id int64) error {
	err := model.DeletePodFromDB(id)
	if err != nil {
		return err
	}
	return nil
}

func deletePodByPodName(podName string) error {
	err := model.DeletePodByPodNameFromDB(podName)
	if err != nil {
		return err
	}
	return nil
}

func updatePodById(pod *gf_cluster.Pod) error {
	err := model.UpdatePodFromDB(pod)
	if err != nil {
		return err
	}
	return nil
}

func updatePodByPodName(pod *v1.Pod) error {
	cpuUsed, memoryUsed, storageUSed := getPodResourceInfo(pod)
	instanceGroupName, instanceGroupId, _ := getPodGroupInfoFromLabels(pod)
	runningTime := getPodRunningTime(pod)
	podInfo := &gf_cluster.Pod{
		NodeName:          pod.Spec.NodeName,
		NodeIp:            pod.Status.HostIP,
		PodName:           pod.GetName(),
		PodIP:             pod.Status.PodIP,
		AllocatedCpuCores: cpuUsed,
		AllocatedMemoryGi: memoryUsed,
		AllocatedDiskGi:   storageUSed,
		InstanceGroupName: instanceGroupName,
		RunningTime:       runningTime,
		Status:            string(pod.Status.Phase),
		InstanceGroupId:   instanceGroupId,
		StartTime:         pod.Status.StartTime.Time.Unix(),
	}
	err := model.UpdatePodByPodNameFromDB(podInfo)
	if err != nil {
		return err
	}
	return nil
}

func ListPodByCreatedUserId(podIp string, nodeIp string, instanceGroupName string, createdUserId int64, pageNumber int, pageSize int) ([]*gf_cluster.PodSummary, int, error) {
	podList, total, err := model.ListPodByCreatedUserFromDB(podIp, nodeIp, instanceGroupName, createdUserId, pageNumber, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return podList, total, nil
}

// getPodResourceInfo 获取k8s pod资源信息： AllocatedCpuCores & AllocatedMemoryGi & AllocatedDiskGi
func getPodResourceInfo(pod *v1.Pod) (float64, float64, float64) {
	cpuResource := resource.NewScaledQuantity(0, resource.Kilo)
	memoryResource := resource.NewScaledQuantity(0, resource.Kilo)
	storageResource := resource.NewScaledQuantity(0, resource.Kilo)
	if pod.Status.Phase == v1.PodRunning {
		for _, container := range pod.Spec.Containers {
			cpuResource.Add(*container.Resources.Limits.Cpu())
			memoryResource.Add(*container.Resources.Limits.Memory())
			storageResource.Add(*container.Resources.Limits.StorageEphemeral())
		}
	}
	cpuUsed := cluster.CpuQuantity2Float(*cpuResource)
	memoryUsed := cluster.StorageQuantity2Float(*memoryResource)
	storageUSed := cluster.StorageQuantity2Float(*storageResource)
	return cpuUsed, memoryUsed, storageUSed
}

// getPodGroupInfo 获取实例组id和实例组名称
func getPodGroupInfoFromLabels(pod *v1.Pod) (string, int64, error) {
	groupName, exist := pod.Labels[gf_cluster.ClusterInstanceGroupKey]
	if !exist {
		groupName = "unknown"
	}

	groupIdStr, exist := pod.Labels[gf_cluster.ClusterInstanceGroupIdKey]
	if !exist {
		groupIdStr = "0"
	}
	groupId, err := strconv.ParseInt(groupIdStr, 10, 64)
	if err != nil {
		logs.Logger.Error("获取GroupId失败", zap.String("value", groupIdStr), zap.Error(err))
	}
	return groupName, groupId, err
}

// getPodRunningTime 获取pod运行时间
func getPodRunningTime(pod *v1.Pod) string {
	runningTime := ""
	if pod.Status.StartTime != nil {
		runningTime = cluster.FormatHumanReadableDuration(time.Now().Sub(pod.Status.StartTime.Time))
	}
	return runningTime
}
