package models

type ObjectRepoLocation struct {
	RepoURL  string `json:"repoURL,omitempty"`
	UserName string `json:"userName,omitempty"`
	Password string `json:"password,omitempty"`
	Token    string `json:"token,omitempty"`
}
