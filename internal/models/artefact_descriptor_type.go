package models

// ArtefactDescriptorType : Type of descriptor present in the artefact.  App provider can either define either a Helm chart or a Terraform script or container spec.
type ArtefactDescriptorType string

// List of ArtefactDescriptorType
const (
	HELM          ArtefactDescriptorType = "HELM"
	TERRAFORM     ArtefactDescriptorType = "TERRAFORM"
	ANSIBLE       ArtefactDescriptorType = "ANSIBLE"
	SHELL         ArtefactDescriptorType = "SHELL"
	COMPONENTSPEC ArtefactDescriptorType = "COMPONENTSPEC"
)
