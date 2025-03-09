package models

// RepoType : Artefact or file repository location. PUBLICREPO is used of public URLs like GitHub, Helm repo, docker registry etc., PRIVATEREPO is used for private repo managed by the application developer, UPLOAD is for the case when artefact/file is uploaded from MEC web portal.  OP should pull the image from ‘repoUrl' immediately after receiving the request and then send back the response. In case the repoURL corresponds to a docker registry, use docker v2 http api to do the pull.
type RepoType string

// List of RepoType
const (
	PRIVATEREPO RepoType = "PRIVATEREPO"
	PUBLICREPO  RepoType = "PUBLICREPO"
	UPLOAD      RepoType = "UPLOAD"
)
