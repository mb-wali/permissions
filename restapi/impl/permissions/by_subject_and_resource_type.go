package permissions

import (
	"database/sql"
	"fmt"

	"github.com/cyverse-de/permissions/clients/grouper"
	"github.com/cyverse-de/permissions/logger"
	"github.com/cyverse-de/permissions/models"
	permsdb "github.com/cyverse-de/permissions/restapi/impl/db"
	"github.com/cyverse-de/permissions/restapi/operations/permissions"

	"github.com/go-openapi/runtime/middleware"
)

func bySubjectAndResourceTypeOk(perms []*models.Permission) middleware.Responder {
	return permissions.NewBySubjectAndResourceTypeOK().WithPayload(
		&models.PermissionList{Permissions: perms},
	)
}

func bySubjectAndResourceTypeInternalServerError(reason string) middleware.Responder {
	return permissions.NewBySubjectAndResourceTypeInternalServerError().WithPayload(
		&models.ErrorOut{Reason: &reason},
	)
}

func bySubjectAndResourceTypeBadRequest(reason string) middleware.Responder {
	return permissions.NewBySubjectAndResourceTypeBadRequest().WithPayload(
		&models.ErrorOut{Reason: &reason},
	)
}

func BuildBySubjectAndResourceTypeHandler(
	db *sql.DB, grouperClient grouper.Grouper,
) func(permissions.BySubjectAndResourceTypeParams) middleware.Responder {

	// Return the handler function.
	return func(params permissions.BySubjectAndResourceTypeParams) middleware.Responder {
		subjectType := params.SubjectType
		subjectId := params.SubjectID
		resourceTypeName := params.ResourceType
		lookup := extractLookupFlag(params.Lookup)
		minLevel := params.MinLevel

		// Create a transaction for the request.
		tx, err := db.Begin()
		if err != nil {
			logger.Log.Error(err)
			return bySubjectAndResourceTypeInternalServerError(err.Error())
		}

		// Verify that the subject type is correct.
		subject, err := permsdb.GetSubjectByExternalId(tx, models.ExternalSubjectID(subjectId))
		if err != nil {
			tx.Rollback()
			logger.Log.Error(err)
			return bySubjectAndResourceTypeInternalServerError(err.Error())
		}
		if subject != nil && string(*subject.SubjectType) != subjectType {
			tx.Rollback()
			reason := fmt.Sprintf("incorrect type for subject, %s: %s", subjectId, subjectType)
			return bySubjectAndResourceTypeBadRequest(reason)
		}

		// Verify that the resource type exists.
		resourceType, err := permsdb.GetResourceTypeByName(tx, &resourceTypeName)
		if err != nil {
			tx.Rollback()
			logger.Log.Error(err)
			return bySubjectAndResourceTypeInternalServerError(err.Error())
		}
		if resourceType == nil {
			tx.Rollback()
			return bySubjectAndResourceTypeOk(make([]*models.Permission, 0))
		}

		// Get the list of subject IDs to use for the query.
		subjectIds, err := buildSubjectIdList(grouperClient, subjectType, subjectId, lookup)
		if err != nil {
			tx.Rollback()
			logger.Log.Error(err)
			return bySubjectInternalServerError(err.Error())
		}

		// Perform the lookup.
		var perms []*models.Permission
		if minLevel == nil {
			perms, err = permsdb.PermissionsForSubjectsAndResourceType(tx, subjectIds, resourceTypeName)
			if err != nil {
				tx.Rollback()
				logger.Log.Error(err)
				return bySubjectAndResourceTypeInternalServerError(err.Error())
			}
		} else {
			perms, err = permsdb.PermissionsForSubjectsAndResourceTypeMinLevel(tx, subjectIds, resourceTypeName, *minLevel)
			if err != nil {
				tx.Rollback()
				logger.Log.Error(err)
				return bySubjectAndResourceTypeInternalServerError(err.Error())
			}
		}

		// Commit the transaction.
		if err := tx.Commit(); err != nil {
			tx.Rollback()
			logger.Log.Error(err)
			return bySubjectAndResourceTypeInternalServerError(err.Error())
		}

		// Add the subject source ID to the response body.
		if err := grouperClient.AddSourceIDToPermissions(perms); err != nil {
			logger.Log.Error(err)
			return bySubjectAndResourceTypeInternalServerError(err.Error())
		}

		return bySubjectAndResourceTypeOk(perms)
	}
}
