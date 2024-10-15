-- name: CreateProject :one
INSERT INTO 
    projects(name, description, user_id, created_at)
    VALUES ($1, $2, $3, $4)
RETURNING name, description, created_at;

-- name: UpdateProject :exec
UPDATE projects SET name = $1, description = $2 WHERE id = $3 AND user_id = $4 AND deleted_at IS NULL;

-- name: FindAllProjects :many
SELECT id, name, description, created_at FROM projects WHERE user_id = $1 AND deleted_at IS NULL;

-- name: FindProjectByID :one
SELECT id, name, description, created_at FROM projects WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL;

-- name: CheckProjectWithinUserID :one
SELECT EXISTS(SELECT 1 FROM projects WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL);

-- name: CountProject :one
SELECT COUNT(*) FROM projects WHERE user_id = $1 AND deleted_at IS NULL;

-- name: DeleteProject :exec
UPDATE projects SET deleted_at = NOW() WHERE user_id = $1 AND id = $2 AND deleted_at IS NULL;
