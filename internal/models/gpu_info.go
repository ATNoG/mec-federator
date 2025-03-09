package models

type GpuInfo struct {
	GpuVendorType string `json:"gpuVendorType"`
	GpuModeName   string `json:"gpuModeName"`
	GpuMemory     int32  `json:"gpuMemory"`
	NumGPU        int32  `json:"numGPU"`
}
