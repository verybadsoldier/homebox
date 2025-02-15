package adapters

import (
	"net/http"

	"github.com/hay-kot/homebox/backend/pkgs/server"
)

// Action is a function that adapts a function to the server.Handler interface.
// It decodes the request body into a value of type T and passes it to the function f.
// The function f is expected to return a value of type Y and an error.
//
// Example:
//
//	type Body struct {
//	    Foo string `json:"foo"`
//	}
//
//	fn := func(ctx context.Context, b Body) (any, error) {
//	    // do something with b
//	    return nil, nil
//	}
//
// r.Post("/foo", adapters.Action(fn, http.StatusCreated))
func Action[T any, Y any](f AdapterFunc[T, Y], ok int) server.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		v, err := decode[T](r)
		if err != nil {
			return err
		}

		res, err := f(r.Context(), v)
		if err != nil {
			return err
		}

		return server.Respond(w, ok, res)
	}
}

// ActionID functions the same as Action, but it also decodes a UUID from the URL path.
//
// Example:
//
//	type Body struct {
//	    Foo string `json:"foo"`
//	}
//
//	fn := func(ctx context.Context, ID uuid.UUID, b Body) (any, error) {
//	    // do something with ID and b
//	    return nil, nil
//	}
//
//	r.Post("/foo/{id}", adapters.ActionID(fn, http.StatusCreated))
func ActionID[T any, Y any](param string, f IDFunc[T, Y], ok int) server.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		ID, err := routeUUID(r, param)
		if err != nil {
			return err
		}

		v, err := decode[T](r)
		if err != nil {
			return err
		}

		res, err := f(r.Context(), ID, v)
		if err != nil {
			return err
		}

		return server.Respond(w, ok, res)
	}
}
