// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: add_post.sql

package search_queries

import (
	"context"
	"database/sql"
	"time"
)

const addPost = `-- name: AddPost :exec
INSERT INTO posts (id, text, parent_post_id, root_post_id, author_did, created_at, has_embedded_media, parent_relationship, sentiment, sentiment_confidence)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
`

type AddPostParams struct {
	ID                  string          `json:"id"`
	Text                string          `json:"text"`
	ParentPostID        sql.NullString  `json:"parent_post_id"`
	RootPostID          sql.NullString  `json:"root_post_id"`
	AuthorDid           string          `json:"author_did"`
	CreatedAt           time.Time       `json:"created_at"`
	HasEmbeddedMedia    bool            `json:"has_embedded_media"`
	ParentRelationship  sql.NullString  `json:"parent_relationship"`
	Sentiment           sql.NullString  `json:"sentiment"`
	SentimentConfidence sql.NullFloat64 `json:"sentiment_confidence"`
}

func (q *Queries) AddPost(ctx context.Context, arg AddPostParams) error {
	_, err := q.exec(ctx, q.addPostStmt, addPost,
		arg.ID,
		arg.Text,
		arg.ParentPostID,
		arg.RootPostID,
		arg.AuthorDid,
		arg.CreatedAt,
		arg.HasEmbeddedMedia,
		arg.ParentRelationship,
		arg.Sentiment,
		arg.SentimentConfidence,
	)
	return err
}
