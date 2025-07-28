package scope

import (
	"context"

	"gitlab.com/tantai-kanban/kanban-api/internal/models"
)

type PayloadCtxKey struct{}
type ScopeCtxKey struct{}
type ThirdPartyScopeKey struct{}
type SessionUserCtxKey struct{}

// SetPayloadToContext sets the payload to context
func SetPayloadToContext(ctx context.Context, payload Payload) context.Context {
	return context.WithValue(ctx, PayloadCtxKey{}, payload)
}

// GetPayloadFromContext gets the payload from context
func GetPayloadFromContext(ctx context.Context) (Payload, bool) {
	payload, ok := ctx.Value(PayloadCtxKey{}).(Payload)
	return payload, ok
}

// GetSubFromContext gets the subject from context
func GetUserIdFromContext(ctx context.Context) (string, bool) {
	payload, ok := GetPayloadFromContext(ctx)
	if !ok {
		return "", false
	}
	return payload.UserID, true
}

// GetShopIDFromContext gets the shop id from context
func GetEmailFromContext(ctx context.Context) (string, bool) {
	payload, ok := GetPayloadFromContext(ctx)
	if !ok {
		return "", false
	}
	return payload.Email, true
}

func SetScopeToContext(ctx context.Context, scope models.Scope) context.Context {
	return context.WithValue(ctx, ScopeCtxKey{}, scope)
}

// GetScopeFromContext gets the scope from context
func GetScopeFromContext(ctx context.Context) (models.Scope, bool) {
	scope, ok := ctx.Value(ScopeCtxKey{}).(models.Scope)
	return scope, ok
}
