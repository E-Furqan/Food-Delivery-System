package utils

import (
	"fmt"

	payload "github.com/E-Furqan/Food-Delivery-System/Payload"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Claims struct {
	payload.RestaurantClaim
	jwt.StandardClaims
}

func Verification(c *gin.Context) (any, error) {
	RestaurantID, exists := c.Get("RestaurantID")
	if !exists {
		return 0, fmt.Errorf("restaurant id does not exist")
	}

	RestaurantID, ok := RestaurantID.(uint)
	if !ok {
		return 0, fmt.Errorf("invalid restaurant id")
	}
	return RestaurantID, nil
}

func GenerateResponse(httpStatusCode int, c *gin.Context, title1 string, message1 string, title2 string, input interface{}) {
	response := gin.H{
		title1: message1,
	}

	if title2 != "" && input != nil {
		response[title2] = input
	}

	c.JSON(httpStatusCode, response)
}
