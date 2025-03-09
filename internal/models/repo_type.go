package models

type RepoType string

const (
	PRIVATEREPO RepoType = "PRIVATEREPO"
	PUBLICREPO  RepoType = "PUBLICREPO"
	UPLOAD      RepoType = "UPLOAD"
)
