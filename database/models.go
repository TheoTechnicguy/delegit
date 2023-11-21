package database

import (
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

// The Feedback structure represents a feedback, comment, or note
// left by users on the page.
type Feedback struct {
	// Each feedback is identified uniquely by their ID.
	// The ID is set by the database, who has full authority over
	// identity value attribution. Once set, the ID is not allowed
	// to change.
	// An ID of `0` (int default) is considered to be a blank ID,
	// all other values refer to one entry in the database.
	ID uint `gorm:"<-:create;primaryKey" json:"ID" validate:"omitempty,min=1"`

	// The course field identifies the course related to the
	// feedback. It is required and must conform to the university's
	// course code formatting, case insensitive.
	Course string `gorm:"<-;size:10;not null" json:"Course" validate:"required,iscourse"`

	// The feedback is the main content of the submission. It is
	// required, contains at least 25 and at most 2000 alphanumeric
	// or unicode characters.
	// NOTE: The upper bound may change in the future.
	Feedback string `gorm:"<-;not null" json:"Feedback" validate:"required,min=25,max=2000,alphanumunicode"`

	// Upvotes are votes cast by people to indicate them being in
	// agreement, and supporting the feedback given.
	// It must be initialized to the default value (0) when creating
	// an entry, and be at most 2000.
	// NOTE: The upper bound may change in the future.
	Upvotes uint `gorm:"<-;default:0;size:11;scale:0;precision:4" json:"Upvotes" validate:"omitempty,min=0,max=2000"`

	// Downvotes are votes cast by people to indicate them being in
	// disagreement, and opposing the feedback given.
	// It must be initialized to the default value (0) when creating
	// an entry, and be at most 2000.
	// NOTE: The upper bound may change in the future.
	Downvotes uint `gorm:"<-;default:0;size:11;scale:0;precision:4" json:"Downvotes" validate:"omitempty,min=0,max=2000"`
}

// isCourse validates that a filed fills the UCLouvain course
// scheme. The course scheme is as follows.
//
//	course := "L" faculty code
//	faculty := letter letter letter letter?
//	code := digit digit digit digit
//
// Letters are alpha ascii letters, generally uppercase. Lowercase
// formatting should be accepted however.
// Codes are at least 1000 and at most 9999.
func isCourse(fl validator.FieldLevel) bool {
	course := strings.ToLower(fl.Field().String())

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Var(course, "alphanum,min=8,max=9,startswith=l"); err != nil {
		return false
	}

	code, err := strconv.Atoi(course[len(course)-4:])
	if err != nil {
		return false
	}
	if code < 1000 {
		return false
	}

	faculty := course[1 : len(course)-4]
	if err := validate.Var(faculty, "alpha,min=3,max=4"); err != nil {
		return false
	}

	return true
}
