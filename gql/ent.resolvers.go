package gql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.48

import (
	"context"

	"entgo.io/contrib/entgql"
	"github.com/strrl/kubernetes-auditing-dashboard/ent"
)

// Node is the resolver for the node field.
func (r *queryResolver) Node(ctx context.Context, id int) (ent.Noder, error) {
	return r.entClient.Noder(ctx, id)
}

// Nodes is the resolver for the nodes field.
func (r *queryResolver) Nodes(ctx context.Context, ids []int) ([]ent.Noder, error) {
	return r.entClient.Noders(ctx, ids)
}

// AuditEvents is the resolver for the auditEvents field.
func (r *queryResolver) AuditEvents(ctx context.Context, after *entgql.Cursor[int], first *int, before *entgql.Cursor[int], last *int, where *ent.AuditEventWhereInput) (*ent.AuditEventConnection, error) {
	return r.entClient.AuditEvent.Query().Paginate(ctx, after, first, before, last,
		ent.WithAuditEventFilter(where.Filter),
	)
}

// ResourceKinds is the resolver for the resourceKinds field.
func (r *queryResolver) ResourceKinds(ctx context.Context, after *entgql.Cursor[int], first *int, before *entgql.Cursor[int], last *int, where *ent.ResourceKindWhereInput) (*ent.ResourceKindConnection, error) {
	return r.entClient.ResourceKind.Query().Paginate(ctx, after, first, before, last,
		ent.WithResourceKindFilter(where.Filter),
	)
}

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
