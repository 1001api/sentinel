-- name: CreateAPIKey :one
INSERT INTO 
    api_keys(name, token, user_id, created_at, expired_at) VALUES ($1, $2, $3, $4, $5) 
RETURNING name, token, created_at, expired_at;

-- name: FindAllAPIKeys :many
SELECT id, name, token, created_at, expired_at FROM api_keys WHERE user_id = $1;

-- name: DeleteAPIKey :exec
DELETE FROM api_keys WHERE user_id = $1 AND id = $2;
