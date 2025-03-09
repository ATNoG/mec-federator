package models

type ObjectRepoLocation struct {
	RepoURL string `json:"repoURL,omitempty"`
	// Username to access the repository
	UserName string `json:"userName,omitempty"`
	// Password to access the repository
	Password string `json:"password,omitempty"`
	// Authorization token to access the repository
	Token string `json:"token,omitempty"`
}
