package gf_cluster

// Pod 表示一个Pod
type Pod struct {
	Id                int64   `json:"id"`
	NodeName          string  `json:"node_name"`
	NodeIp            string  `json:"node_ip"`
	PodName           string  `json:"pod_name"`
	PodIP             string  `json:"pod_ip"`
	AllocatedCpuCores float64 `json:"allocated_cpu_cores"`
	AllocatedMemoryGi float64 `json:"allocated_memory_gi"`
	AllocatedDiskGi   float64 `json:"allocated_disk_gi"`
	GroupName         string  `json:"group_name"`
	RunningTime       string  `json:"running_time"`
	Status            string  `json:"status"`
	GroupId           int64   `json:"group_id"`
	StartTime         int64   `json:"start_time"`
}

func (Pod) TableName() string {
	return "pod"
}
