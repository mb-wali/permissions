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

func BuildAddSubjectHandler(db *sql.DB) func(subjects.AddSubjectParams) middleware.Responder {

	// Return the handler function.
	return func(params subjects.AddSubjectParams) middleware.Responder {
		subjectIn := params.SubjectIn

		// Start a transaction for this request.
		tx, err := db.Begin()
		if err != nil {
			logger.Log.Error(err)
			reason := err.Error()
			return subjects.NewAddSubjectInternalServerError().WithPayload(
				&models.ErrorOut{Reason: &reason},
			)
		}

		// Make sure that a subject with the same ID doesn't exist already.
		exists, err := permsdb.SubjectIdExists(tx, *subjectIn.SubjectID)
		if err != nil {
			tx.Rollback()
			logger.Log.Error(err)
			reason := err.Error()
			return subjects.NewAddSubjectInternalServerError().WithPayload(
				&models.ErrorOut{Reason: &reason},
			)
		}
		if exists {
			tx.Rollback()
			reason := fmt.Sprintf("subject, %s, already exists", string(*subjectIn.SubjectID))
			return subjects.NewAddSubjectBadRequest().WithPayload(
				&models.ErrorOut{Reason: &reason},
			)
		}

		// Add the subject.
		subjectOut, err := permsdb.AddSubject(tx, *subjectIn.SubjectID, *subjectIn.SubjectType)
		if err != nil {
			tx.Rollback()
			logger.Log.Error(err)
			reason := err.Error()
			return subjects.NewAddSubjectInternalServerError().WithPayload(
				&models.ErrorOut{Reason: &reason},
			)
		}

		// Commit the transaction.
		if err := tx.Commit(); err != nil {
			tx.Rollback()
			logger.Log.Error(err)
			reason := err.Error()
			return subjects.NewAddSubjectInternalServerError().WithPayload(
				&models.ErrorOut{Reason: &reason},
			)
		}

		return subjects.NewAddSubjectCreated().WithPayload(subjectOut)
	}
}
