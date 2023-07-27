// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: set_indexed_timestamp.sql

package search_queries

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

const setPostIndexedTimestamp = `-- name: SetPostIndexedTimestamp :exec
UPDATE posts
SET indexed_at = $1
WHERE id = ANY($2::text [])
`

type SetPostIndexedTimestampParams struct {
	IndexedAt sql.NullTime `json:"indexed_at"`
	PostIds   []string     `json:"post_ids"`
}

func (q *Queries) SetPostIndexedTimestamp(ctx context.Context, arg SetPostIndexedTimestampParams) error {
	_, err := q.exec(ctx, q.setPostIndexedTimestampStmt, setPostIndexedTimestamp, arg.IndexedAt, pq.Array(arg.PostIds))
	return err
}
