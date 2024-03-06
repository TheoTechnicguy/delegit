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

	"git.licolas.net/delegit/delegit/database"
	"git.licolas.net/delegit/delegit/models"
	"git.licolas.net/delegit/delegit/validators"
	"github.com/go-playground/validator/v10"
)

var (
	db *database.Database
)

func validateFeedback(f *models.Feedback) error {
	v := validator.New()
	v.RegisterValidation("iscourse", validators.IsCourse, false)
	return v.Struct(f)
}

func sanitizeFeedback(f *models.Feedback) {
	f.ID = 0
	f.Upvotes = 0
	f.Downvotes = 0
}

func GetAllFeedback() ([]*models.Feedback, error) {
	return db.GetAllFeedback()
}

func GetFeedback(id uint) (*models.Feedback, error) {
	return db.GetFeedback(id)
}

func AddFeedback(f *models.Feedback) (*models.Feedback, error) {
	// Sanitize data
	sanitizeFeedback(f)

	if err := validateFeedback(f); err != nil {
		return nil, err
	}

	return db.AddFeedback(f)
}

func UpdateFeedback(f *models.Feedback) (*models.Feedback, error) {
	if err := validateFeedback(f); err != nil {
		return nil, err
	}

	return db.UpdateFeedback(f)
}

func DeleteFeedback(f *models.Feedback) (*models.Feedback, error) {
	if err := validateFeedback(f); err != nil {
		return nil, err
	}

	err := db.DeleteFeedback(f)
	f.ID = 0
	return f, err
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
		return nil, fmt.Errorf("unknown increment")
	}

	return feedback, err
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
		return nil, fmt.Errorf("unknown increment")
	}

	return feedback, err
}
func Setup(database *database.Database) {
	db = database
}