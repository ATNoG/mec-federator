package models

type ServiceApiNameVal string

const (
	QUALITY_ON_DEMAND   ServiceApiNameVal = "QualityOnDemand"
	DEVICE_LOCATION     ServiceApiNameVal = "DeviceLocation"
	DEVICE_STATUS       ServiceApiNameVal = "DeviceStatus"
	SIM_SWAP            ServiceApiNameVal = "SimSwap"
	NUMBER_VERIFICATION ServiceApiNameVal = "NumberVerification"
	DEVICE_IDENTIFIER   ServiceApiNameVal = "DeviceIdentifier"
)
