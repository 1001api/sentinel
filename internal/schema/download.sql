-- name: GetEventTableHeaders :many
SELECT column_name::text FROM information_schema.columns WHERE table_schema = 'public' AND table_name = 'events';

-- name: DownloadIntervalEventData :many
SELECT * FROM events 
WHERE user_id = $1 AND project_id = $2
AND (@interval::int = -1 OR received_at >= NOW() - INTERVAL '1 day' * @interval::int);
