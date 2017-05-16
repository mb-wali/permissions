package resources

import (
	"database/sql"
	"github.com/cyverse-de/permissions/models"
	permsdb "github.com/cyverse-de/permissions/restapi/impl/db"
	"github.com/cyverse-de/permissions/restapi/operations/resources"

	"github.com/cyverse-de/logcabin"
	"github.com/go-openapi/runtime/middleware"
)

func BuildListResourcesHandler(db *sql.DB) func(resources.ListResourcesParams) middleware.Responder {

	// Return the handler function.
	return func(params resources.ListResourcesParams) middleware.Responder {

		// Start a transaction for this request.
		tx, err := db.Begin()
		if err != nil {
			logcabin.Error.Print(err)
			reason := err.Error()
			return resources.NewListResourcesInternalServerError().WithPayload(&models.ErrorOut{&reason})
		}
		defer tx.Commit()

		// List all resources.
		result, err := permsdb.ListResources(tx, params.ResourceTypeName, params.ResourceName)
		if err != nil {
			logcabin.Error.Print(err)
			reason := err.Error()
			return resources.NewListResourcesInternalServerError().WithPayload(&models.ErrorOut{&reason})
		}

		// Return the results.
		return resources.NewListResourcesOK().WithPayload(&models.ResourcesOut{Resources: result})
	}
}
