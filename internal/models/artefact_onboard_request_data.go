package models

type ArtefactOnboardRequest struct {
	ArtefactId             string                 `json:"artefactId" binding:"required"`
	ArtefactName           string                 `json:"artefactName" binding:"required"`
	AppProviderId          string                 `json:"appProviderId" binding:"required"`
	ArtefactVersionInfo    string                 `json:"artefactVersionInfo" binding:"required"`
	ArtefactVirtType       ArtefactVirtType       `json:"artefactVirtType" binding:"required"`
	ArtefactDescriptorType ArtefactDescriptorType `json:"artefactDescriptorType" binding:"required"`
	ComponentSpec          []ComponentSpec        `json:"componentSpec" binding:"required"`

	ArtefactDescription  string             `json:"artefactDescription,omitempty"`
	ArtefactFileName     string             `json:"artefactFileName,omitempty"`
	ArtefactFileFormat   ArtefactFileFormat `json:"artefactFileFormat,omitempty"`
	RepoType             RepoType           `json:"repoType,omitempty"`
	ArtefactRepoLocation ObjectRepoLocation `json:"objectRepoLocation,omitempty"`
	ArtefactFile         File               `json:"artefactFile,omitempty"`
}
