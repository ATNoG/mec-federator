package models

type Artefact struct {
	Id string `json:"artefactId" bson:"artefactId"`

	FederationContextId string `json:"federationContextId" bson:"federationContextId"`

	Name string `json:"artefactName" bson:"artefactName"`

	Description string `json:"artefactDescription" bson:"artefactDescription"`

	VirtType ArtefactVirtType `json:"artefactVirtType" bson:"artefactVirtType"`

	DescriptorType ArtefactDescriptorType `json:"artefactDescriptorType" bson:"artefactDescriptorType"`

	FileName string `json:"artefactFileName,omitempty" bson:"artefactFileName"`

	FileFormat ArtefactFileFormat `json:"artefactFileFormat,omitempty" bson:"artefactFileFormat"`

	ArtefactFile *[]byte `json:"artefactFile,omitempty" bson:"artefactFile"`

	ComponentSpec []ComponentSpec `json:"componentSpec,omitempty" bson:"componentSpec"`

	RepoType RepoType `json:"repoType,omitempty" bson:"repoType"`

	ObjectRepoLocation ObjectRepoLocation `json:"objectRepoLocation,omitempty" bson:"objectRepoLocation"`
}
