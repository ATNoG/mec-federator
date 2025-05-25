package dto

type NewAppPkg struct {
	AppdId     string `json:"appd_id" binding:"required" bson:"appd_id"`
	Name       string `json:"name" binding:"required" bson:"name"`
	Provider   string `json:"provider" binding:"required" bson:"provider"`
	Version    string `json:"version" binding:"required" bson:"version"`
	MecVersion string `json:"mec-version" binding:"required" bson:"mec-version"`
	InfoName   string `json:"info-name" binding:"required" bson:"info-name"`
	Description string `json:"description" binding:"required" bson:"description"`
	AppD       []byte `json:"appd" binding:"required" bson:"appd"`
}

type NewAppPkgMessage struct {
	AppPkgId  string `json:"app_pkg_id" binding:"required"`
	MessageId string `json:"msg_id" binding:"required"`
}
