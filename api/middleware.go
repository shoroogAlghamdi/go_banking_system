package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shoroogAlghamdi/banking_system/token"
)

const (
	authorizationHeaderKey = "authorization"
	authorizationTypeAllowed = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errResponse(err))
			return
		}

		// splits string based on space
		fields := strings.Fields(authorizationHeader)

		// 2 one for bearer one for token
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errResponse(err))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeAllowed {
			err := fmt.Errorf("invalid authorization type %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errResponse(err))
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errResponse(err))
			return
		}

		// Add payload key to the context to be used in handler and othe rmiddleware
		ctx.Set(authorizationPayloadKey, payload)

		// call next to either go back to the handler or other middleware
		ctx.Next()
	}
}
