// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: posts.sql

package store_queries

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

const createPost = `-- name: CreatePost :exec
INSERT INTO posts (
        actor_did,
        rkey,
        content,
        parent_post_actor_did,
        parent_post_rkey,
        quote_post_actor_did,
        quote_post_rkey,
        root_post_actor_did,
        root_post_rkey,
        has_embedded_media,
        created_at
    )
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7,
        $8,
        $9,
        $10,
        $11
    )
`

type CreatePostParams struct {
	ActorDid           string         `json:"actor_did"`
	Rkey               string         `json:"rkey"`
	Content            sql.NullString `json:"content"`
	ParentPostActorDid sql.NullString `json:"parent_post_actor_did"`
	ParentPostRkey     sql.NullString `json:"parent_post_rkey"`
	QuotePostActorDid  sql.NullString `json:"quote_post_actor_did"`
	QuotePostRkey      sql.NullString `json:"quote_post_rkey"`
	RootPostActorDid   sql.NullString `json:"root_post_actor_did"`
	RootPostRkey       sql.NullString `json:"root_post_rkey"`
	HasEmbeddedMedia   bool           `json:"has_embedded_media"`
	CreatedAt          sql.NullTime   `json:"created_at"`
}

func (q *Queries) CreatePost(ctx context.Context, arg CreatePostParams) error {
	_, err := q.exec(ctx, q.createPostStmt, createPost,
		arg.ActorDid,
		arg.Rkey,
		arg.Content,
		arg.ParentPostActorDid,
		arg.ParentPostRkey,
		arg.QuotePostActorDid,
		arg.QuotePostRkey,
		arg.RootPostActorDid,
		arg.RootPostRkey,
		arg.HasEmbeddedMedia,
		arg.CreatedAt,
	)
	return err
}

const deletePost = `-- name: DeletePost :exec
DELETE FROM posts
WHERE actor_did = $1
    AND rkey = $2
`

type DeletePostParams struct {
	ActorDid string `json:"actor_did"`
	Rkey     string `json:"rkey"`
}

func (q *Queries) DeletePost(ctx context.Context, arg DeletePostParams) error {
	_, err := q.exec(ctx, q.deletePostStmt, deletePost, arg.ActorDid, arg.Rkey)
	return err
}

const getMyPostsByFuzzyContent = `-- name: GetMyPostsByFuzzyContent :many
SELECT actor_did, rkey, content, parent_post_actor_did, quote_post_actor_did, quote_post_rkey, parent_post_rkey, root_post_actor_did, root_post_rkey, has_embedded_media, created_at, inserted_at
FROM posts
WHERE actor_did = $1
    AND content ILIKE concat('%', $4::text, '%')::text
    AND content not ilike '%!jazbot%'
ORDER BY created_at DESC
LIMIT $2 OFFSET $3
`

type GetMyPostsByFuzzyContentParams struct {
	ActorDid string `json:"actor_did"`
	Limit    int32  `json:"limit"`
	Offset   int32  `json:"offset"`
	Query    string `json:"query"`
}

func (q *Queries) GetMyPostsByFuzzyContent(ctx context.Context, arg GetMyPostsByFuzzyContentParams) ([]Post, error) {
	rows, err := q.query(ctx, q.getMyPostsByFuzzyContentStmt, getMyPostsByFuzzyContent,
		arg.ActorDid,
		arg.Limit,
		arg.Offset,
		arg.Query,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Post
	for rows.Next() {
		var i Post
		if err := rows.Scan(
			&i.ActorDid,
			&i.Rkey,
			&i.Content,
			&i.ParentPostActorDid,
			&i.QuotePostActorDid,
			&i.QuotePostRkey,
			&i.ParentPostRkey,
			&i.RootPostActorDid,
			&i.RootPostRkey,
			&i.HasEmbeddedMedia,
			&i.CreatedAt,
			&i.InsertedAt,
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

const getPost = `-- name: GetPost :one
SELECT actor_did, rkey, content, parent_post_actor_did, quote_post_actor_did, quote_post_rkey, parent_post_rkey, root_post_actor_did, root_post_rkey, has_embedded_media, created_at, inserted_at
FROM posts
WHERE actor_did = $1
    AND rkey = $2
`

type GetPostParams struct {
	ActorDid string `json:"actor_did"`
	Rkey     string `json:"rkey"`
}

func (q *Queries) GetPost(ctx context.Context, arg GetPostParams) (Post, error) {
	row := q.queryRow(ctx, q.getPostStmt, getPost, arg.ActorDid, arg.Rkey)
	var i Post
	err := row.Scan(
		&i.ActorDid,
		&i.Rkey,
		&i.Content,
		&i.ParentPostActorDid,
		&i.QuotePostActorDid,
		&i.QuotePostRkey,
		&i.ParentPostRkey,
		&i.RootPostActorDid,
		&i.RootPostRkey,
		&i.HasEmbeddedMedia,
		&i.CreatedAt,
		&i.InsertedAt,
	)
	return i, err
}

const getPostWithReplies = `-- name: GetPostWithReplies :many
WITH RootPost AS (
    SELECT p.actor_did, p.rkey, p.content, p.parent_post_actor_did, p.quote_post_actor_did, p.quote_post_rkey, p.parent_post_rkey, p.root_post_actor_did, p.root_post_rkey, p.has_embedded_media, p.created_at, p.inserted_at,
        array_agg(COALESCE(i.cid, ''))::TEXT [] as image_cids,
        array_agg(COALESCE(i.alt_text, ''))::TEXT [] as image_alts
    FROM posts p
        LEFT JOIN images i ON p.actor_did = i.post_actor_did
        AND p.rkey = i.post_rkey
    WHERE p.actor_did = $1
        AND p.rkey = $2
    GROUP BY p.actor_did,
        p.rkey
),
Replies AS (
    SELECT p.actor_did, p.rkey, p.content, p.parent_post_actor_did, p.quote_post_actor_did, p.quote_post_rkey, p.parent_post_rkey, p.root_post_actor_did, p.root_post_rkey, p.has_embedded_media, p.created_at, p.inserted_at,
        array_agg(COALESCE(i.cid, ''))::TEXT [] as image_cids,
        array_agg(COALESCE(i.alt_text, ''))::TEXT [] as image_alts
    FROM posts p
        LEFT JOIN images i ON p.actor_did = i.post_actor_did
        AND p.rkey = i.post_rkey
    WHERE p.parent_post_actor_did = (
            SELECT actor_did
            FROM RootPost
        )
        AND p.parent_post_rkey = (
            SELECT rkey
            FROM RootPost
        )
    GROUP BY p.actor_did,
        p.rkey
),
RootLikeCount AS (
    SELECT lc.subject_id,
        lc.num_likes
    FROM subjects s
        JOIN like_counts lc ON s.id = lc.subject_id
    WHERE s.actor_did = (
            SELECT actor_did
            FROM RootPost
        )
        AND s.rkey = (
            SELECT rkey
            FROM RootPost
        )
),
ReplyLikeCounts AS (
    SELECT s.actor_did,
        s.rkey,
        lc.num_likes
    FROM subjects s
        JOIN like_counts lc ON s.id = lc.subject_id
    WHERE s.actor_did IN (
            SELECT actor_did
            FROM Replies
        )
        AND s.rkey IN (
            SELECT rkey
            FROM Replies
        )
)
SELECT rp.actor_did, rp.rkey, rp.content, rp.parent_post_actor_did, rp.quote_post_actor_did, rp.quote_post_rkey, rp.parent_post_rkey, rp.root_post_actor_did, rp.root_post_rkey, rp.has_embedded_media, rp.created_at, rp.inserted_at, rp.image_cids, rp.image_alts,
    rlc.num_likes AS like_count
FROM RootPost rp
    LEFT JOIN RootLikeCount rlc ON rlc.subject_id = (
        SELECT id
        FROM subjects
        WHERE actor_did = rp.actor_did
            AND rkey = rp.rkey
    )
UNION ALL
SELECT r.actor_did, r.rkey, r.content, r.parent_post_actor_did, r.quote_post_actor_did, r.quote_post_rkey, r.parent_post_rkey, r.root_post_actor_did, r.root_post_rkey, r.has_embedded_media, r.created_at, r.inserted_at, r.image_cids, r.image_alts,
    rlc.num_likes AS like_count
FROM Replies r
    LEFT JOIN ReplyLikeCounts rlc ON r.actor_did = rlc.actor_did
    AND r.rkey = rlc.rkey
`

type GetPostWithRepliesParams struct {
	ActorDid string `json:"actor_did"`
	Rkey     string `json:"rkey"`
}

type GetPostWithRepliesRow struct {
	ActorDid           string         `json:"actor_did"`
	Rkey               string         `json:"rkey"`
	Content            sql.NullString `json:"content"`
	ParentPostActorDid sql.NullString `json:"parent_post_actor_did"`
	QuotePostActorDid  sql.NullString `json:"quote_post_actor_did"`
	QuotePostRkey      sql.NullString `json:"quote_post_rkey"`
	ParentPostRkey     sql.NullString `json:"parent_post_rkey"`
	RootPostActorDid   sql.NullString `json:"root_post_actor_did"`
	RootPostRkey       sql.NullString `json:"root_post_rkey"`
	HasEmbeddedMedia   bool           `json:"has_embedded_media"`
	CreatedAt          sql.NullTime   `json:"created_at"`
	InsertedAt         time.Time      `json:"inserted_at"`
	ImageCids          []string       `json:"image_cids"`
	ImageAlts          []string       `json:"image_alts"`
	LikeCount          sql.NullInt64  `json:"like_count"`
}

func (q *Queries) GetPostWithReplies(ctx context.Context, arg GetPostWithRepliesParams) ([]GetPostWithRepliesRow, error) {
	rows, err := q.query(ctx, q.getPostWithRepliesStmt, getPostWithReplies, arg.ActorDid, arg.Rkey)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPostWithRepliesRow
	for rows.Next() {
		var i GetPostWithRepliesRow
		if err := rows.Scan(
			&i.ActorDid,
			&i.Rkey,
			&i.Content,
			&i.ParentPostActorDid,
			&i.QuotePostActorDid,
			&i.QuotePostRkey,
			&i.ParentPostRkey,
			&i.RootPostActorDid,
			&i.RootPostRkey,
			&i.HasEmbeddedMedia,
			&i.CreatedAt,
			&i.InsertedAt,
			pq.Array(&i.ImageCids),
			pq.Array(&i.ImageAlts),
			&i.LikeCount,
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

const getPostsByActor = `-- name: GetPostsByActor :many
SELECT actor_did, rkey, content, parent_post_actor_did, quote_post_actor_did, quote_post_rkey, parent_post_rkey, root_post_actor_did, root_post_rkey, has_embedded_media, created_at, inserted_at
FROM posts
WHERE actor_did = $1
ORDER BY created_at DESC
LIMIT $2
`

type GetPostsByActorParams struct {
	ActorDid string `json:"actor_did"`
	Limit    int32  `json:"limit"`
}

func (q *Queries) GetPostsByActor(ctx context.Context, arg GetPostsByActorParams) ([]Post, error) {
	rows, err := q.query(ctx, q.getPostsByActorStmt, getPostsByActor, arg.ActorDid, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Post
	for rows.Next() {
		var i Post
		if err := rows.Scan(
			&i.ActorDid,
			&i.Rkey,
			&i.Content,
			&i.ParentPostActorDid,
			&i.QuotePostActorDid,
			&i.QuotePostRkey,
			&i.ParentPostRkey,
			&i.RootPostActorDid,
			&i.RootPostRkey,
			&i.HasEmbeddedMedia,
			&i.CreatedAt,
			&i.InsertedAt,
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

const getPostsByActorsFollowingTarget = `-- name: GetPostsByActorsFollowingTarget :many
WITH followers AS (
    SELECT actor_did
    FROM follows
    WHERE target_did = $1
)
SELECT p.actor_did, p.rkey, p.content, p.parent_post_actor_did, p.quote_post_actor_did, p.quote_post_rkey, p.parent_post_rkey, p.root_post_actor_did, p.root_post_rkey, p.has_embedded_media, p.created_at, p.inserted_at
FROM posts p
    JOIN followers f ON f.actor_did = p.actor_did
WHERE (p.created_at, p.actor_did, p.rkey) < (
        $3::TIMESTAMPTZ,
        $4::TEXT,
        $5::TEXT
    )
    AND (p.root_post_rkey IS NULL)
    AND (
        (p.parent_relationship IS NULL)
        OR (p.parent_relationship <> 'r'::text)
    )
ORDER BY p.created_at DESC,
    p.actor_did DESC,
    p.rkey DESC
LIMIT $2
`

type GetPostsByActorsFollowingTargetParams struct {
	TargetDid       string    `json:"target_did"`
	Limit           int32     `json:"limit"`
	CursorCreatedAt time.Time `json:"cursor_created_at"`
	CursorActorDid  string    `json:"cursor_actor_did"`
	CursorRkey      string    `json:"cursor_rkey"`
}

func (q *Queries) GetPostsByActorsFollowingTarget(ctx context.Context, arg GetPostsByActorsFollowingTargetParams) ([]Post, error) {
	rows, err := q.query(ctx, q.getPostsByActorsFollowingTargetStmt, getPostsByActorsFollowingTarget,
		arg.TargetDid,
		arg.Limit,
		arg.CursorCreatedAt,
		arg.CursorActorDid,
		arg.CursorRkey,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Post
	for rows.Next() {
		var i Post
		if err := rows.Scan(
			&i.ActorDid,
			&i.Rkey,
			&i.Content,
			&i.ParentPostActorDid,
			&i.QuotePostActorDid,
			&i.QuotePostRkey,
			&i.ParentPostRkey,
			&i.RootPostActorDid,
			&i.RootPostRkey,
			&i.HasEmbeddedMedia,
			&i.CreatedAt,
			&i.InsertedAt,
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

const getPostsFromNonMoots = `-- name: GetPostsFromNonMoots :many
WITH my_follows AS (
    SELECT target_did
    FROM follows
    WHERE follows.actor_did = $1
),
non_moots AS (
    SELECT actor_did
    FROM follows f
        LEFT JOIN my_follows ON f.actor_did = my_follows.target_did
    WHERE f.target_did = $1
        AND my_follows.target_did IS NULL
),
non_moots_and_non_spam AS (
    SELECT nm.actor_did
    FROM non_moots nm
        LEFT JOIN following_counts fc ON nm.actor_did = fc.actor_did
    WHERE fc.num_following < 4000
)
SELECT p.actor_did, p.rkey, p.content, p.parent_post_actor_did, p.quote_post_actor_did, p.quote_post_rkey, p.parent_post_rkey, p.root_post_actor_did, p.root_post_rkey, p.has_embedded_media, p.created_at, p.inserted_at
FROM posts p
    JOIN non_moots_and_non_spam f ON f.actor_did = p.actor_did
WHERE (p.created_at, p.actor_did, p.rkey) < (
        $3::TIMESTAMPTZ,
        $4::TEXT,
        $5::TEXT
    )
    AND p.root_post_rkey IS NULL
    AND p.parent_post_rkey IS NULL
    AND p.created_at > NOW() - make_interval(hours := 24)
ORDER BY p.created_at DESC,
    p.actor_did DESC,
    p.rkey DESC
LIMIT $2
`

type GetPostsFromNonMootsParams struct {
	ActorDid        string    `json:"actor_did"`
	Limit           int32     `json:"limit"`
	CursorCreatedAt time.Time `json:"cursor_created_at"`
	CursorActorDid  string    `json:"cursor_actor_did"`
	CursorRkey      string    `json:"cursor_rkey"`
}

func (q *Queries) GetPostsFromNonMoots(ctx context.Context, arg GetPostsFromNonMootsParams) ([]Post, error) {
	rows, err := q.query(ctx, q.getPostsFromNonMootsStmt, getPostsFromNonMoots,
		arg.ActorDid,
		arg.Limit,
		arg.CursorCreatedAt,
		arg.CursorActorDid,
		arg.CursorRkey,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Post
	for rows.Next() {
		var i Post
		if err := rows.Scan(
			&i.ActorDid,
			&i.Rkey,
			&i.Content,
			&i.ParentPostActorDid,
			&i.QuotePostActorDid,
			&i.QuotePostRkey,
			&i.ParentPostRkey,
			&i.RootPostActorDid,
			&i.RootPostRkey,
			&i.HasEmbeddedMedia,
			&i.CreatedAt,
			&i.InsertedAt,
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

const getTopPosts = `-- name: GetTopPosts :many
WITH TopSubjects AS (
    SELECT s.id, s.actor_did, s.rkey, s.col,
        lc.num_likes
    FROM subjects s
        JOIN like_counts lc ON lc.subject_id = s.id
    WHERE s.col = 1
        AND lc.num_likes > 100
    ORDER BY lc.num_likes DESC
    LIMIT $1 + 30 OFFSET $2
)
SELECT p.actor_did, p.rkey, p.content, p.parent_post_actor_did, p.quote_post_actor_did, p.quote_post_rkey, p.parent_post_rkey, p.root_post_actor_did, p.root_post_rkey, p.has_embedded_media, p.created_at, p.inserted_at
FROM posts p
    JOIN TopSubjects s ON p.actor_did = s.actor_did
    AND p.rkey = s.rkey
ORDER BY s.num_likes DESC
LIMIT $1
`

type GetTopPostsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) GetTopPosts(ctx context.Context, arg GetTopPostsParams) ([]Post, error) {
	rows, err := q.query(ctx, q.getTopPostsStmt, getTopPosts, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Post
	for rows.Next() {
		var i Post
		if err := rows.Scan(
			&i.ActorDid,
			&i.Rkey,
			&i.Content,
			&i.ParentPostActorDid,
			&i.QuotePostActorDid,
			&i.QuotePostRkey,
			&i.ParentPostRkey,
			&i.RootPostActorDid,
			&i.RootPostRkey,
			&i.HasEmbeddedMedia,
			&i.CreatedAt,
			&i.InsertedAt,
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

const getTopPostsForActor = `-- name: GetTopPostsForActor :many
WITH TopSubjects AS (
    SELECT s.id, s.actor_did, s.rkey, s.col,
        lc.num_likes
    FROM subjects s
        JOIN like_counts lc ON lc.subject_id = s.id
    WHERE s.col = 1
        AND lc.num_likes > 1
        AND s.actor_did = $1
    ORDER BY lc.num_likes DESC
    LIMIT $2 OFFSET $3
)
SELECT p.actor_did, p.rkey, p.content, p.parent_post_actor_did, p.quote_post_actor_did, p.quote_post_rkey, p.parent_post_rkey, p.root_post_actor_did, p.root_post_rkey, p.has_embedded_media, p.created_at, p.inserted_at
FROM posts p
    JOIN TopSubjects s ON p.actor_did = s.actor_did
    AND p.rkey = s.rkey
ORDER BY s.num_likes DESC
`

type GetTopPostsForActorParams struct {
	ActorDid string `json:"actor_did"`
	Limit    int32  `json:"limit"`
	Offset   int32  `json:"offset"`
}

func (q *Queries) GetTopPostsForActor(ctx context.Context, arg GetTopPostsForActorParams) ([]Post, error) {
	rows, err := q.query(ctx, q.getTopPostsForActorStmt, getTopPostsForActor, arg.ActorDid, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Post
	for rows.Next() {
		var i Post
		if err := rows.Scan(
			&i.ActorDid,
			&i.Rkey,
			&i.Content,
			&i.ParentPostActorDid,
			&i.QuotePostActorDid,
			&i.QuotePostRkey,
			&i.ParentPostRkey,
			&i.RootPostActorDid,
			&i.RootPostRkey,
			&i.HasEmbeddedMedia,
			&i.CreatedAt,
			&i.InsertedAt,
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
