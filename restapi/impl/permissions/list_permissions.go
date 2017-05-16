package permissions

import (
	"database/sql"
	"github.com/cyverse-de/permissions/models"
	permsdb "github.com/cyverse-de/permissions/restapi/impl/db"
	"github.com/cyverse-de/permissions/restapi/operations/permissions"

	"github.com/cyverse-de/logcabin"
	"github.com/go-openapi/runtime/middleware"
)

func BuildListPermissionsHandler(db *sql.DB) func(permissions.ListPermissionsParams) middleware.Responder {

	// Return the handler function.
	return func(permissions.ListPermissionsParams) middleware.Responder {

		// Start a transaction for this request.
		tx, err := db.Begin()
		if err != nil {
			logcabin.Error.Print(err)
			reason := err.Error()
			return permissions.NewListPermissionsInternalServerError().WithPayload(&models.ErrorOut{&reason})
		}
		defer tx.Commit()

		// List all permissions.
		result, err := permsdb.ListPermissions(tx)
		if err != nil {
			logcabin.Error.Print(err)
			reason := err.Error()
			return permissions.NewListPermissionsInternalServerError().WithPayload(&models.ErrorOut{&reason})
		}

		// Return the results.
		return permissions.NewListPermissionsOK().WithPayload(&models.PermissionList{Permissions: result})
	}
}
