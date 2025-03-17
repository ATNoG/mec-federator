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

	tokenEndpoint := config.AppConfig.KeycloakTokenEndpoint
	clientId := c.PostForm("client_id")
	clientSecrect := c.PostForm("client_secret")
	if tokenEndpoint == "" || clientId == "" || clientSecrect == "" {
		log.Fatal("KeycloakAuthentication", "Missing configuration variables")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Missing configuration variables"})
		return
	}

	log.Print("BeginAuthController - Building request to Keycloak")

	reqBody := strings.NewReader("grant_type=client_credentials&client_id=" + clientId + "&client_secret=" + clientSecrect)
	req, err := http.NewRequest("POST", tokenEndpoint, reqBody)
	if err != nil {
		slog.Error("KeycloakAuthentication", "Error creating request: %v", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating authentication request"})
		return
	}

	log.Print("BeginAuthController - Sending request")

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
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

	log.Print("BeginAuthController - Decoding response")
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		slog.Error("KeycloakAuthentication", "Error decoding response: %v", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error decoding authentication response"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving access token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accessToken": accessToken.Token,
		"expiresAt":   accessToken.ExpiresAt.Format(time.RFC3339),
	})
}
