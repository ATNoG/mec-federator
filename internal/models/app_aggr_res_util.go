package models

type AppAggrResUtil struct {
	AppId string `json:"appId"`

	AppProvId string `json:"appProvId"`
	// No of application instances of appId in a zone
	NoOfAppInstances int32 `json:"noOfAppInstances"`

	AppInstances []string `json:"appInstances"`

	CpuUtil *CpuUtilization `json:"cpuUtil"`

	MemUtil *MemUtilization `json:"memUtil"`

	DiskUtil *DiskUtilization `json:"diskUtil"`

	NetworkUtil *NetworkUtilization `json:"networkUtil"`

	FlavourUtil *[]FlavourMetrics `json:"flavourUtil"`
}
