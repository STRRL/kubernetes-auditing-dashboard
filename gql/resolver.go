package gql

import "github.com/strrl/kubernetes-auditing-dashboard/ent"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	entClient *ent.Client
}

func NewResolver(entClient *ent.Client) *Resolver {
	return &Resolver{entClient: entClient}
}
