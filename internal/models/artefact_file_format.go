package models

type ArtefactFileFormat string

const (
	ZIP   ArtefactFileFormat = "ZIP"
	TAR   ArtefactFileFormat = "TAR"
	TEXT  ArtefactFileFormat = "TEXT"
	TARGZ ArtefactFileFormat = "TARGZ"
)
