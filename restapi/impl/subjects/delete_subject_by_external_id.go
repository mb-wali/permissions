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

func deleteSubjectByExternalIdInternalServerError(reason string) middleware.Responder {
	return subjects.NewDeleteSubjectByExternalIDInternalServerError().WithPayload(
		&models.ErrorOut{Reason: &reason},
	)
}

func deleteSubjectByExternalIdNotFound(reason string) middleware.Responder {
	return subjects.NewDeleteSubjectByExternalIDNotFound().WithPayload(
		&models.ErrorOut{Reason: &reason},
	)
}

func deleteSubjectByExternalIdOk() middleware.Responder {
	return subjects.NewDeleteSubjectByExternalIDOK()
}

func BuildDeleteSubjectByExternalIdHandler(
	db *sql.DB,
) func(subjects.DeleteSubjectByExternalIDParams) middleware.Responder {

	// Return the handler function.
	return func(params subjects.DeleteSubjectByExternalIDParams) middleware.Responder {
		subjectId := models.ExternalSubjectID(params.SubjectID)
		subjectType := models.SubjectType(params.SubjectType)

		// Start a transaction for the request.
		tx, err := db.Begin()
		if err != nil {
			logger.Log.Error(err)
			return deleteSubjectByExternalIdInternalServerError(err.Error())
		}

		// Look up the subject.
		subject, err := permsdb.GetSubject(tx, subjectId, subjectType)
		if err != nil {
			tx.Rollback()
			logger.Log.Error(err)
			return deleteSubjectByExternalIdInternalServerError(err.Error())
		}
		if subject == nil {
			tx.Rollback()
			reason := fmt.Sprintf("subject not found: %s:%s", string(subjectType), string(subjectId))
			return deleteSubjectByExternalIdNotFound(reason)
		}

		// Delete the subject.
		if err := permsdb.DeleteSubject(tx, *subject.ID); err != nil {
			tx.Rollback()
			return deleteSubjectByExternalIdInternalServerError(err.Error())
		}

		// Commit the transaction.
		if err := tx.Commit(); err != nil {
			tx.Rollback()
			logger.Log.Error(err)
			return deleteSubjectByExternalIdInternalServerError(err.Error())
		}

		return deleteSubjectByExternalIdOk()
	}
}
