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

-- name: CheckProjectAggrEligibility :one
WITH latest_aggr AS (
    SELECT COALESCE(MAX(aggregated_at), '1970-01-01'::timestamp) AS last_aggregated
    FROM project_aggregations AS pa
    WHERE pa.project_id = $1
),
latest_events AS (
    SELECT COUNT(1) AS event_count
    FROM events AS e
    WHERE e.received_at > (SELECT la.last_aggregated FROM latest_aggr AS la)
    AND e.project_id = $1
)
SELECT * FROM latest_events;

-- name: CreateProjectAggr :exec
INSERT INTO project_aggregations (
    project_id, 
    user_id, 
    total_events, 
    total_event_types, 
    total_unique_users, 
    total_locations, 
    total_unique_page_urls, 
    most_visited_urls, 
    most_visited_countries, 
    most_visited_cities, 
    most_visited_regions, 
    most_firing_elements, 
    last_visited_users, 
    most_used_browsers, 
    most_fired_event_types, 
    most_fired_event_labels,
    aggregated_at,
    aggregated_at_str
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18
);

-- name: FindProjectAggr :many
SELECT * FROM project_aggregations
WHERE project_id = $1 AND user_id = $2
ORDER BY aggregated_at DESC LIMIT $3;
