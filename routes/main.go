/**
 * file: routes/main.go
 * author: theo technicguy
 * license: apache-2.0
 */

package routes

import (
	"git.licolas.net/delegit/delegit/database"
	"git.licolas.net/delegit/delegit/routes/api"
	"git.licolas.net/delegit/delegit/routes/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, db *database.Database) {
	router.Use(middleware.CommonHeaders)

	apiRouter := router.Group("/api")
	api.RegisterFeedbackEndpoints(db, apiRouter)

	RegisterFrontendRoutes(router)
}
