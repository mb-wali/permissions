package grouper

import (
	"database/sql"

	"github.com/cyverse-de/dbutil"
	"github.com/cyverse-de/permissions/models"

	"github.com/lib/pq"
)

type GroupInfo struct {
	ID   string
	Name string
}

type Grouper interface {
	GroupsForSubject(string) ([]*GroupInfo, error)
	AddSourceIDToPermissions([]*models.Permission) error
	AddSourceIDToPermission(*models.Permission) error
}

// Note: the grouper client is intended to be a read-only client. Explicit transactions are not
// used here for that reason.
type GrouperClient struct {
	db     *sql.DB
	prefix string
}

func NewGrouperClient(dburi, prefix string) (*GrouperClient, error) {
	connector, err := dbutil.NewDefaultConnector("1m")
	if err != nil {
		return nil, err
	}

	db, err := connector.Connect("postgres", dburi)
	if err != nil {
		return nil, err
	}

	return &GrouperClient{
		db:     db,
		prefix: prefix,
	}, nil
}

func (gc *GrouperClient) GroupsForSubject(subjectId string) ([]*GroupInfo, error) {

	// Query the database.
	query := `SELECT group_id, group_name FROM grouper_memberships_v
            WHERE subject_id = $1 AND group_name LIKE $2`
	rows, err := gc.db.Query(query, subjectId, gc.prefix+"%")
	if err != nil {
		return nil, err
	}

	// Extract the groups from the database.
	groups := make([]*GroupInfo, 0)
	for rows.Next() {
		var group GroupInfo
		if err := rows.Scan(&group.ID, &group.Name); err != nil {
			return nil, err
		}
		groups = append(groups, &group)
	}

	return groups, nil
}

func (gc *GrouperClient) AddSourceIDToPermissions(permissions []*models.Permission) error {

	// Get a list of subject identifiers.
	subjectIDs := make([]string, 0)
	for _, permission := range permissions {
		subjectIDs = append(subjectIDs, string(permission.Subject.SubjectID))
	}

	// Query the database.
	query := `SELECT subject_id, subject_source FROM grouper_members
            WHERE subject_id = ANY($1)`
	rows, err := gc.db.Query(query, pq.Array(subjectIDs))
	if err != nil {
		return err
	}
	defer rows.Close()

	// Build a map from subject ID to source ID.
	m := make(map[string]string)
	for rows.Next() {
		var subjectID, sourceID string
		if err := rows.Scan(&subjectID, &sourceID); err != nil {
			return err
		}
		m[subjectID] = sourceID
	}

	// Add the subject IDs to the permission objects.
	for _, permission := range permissions {
		sourceID := m[string(permission.Subject.SubjectID)]
		permission.Subject.SubjectSourceID = models.SubjectSourceID(sourceID)
	}

	return nil
}

func (gc *GrouperClient) AddSourceIDToPermission(permission *models.Permission) error {
	return gc.AddSourceIDToPermissions([]*models.Permission{permission})
}
