package models

type GpuInfo struct {
	// GPU vendor name e.g. NVIDIA, AMD etc.
	GpuVendorType string `json:"gpuVendorType"`
	// Model name corresponding to vendorType may include info e.g. for NVIDIA, model name could be “Tesla M60”, “Tesla V100” etc.
	GpuModeName string `json:"gpuModeName"`
	// GPU memory in Mbytes
	GpuMemory int32 `json:"gpuMemory"`
	// Number of GPUs
	NumGPU int32 `json:"numGPU"`
}
