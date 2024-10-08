// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"math"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/strrl/kubernetes-auditing-dashboard/ent/predicate"
	"github.com/strrl/kubernetes-auditing-dashboard/ent/resourcekind"
)

// ResourceKindQuery is the builder for querying ResourceKind entities.
type ResourceKindQuery struct {
	config
	ctx        *QueryContext
	order      []resourcekind.OrderOption
	inters     []Interceptor
	predicates []predicate.ResourceKind
	modifiers  []func(*sql.Selector)
	loadTotal  []func(context.Context, []*ResourceKind) error
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the ResourceKindQuery builder.
func (rkq *ResourceKindQuery) Where(ps ...predicate.ResourceKind) *ResourceKindQuery {
	rkq.predicates = append(rkq.predicates, ps...)
	return rkq
}

// Limit the number of records to be returned by this query.
func (rkq *ResourceKindQuery) Limit(limit int) *ResourceKindQuery {
	rkq.ctx.Limit = &limit
	return rkq
}

// Offset to start from.
func (rkq *ResourceKindQuery) Offset(offset int) *ResourceKindQuery {
	rkq.ctx.Offset = &offset
	return rkq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (rkq *ResourceKindQuery) Unique(unique bool) *ResourceKindQuery {
	rkq.ctx.Unique = &unique
	return rkq
}

// Order specifies how the records should be ordered.
func (rkq *ResourceKindQuery) Order(o ...resourcekind.OrderOption) *ResourceKindQuery {
	rkq.order = append(rkq.order, o...)
	return rkq
}

// First returns the first ResourceKind entity from the query.
// Returns a *NotFoundError when no ResourceKind was found.
func (rkq *ResourceKindQuery) First(ctx context.Context) (*ResourceKind, error) {
	nodes, err := rkq.Limit(1).All(setContextOp(ctx, rkq.ctx, ent.OpQueryFirst))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{resourcekind.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (rkq *ResourceKindQuery) FirstX(ctx context.Context) *ResourceKind {
	node, err := rkq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first ResourceKind ID from the query.
// Returns a *NotFoundError when no ResourceKind ID was found.
func (rkq *ResourceKindQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = rkq.Limit(1).IDs(setContextOp(ctx, rkq.ctx, ent.OpQueryFirstID)); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{resourcekind.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (rkq *ResourceKindQuery) FirstIDX(ctx context.Context) int {
	id, err := rkq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single ResourceKind entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one ResourceKind entity is found.
// Returns a *NotFoundError when no ResourceKind entities are found.
func (rkq *ResourceKindQuery) Only(ctx context.Context) (*ResourceKind, error) {
	nodes, err := rkq.Limit(2).All(setContextOp(ctx, rkq.ctx, ent.OpQueryOnly))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{resourcekind.Label}
	default:
		return nil, &NotSingularError{resourcekind.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (rkq *ResourceKindQuery) OnlyX(ctx context.Context) *ResourceKind {
	node, err := rkq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only ResourceKind ID in the query.
// Returns a *NotSingularError when more than one ResourceKind ID is found.
// Returns a *NotFoundError when no entities are found.
func (rkq *ResourceKindQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = rkq.Limit(2).IDs(setContextOp(ctx, rkq.ctx, ent.OpQueryOnlyID)); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{resourcekind.Label}
	default:
		err = &NotSingularError{resourcekind.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (rkq *ResourceKindQuery) OnlyIDX(ctx context.Context) int {
	id, err := rkq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of ResourceKinds.
func (rkq *ResourceKindQuery) All(ctx context.Context) ([]*ResourceKind, error) {
	ctx = setContextOp(ctx, rkq.ctx, ent.OpQueryAll)
	if err := rkq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*ResourceKind, *ResourceKindQuery]()
	return withInterceptors[[]*ResourceKind](ctx, rkq, qr, rkq.inters)
}

// AllX is like All, but panics if an error occurs.
func (rkq *ResourceKindQuery) AllX(ctx context.Context) []*ResourceKind {
	nodes, err := rkq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of ResourceKind IDs.
func (rkq *ResourceKindQuery) IDs(ctx context.Context) (ids []int, err error) {
	if rkq.ctx.Unique == nil && rkq.path != nil {
		rkq.Unique(true)
	}
	ctx = setContextOp(ctx, rkq.ctx, ent.OpQueryIDs)
	if err = rkq.Select(resourcekind.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (rkq *ResourceKindQuery) IDsX(ctx context.Context) []int {
	ids, err := rkq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (rkq *ResourceKindQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, rkq.ctx, ent.OpQueryCount)
	if err := rkq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, rkq, querierCount[*ResourceKindQuery](), rkq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (rkq *ResourceKindQuery) CountX(ctx context.Context) int {
	count, err := rkq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (rkq *ResourceKindQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, rkq.ctx, ent.OpQueryExist)
	switch _, err := rkq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (rkq *ResourceKindQuery) ExistX(ctx context.Context) bool {
	exist, err := rkq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the ResourceKindQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (rkq *ResourceKindQuery) Clone() *ResourceKindQuery {
	if rkq == nil {
		return nil
	}
	return &ResourceKindQuery{
		config:     rkq.config,
		ctx:        rkq.ctx.Clone(),
		order:      append([]resourcekind.OrderOption{}, rkq.order...),
		inters:     append([]Interceptor{}, rkq.inters...),
		predicates: append([]predicate.ResourceKind{}, rkq.predicates...),
		// clone intermediate query.
		sql:  rkq.sql.Clone(),
		path: rkq.path,
	}
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		Name string `json:"name,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.ResourceKind.Query().
//		GroupBy(resourcekind.FieldName).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (rkq *ResourceKindQuery) GroupBy(field string, fields ...string) *ResourceKindGroupBy {
	rkq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &ResourceKindGroupBy{build: rkq}
	grbuild.flds = &rkq.ctx.Fields
	grbuild.label = resourcekind.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		Name string `json:"name,omitempty"`
//	}
//
//	client.ResourceKind.Query().
//		Select(resourcekind.FieldName).
//		Scan(ctx, &v)
func (rkq *ResourceKindQuery) Select(fields ...string) *ResourceKindSelect {
	rkq.ctx.Fields = append(rkq.ctx.Fields, fields...)
	sbuild := &ResourceKindSelect{ResourceKindQuery: rkq}
	sbuild.label = resourcekind.Label
	sbuild.flds, sbuild.scan = &rkq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a ResourceKindSelect configured with the given aggregations.
func (rkq *ResourceKindQuery) Aggregate(fns ...AggregateFunc) *ResourceKindSelect {
	return rkq.Select().Aggregate(fns...)
}

func (rkq *ResourceKindQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range rkq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, rkq); err != nil {
				return err
			}
		}
	}
	for _, f := range rkq.ctx.Fields {
		if !resourcekind.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if rkq.path != nil {
		prev, err := rkq.path(ctx)
		if err != nil {
			return err
		}
		rkq.sql = prev
	}
	return nil
}

func (rkq *ResourceKindQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*ResourceKind, error) {
	var (
		nodes = []*ResourceKind{}
		_spec = rkq.querySpec()
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*ResourceKind).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &ResourceKind{config: rkq.config}
		nodes = append(nodes, node)
		return node.assignValues(columns, values)
	}
	if len(rkq.modifiers) > 0 {
		_spec.Modifiers = rkq.modifiers
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, rkq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	for i := range rkq.loadTotal {
		if err := rkq.loadTotal[i](ctx, nodes); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (rkq *ResourceKindQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := rkq.querySpec()
	if len(rkq.modifiers) > 0 {
		_spec.Modifiers = rkq.modifiers
	}
	_spec.Node.Columns = rkq.ctx.Fields
	if len(rkq.ctx.Fields) > 0 {
		_spec.Unique = rkq.ctx.Unique != nil && *rkq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, rkq.driver, _spec)
}

func (rkq *ResourceKindQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(resourcekind.Table, resourcekind.Columns, sqlgraph.NewFieldSpec(resourcekind.FieldID, field.TypeInt))
	_spec.From = rkq.sql
	if unique := rkq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if rkq.path != nil {
		_spec.Unique = true
	}
	if fields := rkq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, resourcekind.FieldID)
		for i := range fields {
			if fields[i] != resourcekind.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := rkq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := rkq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := rkq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := rkq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (rkq *ResourceKindQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(rkq.driver.Dialect())
	t1 := builder.Table(resourcekind.Table)
	columns := rkq.ctx.Fields
	if len(columns) == 0 {
		columns = resourcekind.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if rkq.sql != nil {
		selector = rkq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if rkq.ctx.Unique != nil && *rkq.ctx.Unique {
		selector.Distinct()
	}
	for _, p := range rkq.predicates {
		p(selector)
	}
	for _, p := range rkq.order {
		p(selector)
	}
	if offset := rkq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := rkq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// ResourceKindGroupBy is the group-by builder for ResourceKind entities.
type ResourceKindGroupBy struct {
	selector
	build *ResourceKindQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (rkgb *ResourceKindGroupBy) Aggregate(fns ...AggregateFunc) *ResourceKindGroupBy {
	rkgb.fns = append(rkgb.fns, fns...)
	return rkgb
}

// Scan applies the selector query and scans the result into the given value.
func (rkgb *ResourceKindGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, rkgb.build.ctx, ent.OpQueryGroupBy)
	if err := rkgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*ResourceKindQuery, *ResourceKindGroupBy](ctx, rkgb.build, rkgb, rkgb.build.inters, v)
}

func (rkgb *ResourceKindGroupBy) sqlScan(ctx context.Context, root *ResourceKindQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(rkgb.fns))
	for _, fn := range rkgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*rkgb.flds)+len(rkgb.fns))
		for _, f := range *rkgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*rkgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := rkgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// ResourceKindSelect is the builder for selecting fields of ResourceKind entities.
type ResourceKindSelect struct {
	*ResourceKindQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (rks *ResourceKindSelect) Aggregate(fns ...AggregateFunc) *ResourceKindSelect {
	rks.fns = append(rks.fns, fns...)
	return rks
}

// Scan applies the selector query and scans the result into the given value.
func (rks *ResourceKindSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, rks.ctx, ent.OpQuerySelect)
	if err := rks.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*ResourceKindQuery, *ResourceKindSelect](ctx, rks.ResourceKindQuery, rks, rks.inters, v)
}

func (rks *ResourceKindSelect) sqlScan(ctx context.Context, root *ResourceKindQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(rks.fns))
	for _, fn := range rks.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*rks.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := rks.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
