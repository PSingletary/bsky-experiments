// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: get_post_page_cursor.sql

package search_queries

import (
	"context"
	"database/sql"
	"time"
)

const getPostPageCursor = `-- name: GetPostPageCursor :many
WITH filtered_posts AS (
    SELECT id, text, parent_post_id, root_post_id, author_did, created_at, has_embedded_media, parent_relationship, sentiment, sentiment_confidence, indexed_at
    FROM posts
    WHERE created_at < $1
    ORDER BY created_at DESC
    LIMIT $2
)
SELECT fp.id,
    fp.text,
    fp.parent_post_id,
    fp.root_post_id,
    fp.author_did,
    fp.created_at,
    fp.has_embedded_media,
    fp.parent_relationship,
    fp.sentiment,
    fp.sentiment_confidence,
    fp.indexed_at,
    (
        SELECT COALESCE(
                json_agg(l.label) FILTER (
                    WHERE l.label IS NOT NULL
                ),
                '[]'
            )
        FROM post_labels l
        WHERE l.post_id = fp.id
    ) as labels
FROM filtered_posts fp
`

type GetPostPageCursorParams struct {
	CreatedAt time.Time `json:"created_at"`
	Limit     int32     `json:"limit"`
}

type GetPostPageCursorRow struct {
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

func (q *Queries) GetPostPageCursor(ctx context.Context, arg GetPostPageCursorParams) ([]GetPostPageCursorRow, error) {
	rows, err := q.query(ctx, q.getPostPageCursorStmt, getPostPageCursor, arg.CreatedAt, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPostPageCursorRow
	for rows.Next() {
		var i GetPostPageCursorRow
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
