package models

type ComputeResourceInfo struct {
	// CPU Instruction Set Architecture (ISA) E.g., Intel, Arm etc.
	CpuArchType string `json:"cpuArchType"`

	NumCPU string `json:"numCPU"`
	// Amount of RAM in Mbytes
	Memory int64 `json:"memory"`
	// Amount of disk storage in Gbytes for a given ISA type
	DiskStorage int32 `json:"diskStorage,omitempty"`

	Gpu []GpuInfo `json:"gpu,omitempty"`
	// Number of Intel VPUs available for a given ISA type
	Vpu int32 `json:"vpu,omitempty"`
	// Number of FPGAs available for a given ISA type
	Fpga int32 `json:"fpga,omitempty"`

	Hugepages []HugePage `json:"hugepages,omitempty"`
	// Support for exclusive CPUs
	CpuExclusivity bool `json:"cpuExclusivity,omitempty"`
}
