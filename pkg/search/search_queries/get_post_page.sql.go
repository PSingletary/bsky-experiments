// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: get_post_page.sql

package search_queries

import (
	"context"
	"database/sql"
	"time"
)

const getPostPage = `-- name: GetPostPage :many
SELECT p.id,
    p.text,
    p.parent_post_id,
    p.root_post_id,
    p.author_did,
    p.created_at,
    p.has_embedded_media,
    p.parent_relationship,
    p.sentiment,
    p.sentiment_confidence,
    p.indexed_at,
    COALESCE(
        json_agg(l.label) FILTER (
            WHERE l.label IS NOT NULL
        ),
        '[]'
    ) as labels
FROM posts p
    LEFT JOIN post_labels l on l.post_id = p.id
GROUP BY p.id
ORDER BY p.id
LIMIT $1 OFFSET $2
`

type GetPostPageParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type GetPostPageRow struct {
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
	IndexedAt           sql.NullTime    `json:"indexed_at"`
	Labels              interface{}     `json:"labels"`
}

func (q *Queries) GetPostPage(ctx context.Context, arg GetPostPageParams) ([]GetPostPageRow, error) {
	rows, err := q.query(ctx, q.getPostPageStmt, getPostPage, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPostPageRow
	for rows.Next() {
		var i GetPostPageRow
		if err := rows.Scan(
			&i.ID,
			&i.Text,
			&i.ParentPostID,
			&i.RootPostID,
			&i.AuthorDid,
			&i.CreatedAt,
			&i.HasEmbeddedMedia,
			&i.ParentRelationship,
			&i.Sentiment,
			&i.SentimentConfidence,
			&i.IndexedAt,
			&i.Labels,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
