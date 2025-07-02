package dto

import "go.mongodb.org/mongo-driver/v2/bson"

type NewAppPkg struct {
	AppdId      string `json:"appd_id" binding:"required" bson:"appd_id"`
	Name        string `json:"name" binding:"required" bson:"name"`
	Provider    string `json:"provider" binding:"required" bson:"provider"`
	Version     string `json:"version" binding:"required" bson:"version"`
	MecVersion  string `json:"mec-version" binding:"required" bson:"mec-version"`
	InfoName    string `json:"info-name" binding:"required" bson:"info-name"`
	Description string `json:"description" binding:"required" bson:"description"`
	AppD        []byte `json:"appd" binding:"required" bson:"appd"`
}

type OrchAppPkg struct {
	Id          bson.ObjectID `json:"_id" binding:"required" bson:"_id"`
	AppdId      string        `json:"appd_id" binding:"required" bson:"appd_id"`
	Name        string        `json:"name" binding:"required" bson:"name"`
	Provider    string        `json:"provider" binding:"required" bson:"provider"`
	Version     string        `json:"version" binding:"required" bson:"version"`
	MecVersion  string        `json:"mec-version" binding:"required" bson:"mec-version"`
	InfoName    string        `json:"info-name" binding:"required" bson:"info-name"`
	Description string        `json:"description" binding:"required" bson:"description"`
	AppD        []byte        `json:"appd" binding:"required" bson:"appd"`
	NsPkgId     string        `json:"ns_pkg_id" binding:"required" bson:"ns_pkg_id"`
	VnfPkgId    string        `json:"vnf_pkg_id" binding:"required" bson:"vnf_pkg_id"`
}

type NewAppPkgMessage struct {
	AppPkgId string `json:"app_pkg_id" binding:"required"`
}

type DeleteAppPkgMessage struct {
	AppPkgId string `json:"app_pkg_id" binding:"required"`
}
