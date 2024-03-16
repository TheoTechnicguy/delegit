/**
 * file: routes/middleware/common.go
 * author: theo technicguy
 * license: apache-2.0
 *
 * This file contains common middleware.
 */

package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CommonHeaders is a common middleware inserting common headers
// that should be included in every response from the server.
func CommonHeaders(ctx *gin.Context) {
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	ctx.Writer.Header().Set("Access-Control-Max-Age", "300")
	ctx.Writer.Header().Set("X-Content-Type-Options", "nosniff")
	ctx.Next()
}

func Terminate(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusNoContent)
}
