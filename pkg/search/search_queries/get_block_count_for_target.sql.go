// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: get_block_count_for_target.sql

package search_queries

import (
	"context"
)

const getBlockedByCountForTarget = `-- name: GetBlockedByCountForTarget :one
SELECT COUNT(*) AS count
FROM author_blocks
WHERE target_did = $1
`

func (q *Queries) GetBlockedByCountForTarget(ctx context.Context, targetDid string) (int64, error) {
	row := q.queryRow(ctx, q.getBlockedByCountForTargetStmt, getBlockedByCountForTarget, targetDid)
	var count int64
	err := row.Scan(&count)
	return count, err
}
