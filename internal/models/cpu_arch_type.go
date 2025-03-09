package models

// CpuArchType : CPU Instruction Set Architecture (ISA) E.g., Intel, Arm etc.
type CpuArchType string

// List of CPUArchType
const (
	X86    CpuArchType = "ISA_X86"
	X86_64 CpuArchType = "ISA_X86_64"
	ARM_64 CpuArchType = "ISA_ARM_64"
)
