package permissions

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-swagger/go-swagger/httpkit/middleware"
)

// BySubjectAndResourceTypeHandlerFunc turns a function with the right signature into a by subject and resource type handler
type BySubjectAndResourceTypeHandlerFunc func(BySubjectAndResourceTypeParams) middleware.Responder

// Handle executing the request and returning a response
func (fn BySubjectAndResourceTypeHandlerFunc) Handle(params BySubjectAndResourceTypeParams) middleware.Responder {
	return fn(params)
}

// BySubjectAndResourceTypeHandler interface for that can handle valid by subject and resource type params
type BySubjectAndResourceTypeHandler interface {
	Handle(BySubjectAndResourceTypeParams) middleware.Responder
}

// NewBySubjectAndResourceType creates a new http.Handler for the by subject and resource type operation
func NewBySubjectAndResourceType(ctx *middleware.Context, handler BySubjectAndResourceTypeHandler) *BySubjectAndResourceType {
	return &BySubjectAndResourceType{Context: ctx, Handler: handler}
}

/*BySubjectAndResourceType swagger:route GET /permissions/subjects/{subject_type}/{subject_id}/{resource_type} permissions bySubjectAndResourceType

Look Up by Subject and Resource Type

Looks up all permissions granted to a subject for resources of the given type. If lookup mode is enabled and the subject is a user, the most lenient permissions granted to the subject or any groups the subject belongs to will be listed. If lookup mode is not enabled or the subject is a group then only permissions assigned directly to the subject will be listed. This endpoint will return an error status if the subject ID is in use and associated with a different subject type.

*/
type BySubjectAndResourceType struct {
	Context *middleware.Context
	Handler BySubjectAndResourceTypeHandler
}

func (o *BySubjectAndResourceType) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, _ := o.Context.RouteInfo(r)
	var Params = NewBySubjectAndResourceTypeParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
