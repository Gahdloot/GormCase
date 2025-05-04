package query

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

// Query represents a database query
type Query struct {
	table      string
	fields     []string
	conditions []string
	args       []interface{}
	orderBy    []string
	groupBy    []string
	having     []string
	limit      int
	offset     int
	joins      []string
	distinct   bool
}

// NewQuery creates a new Query
func NewQuery(table string) *Query {
	return &Query{
		table:      table,
		fields:     []string{"*"},
		conditions: []string{},
		args:       []interface{}{},
		orderBy:    []string{},
		groupBy:    []string{},
		having:     []string{},
		joins:      []string{},
	}
}

// Select sets the fields to select
func (q *Query) Select(fields ...string) *Query {
	q.fields = fields
	return q
}

// Where adds a WHERE condition
func (q *Query) Where(condition string, args ...interface{}) *Query {
	q.conditions = append(q.conditions, condition)
	q.args = append(q.args, args...)
	return q
}

// OrderBy adds an ORDER BY clause
func (q *Query) OrderBy(fields ...string) *Query {
	q.orderBy = append(q.orderBy, fields...)
	return q
}

// GroupBy adds a GROUP BY clause
func (q *Query) GroupBy(fields ...string) *Query {
	q.groupBy = append(q.groupBy, fields...)
	return q
}

// Having adds a HAVING clause
func (q *Query) Having(condition string, args ...interface{}) *Query {
	q.having = append(q.having, condition)
	q.args = append(q.args, args...)
	return q
}

// Limit sets the LIMIT clause
func (q *Query) Limit(limit int) *Query {
	q.limit = limit
	return q
}

// Offset sets the OFFSET clause
func (q *Query) Offset(offset int) *Query {
	q.offset = offset
	return q
}

// Join adds a JOIN clause
func (q *Query) Join(join string) *Query {
	q.joins = append(q.joins, join)
	return q
}

// Distinct sets the DISTINCT clause
func (q *Query) Distinct() *Query {
	q.distinct = true
	return q
}

// Build constructs the SQL query
func (q *Query) Build() (string, []interface{}) {
	var parts []string

	// SELECT
	if q.distinct {
		parts = append(parts, "SELECT DISTINCT")
	} else {
		parts = append(parts, "SELECT")
	}

	// Fields
	parts = append(parts, strings.Join(q.fields, ", "))

	// FROM
	parts = append(parts, "FROM", q.table)

	// JOINs
	if len(q.joins) > 0 {
		parts = append(parts, strings.Join(q.joins, " "))
	}

	// WHERE
	if len(q.conditions) > 0 {
		parts = append(parts, "WHERE", strings.Join(q.conditions, " AND "))
	}

	// GROUP BY
	if len(q.groupBy) > 0 {
		parts = append(parts, "GROUP BY", strings.Join(q.groupBy, ", "))
	}

	// HAVING
	if len(q.having) > 0 {
		parts = append(parts, "HAVING", strings.Join(q.having, " AND "))
	}

	// ORDER BY
	if len(q.orderBy) > 0 {
		parts = append(parts, "ORDER BY", strings.Join(q.orderBy, ", "))
	}

	// LIMIT
	if q.limit > 0 {
		parts = append(parts, fmt.Sprintf("LIMIT %d", q.limit))
	}

	// OFFSET
	if q.offset > 0 {
		parts = append(parts, fmt.Sprintf("OFFSET %d", q.offset))
	}

	return strings.Join(parts, " "), q.args
}

// Execute executes the query and returns the result
func (q *Query) Execute(ctx context.Context, db *sql.DB) (*sql.Rows, error) {
	query, args := q.Build()
	return db.QueryContext(ctx, query, args...)
}

// ExecuteRow executes the query and returns a single row
func (q *Query) ExecuteRow(ctx context.Context, db *sql.DB) *sql.Row {
	query, args := q.Build()
	return db.QueryRowContext(ctx, query, args...)
}

// Count executes a COUNT query
func (q *Query) Count(ctx context.Context, db *sql.DB) (int, error) {
	// Save original fields
	originalFields := q.fields
	q.fields = []string{"COUNT(*)"}

	// Execute query
	var count int
	err := q.ExecuteRow(ctx, db).Scan(&count)

	// Restore original fields
	q.fields = originalFields

	return count, err
}

// Exists checks if any rows match the query
func (q *Query) Exists(ctx context.Context, db *sql.DB) (bool, error) {
	// Save original fields and limit
	originalFields := q.fields
	originalLimit := q.limit
	q.fields = []string{"1"}
	q.limit = 1

	// Execute query
	rows, err := q.Execute(ctx, db)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	// Restore original fields and limit
	q.fields = originalFields
	q.limit = originalLimit

	return rows.Next(), nil
}
