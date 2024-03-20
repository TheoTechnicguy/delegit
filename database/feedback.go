/**
 * file: database/feedback.go
 * author: theo technicguy
 * license: apache-2.0
 *
 * This file contains the feedback database
 * logic for the data persistance plane.
 */

package database

import (
	"git.licolas.net/delegit/delegit/models"
)

func (db *Database) GetAllFeedback() (f []*models.Feedback, err error) {
	err = db.db.Find(&f).Error
	return
}

func (db *Database) GetFeedback(id uint) (*models.Feedback, error) {
	f := new(models.Feedback)
	if r := db.db.First(&f, id); r.Error != nil {
		return nil, r.Error
	}

	return f, nil
}

func (db *Database) AddFeedback(feedback *models.Feedback) (*models.Feedback, error) {
	if r := db.db.Create(feedback); r.Error != nil {
		return nil, r.Error
	}

	return feedback, nil
}

func (db *Database) UpdateFeedback(feedback *models.Feedback) (*models.Feedback, error) {
	if r := db.db.Save(feedback); r.Error != nil {
		return nil, r.Error
	}

	return feedback, nil
}

func (db *Database) DeleteFeedback(feedback *models.Feedback) error {
	if r := db.db.
		Where("course = ?", feedback.Course).
		Where("feedback = ?", feedback.Feedback).
		Where("upvotes = ?", feedback.Upvotes).
		Where("downvotes = ?", feedback.Downvotes).
		Delete(feedback); r.Error != nil {
		return r.Error
	}
	return nil
}

// updateFeedbackAppreciation is a helper function to update the feedback
// appreciation counters. It is used to increment and decrement the upvotes
// and downvotes of a feedback.
// It takes the id of the feedback to update and a mutator function that
// takes a pointer to the feedback and modifies it in place.
func (db *Database) updateFeedbackAppreciation(id uint, mutator func(*models.Feedback)) (*models.Feedback, error) {
	f := new(models.Feedback)
	tx := db.db.Begin()
	defer tx.Rollback()

	if r := tx.First(&f, id); r.Error != nil {
		return nil, r.Error
	}

	mutator(f)

	if r := tx.Save(f); r.Error != nil {
		return nil, r.Error
	}

	tx.Commit()

	return f, tx.Error
}

func (db *Database) IncrementFeedbackUpvotes(id uint) (*models.Feedback, error) {
	mutator := func(f *models.Feedback) {
		f.Upvotes++
	}

	return db.updateFeedbackAppreciation(id, mutator)
}

func (db *Database) DecrementFeedbackUpvotes(id uint) (*models.Feedback, error) {
	mutator := func(f *models.Feedback) {
		f.Upvotes--
	}

	return db.updateFeedbackAppreciation(id, mutator)
}

func (db *Database) IncrementFeedbackDownvotes(id uint) (*models.Feedback, error) {
	mutator := func(f *models.Feedback) {
		f.Downvotes++
	}

	return db.updateFeedbackAppreciation(id, mutator)
}

func (db *Database) DecrementFeedbackDownvotes(id uint) (*models.Feedback, error) {
	mutator := func(f *models.Feedback) {
		f.Downvotes--
	}

	return db.updateFeedbackAppreciation(id, mutator)
}
