package middlewares

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/qhh/prjEcom/pkg/utils"
)

const (
	authorizationHeaderKey  = "Authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

// AuthMiddleware creates a middleware for authorization
func AuthMiddleware(jwtMaker *utils.JWTMaker) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header is not provided"})
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unsupported authorization type"})
			return
		}

		accessToken := fields[1]
		payload, err := jwtMaker.VerifyToken(accessToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.Set(authorizationPayloadKey, payload)
		c.Next()
	}
}

// RoleMiddleware creates a middleware for role-based authorization
func RoleMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		payload, exists := c.Get(authorizationPayloadKey)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		userPayload, ok := payload.(*utils.Payload)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "invalid token payload"})
			return
		}

		// Check if the user has the required role
		roleMatch := false
		for _, role := range roles {
			if userPayload.Role == role {
				roleMatch = true
				break
			}
		}

		if !roleMatch {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "access denied: required role not found"})
			return
		}

		c.Next()
	}
}

// GetAuthPayload extracts the token payload from the context
func GetAuthPayload(c *gin.Context) (*utils.Payload, error) {
	payload, exists := c.Get(authorizationPayloadKey)
	if !exists {
		return nil, errors.New("authorization payload not found")
	}

	userPayload, ok := payload.(*utils.Payload)
	if !ok {
		return nil, errors.New("invalid authorization payload")
	}

	return userPayload, nil
}
