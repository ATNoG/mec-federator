package models

type Flavour struct {
	FlavourId string `json:"flavourId"`

	CpuArchType *CpuArchType `json:"cpuArchType"`
	// A list of operating systems which this flavour configuration can support e.g., RHEL Linux, Ubuntu 18.04 LTS, MS Windows 2012 R2.
	SupportedOSTypes []OsType `json:"supportedOSTypes"`
	// Number of available vCPUs
	NumCPU int32 `json:"numCPU"`
	// Amount of RAM in Mbytes
	MemorySize int32 `json:"memorySize"`
	// Amount of disk storage in Gbytes
	StorageSize int32 `json:"storageSize"`

	Gpu []GpuInfo `json:"gpu,omitempty"`
	// Number of FPGAs
	Fpga int32 `json:"fpga,omitempty"`
	// Number of Intel VPUs available
	Vpu int32 `json:"vpu,omitempty"`

	Hugepages []HugePage `json:"hugepages,omitempty"`
	// Support for exclusive CPUs
	CpuExclusivity bool `json:"cpuExclusivity,omitempty"`
}
