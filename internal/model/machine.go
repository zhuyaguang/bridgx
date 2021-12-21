package model

// Machine 机器
type Machine struct {
	Ip       string `json:"ip"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type MachineRequest struct {
	Ip       string `json:"ip"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type MachineResponse struct {
	Ip     string `json:"ip"`
	IsPass bool   `json:"is_pass"`
}
