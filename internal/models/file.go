package models

type File struct {
	Id string `json:"fileId" bson:"fileId"`

	FederationContextId string `json:"federationContextId" bson:"federationContextId"`

	FileName string `json:"fileName" bson:"fileName"`

	FileDescription string `json:"fileDescription,omitempty" bson:"fileDescription"`

	FileVersionInfo string `json:"fileVersionInfo,omitempty" bson:"fileVersionInfo"`

	FileType VirtImageType `json:"fileType" bson:"fileType"`

	Checksum string `json:"checksum" bson:"checksum"`

	ImgOsType OsType `json:"imgOsType" bson:"imgOsType"`

	ImgInsSetArch CpuArchType `json:"imgInsSetArch" bson:"imgInsSetArch"`

	RepoType RepoType `json:"repoType,omitempty" bson:"repoType"`

	ObjectRepoLocation ObjectRepoLocation `json:"objectRepoLocation,omitempty" bson:"objectRepoLocation"`

	ArtefactFile *[]byte `json:"artefactFile,omitempty" bson:"artefactFile"`
}
