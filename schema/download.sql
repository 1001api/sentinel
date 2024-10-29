-- name: GetEventTableHeaders :many
SELECT column_name::text FROM information_schema.columns WHERE table_schema = 'public' AND table_name = 'events';

-- name: DownloadLastMonthData :many
SELECT * FROM events 
WHERE user_id = $1 AND project_id = $2
AND received_at > date_trunc('month', NOW());
