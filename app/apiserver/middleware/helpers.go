package middleware

import (
	"net/http"

	"github.com/inhumanLightBackend/app/models/roles"
	"github.com/spf13/cast"
)

// Get user map in context of request
func UserContextMap(ctx interface{}) map[string]string {
	return cast.ToStringMapString(ctx)
}

// Check if user in context have admin privilege
func IsAdmin(r *http.Request) bool {
	userCtx := UserContextMap(r.Context().Value(CtxUserKey))
	return userCtx["access"] == roles.ADMIN
}