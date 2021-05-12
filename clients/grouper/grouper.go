package grouper

import (
	"database/sql"

	"github.com/cyverse-de/dbutil"
	"github.com/cyverse-de/permissions/models"

	"github.com/lib/pq"
)

// GroupInfo represents information about a group in Grouper.
type GroupInfo struct {
	ID   string
	Name string
}

// Grouper is the interface implemented by a Grouper client instance.
type Grouper interface {
	GroupsForSubject(string) ([]*GroupInfo, error)
	AddSourceIDToPermissions([]*models.Permission) error
	AddSourceIDToPermission(*models.Permission) error
}

// Client represents a Grouper client instance.
//
// Note: the grouper client is intended to be a read-only client. Explicit transactions are not
// used here for that reason.
type Client struct {
	db     *sql.DB
	prefix string
}

// NewGrouperClient creates and returns a new Grouper client for the given database URI and group name prefix.
func NewGrouperClient(dbURI, prefix string) (*Client, error) {
	connector, err := dbutil.NewDefaultConnector("1m")
	if err != nil {
		return nil, err
	}

	db, err := connector.Connect("postgres", dbURI)
	if err != nil {
		return nil, err
	}

	return &Client{
		db:     db,
		prefix: prefix,
	}, nil
}

// GroupsForSubject returns the list of groups that the subject with the given ID belongs to.
func (gc *Client) GroupsForSubject(subjectID string) ([]*GroupInfo, error) {

	// Query the database.
	query := `SELECT group_id, group_name FROM grouper_memberships_v
            WHERE subject_id = $1 AND group_name LIKE $2 AND list_name = 'members'`
	rows, err := gc.db.Query(query, subjectID, gc.prefix+"%")
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

// AddSourceIDToPermissions adds the subject source IDs to a slice of Permission objects.
func (gc *Client) AddSourceIDToPermissions(permissions []*models.Permission) error {

	// Get a list of subject identifiers.
	subjectIDs := make([]string, 0)
	for _, permission := range permissions {
		subjectIDs = append(subjectIDs, string(*permission.Subject.SubjectID))
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
		var sourceID models.SubjectSourceID
		sourceID = models.SubjectSourceID(m[string(*permission.Subject.SubjectID)])
		permission.Subject.SubjectSourceID = &sourceID
	}

	return nil
}

// AddSourceIDToPermission adds the subject source ID to a permission object.
func (gc *Client) AddSourceIDToPermission(permission *models.Permission) error {
	return gc.AddSourceIDToPermissions([]*models.Permission{permission})
}
