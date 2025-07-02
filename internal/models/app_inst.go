package models

type AppInstance struct {
	Id                  string `json:"appInstanceId" bson:"appInstanceId"`
	FederationContextId string `json:"federationContextId" bson:"federationContextId"`
	Name                string `json:"name" bson:"name"`
	Description         string `json:"description" bson:"description"`
	ArtefactId          string `json:"artefactId" bson:"artefactId"`
	AppiId              string `json:"appiId" bson:"appiId"`
	AppPkgId            string `json:"appPkgId,omitempty" bson:"appPkgId"`
}
