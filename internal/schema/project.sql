-- name: CreateProject :one
INSERT INTO 
    projects(name, description, url, user_id, created_at)
    VALUES ($1, $2, $3, $4, $5)
RETURNING name, description, created_at;

-- name: UpdateProject :exec
UPDATE projects SET name = $1, description = $2, url = $3 WHERE id = $4 AND user_id = $5 AND deleted_at IS NULL;

-- name: FindAllProjects :many
SELECT id, name, description, url, created_at FROM projects WHERE user_id = $1 AND deleted_at IS NULL;

-- name: FindProjectByID :one
SELECT id, name, description, url, created_at FROM projects WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL;

-- name: CheckProjectWithinUserID :one
SELECT EXISTS(SELECT 1 FROM projects WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL);

-- name: CountProject :one
SELECT COUNT(*) FROM projects WHERE user_id = $1 AND deleted_at IS NULL;

-- name: DeleteProject :exec
UPDATE projects SET deleted_at = NOW() WHERE user_id = $1 AND id = $2 AND deleted_at IS NULL;

-- name: CountProjectSize :one
SELECT 
    COALESCE(SUM(pg_column_size(events.*)) / 1024, 0)::bigint AS total_project_size
FROM events 
WHERE user_id = $1 AND project_id = $2;

-- name: LastProjectDataReceived :one
SELECT received_at FROM events 
WHERE user_id = $1 AND project_id = $2
ORDER BY received_at DESC LIMIT 1;
