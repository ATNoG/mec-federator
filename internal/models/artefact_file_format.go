package models

// ArtefactFileFormat : Artefacts like Helm charts or Terraform scripts may need compressed format.
type ArtefactFileFormat string

// List of ArtefactFileFormat
const (
	ZIP   ArtefactFileFormat = "ZIP"
	TAR   ArtefactFileFormat = "TAR"
	TEXT  ArtefactFileFormat = "TEXT"
	TARGZ ArtefactFileFormat = "TARGZ"
)
