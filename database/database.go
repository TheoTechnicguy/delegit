/**
 * file: database/database.go
 * author: theo technicguy
 * license: apache-2.0
 *
 * The database package provides data persistance
 * for the entire application.
 */

// The database package provides persistance for application data
package database

import (
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	ErrInvalidFeedback error = errors.New("invalid feedback")

	ErrInvalidDatabaseKind error = errors.New("invalid database type")
)

type Database struct {
	db *gorm.DB
}

func NewDatabase(kind, dsn string) (*Database, error) {
	var dialect gorm.Dialector
	switch kind {
	case "sqlite":
		dialect = sqlite.Open(dsn)
	default:
		return nil, ErrInvalidDatabaseKind
	}

	db, err := NewDatabaseFromDialector(dialect, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate()
	return db, err
}

func NewDatabaseFromDialector(dialect gorm.Dialector, config *gorm.Config) (*Database, error) {
	db := new(Database)

	var err error
	db.db, err = gorm.Open(dialect, config)
	if err != nil {
		return nil, err
	}

	return db, err
}

func (db *Database) AutoMigrate() error {
	t := []any{
		Feedback{},
	}

	for _, v := range t {
		if err := db.db.AutoMigrate(&v); err != nil {
			return err
		}
	}

	return nil
}

func (db *Database) GetAllFeedback() (f []*Feedback, err error) {
	err = db.db.Find(&f).Error
	return
}

func (db *Database) GetFeedback(id uint) (*Feedback, error) {
	f := &Feedback{}
	if r := db.db.First(&f, id); r.Error != nil {
		return nil, r.Error
	}

	return f, nil
}

func (db *Database) AddFeedback(feedback *Feedback) (*Feedback, error) {
	log.Debug().Any("feedback", feedback).Msg("adding feedback")
	// Check mandatory fields
	if feedback.Course == "" {
		log.Debug().Err(ErrInvalidFeedback).Msg("empty course field")
		return nil, ErrInvalidFeedback
	}
	if feedback.Feedback == "" {
		log.Debug().Err(ErrInvalidFeedback).Msg("empty feedback field")
		return nil, ErrInvalidFeedback
	}

	// Sanitize data
	feedback.ID = 0
	feedback.Upvotes = 0
	feedback.Downvotes = 0

	if r := db.db.Create(feedback); r.Error != nil {
		return nil, r.Error
	}

	return feedback, nil
}

func (db *Database) UpdateFeedback(feedback *Feedback) (*Feedback, error) {
	// Check mandatory fields
	if feedback.Course == "" {
		return nil, ErrInvalidFeedback
	}
	if feedback.Feedback == "" {
		return nil, ErrInvalidFeedback
	}

	if r := db.db.Save(feedback); r.Error != nil {
		return nil, r.Error
	}

	return feedback, nil
}

func (db *Database) DeleteFeedback(feedback *Feedback) error {
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

func (db *Database) updateFeedbackAppreciation(id uint, appr string, change int) (*Feedback, error) {
	f := &Feedback{}
	tx := db.db.Begin()
	defer tx.Rollback()

	if r := tx.First(&f, id); r.Error != nil {
		return nil, r.Error
	}

	var votes *uint
	switch appr {
	case "upvotes":
		votes = &f.Upvotes
	case "downvotes":
		votes = &f.Downvotes
	default:
		return nil, fmt.Errorf("unknown appreciation")
	}

	switch change {
	case 1:
		if *votes > 2000 {
			return nil, fmt.Errorf("out of bounds")
		} else {
			*votes += 1
		}
	case -1:
		if *votes == 0 {
			return nil, fmt.Errorf("out of bounds")
		} else {
			*votes -= 1
		}
	}

	if r := tx.Save(f); r.Error != nil {
		return nil, r.Error
	}

	tx.Commit()

	return f, tx.Error
}

func (db *Database) IncrementFeedbackUpvotes(id uint) (*Feedback, error) {
	return db.updateFeedbackAppreciation(id, "upvotes", 1)
}

func (db *Database) DecrementFeedbackUpvotes(id uint) (*Feedback, error) {
	return db.updateFeedbackAppreciation(id, "upvotes", -1)
}

func (db *Database) IncrementFeedbackDownvotes(id uint) (*Feedback, error) {
	return db.updateFeedbackAppreciation(id, "downvotes", 1)
}

func (db *Database) DecrementFeedbackDownvotes(id uint) (*Feedback, error) {
	return db.updateFeedbackAppreciation(id, "downvotes", -1)
}
