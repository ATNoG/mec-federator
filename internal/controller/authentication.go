package controller

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mankings/mec-federator/internal/config"
	"github.com/mankings/mec-federator/internal/models"
	"github.com/mankings/mec-federator/internal/services"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AuthController struct {
	authService *services.AuthServiceImpl
	mongoClient *mongo.Client
}

func NewAuthController(authService *services.AuthServiceImpl, mongoClient *mongo.Client) *AuthController {
	return &AuthController{
		authService: authService,
		mongoClient: mongoClient,
	}
}

// Retrieve the access token from the cookie and check if it is valid
func (ac *AuthController) IsAuthenticated(c *gin.Context) bool {
	tokenStr, err := c.Cookie("accessToken")
	if err != nil {
		return false
	}
	_, err = ac.authService.QueryAccessToken(tokenStr)
	return err == nil
}

func (ac *AuthController) BeginAuthController(c *gin.Context) {
	slog.Info("Gathering OAuth2.0 configuration variables")
	tokenEndpoint := config.AppConfig.OAuth2TokenEndpoint
	clientId := config.AppConfig.OAuth2ClientId    // this might need to come from the request itself
	clientSecrect := config.AppConfig.OAuth2Secret // this might need to come from the request itself

	slog.Info("Building request to Keycloak")
	reqBody := strings.NewReader("grant_type=client_credentials&client_id=" + clientId + "&client_secret=" + clientSecrect)
	req, err := http.NewRequest("POST", tokenEndpoint, reqBody)
	if err != nil {
		slog.Error("KeycloakAuthentication", "Error creating request: %v", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating authentication request"})
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	slog.Info("Sending request")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error("KeycloakAuthentication", "Error sending request: %v", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error sending authentication request"})
		return
	}

	var tokenResponse struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}

	slog.Info("Decoding response")
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		slog.Error("KeycloakAuthentication", "Error decoding response: %v", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error decoding authentication response"})
		return
	}

	expiresAt := time.Now().Add(time.Duration(tokenResponse.ExpiresIn) * time.Second)

	slog.Info("Building AccessToken")
	accessToken := models.AccessToken{
		Token:     tokenResponse.AccessToken,
		ExpiresAt: expiresAt,
	}

	slog.Info("Saving AccessToken")
	err = ac.authService.SaveAccessToken(accessToken)
	if err != nil {
		slog.Error("KeycloakAuthentication", "SaveAccessToken error: %v", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving access token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accessToken": accessToken.Token,
		"expiresAt":   accessToken.ExpiresAt.Format(time.RFC3339),
	})
}
