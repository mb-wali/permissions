package db

import (
	"github.com/cyverse-de/permissions/models"
)

type SubjectDto struct {
	ID          *models.InternalSubjectID
	SubjectID   *models.ExternalSubjectID
	SubjectType *models.SubjectType
}

func (s *SubjectDto) ToSubjectOut() *models.SubjectOut {
	var subjectOut models.SubjectOut

	subjectOut.ID = s.ID
	subjectOut.SubjectID = s.SubjectID
	subjectOut.SubjectType = s.SubjectType

	return &subjectOut
}
