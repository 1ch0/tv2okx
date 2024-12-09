package utils

import (
	"context"
)

type contextKey int

const (
	projectKey contextKey = iota
	usernameKey
	permissionKey
)

// WithProject carries project in context
func WithProject(parent context.Context, project string) context.Context {
	return context.WithValue(parent, projectKey, project)
}

// ProjectFrom extract project from context
func ProjectFrom(ctx context.Context) (string, bool) {
	project, ok := ctx.Value(projectKey).(string)
	return project, ok
}

// WithUsername carries username in context
func WithUsername(parent context.Context, username string) context.Context {
	return context.WithValue(parent, usernameKey, username)
}

// UsernameFrom extract username from context
func UsernameFrom(ctx context.Context) (string, bool) {
	username, ok := ctx.Value(usernameKey).(string)
	return username, ok
}

// WithUserRole carries user role in context
func WithUserRole(parent context.Context, roles []string) context.Context {
	return context.WithValue(parent, permissionKey, roles)
}

// UserRoleFrom extract user role from context
func UserRoleFrom(ctx context.Context) ([]string, bool) {
	roles, ok := ctx.Value(permissionKey).([]string)
	return roles, ok
}
