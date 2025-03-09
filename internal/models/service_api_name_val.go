package models

// ServiceApiNameVal : Name of the Service API
type ServiceApiNameVal string

// List of serviceAPINameVal
const (
	QUALITY_ON_DEMAND   ServiceApiNameVal = "QualityOnDemand"
	DEVICE_LOCATION     ServiceApiNameVal = "DeviceLocation"
	DEVICE_STATUS       ServiceApiNameVal = "DeviceStatus"
	SIM_SWAP            ServiceApiNameVal = "SimSwap"
	NUMBER_VERIFICATION ServiceApiNameVal = "NumberVerification"
	DEVICE_IDENTIFIER   ServiceApiNameVal = "DeviceIdentifier"
)
