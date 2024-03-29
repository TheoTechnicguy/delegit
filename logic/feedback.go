/**
 * file: logic/feedback.go
 * author: theo technicguy
 * license: apache-2.0
 *
 * This file contains the logic (at least most of it)
 * concerning feedback.
 */

package logic

import (
	"fmt"
	"net/http"

	"git.licolas.net/delegit/delegit/database"
	"git.licolas.net/delegit/delegit/models"
	"git.licolas.net/delegit/delegit/uxerrors"
	"git.licolas.net/delegit/delegit/validators"
)

var (
	db *database.Database
)

func sanitizeFeedback(f *models.Feedback) {
	f.ID = 0
	f.Upvotes = 0
	f.Downvotes = 0
}

func GetAllFeedback() ([]*models.Feedback, error) {
	fs, err := db.GetAllFeedback()
	if err != nil {
		return nil, handleDatabaseError(err)
	}

	return fs, nil
}

func GetFeedback(id uint) (*models.Feedback, error) {
	f, err := db.GetFeedback(id)
	if err != nil {
		return nil, handleDatabaseError(err)
	}

	return f, nil
}

func AddFeedback(f *models.Feedback) (*models.Feedback, error) {
	sanitizeFeedback(f)

	if err := validators.ValidateFeedback(f); err != nil {
		return nil, err
	}

	r, err := db.AddFeedback(f)
	if err != nil {
		return nil, handleDatabaseError(err)
	}

	return r, nil
}

func UpdateFeedback(f *models.Feedback) (*models.Feedback, error) {
	if err := validators.ValidateFeedback(f); err != nil {
		return nil, err
	}

	r, err := db.UpdateFeedback(f)
	if err != nil {
		return nil, handleDatabaseError(err)
	}

	return r, nil
}

func DeleteFeedback(f *models.Feedback) (*models.Feedback, error) {
	if err := validators.ValidateFeedback(f); err != nil {
		return nil, err
	}

	err := db.DeleteFeedback(f)
	f.ID = 0
	return f, handleDatabaseError(err)
}

func UpdateFeedbackUpvotes(id uint, votes int) (*models.Feedback, error) {
	var feedback *models.Feedback
	var err error
	switch votes {
	case 1:
		feedback, err = db.IncrementFeedbackUpvotes(id)
	case -1:
		feedback, err = db.DecrementFeedbackUpvotes(id)
	default:
		uxe := uxerrors.New(fmt.Errorf("unknown increment"))
		uxe.Summary = "The increment you are attempting to do is invalid"
		uxe.Detail = fmt.Sprintf("You are trying to increment feedback upvotes by %d, but only 1 or -1 is allowed. Correct the values and try again.", votes)
		return nil, uxerrors.NewErrors(http.StatusBadRequest).Append(uxe)
	}

	return feedback, handleDatabaseError(err)
}

func UpdateFeedbackDownvotes(id uint, votes int) (*models.Feedback, error) {
	var feedback *models.Feedback
	var err error
	switch votes {
	case 1:
		feedback, err = db.IncrementFeedbackDownvotes(id)
	case -1:
		feedback, err = db.DecrementFeedbackDownvotes(id)
	default:
		uxe := uxerrors.New(fmt.Errorf("unknown increment"))
		uxe.Summary = "The increment you are attempting to do is invalid"
		uxe.Detail = fmt.Sprintf("You are trying to increment feedback downvotes by %d, but only 1 or -1 is allowed. Correct the values and try again.", votes)
		return nil, uxerrors.NewErrors(http.StatusBadRequest).Append(uxe)
	}

	return feedback, handleDatabaseError(err)
}

func Setup(database *database.Database) {
	db = database
}
