package gql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/strrl/kubernetes-auditing-dashboard/ent"
	"github.com/strrl/kubernetes-auditing-dashboard/ent/auditevent"
)

// CompletedRequestResponseAuditEvents is the resolver for the completedRequestResponseAuditEvents field.
func (r *queryResolver) CompletedRequestResponseAuditEvents(ctx context.Context, page *int, pageSize *int) (*AuditEventPagination, error) {
	count, err := r.Resolver.entClient.AuditEvent.Query().
		Where(auditevent.LevelEQ("RequestResponse")).
		Where(auditevent.StageEQ("ResponseComplete")).
		Count(ctx)
	if err != nil {
		return nil, err
	}
	offset, limit := paginationToSQL(page, pageSize)
	rows, err := r.Resolver.entClient.AuditEvent.Query().
		Where(auditevent.LevelEQ("RequestResponse")).
		Where(auditevent.StageEQ("ResponseComplete")).
		Order(ent.Asc(auditevent.FieldID)).
		Offset(offset).
		Limit(limit).
		All(ctx)

	if err != nil {
		return nil, err
	}

	realPage := DefaultPage
	realPageSize := DefaultPageSize
	if page != nil {
		realPage = *page
	}
	if pageSize != nil {
		realPageSize = *pageSize
	}
	result := &AuditEventPagination{
		Total:           count,
		Page:            realPage,
		PageSize:        realPageSize,
		TotalPages:      count/realPageSize + 1,
		HasNextPage:     count > (realPage+1)*realPageSize,
		HasPreviousPage: realPage > 0,
		Rows:            rows,
	}
	return result, nil
}
