package dto

type InstantiateAppPkgMessage struct {
	AppPkgId    string `json:"app_pkg_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	VimId       string `json:"vim_id"`
}

type TerminateAppPkgMessage struct {
	AppInstanceId string `json:"appi_id"`
}
