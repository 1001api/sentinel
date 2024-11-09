-- name: CheckAdminExist :one
SELECT EXISTS(SELECT 1 FROM users WHERE root_user = true);

-- name: FindUserByEmail :one
SELECT id, fullname, email, profile_url FROM users WHERE email = $1;

-- name: FindUserByEmailWithHash :one
SELECT id, email, password_hashed FROM users WHERE email = $1;

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
    fullname, email, password_hashed, profile_url, root_user, public_key
) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, fullname, email, profile_url;
