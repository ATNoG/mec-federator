package models

type AppInstance struct {
	Id                  string `json:"appInstanceId" bson:"appInstanceId"`
	FederationContextId string `json:"federationContextId" bson:"federationContextId"`
	Name                string `json:"name" bson:"name"`
	Description         string `json:"description" bson:"description"`
	ArtefactId          string `json:"artefactId" bson:"artefactId"`
	AppiId              string `json:"appiId" bson:"appiId"`

	// is empty if the app instance is running on a partner zone
	AppPkgId string `json:"appPkgId,omitempty" bson:"appPkgId"`
	// refers to the ns id of the app instance if it is running on a partner zone
	NsId string `json:"nsId,omitempty" bson:"nsId"`
	// refers to the vnf id of the app instance if it is running on a partner zone
	VnfId string `json:"vnfId,omitempty" bson:"vnfId"`
}
