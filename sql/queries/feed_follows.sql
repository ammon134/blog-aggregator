-- name: CreateFeedFollow :one
INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows
WHERE user_id = $1 AND id = $2;

-- name: ListFeedFollowsByUserId :many
SELECT * FROM feed_follows
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: GetFeedFollowByUserIdFeedID :one
SELECT * FROM feed_follows
WHERE user_id = $1 AND feed_id = $2;
