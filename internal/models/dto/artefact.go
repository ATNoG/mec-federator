package dto

import (
	"errors"

	"github.com/mankings/mec-federator/internal/models"
)

type ArtefactOnboardRequest struct {
	ArtefactId             string                        `json:"artefactId" form:"artefactId" binding:"required"`
	ArtefactName           string                        `json:"artefactName" form:"artefactName" binding:"required"`
	AppProviderId          string                        `json:"appProviderId" form:"appProviderId" binding:"required"`
	ArtefactVersionInfo    string                        `json:"artefactVersionInfo" form:"artefactVersionInfo" binding:"required"`
	ArtefactVirtType       models.ArtefactVirtType       `json:"artefactVirtType" form:"artefactVirtType" binding:"required"`
	ArtefactDescriptorType models.ArtefactDescriptorType `json:"artefactDescriptorType" form:"artefactDescriptorType" binding:"required"`

	ComponentSpec        []models.ComponentSpec    `json:"componentSpec,omitempty" form:"componentSpec,omitempty"`
	ArtefactDescription  string                    `json:"artefactDescription,omitempty" form:"artefactDescription,omitempty"`
	ArtefactFileName     string                    `json:"artefactFileName,omitempty" form:"artefactFileName,omitempty"`
	ArtefactFileFormat   models.ArtefactFileFormat `json:"artefactFileFormat,omitempty" form:"artefactFileFormat,omitempty"`
	RepoType             models.RepoType           `json:"repoType,omitempty" form:"repoType,omitempty"`
	ArtefactRepoLocation models.ObjectRepoLocation `json:"objectRepoLocation,omitempty" form:"objectRepoLocation,omitempty"`
	// ArtefactFile         models.File               `json:"artefactFile,omitempty"`
}

func (a *ArtefactOnboardRequest) Validate() error {
	if a.ArtefactId == "" {
		return errors.New("artefactId is required")
	}
	if a.ArtefactName == "" {
		return errors.New("artefactName is required")
	}

	if a.AppProviderId == "" {
		return errors.New("appProviderId is required")
	}

	if a.ArtefactVersionInfo == "" {
		return errors.New("artefactVersionInfo is required")
	}

	if a.ArtefactVirtType == "" {
		return errors.New("artefactVirtType is required")
	}

	if a.ArtefactDescriptorType == "" {
		return errors.New("artefactDescriptorType is required")
	}

	if a.ArtefactFileFormat == "" {
		return errors.New("artefactFileFormat is required")
	}

	if a.RepoType == "" {
		return errors.New("repoType is required")
	} else if a.RepoType == models.PUBLICREPO || a.RepoType == models.PRIVATEREPO {
		return errors.New("invalid repoType")
	}

	return nil
}
