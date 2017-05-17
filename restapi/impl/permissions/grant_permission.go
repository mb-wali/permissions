package permissions

import (
	"database/sql"
	"github.com/cyverse-de/permissions/clients/grouper"
	"github.com/cyverse-de/permissions/logger"
	"github.com/cyverse-de/permissions/models"
	permsdb "github.com/cyverse-de/permissions/restapi/impl/db"
	"github.com/cyverse-de/permissions/restapi/operations/permissions"

	"github.com/go-openapi/runtime/middleware"
)

func grantPermissionInternalServerError(reason string) middleware.Responder {
	return permissions.NewGrantPermissionInternalServerError().WithPayload(
		&models.ErrorOut{Reason: &reason},
	)
}

func grantPermissionBadRequest(reason string) middleware.Responder {
	return permissions.NewGrantPermissionBadRequest().WithPayload(
		&models.ErrorOut{Reason: &reason},
	)
}

func BuildGrantPermissionHandler(
	db *sql.DB, grouperClient grouper.Grouper,
) func(permissions.GrantPermissionParams) middleware.Responder {

	erf := &ErrorResponseFns{
		InternalServerError: grantPermissionInternalServerError,
		BadRequest:          grantPermissionBadRequest,
	}

	// Return the hnadler function.
	return func(params permissions.GrantPermissionParams) middleware.Responder {
		req := params.PermissionGrantRequest

		// Create a transaction for the request.
		tx, err := db.Begin()
		if err != nil {
			logger.Log.Error(err)
			return grantPermissionInternalServerError(err.Error())
		}

		// Either get or add the subject.
		subject, errorResponder := getOrAddSubject(tx, req.Subject, erf)
		if errorResponder != nil {
			tx.Rollback()
			return errorResponder
		}

		// Either get or add the resource.
		resource, errorResponder := getOrAddResource(tx, req.Resource, erf)
		if errorResponder != nil {
			tx.Rollback()
			return errorResponder
		}

		// Look up the permission level.
		permissionLevelId, errorResponder := getPermissionLevel(tx, req.PermissionLevel, erf)
		if errorResponder != nil {
			tx.Rollback()
			return errorResponder
		}

		// Either update or add the permission.
		permission, err := permsdb.UpsertPermission(tx, subject.ID, *resource.ID, *permissionLevelId)
		if err != nil {
			tx.Rollback()
			logger.Log.Error(err)
			return grantPermissionInternalServerError(err.Error())
		}

		// Commit the transaction.
		if err := tx.Commit(); err != nil {
			tx.Rollback()
			logger.Log.Error(err)
			return grantPermissionInternalServerError(err.Error())
		}

		// Add the subject source ID to the permission object.
		if err := grouperClient.AddSourceIDToPermission(permission); err != nil {
			logger.Log.Error(err)
			return grantPermissionInternalServerError(err.Error())
		}

		return permissions.NewGrantPermissionOK().WithPayload(permission)
	}
}
