package ewbi

import (
	"time"

	"github.com/mankings/mec-federator/internal/models"
)

// Federation Management

type FederationRenewalResponse struct {
	FederationContextId string `json:"federationContextId"`

	FederationExpiryDate time.Time `json:"federationExpiryDate"`

	FederationRenewalDate time.Time `json:"federationRenewalDate"`
}

// Artefact Management

type OnboardArtefactRequest struct {
	ArtefactId             string                        `json:"artefactId" binding:"required"`
	ArtefactName           string                        `json:"artefactName" binding:"required"`
	AppProviderId          string                        `json:"appProviderId" binding:"required"`
	ArtefactVersionInfo    string                        `json:"artefactVersionInfo" binding:"required"`
	ArtefactVirtType       models.ArtefactVirtType       `json:"artefactVirtType" binding:"required"`
	ArtefactDescriptorType models.ArtefactDescriptorType `json:"artefactDescriptorType" binding:"required"`
	ComponentSpec          []models.ComponentSpec        `json:"componentSpec" binding:"required"`

	ArtefactDescription  string                    `json:"artefactDescription,omitempty"`
	ArtefactFileName     string                    `json:"artefactFileName,omitempty"`
	ArtefactFileFormat   models.ArtefactFileFormat `json:"artefactFileFormat,omitempty"`
	RepoType             models.RepoType           `json:"repoType,omitempty"`
	ArtefactRepoLocation models.ObjectRepoLocation `json:"objectRepoLocation,omitempty"`
	ArtefactFile         models.File               `json:"artefactFile,omitempty"`
}
