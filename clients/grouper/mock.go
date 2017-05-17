package grouper

import (
	"github.com/cyverse-de/permissions/models"
)

type MockGrouperClient struct {
	groups map[string][]*GroupInfo
}

func NewMockGrouperClient(groups map[string][]*GroupInfo) *MockGrouperClient {
	return &MockGrouperClient{groups: groups}
}

func (gc *MockGrouperClient) GroupsForSubject(subjectId string) ([]*GroupInfo, error) {
	return gc.groups[subjectId], nil
}

func (gc *MockGrouperClient) AddSourceIDToPermissions(_ []*models.Permission) error {
	return nil
}
