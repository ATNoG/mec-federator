package models

type ComputeResourceInfo struct {
	CpuArchType string `json:"cpuArchType"`

	NumCPU      string `json:"numCPU"`
	Memory      int64  `json:"memory"`
	DiskStorage int32  `json:"diskStorage,omitempty"`

	Gpu  []GpuInfo `json:"gpu,omitempty"`
	Vpu  int32     `json:"vpu,omitempty"`
	Fpga int32     `json:"fpga,omitempty"`

	Hugepages      []HugePage `json:"hugepages,omitempty"`
	CpuExclusivity bool       `json:"cpuExclusivity,omitempty"`
}
