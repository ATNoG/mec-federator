package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mankings/mec-federator/internal/config"
	"github.com/mankings/mec-federator/internal/models"
)

func PostFederationPartnerController(c *gin.Context) {
	logger := config.GetLogger("PostFederationPartnerController")

	var federationRequestData models.FederationRequestData
	if err := c.ShouldBindJSON(&federationRequestData); err != nil {
		logger.Errorf("Error while binding JSON: %v", err)

		problemDetails := models.ProblemDetails{
			Title:   "Insufficient parameters",
			Details: "Incorrect values received in federation request",
			Cause:   "INVALID_FED_RQST_PARAMETERS",
		}
		c.JSON(http.StatusBadRequest, problemDetails)
		return
	}

	logger.Infof("Received federation request data: %v", federationRequestData)

	federationResponseData := models.FederationResponseData{
		PartnerOPFederationId: "partnerOPFederationId",
		PartnerOPCountryCode:  "partnerOPCountryCode",
		FederationContextId:   "federationContextId",
		PartnerOPMobileNetworkCodes: models.MobileNetworkCodes{
			Mcc: "mcc",
			Mncs: []string{
				"mnc1",
				"mnc2",
			},
		},
		PartnerOPFixedNetworkCodes: []string{
			"partnerOPFixedNetworkCode1",
			"partnerOPFixedNetworkCode2",
		},
		EdgeDiscoveryServiceEndpoint: models.ServiceEndpoint{
			Port: 8080,
		},
		OfferedAvailabilityZones: []string{"AZ1", "AZ2"},
		PlatformCaps:             []string{"cap1", "cap2"},
		FederationRenewalDate:    time.Now().AddDate(0, 6, 0),
		FederationExpiryDate:     time.Now().AddDate(1, 0, 0),
		DefaultSubscriptions:     []string{"sub1", "sub2"},
		SubscriptionsPeriodicity: "monthly",
		OpsInfoExposureEndpoint:  "opsInfoExposureEndpoint",
	}

	c.JSON(http.StatusOK, federationResponseData)
}
