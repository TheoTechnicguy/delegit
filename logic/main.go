/**
 * file: logic/main.go
 * author: theo technicguy
 * license: apache-2.0
 *
 * This main logic file contains generic logic
 * utility parts for the application.
 */

package logic

import (
	"net/http"

	"git.licolas.net/delegit/delegit/uxerrors"
	"gorm.io/gorm"
)

func handleDatabaseError(err error) error {
	switch err {
	case nil:
		return nil
	case gorm.ErrRecordNotFound:
		uxe := uxerrors.New(err)
		uxe.Summary = "Feedback not found"
		uxe.Detail = "The feedback you requested could not be found. Check the ID and try again."
		return uxerrors.NewErrors(http.StatusNotFound).Append(uxe)
	default:
		return uxerrors.NewErrors(http.StatusInternalServerError).AppendNew(err)
	}
}
