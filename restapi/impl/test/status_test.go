package test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/cyverse-de/permissions/restapi/operations/status"

	impl "github.com/cyverse-de/permissions/restapi/impl/status"
)

var testJSON = json.RawMessage([]byte(`{
    "info": {
        "description": "The Description",
        "title":       "The Title",
        "version":     "The Version"
    }
}`))

func TestGetStatus(t *testing.T) {
	f := impl.BuildStatusHandler(testJSON)
	r := f(status.NewGetParams())

	// Verify that we got the expected return value.
	if reflect.TypeOf(r) != reflect.TypeOf((*status.GetOK)(nil)) {
		t.Errorf("unexpected return type from status handler: %s", reflect.TypeOf(r))
	}

	// Convert the return value for convenience.
	s := r.(*status.GetOK)

	// Verify that we got the expected payload.
	if *s.Payload.Description != "The Description" {
		t.Errorf("unexpected description: %s", *s.Payload.Description)
	}
	if *s.Payload.Service != "The Title" {
		t.Errorf("unexpected service name: %s", *s.Payload.Service)
	}
	if *s.Payload.Version != "The Version" {
		t.Errorf("unexpected version: %s", *s.Payload.Version)
	}
}
