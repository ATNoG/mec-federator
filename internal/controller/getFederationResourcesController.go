package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Retrieves REST APIs supported by an OP for federation services.
// @Description Retrieves REST APIs supported by an OP for federation services.
// @Tags FederationAPIManagement
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 422 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure 503 {object} map[string]interface{}
// @Failure 520 {object} map[string]interface{}
// @Router /federation-resources [get]
func GetFederationResourcesController(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"federationSupportedAPIs": []string{
			"API1",
			"API2",
			"API3",
		},
	})
}
