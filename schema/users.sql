-- name: FindUserByEmail :one
SELECT id, fullname, email, profile_url FROM users WHERE email = $1;

-- name: FindUserByID :one
SELECT id, fullname, email, profile_url FROM users WHERE id = $1;

-- name: FindUserPublicKey :one
SELECT public_key FROM users WHERE id = $1;

-- name: FindUserByPublicKey :one
SELECT id, fullname, email, profile_url FROM users WHERE public_key = $1;

-- name: CheckUserIDExist :one
SELECT EXISTS(SELECT 1 FROM users WHERE id = $1);

-- name: CreateUser :one
INSERT INTO users(
    fullname, email, oauth_provider, oauth_id, profile_url, public_key
) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, fullname, email, profile_url;
