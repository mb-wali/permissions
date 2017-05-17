package permissions

import (
	"database/sql"
	"github.com/cyverse-de/permissions/models"
	permsdb "github.com/cyverse-de/permissions/restapi/impl/db"
	"github.com/cyverse-de/permissions/restapi/operations/permissions"

	"github.com/cyverse-de/logcabin"
	"github.com/cyverse-de/permissions/clients/grouper"
	"github.com/go-openapi/runtime/middleware"
)

func internalServerError(reason string) *permissions.ListPermissionsInternalServerError {
	return permissions.NewListPermissionsInternalServerError().WithPayload(
		&models.ErrorOut{Reason: &reason},
	)
}

func BuildListPermissionsHandler(
	db *sql.DB, grouper grouper.Grouper,
) func(permissions.ListPermissionsParams) middleware.Responder {

	// Return the handler function.
	return func(permissions.ListPermissionsParams) middleware.Responder {

		// Start a transaction for this request.
		tx, err := db.Begin()
		if err != nil {
			logcabin.Error.Print(err)
			return internalServerError(err.Error())
		}
		defer tx.Commit()

		// List all permissions.
		result, err := permsdb.ListPermissions(tx)
		if err != nil {
			logcabin.Error.Print(err)
			return internalServerError(err.Error())
		}

		// Add subject sources to the permission list.
		if err = grouper.AddSourceIDToPermissions(result); err != nil {
			logcabin.Error.Print(err)
			return internalServerError(err.Error())
		}

		// Return the results.
		return permissions.NewListPermissionsOK().WithPayload(&models.PermissionList{Permissions: result})
	}
}
