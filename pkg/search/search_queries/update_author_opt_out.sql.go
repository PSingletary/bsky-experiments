// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: update_author_opt_out.sql

package search_queries

import (
	"context"
)

const updateAuthorOptOut = `-- name: UpdateAuthorOptOut :exec
UPDATE authors
SET cluster_opt_out = $2
WHERE did = $1
`

type UpdateAuthorOptOutParams struct {
	Did           string `json:"did"`
	ClusterOptOut bool   `json:"cluster_opt_out"`
}

func (q *Queries) UpdateAuthorOptOut(ctx context.Context, arg UpdateAuthorOptOutParams) error {
	_, err := q.exec(ctx, q.updateAuthorOptOutStmt, updateAuthorOptOut, arg.Did, arg.ClusterOptOut)
	return err
}
