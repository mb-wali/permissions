package grouper

import (
	"github.com/cyverse-de/permissions/models"
)

// MockGrouperClient represents a mock Grouper client.
type MockGrouperClient struct {
	groups map[string][]*GroupInfo
}

// NewMockGrouperClient returns a new Mock grouper client.
func NewMockGrouperClient(groups map[string][]*GroupInfo) *MockGrouperClient {
	return &MockGrouperClient{groups: groups}
}

// GroupsForSubject returns a mock list of groups for a subject.
func (gc *MockGrouperClient) GroupsForSubject(subjectID string) ([]*GroupInfo, error) {
	return gc.groups[subjectID], nil
}

// AddSourceIDToPermissions is a no-op for now.
func (gc *MockGrouperClient) AddSourceIDToPermissions(_ []*models.Permission) error {
	return nil
}

// AddSourceIDToPermission is a no-op for now.
func (gc *MockGrouperClient) AddSourceIDToPermission(_ *models.Permission) error {
	return nil
}
