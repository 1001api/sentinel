-- name: CreateEvent :exec
INSERT INTO events (
    event_type,
    event_label,
    page_url,
    element_path,
    element_type,
    ip_addr,
    user_agent,
    browser_name,
    country,
    region,
    city,
    session_id,
    device_type,
    time_on_page,
    screen_resolution,
    fired_at,
    received_at,
    user_id,
    project_id
) VALUES (
    $1,  -- event_type
    $2,  -- event_label
    $3,  -- page_url
    $4,  -- element_path
    $5,  -- element_type
    $6,  -- ip_addr
    $7,  -- user_agent
    $8,  -- browser_name
    $9,  -- country
    $10, -- region
    $11, -- city
    $12, -- session_id
    $13, -- device_type
    $14, -- time_on_page
    $15, -- screen_resolution
    $16, -- fired_at
    $17, -- received_at
    $18, -- user_id
    $19  -- project_id
);

-- name: GetLiveEvents :many
SELECT
    p.name,
    e.event_type,
    e.event_label,
    e.page_url,
    e.element_path,
    e.country,
    e.fired_at,
    e.received_at
FROM events AS e
JOIN projects AS p ON e.project_id = p.id
WHERE e.user_id = $1 AND e.received_at >= NOW() - INTERVAL '1 hour'
ORDER BY e.received_at DESC
LIMIT 100;

-- name: GetEvents :many
SELECT
    p.name AS project_name,
    e.event_type,
    e.event_label,
    e.page_url,
    e.element_path,
    e.element_type,
    e.ip_addr,
    e.user_agent,
    e.browser_name,
    e.country,
    e.region,
    e.city,
    e.session_id,
    e.device_type,
    e.time_on_page,
    e.screen_resolution,
    e.fired_at,
    e.received_at,
    e.project_id
FROM events AS e
JOIN projects AS p ON e.project_id = p.id
WHERE e.user_id = $1
AND (@interval::int = -1 OR received_at >= NOW() - INTERVAL '1 day' * @interval::int)
-- check if project id is provided and is not default empty UUID 
AND (@project_id::uuid = '00000000-0000-0000-0000-000000000000' OR e.project_id = @project_id) 
ORDER BY e.received_at DESC
LIMIT COALESCE(@limit_count::integer, 100);

-- name: GetLiveEventsDetail :many
SELECT 
    event_type,
    event_label,
    page_url,
    element_path,
    element_type,
    ip_addr,
    user_agent,
    browser_name,
    country,
    region,
    city,
    device_type,
    time_on_page,
    screen_resolution,
    fired_at,
    received_at
FROM events
WHERE
    user_id = $2 AND project_id = $1 
    AND (@byLastHour::bool IS NOT TRUE OR received_at >= NOW() - INTERVAL '1 minute')
ORDER BY received_at DESC
LIMIT COALESCE(@limit_count::integer, 100);

-- name: GetWeeklyEvents :many
SELECT
  DATE_TRUNC('day', received_at)::timestamp AS timestamp,
  COUNT(*) AS total
FROM events
WHERE user_id = $2 AND project_id = $1 AND received_at >= NOW() - INTERVAL '7 days'
GROUP BY timestamp ORDER BY timestamp ASC;

-- name: GetWeeklyEventsTotal :one
SELECT
COUNT(id) AS total
FROM events WHERE received_at >= NOW() - INTERVAL '7 days'
AND user_id = $2 AND project_id = $1;

-- name: DeleteEventByProjectID :exec
DELETE FROM events WHERE user_id = $1 AND project_id = $2;

-- name: GetPercentageEventsType :many
SELECT event_type, COUNT(*) AS total
FROM events
WHERE user_id = $2 AND project_id = $1
GROUP BY event_type
ORDER BY total DESC
LIMIT 10;

-- name: GetPercentageEventsLabel :many
SELECT event_label, count(event_label) AS total
FROM events
WHERE user_id = $2 AND project_id = $1
GROUP BY event_label
ORDER BY total DESC
LIMIT 10;

-- name: CountUserMonthlyEvents :one
SELECT COUNT(id) FROM events 
WHERE user_id = $1
AND received_at > date_trunc('month', NOW());
