package models

type CpuArchType string

const (
	X86    CpuArchType = "ISA_X86"
	X86_64 CpuArchType = "ISA_X86_64"
	ARM_64 CpuArchType = "ISA_ARM_64"
)
