package subjects

import (
	"database/sql"
	"fmt"

	"github.com/cyverse-de/permissions/logger"
	"github.com/cyverse-de/permissions/models"
	permsdb "github.com/cyverse-de/permissions/restapi/impl/db"
	"github.com/cyverse-de/permissions/restapi/operations/subjects"

	"github.com/go-openapi/runtime/middleware"
)

// BuildDeleteSubjectHandler builds the request handler for the delete subject endpoint.
func BuildDeleteSubjectHandler(db *sql.DB) func(subjects.DeleteSubjectParams) middleware.Responder {

	// Return the handler function.
	return func(params subjects.DeleteSubjectParams) middleware.Responder {
		id := models.InternalSubjectID(params.ID)

		// Start a transaction for this request.
		tx, err := db.Begin()
		if err != nil {
			logger.Log.Error(err)
			reason := err.Error()
			return subjects.NewDeleteSubjectInternalServerError().WithPayload(
				&models.ErrorOut{Reason: &reason},
			)
		}

		// Verify that the subject exists.
		exists, err := permsdb.SubjectExists(tx, id)
		if err != nil {
			tx.Rollback()
			logger.Log.Error(err)
			reason := err.Error()
			return subjects.NewDeleteSubjectInternalServerError().WithPayload(
				&models.ErrorOut{Reason: &reason},
			)
		}
		if !exists {
			tx.Rollback()
			reason := fmt.Sprintf("subject, %s, not found", string(id))
			return subjects.NewDeleteSubjectNotFound().WithPayload(
				&models.ErrorOut{Reason: &reason},
			)
		}

		// Delete the subject.
		if err := permsdb.DeleteSubject(tx, id); err != nil {
			tx.Rollback()
			logger.Log.Error(err)
			reason := err.Error()
			return subjects.NewDeleteSubjectInternalServerError().WithPayload(
				&models.ErrorOut{Reason: &reason},
			)
		}

		// Commit the transaction.
		if err := tx.Commit(); err != nil {
			tx.Rollback()
			logger.Log.Error(err)
			reason := err.Error()
			return subjects.NewDeleteSubjectInternalServerError().WithPayload(
				&models.ErrorOut{Reason: &reason},
			)
		}

		return subjects.NewDeleteSubjectOK()
	}
}
