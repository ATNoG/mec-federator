package models

type Application struct {
	Id                  string `json:"appId" bson:"appId"`
	FederationContextId string `json:"federationContextId" bson:"federationContextId"`
	Name                string `json:"name" bson:"name"`
	Description         string `json:"description" bson:"description"`
	ArtefactId          string `json:"artefactId" bson:"artefactId"`

	// is empty if the app instance is running on a partner zone
	AppPkgId string `json:"appPkgId,omitempty" bson:"appPkgId"`
}
