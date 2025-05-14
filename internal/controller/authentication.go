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
	"github.com/mankings/mec-federator/internal/models/dto"
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

// @Summary Issue Access Token
// @Description Issue an access token from Keycloak
// @Tags Authentication
// @Accept json
// @Produce json
// @Param body body dto.AccessTokenRequestData true "Client ID and Client Secret"
// @Success 200 {object} models.AccessToken
// @Failure 400 {object} models.ProblemDetails
// @Failure 500 {object} models.ProblemDetails
// @Router /auth/token [post]
func (ac *AuthController) IssueAccessTokenController(c *gin.Context) {
	log.Print("BeginAuthController - Gathering OAuth2.0 configuration variables")

	var expectedBody dto.AccessTokenRequestData
	if err := c.ShouldBindJSON(&expectedBody); err != nil {
		utils.HandleProblem(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	tokenEndpoint := config.AppConfig.KeycloakTokenEndpoint
	clientId := expectedBody.ClientId
	clientSecret := expectedBody.ClientSecret

	switch {
	case tokenEndpoint == "":
		log.Fatal("KeycloakAuthentication", "Missing tokenEndpoint configuration variable")
		utils.HandleProblem(c, http.StatusInternalServerError, "Missing tokenEndpoint configuration variable")
		return
	case clientId == "":
		log.Fatal("KeycloakAuthentication", "Missing clientId body parameter")
		utils.HandleProblem(c, http.StatusInternalServerError, "Missing clientId body parameter")
		return
	case clientSecret == "":
		log.Fatal("KeycloakAuthentication", "Missing clientSecret body parameter")
		utils.HandleProblem(c, http.StatusInternalServerError, "Missing clientSecret body parameter")
		return
	}

	log.Print("BeginAuthController - Building request to Keycloak")
	reqBody := strings.NewReader("grant_type=client_credentials&client_id=" + clientId + "&client_secret=" + clientSecret)
	req, err := http.NewRequest("POST", tokenEndpoint, reqBody)
	if err != nil {
		slog.Error("KeycloakAuthentication", "Error creating request: %v", err.Error())
		utils.HandleProblem(c, http.StatusInternalServerError, "Error creating authentication request")
		return
	}

	log.Print("BeginAuthController - Sending request")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		utils.HandleProblem(c, http.StatusInternalServerError, "Error sending authentication request")
		slog.Error("KeycloakAuthentication", "Error sending request: %v", err.Error())
		return
	}

	log.Print("BeginAuthController - Decoding response")

	var tokenResponse struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		slog.Error("KeycloakAuthentication", "Error decoding response: %v", err.Error())
		utils.HandleProblem(c, http.StatusInternalServerError, "Error decoding authentication response")
		return
	}

	log.Print("BeginAuthController - Building AccessToken")
	expiresAt := time.Now().Add(time.Duration(tokenResponse.ExpiresIn) * time.Second)
	accessToken := models.AccessToken{
		AccessToken: tokenResponse.AccessToken,
		ExpiresAt:   expiresAt,
	}

	log.Print("BeginAuthController - Saving AccessToken")
	err = ac.authService.SaveAccessToken(accessToken)
	if err != nil {
		slog.Error("KeycloakAuthentication", "SaveAccessToken error: %v", err.Error())
		utils.HandleProblem(c, http.StatusInternalServerError, "Error saving access token")
		return
	}

	log.Print("BeginAuthController - Responding OK with AccessToken")
	c.JSON(http.StatusOK, accessToken)
}
