package AuthController

import (
	"net/http"
	"time"

	model "github.com/E-Furqan/Food-Delivery-System/Authentication_Service/Model"
	utils "github.com/E-Furqan/Food-Delivery-System/Authentication_Service/Utils"
	"github.com/gin-gonic/gin"
)

func GenerateTokens(c *gin.Context) {
	var input model.Input
	var accessClaims model.Claims
	var refreshClaims model.Claims

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessClaims, refreshClaims = utils.CreateClaim(input)
	accessTokenString, refreshTokenString, err := utils.GenerateTokens(accessClaims, refreshClaims)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	var tokens model.Tokens
	tokens.AccessToken = accessTokenString
	tokens.RefreshToken = refreshTokenString
	tokens.Expiration = time.Now().Add(24 * time.Hour).Unix()

	c.JSON(http.StatusOK, tokens)

}

func ReFreshToken(c *gin.Context) {
	var input model.Tokens

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, err := utils.RefreshToken(input.RefreshToken, c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	var tokens model.Tokens
	tokens.AccessToken = accessToken
	tokens.RefreshToken = input.RefreshToken
	tokens.Expiration = time.Now().Add(30 * time.Minute).Unix()

	c.JSON(http.StatusOK, tokens)
}
