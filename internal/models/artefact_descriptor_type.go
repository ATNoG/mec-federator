package models

type ArtefactDescriptorType string

const (
	HELM          ArtefactDescriptorType = "HELM"
	TERRAFORM     ArtefactDescriptorType = "TERRAFORM"
	ANSIBLE       ArtefactDescriptorType = "ANSIBLE"
	SHELL         ArtefactDescriptorType = "SHELL"
	COMPONENTSPEC ArtefactDescriptorType = "COMPONENTSPEC"
)
