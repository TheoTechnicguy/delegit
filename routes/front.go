package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterFrontendRoutes(router *gin.Engine) {

	router.Static("/styles", "public/styles")
	router.Static("/js", "public/js")
	router.Static("/images", "public/images")
	router.Static("/fonts", "public/fonts")
	router.LoadHTMLFiles("public/index.html")

	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", nil)
	})
}
