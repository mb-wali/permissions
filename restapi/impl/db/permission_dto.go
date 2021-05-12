package db

import (
	"github.com/cyverse-de/permissions/models"
)

type PermissionDto struct {
	ID                *models.PermissionID
	InternalSubjectID *models.InternalSubjectID
	SubjectID         *models.ExternalSubjectID
	SubjectType       *models.SubjectType
	ResourceID        string
	ResourceName      string
	ResourceType      string
	PermissionLevel   *models.PermissionLevel
}

func (p *PermissionDto) ToPermission() *models.Permission {

	// Extract the subject.
	subject := &models.SubjectOut{
		ID:          p.InternalSubjectID,
		SubjectID:   p.SubjectID,
		SubjectType: p.SubjectType,
	}

	// Extract the resource.
	resource := &models.ResourceOut{
		ID:           &p.ResourceID,
		Name:         &p.ResourceName,
		ResourceType: &p.ResourceType,
	}

	// Extract the permission itself.
	permission := &models.Permission{
		ID:              p.ID,
		PermissionLevel: p.PermissionLevel,
		Resource:        resource,
		Subject:         subject,
	}

	return permission
}
