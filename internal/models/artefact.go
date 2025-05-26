package models

type Artefact struct {
	Id                  string                 `json:"artefactId" bson:"artefactId"`
	FederationContextId string                 `json:"federationContextId" bson:"federationContextId"`
	AppProviderId       string                 `json:"appProviderId" bson:"appProviderId"`
	Name                string                 `json:"artefactName" bson:"artefactName"`
	VersionInfo         string                 `json:"artefactVersionInfo" bson:"artefactVersionInfo"`
	Description         string                 `json:"artefactDescription" bson:"artefactDescription"`
	VirtType            ArtefactVirtType       `json:"artefactVirtType" bson:"artefactVirtType"`
	DescriptorType      ArtefactDescriptorType `json:"artefactDescriptorType" bson:"artefactDescriptorType"`
	FileName            string                 `json:"artefactFileName,omitempty" bson:"artefactFileName"`
	FileFormat          ArtefactFileFormat     `json:"artefactFileFormat,omitempty" bson:"artefactFileFormat"`
	ArtefactFile        *[]byte                `json:"artefactFile,omitempty" bson:"artefactFile"`

	ComponentSpec      []ComponentSpec    `json:"componentSpec,omitempty" bson:"componentSpec"`
	ObjectRepoLocation ObjectRepoLocation `json:"objectRepoLocation,omitempty" bson:"objectRepoLocation"`

	AppPkgId string `json:"appPkgId,omitempty" bson:"appPkgId"`
}
