// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: remove_post_like.sql

package search_queries

import (
	"context"
)

const removeLikeFromPost = `-- name: RemoveLikeFromPost :exec
UPDATE post_likes
SET like_count = GREATEST(0, like_count - 1)
WHERE post_id = $1
`

func (q *Queries) RemoveLikeFromPost(ctx context.Context, postID string) error {
	_, err := q.exec(ctx, q.removeLikeFromPostStmt, removeLikeFromPost, postID)
	return err
}
