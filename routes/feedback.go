/**
 * file: router/feedback.go
 * author: theo technicguy
 * license: apache-2.0
 *
 * This file contains all routes leading to
 * the feedback endpoints.
 */

package routes

import (
	"net/http"
	"strconv"

	"git.licolas.net/delegit/delegit/database"
	"github.com/gin-gonic/gin"
)

var (
	db *database.Database
)

func getAllFeedback(ctx *gin.Context) {
	feedback, err := db.GetAllFeedback()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err})
	}
	ctx.JSON(http.StatusOK, feedback)
}

func postFeedback(ctx *gin.Context) {
	var feedback database.Feedback
	if err := ctx.ShouldBind(&feedback); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	f, err := db.AddFeedback(&feedback)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, f)
}

func getFeedback(ctx *gin.Context) {
	_id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	id := uint(_id)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	feedback, err := db.GetFeedback(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, feedback)
}

func putFeedback(ctx *gin.Context) {
	var feedback *database.Feedback
	if err := ctx.ShouldBind(&feedback); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	feedback, err := db.UpdateFeedback(feedback)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, feedback)
}

func deleteFeedback(ctx *gin.Context) {
	var feedback *database.Feedback
	if err := ctx.ShouldBind(&feedback); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := db.DeleteFeedback(feedback); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	feedback.ID = 0
	ctx.JSON(http.StatusOK, feedback)
}

func optionsFeedbackList(ctx *gin.Context) {
	ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	ctx.AbortWithStatus(http.StatusNoContent)
}

func optionsFeedbackEntry(ctx *gin.Context) {
	ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET, PUT, DELETE, OPTIONS")
	ctx.AbortWithStatus(http.StatusNoContent)
}

func RegisterFeedbackEndpoints(database *database.Database, router *gin.Engine) {
	db = database

	router.Use(CommonHeaders)

	router.GET("/feedback", getAllFeedback)
	router.GET("/feedback/:id", getFeedback)
	router.POST("/feedback", postFeedback)
	router.PUT("/feedback/:id", putFeedback)
	router.DELETE("/feedback/:id", deleteFeedback)
	router.OPTIONS("/feedback", optionsFeedbackList)
	router.OPTIONS("/feedback/:id", optionsFeedbackEntry)
}
