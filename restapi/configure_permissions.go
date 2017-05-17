package restapi

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/cyverse-de/configurate"
	"github.com/cyverse-de/dbutil"
	"github.com/cyverse-de/logcabin"
	"github.com/cyverse-de/version"

	errors "github.com/go-openapi/errors"
	httpkit "github.com/go-openapi/runtime"
	swag "github.com/go-openapi/swag"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"github.com/tylerb/graceful"

	"github.com/cyverse-de/permissions/clients/grouper"
	"github.com/cyverse-de/permissions/restapi/operations"
	"github.com/cyverse-de/permissions/restapi/operations/permissions"
	"github.com/cyverse-de/permissions/restapi/operations/resource_types"
	"github.com/cyverse-de/permissions/restapi/operations/resources"
	"github.com/cyverse-de/permissions/restapi/operations/status"
	"github.com/cyverse-de/permissions/restapi/operations/subjects"

	permissions_impl "github.com/cyverse-de/permissions/restapi/impl/permissions"
	resource_types_impl "github.com/cyverse-de/permissions/restapi/impl/resource_types"
	resources_impl "github.com/cyverse-de/permissions/restapi/impl/resources"
	status_impl "github.com/cyverse-de/permissions/restapi/impl/status"
	subjects_impl "github.com/cyverse-de/permissions/restapi/impl/subjects"
)

// This file is safe to edit. Once it exists it will not be overwritten

// The default permissions configuration.
const DefaultConfig = `
db:
  uri: "postgresql://de:notprod@dedb:5432/permissions?sslmode=disable"

grouperdb:
  uri: "postgresql://de:notprod@adedb:5432/grouper?sslmode=disable"
  folder_name_prefix: "iplant:de:docker-compose"
`

// Command line options that aren't managed by go-swagger.
var options struct {
	CfgPath     string `long:"config" default:"/etc/iplant/de/permissions.yaml" description:"The path to the config file"`
	ShowVersion bool   `short:"v" long:"version" description:"Print the app version and exit"`
}

// Register the command-line options.
func configureFlags(api *operations.PermissionsAPI) {
	api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{
		swag.CommandLineOptionsGroup{"Service Options", "", &options},
	}
}

// Validate the custom command-line options.
func validateOptions() error {
	if options.CfgPath == "" {
		return fmt.Errorf("--config must be set")
	}

	return nil
}

// The database connection.
var db *sql.DB
var grouperClient *grouper.GrouperClient

// Initialize the service.
func initService() error {
	if options.ShowVersion {
		version.AppVersion()
		os.Exit(0)
	}

	var (
		err error
		cfg *viper.Viper
	)
	if cfg, err = configurate.InitDefaults(options.CfgPath, DefaultConfig); err != nil {
		return err
	}

	connector, err := dbutil.NewDefaultConnector("1m")
	if err != nil {
		return err
	}

	dburi := cfg.GetString("db.uri")
	db, err = connector.Connect("postgres", dburi)
	if err != nil {
		return err
	}

	grouperDburi := cfg.GetString("grouperdb.uri")
	grouperFolderNamePrefix := cfg.GetString("grouperdb.folder_name_prefix")
	grouperClient, err = grouper.NewGrouperClient(grouperDburi, grouperFolderNamePrefix)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	return nil
}

// Clean up when the service exits.
func cleanup() {
	logcabin.Info.Printf("Closing the database connection.")
	db.Close()
}

func configureAPI(api *operations.PermissionsAPI) http.Handler {
	if err := validateOptions(); err != nil {
		logcabin.Error.Fatal(err)
	}

	if err := initService(); err != nil {
		logcabin.Error.Fatal(err)
	}

	api.ServeError = errors.ServeError

	api.JSONConsumer = httpkit.JSONConsumer()

	api.JSONProducer = httpkit.JSONProducer()

	api.StatusGetHandler = status.GetHandlerFunc(status_impl.BuildStatusHandler(SwaggerJSON))

	api.ResourceTypesGetResourceTypesHandler = resource_types.GetResourceTypesHandlerFunc(
		resource_types_impl.BuildResourceTypesGetHandler(db),
	)

	api.ResourceTypesDeleteResourceTypeByNameHandler = resource_types.DeleteResourceTypeByNameHandlerFunc(
		resource_types_impl.BuildDeleteResourceTypeByNameHandler(db),
	)

	api.ResourceTypesPostResourceTypesHandler = resource_types.PostResourceTypesHandlerFunc(
		resource_types_impl.BuildResourceTypesPostHandler(db),
	)

	api.ResourceTypesPutResourceTypesIDHandler = resource_types.PutResourceTypesIDHandlerFunc(
		resource_types_impl.BuildResourceTypesIDPutHandler(db),
	)

	api.ResourceTypesDeleteResourceTypesIDHandler = resource_types.DeleteResourceTypesIDHandlerFunc(
		resource_types_impl.BuildResourceTypesIDDeleteHandler(db),
	)

	api.ResourcesAddResourceHandler = resources.AddResourceHandlerFunc(
		resources_impl.BuildAddResourceHandler(db),
	)

	api.ResourcesDeleteResourceByNameHandler = resources.DeleteResourceByNameHandlerFunc(
		resources_impl.BuildDeleteResourceByNameHandler(db),
	)

	api.ResourcesListResourcesHandler = resources.ListResourcesHandlerFunc(
		resources_impl.BuildListResourcesHandler(db),
	)

	api.ResourcesUpdateResourceHandler = resources.UpdateResourceHandlerFunc(
		resources_impl.BuildUpdateResourceHandler(db),
	)

	api.ResourcesDeleteResourceHandler = resources.DeleteResourceHandlerFunc(
		resources_impl.BuildDeleteResourceHandler(db),
	)

	api.SubjectsAddSubjectHandler = subjects.AddSubjectHandlerFunc(
		subjects_impl.BuildAddSubjectHandler(db),
	)

	api.SubjectsDeleteSubjectByExternalIDHandler = subjects.DeleteSubjectByExternalIDHandlerFunc(
		subjects_impl.BuildDeleteSubjectByExternalIdHandler(db),
	)

	api.SubjectsListSubjectsHandler = subjects.ListSubjectsHandlerFunc(
		subjects_impl.BuildListSubjectsHandler(db),
	)

	api.SubjectsUpdateSubjectHandler = subjects.UpdateSubjectHandlerFunc(
		subjects_impl.BuildUpdateSubjectHandler(db),
	)

	api.SubjectsDeleteSubjectHandler = subjects.DeleteSubjectHandlerFunc(
		subjects_impl.BuildDeleteSubjectHandler(db),
	)

	api.PermissionsListPermissionsHandler = permissions.ListPermissionsHandlerFunc(
		permissions_impl.BuildListPermissionsHandler(db, grouperClient),
	)

	api.PermissionsGrantPermissionHandler = permissions.GrantPermissionHandlerFunc(
		permissions_impl.BuildGrantPermissionHandler(db, grouperClient),
	)

	api.PermissionsRevokePermissionHandler = permissions.RevokePermissionHandlerFunc(
		permissions_impl.BuildRevokePermissionHandler(db),
	)

	api.PermissionsPutPermissionHandler = permissions.PutPermissionHandlerFunc(
		permissions_impl.BuildPutPermissionHandler(db, grouperClient),
	)

	api.PermissionsBySubjectHandler = permissions.BySubjectHandlerFunc(
		permissions_impl.BuildBySubjectHandler(db, grouperClient),
	)

	api.PermissionsBySubjectAndResourceTypeHandler = permissions.BySubjectAndResourceTypeHandlerFunc(
		permissions_impl.BuildBySubjectAndResourceTypeHandler(db, grouperClient),
	)

	api.PermissionsBySubjectAndResourceHandler = permissions.BySubjectAndResourceHandlerFunc(
		permissions_impl.BuildBySubjectAndResourceHandler(db, grouperClient),
	)

	api.PermissionsListResourcePermissionsHandler = permissions.ListResourcePermissionsHandlerFunc(
		permissions_impl.BuildListResourcePermissionsHandler(db, grouperClient),
	)

	api.ServerShutdown = cleanup

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *graceful.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json
// document. So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return uiMiddleware(handler)
}

// The middleware to serve up the interactive Swagger UI.
func uiMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/docs" {
			http.Redirect(w, r, "/docs/", http.StatusFound)
			return
		}
		if strings.Index(r.URL.Path, "/docs/") == 0 {
			http.StripPrefix("/docs/", http.FileServer(http.Dir("docs"))).ServeHTTP(w, r)
			return
		}
		handler.ServeHTTP(w, r)
	})
}
