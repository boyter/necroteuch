package data

import "context"

func (q *Queries) CreateSchema(ctx context.Context, schema string) error {
	_, err := q.db.ExecContext(ctx, schema)
	return err
}
