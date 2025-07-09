package dto

// entities to store the app instance details in the orchestrator database
type OrchAppI struct {
	AppiId    string                         `json:"appi_id" bson:"appi_id"`
	Domain    string                         `json:"domain" bson:"domain"`
	Instances map[string]map[string]Instance `json:"instances" bson:"instances"`
	AppPkgId  string                         `json:"app_pkg_id" bson:"app_pkg_id"`
	NsPkgId   string                         `json:"ns_pkg_id" bson:"ns_pkg_id"`
	VnfPkgId  string                         `json:"vnf_pkg_id" bson:"vnf_pkg_id"`
}

type Instance struct {
	NSID  string            `json:"ns_id" bson:"ns_id"`
	VNFID string            `json:"vnf_id" bson:"vnf_id"`
	KDUs  map[string]string `json:"kdus" bson:"kdus"`
}

//
// kafka messages to be sent
//

// instantiate a new app instance
type InstantiateAppPkgMessage struct {
	AppPkgId     string `json:"app_pkg_id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	VimId        string `json:"vim_id"`
	Config       string `json:"config"`
	OriginDomain string `json:"original_domain"`
}

// terminate an app instance
type TerminateAppiMessage struct {
	AppiId string `json:"appi_id"`
}

// periodic message to list federated app instances
type FederatedAppisMessage struct {
	Appis map[string]AppiDetails `json:"appis"`
}

type AppiDetails struct {
	Domain    string              `json:"domain"`
	Instances map[string]Instance `json:"instances"`
}

// enable a kdu of an app instance
type EnableAppInstanceKDUMessage struct {
	AppdId string `json:"mec_appd_id"`
	KDUId  string `json:"kdu_id"`
	NSId   string `json:"ns_id"`
	Node   string `json:"node"`
}

// disable a kdu of an app instance
type DisableAppInstanceKDUMessage struct {
	AppdId string `json:"mec_appd_id"`
	KDUId  string `json:"kdu_id"`
	NSId   string `json:"ns_id"`
}

// migrate an app instance to a specific node
type MigrateAppInstanceNodeMessage struct {
	NsId  string `json:"ns_id"`
	VnfId string `json:"vnf_id"`
	KduId string `json:"kdu_id"`
	Node  string `json:"node"`
}

//
// kafka messages to be received
//

// new app instance
type NewAppInstanceMessage struct {
	AppInstanceId string `json:"appi_id"`
}
