package AuthController

import (
	"net/http"
	"time"

	payload "github.com/E-Furqan/Food-Delivery-System/Authentication_Service/Payload"
	utils "github.com/E-Furqan/Food-Delivery-System/Authentication_Service/Utils"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var input payload.Input
	var accessClaims payload.Claims
	var refreshClaims payload.Claims

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

	var tokens payload.Tokens
	tokens.AccessToken = accessTokenString
	tokens.RefreshToken = refreshTokenString
	tokens.Expiration = time.Now().Add(24 * time.Hour).Unix()

	c.JSON(http.StatusOK, tokens)

}

func ReFreshToken(c *gin.Context) {
	var input payload.Tokens

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, err := utils.RefreshToken(input.RefreshToken, c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate the token"})
		return
	}

	var tokens payload.Tokens
	tokens.AccessToken = accessToken
	tokens.RefreshToken = input.RefreshToken
	tokens.Expiration = time.Now().Add(30 * time.Minute).Unix()

	c.JSON(http.StatusOK, tokens)
}
