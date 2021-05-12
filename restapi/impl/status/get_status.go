package status

import (
	"encoding/json"
	"fmt"

	"github.com/cyverse-de/permissions/models"
	"github.com/cyverse-de/permissions/restapi/operations/status"
	"github.com/go-openapi/runtime/middleware"
)

// Info represents the service information retrieved from the Swagger specification.
type Info struct {
	Description string `json:"description"`
	Title       string `json:"title"`
	Version     string `json:"version"`
}

// SwaggerSpec represents just the service imformation portion of the Swagger specification.
type SwaggerSpec struct {
	Info Info `json:"info"`
}

func serviceInfo(swaggerJSON json.RawMessage) (*models.ServiceInfo, error) {
	var decoded SwaggerSpec

	// Extract the service info from the Swagger JSON.
	if err := json.Unmarshal(swaggerJSON, &decoded); err != nil {
		return nil, fmt.Errorf("unable to decode the Swagger JSON: %s", err)
	}

	// Build the service info struct.
	info := decoded.Info
	return &models.ServiceInfo{
		Description: &info.Description,
		Service:     &info.Title,
		Version:     &info.Version,
	}, nil
}

// BuildStatusHandler builds the request handler for the service status endpoint.
func BuildStatusHandler(swaggerJSON json.RawMessage) func(status.GetParams) middleware.Responder {

	// Load the service info. Failure to do so will cause the service to abort.
	info, err := serviceInfo(swaggerJSON)
	if err != nil {
		panic(err)
	}

	// Return the handler function.
	return func(status.GetParams) middleware.Responder {
		return status.NewGetOK().WithPayload(info)
	}
}
