/**
 * file: routes/api/feedback.go
 * author: theo technicguy
 * license: apache-2.0
 *
 * This file contains all routes leading to
 * the feedback endpoints.
 */

package api

import (
	"fmt"
	"net/http"
	"strconv"

	"git.licolas.net/delegit/delegit/database"
	"git.licolas.net/delegit/delegit/logic"
	"git.licolas.net/delegit/delegit/models"
	"git.licolas.net/delegit/delegit/routes/middleware"
	"git.licolas.net/delegit/delegit/uxerrors"
	"github.com/gin-gonic/gin"
)

var (
	db *database.Database
)

func handleError(ctx *gin.Context, err error) {
	var o any
	var status int
	switch v := err.(type) {
	case uxerrors.Errors:
		status = v.Status
		o = v.ToMap(false)
	case uxerrors.Error:
		es := uxerrors.NewErrors(http.StatusInternalServerError).Append(v)
		handleError(ctx, es)
		return
	case error:
		es := uxerrors.NewErrors(http.StatusInternalServerError).AppendNew(v)
		handleError(ctx, es)
		return
	default:
		es := uxerrors.NewErrors(http.StatusInternalServerError).AppendNew(fmt.Errorf("Another, unknown error occurred"))
		handleError(ctx, es)
		return
	}

	ctx.AbortWithStatusJSON(status, o)
}

func feedbackBindError(err error) error {
	uxe := uxerrors.New(err)
	uxe.Summary = "Could not parse your feedback"
	uxe.Detail = "The feedback you gave could not be parsed. This usually means that you did not respect the specification. Check your input and try again."
	return uxerrors.NewErrors(http.StatusBadRequest).Append(uxe)
}

func getAllFeedback(ctx *gin.Context) {
	feedback, err := logic.GetAllFeedback()
	if err != nil {
		handleError(ctx, err)
	}

	ctx.JSON(http.StatusOK, feedback)
}

func postFeedback(ctx *gin.Context) {
	var feedback models.Feedback
	if err := ctx.ShouldBind(&feedback); err != nil {
		handleError(ctx, feedbackBindError(err))
		return
	}

	f, err := logic.AddFeedback(&feedback)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, f)
}

func getFeedback(ctx *gin.Context) {
	_id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	id := uint(_id)

	if err != nil {

		handleError(ctx, feedbackBindError(err))
		return
	}

	feedback, err := logic.GetFeedback(id)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, feedback)
}

func putFeedback(ctx *gin.Context) {
	var feedback *models.Feedback
	if err := ctx.ShouldBind(&feedback); err != nil {
		handleError(ctx, feedbackBindError(err))
		return
	}

	feedback, err := logic.UpdateFeedback(feedback)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, feedback)
}

func deleteFeedback(ctx *gin.Context) {
	var feedback *models.Feedback
	if err := ctx.ShouldBind(&feedback); err != nil {
		handleError(ctx, feedbackBindError(err))
		return
	}

	feedback, err := logic.DeleteFeedback(feedback)
	if err != nil {
		handleError(ctx, err)
		return
	}

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
		handleError(ctx, feedbackBindError(err))
		return
	}

	var votes int
	if err := ctx.ShouldBind(&votes); err != nil {
		handleError(ctx, feedbackBindError(err))
		return
	}

	feedback, err := logic.UpdateFeedbackUpvotes(id, votes)
	if err != nil {
		handleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, feedback)
}

func updateFeedbackDownvotes(ctx *gin.Context) {
	_id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	id := uint(_id)

	if err != nil {
		handleError(ctx, feedbackBindError(err))
		return
	}

	var votes int
	if err := ctx.ShouldBind(&votes); err != nil {
		handleError(ctx, feedbackBindError(err))
		return
	}

	feedback, err := logic.UpdateFeedbackDownvotes(id, votes)
	if err != nil {
		handleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, feedback)
}

func RegisterFeedbackEndpoints(database *database.Database, router *gin.RouterGroup) {
	db = database

	list := router.Group("/feedback")
	list.Use(optionsFeedbackList)
	list.GET("/", getAllFeedback)
	list.POST("/", postFeedback)
	list.OPTIONS("/", middleware.Terminate)

	entry := router.Group("/feedback/:id")
	entry.Use(optionsFeedbackEntry)
	entry.GET("/", getFeedback)
	entry.PATCH("/upvote", updateFeedbackUpvotes)
	entry.PATCH("/downvote", updateFeedbackDownvotes)
	entry.PUT("/", putFeedback)
	entry.DELETE("/", deleteFeedback)
	entry.OPTIONS("/", middleware.Terminate)
}
