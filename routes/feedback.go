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
	"git.licolas.net/delegit/delegit/models"
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
	var feedback models.Feedback
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
	var feedback *models.Feedback
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
	var feedback *models.Feedback
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
}

func optionsFeedbackEntry(ctx *gin.Context) {
	ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET, PUT, PATCH, DELETE, OPTIONS")
}

func updateFeedbackUpvotes(ctx *gin.Context) {
	_id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	id := uint(_id)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var votes int
	if err := ctx.ShouldBind(&votes); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var feedback *models.Feedback
	switch votes {
	case 1:
		feedback, err = db.IncrementFeedbackUpvotes(id)
	case -1:
		feedback, err = db.DecrementFeedbackUpvotes(id)
	default:
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "unknown increment"})
		return
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, feedback)
}

func updateFeedbackDownvotes(ctx *gin.Context) {
	_id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	id := uint(_id)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var votes int
	if err := ctx.ShouldBind(&votes); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var feedback *models.Feedback
	switch votes {
	case 1:
		feedback, err = db.IncrementFeedbackDownvotes(id)
	case -1:
		feedback, err = db.DecrementFeedbackDownvotes(id)
	default:
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "unknown increment"})
		return
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, feedback)
}

func RegisterFeedbackEndpoints(database *database.Database, router *gin.Engine) {
	db = database

	list := router.Group("/feedback")
	list.Use(CommonHeaders, optionsFeedbackList)
	list.GET("/", getAllFeedback)
	list.POST("/", postFeedback)
	list.OPTIONS("/", Terminate)

	entry := router.Group("/feedback/:id")
	entry.Use(optionsFeedbackEntry)
	entry.GET("/", getFeedback)
	entry.PATCH("/upvote", updateFeedbackUpvotes)
	entry.PATCH("/downvote", updateFeedbackDownvotes)
	entry.PUT("/", putFeedback)
	entry.DELETE("/", deleteFeedback)
	entry.OPTIONS("/", Terminate)
}
