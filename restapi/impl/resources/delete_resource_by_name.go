package resources

import (
	"database/sql"
	"fmt"
	"github.com/cyverse-de/permissions/models"
	permsdb "github.com/cyverse-de/permissions/restapi/impl/db"
	"github.com/cyverse-de/permissions/restapi/operations/resources"

	"github.com/cyverse-de/logcabin"
	"github.com/go-openapi/runtime/middleware"
)

func deleteResourceByNameInternalServerError(reason string) middleware.Responder {
	return resources.NewDeleteResourceByNameInternalServerError().WithPayload(
		&models.ErrorOut{Reason: &reason},
	)
}

func deleteResourceByNameOK() middleware.Responder {
	return resources.NewDeleteResourceByNameOK()
}

func deleteResourceByNameNotFound(reason string) middleware.Responder {
	return resources.NewDeleteResourceByNameNotFound().WithPayload(
		&models.ErrorOut{Reason: &reason},
	)
}

func BuildDeleteResourceByNameHandler(db *sql.DB) func(resources.DeleteResourceByNameParams) middleware.Responder {

	// Return the handler function.
	return func(params resources.DeleteResourceByNameParams) middleware.Responder {

		// Start a transaction for the request.
		tx, err := db.Begin()
		if err != nil {
			logcabin.Error.Print(err)
			return deleteResourceByNameInternalServerError(err.Error())
		}

		// Look up the resource.
		resource, err := permsdb.GetResourceByNameAndType(tx, params.ResourceName, params.ResourceTypeName)
		if err != nil {
			tx.Rollback()
			logcabin.Error.Print(err)
			return deleteResourceByNameInternalServerError(err.Error())
		}
		if resource == nil {
			tx.Rollback()
			reason := fmt.Sprintf("resource not found: %s:%s", params.ResourceTypeName, params.ResourceName)
			return deleteResourceByNameNotFound(reason)
		}

		// Delete the resource.
		if err := permsdb.DeleteResource(tx, resource.ID); err != nil {
			tx.Rollback()
			logcabin.Error.Print(err)
			return deleteResourceByNameInternalServerError(err.Error())
		}

		// Commit the transaction.
		if err := tx.Commit(); err != nil {
			tx.Rollback()
			logcabin.Error.Print(err)
			return deleteResourceByNameInternalServerError(err.Error())
		}

		return deleteResourceByNameOK()
	}
}
