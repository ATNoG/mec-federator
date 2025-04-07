package models

type OrchestratorInfo struct {
	OperatorId string `json:"operatorId" bson:"operatorId"`

	OperatorName string `json:"operatorName" bson:"operatorName"`

	OperatorCountryCode string `json:"operatorAddress" bson:"operatorAddress"`

	MobileCountryCode string `json:"operatorPublicKey" bson:"operatorPublicKey"`

	MobileNetworkCodes []string `json:"mobileNetworkCodes" bson:"mobileNetworkCodes"`

	KafkaUrl string `json:"kafkaUrl" bson:"kafkaUrl"`
}
