package resourcetypes

import (
	"database/sql"
	"fmt"

	"github.com/cyverse-de/permissions/models"
	permsdb "github.com/cyverse-de/permissions/restapi/impl/db"
	"github.com/cyverse-de/permissions/restapi/operations/resource_types"
	"github.com/go-openapi/runtime/middleware"
)

// BuildResourceTypesIDDeleteHandler builds the request handler for the resource type deletion endpoint.
func BuildResourceTypesIDDeleteHandler(
	db *sql.DB,
) func(resource_types.DeleteResourceTypesIDParams) middleware.Responder {

	// Return the handler function.
	return func(params resource_types.DeleteResourceTypesIDParams) middleware.Responder {

		// Start a transaction for this request.
		tx, err := db.Begin()
		if err != nil {
			reason := err.Error()
			return resource_types.NewDeleteResourceTypesIDInternalServerError().WithPayload(
				&models.ErrorOut{Reason: &reason},
			)
		}

		// Verify that the resource type exists.
		exists, err := permsdb.ResourceTypeExists(tx, &params.ID)
		if err != nil {
			tx.Rollback()
			reason := err.Error()
			return resource_types.NewDeleteResourceTypesIDInternalServerError().WithPayload(
				&models.ErrorOut{Reason: &reason},
			)
		}
		if !exists {
			tx.Rollback()
			reason := fmt.Sprintf("resource type %s not found", params.ID)
			return resource_types.NewDeleteResourceTypesIDNotFound().WithPayload(
				&models.ErrorOut{Reason: &reason},
			)
		}

		// Verify that the resource type has no resources associated with it.
		numResources, err := permsdb.CountResourcesOfType(tx, &params.ID)
		if err != nil {
			tx.Rollback()
			reason := err.Error()
			return resource_types.NewDeleteResourceTypesIDInternalServerError().WithPayload(
				&models.ErrorOut{Reason: &reason},
			)
		}
		if numResources != 0 {
			tx.Rollback()
			reason := fmt.Sprintf("resource type %s has resources associated with it", params.ID)
			return resource_types.NewDeleteResourceTypesIDBadRequest().WithPayload(
				&models.ErrorOut{Reason: &reason},
			)
		}

		// Delete the resource type.
		err = permsdb.DeleteResourceType(tx, &params.ID)
		if err != nil {
			tx.Rollback()
			reason := err.Error()
			return resource_types.NewDeleteResourceTypesIDInternalServerError().WithPayload(
				&models.ErrorOut{Reason: &reason},
			)
		}

		// Commit the transaction.
		if err := tx.Commit(); err != nil {
			tx.Rollback()
			reason := err.Error()
			return resource_types.NewDeleteResourceTypesIDInternalServerError().WithPayload(
				&models.ErrorOut{Reason: &reason},
			)
		}

		return resource_types.NewDeleteResourceTypesIDOK()
	}
}
