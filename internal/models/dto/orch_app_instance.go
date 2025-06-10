package dto

// entities to store the app instance details in the orchestrator database
type OrchAppI struct {
	AppiId    string                         `json:"appi_id" bson:"appi_id"`
	Domain    string                         `json:"domain" bson:"domain"`
	Instances map[string]map[string]Instance `json:"instances" bson:"instances"`
	NsPkgId   string                         `json:"ns_pkg_id" bson:"ns_pkg_id"`
	VnfPkgId  string                         `json:"vnf_pkg_id" bson:"vnf_pkg_id"`
}

type Instance struct {
	NSID  string            `json:"ns_id" bson:"ns_id"`
	VNFID string            `json:"vnf_id" bson:"vnf_id"`
	KDUs  map[string]string `json:"kdus" bson:"kdus"`
}

//
// kafka messages
//

// instantiate a new app instance
type InstantiateAppPkgMessage struct {
	AppPkgId    string `json:"app_pkg_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	VimId       string `json:"vim_id"`
}

// terminate an app instance
type TerminateAppiMessage struct {
	AppInstanceId string `json:"appi_id"`
}

// periodic message to list federated app instances
type FederatedAppisMessage struct {
	Appis  map[string]AppiDetails `json:"appis"`
	Domain string                 `json:"domain"`
}

type AppiDetails struct {
	Instances map[string]Instance `json:"instances"`
}
