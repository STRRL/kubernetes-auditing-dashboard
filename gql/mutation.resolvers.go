package gql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
)

// ImportResourceKindTsv is the resolver for the importResourceKindTSV field.
func (r *mutationResolver) ImportResourceKindTsv(ctx context.Context, tsv string) (int, error) {
	return 0, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
