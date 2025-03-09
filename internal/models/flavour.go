package models

type Flavour struct {
	FlavourId string `json:"flavourId"`

	CpuArchType      *CpuArchType `json:"cpuArchType"`
	SupportedOSTypes []OsType     `json:"supportedOSTypes"`
	NumCPU           int32        `json:"numCPU"`
	MemorySize       int32        `json:"memorySize"`
	StorageSize      int32        `json:"storageSize"`

	Gpu  []GpuInfo `json:"gpu,omitempty"`
	Fpga int32     `json:"fpga,omitempty"`
	Vpu  int32     `json:"vpu,omitempty"`

	Hugepages      []HugePage `json:"hugepages,omitempty"`
	CpuExclusivity bool       `json:"cpuExclusivity,omitempty"`
}
