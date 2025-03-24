package controller

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mankings/mec-federator/internal/config"
	"github.com/mankings/mec-federator/internal/models"
	"github.com/mankings/mec-federator/internal/services"
	"github.com/mankings/mec-federator/internal/utils"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AuthController struct {
	authService *services.AuthService
	mongoClient *mongo.Client
}

func NewAuthController(authService *services.AuthService, mongoClient *mongo.Client) *AuthController {
	return &AuthController{
		authService: authService,
		mongoClient: mongoClient,
	}
}

func (ac *AuthController) IssueAccessTokenController(c *gin.Context) {
	log.Print("BeginAuthController - Gathering OAuth2.0 configuration variables")

	type requestBody struct {
		ClientId     string `json:"clientId"`
		ClientSecret string `json:"clientSecret"`
	}

	var expectedBody requestBody
	if err := c.ShouldBindJSON(&expectedBody); err != nil {
		problemDetails := utils.NewProblemDetails(http.StatusBadRequest)
		c.AbortWithStatusJSON(http.StatusBadRequest, problemDetails)
		return
	}

	tokenEndpoint := config.AppConfig.KeycloakTokenEndpoint
	clientId := expectedBody.ClientId
	clientSecret := expectedBody.ClientSecret

	switch {
	case tokenEndpoint == "":
		log.Fatal("KeycloakAuthentication", "Missing tokenEndpoint configuration variable")
		problemDetails := utils.NewProblemDetails(http.StatusInternalServerError)
		problemDetails.Detail = "Missing tokenEndpoint configuration variable"
		c.JSON(http.StatusInternalServerError, problemDetails)
		return
	case clientId == "":
		log.Fatal("KeycloakAuthentication", "Missing clientId body parameter")
		problemDetails := utils.NewProblemDetails(http.StatusInternalServerError)
		problemDetails.Detail = "Missing clientId body parameter"
		c.JSON(http.StatusInternalServerError, problemDetails)
		return
	case clientSecret == "":
		log.Fatal("KeycloakAuthentication", "Missing clientSecret body parameter")
		problemDetails := utils.NewProblemDetails(http.StatusInternalServerError)
		problemDetails.Detail = "Missing clientSecret body parameter"
		c.JSON(http.StatusInternalServerError, problemDetails)
		return
	}

	log.Print("BeginAuthController - Building request to Keycloak")

	reqBody := strings.NewReader("grant_type=client_credentials&client_id=" + clientId + "&client_secret=" + clientSecret)
	req, err := http.NewRequest("POST", tokenEndpoint, reqBody)
	if err != nil {
		slog.Error("KeycloakAuthentication", "Error creating request: %v", err.Error())
		problemDetails := utils.NewProblemDetails(http.StatusInternalServerError)
		problemDetails.Detail = "Error creating authentication request"
		c.JSON(http.StatusInternalServerError, problemDetails)
		return
	}

	log.Print("BeginAuthController - Sending request")

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error("KeycloakAuthentication", "Error sending request: %v", err.Error())
		problemDetails := utils.NewProblemDetails(http.StatusInternalServerError)
		problemDetails.Detail = "Error sending authentication request"
		c.JSON(http.StatusInternalServerError, problemDetails)
		return
	}

	log.Print("BeginAuthController - Decoding response")

	var tokenResponse struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		slog.Error("KeycloakAuthentication", "Error decoding response: %v", err.Error())
		problemDetails := utils.NewProblemDetails(http.StatusInternalServerError)
		problemDetails.Detail = "Error decoding authentication response"
		c.JSON(http.StatusInternalServerError, problemDetails)
		return
	}

	expiresAt := time.Now().Add(time.Duration(tokenResponse.ExpiresIn) * time.Second)

	log.Print("BeginAuthController - Building AccessToken")

	accessToken := models.AccessToken{
		Token:     tokenResponse.AccessToken,
		ExpiresAt: expiresAt,
	}

	log.Print("BeginAuthController - Saving AccessToken")

	err = ac.authService.SaveAccessToken(accessToken)
	if err != nil {
		slog.Error("KeycloakAuthentication", "SaveAccessToken error: %v", err.Error())
		problemDetails := utils.NewProblemDetails(http.StatusInternalServerError)
		problemDetails.Detail = "Error saving access token"
		c.JSON(http.StatusInternalServerError, problemDetails)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accessToken": accessToken.Token,
		"expiresAt":   accessToken.ExpiresAt.Format(time.RFC3339),
	})
}
