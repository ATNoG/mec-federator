package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
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

	c.JSON(http.StatusOK, gin.H{"message": "Begin Auth"})
}
