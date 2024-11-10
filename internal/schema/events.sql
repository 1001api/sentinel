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
WHERE user_id = $2 AND project_id = $1 AND received_at >= NOW() - INTERVAL '1 hour'
ORDER BY received_at DESC
LIMIT 50;

-- name: GetEventSummary :one
SELECT 
COUNT(e.id) AS total_events,
COUNT(DISTINCT e.ip_addr) AS total_unique_users,
COUNT(DISTINCT e.event_type) AS total_event_type,
COUNT(DISTINCT e.country) AS total_country_visited,
(
    SELECT sub.page_url FROM events AS sub
    WHERE sub.user_id = $2 AND sub.project_id = $1
    GROUP BY sub.page_url ORDER BY COUNT(sub.page_url) DESC LIMIT 1
) AS most_visited_url,
(
    SELECT sub.country FROM events AS sub
    WHERE sub.user_id = $2 AND sub.project_id = $1
    GROUP BY sub.country ORDER BY COUNT(sub.country) DESC LIMIT 1
) AS most_country_visited
FROM events AS e WHERE e.user_id = $2 AND e.project_id = $1;

-- name: GetTotalEventSummary :one
SELECT 
    COUNT(id) AS total_events,
    COUNT(DISTINCT event_type) AS total_event_type,
    COUNT(DISTINCT ip_addr) AS total_unique_users,
    COUNT(DISTINCT country) AS total_country_visited,
    COUNT(DISTINCT page_url) AS total_page_url
FROM events WHERE user_id = $2 AND project_id = $1;

-- name: GetEventDetailSummary :many
WITH 
most_visited_url AS (
    SELECT 'most_visited_url' AS query_type, sub.page_url AS name, COUNT(sub.page_url) AS total
    FROM events sub
    WHERE sub.page_url IS NOT NULL AND sub.page_url <> ''
    AND sub.user_id = $2 AND sub.project_id = $1
    GROUP BY sub.page_url 
    ORDER BY COUNT(*) DESC 
    LIMIT 5
),
most_visited_country AS (
    SELECT 'most_visited_country' AS query_type, sub.country AS name, COUNT(sub.country) AS total
    FROM events sub
    WHERE sub.country IS NOT NULL AND sub.country <> ''
    AND sub.user_id = $2 AND sub.project_id = $1
    GROUP BY sub.country 
    ORDER BY COUNT(*) DESC 
    LIMIT 5
),
most_visited_city AS (
    SELECT 'most_visited_city' AS query_type, sub.city AS name, COUNT(sub.city) AS total
    FROM events sub
    WHERE sub.city IS NOT NULL AND sub.city <> '' 
    AND sub.user_id = $2 AND sub.project_id = $1
    GROUP BY sub.city 
    ORDER BY COUNT(*) DESC 
    LIMIT 5
),
most_hit_element AS (
    SELECT 'most_hit_element' AS query_type, sub.element_path AS name, COUNT(sub.element_path) AS total
    FROM events sub
    WHERE sub.element_path IS NOT NULL AND sub.element_path <> '' 
    AND sub.user_id = $2 AND sub.project_id = $1
    GROUP BY sub.element_path 
    ORDER BY COUNT(*) DESC 
    LIMIT 5
),
last_visited_user AS (
    SELECT 'last_visited_user' AS query_type, 
           sub.ip_addr::text AS name, 
           sub.received_at AS timestamp
    FROM events sub
    WHERE sub.ip_addr IS NOT NULL
    AND sub.user_id = $2 AND sub.project_id = $1
    ORDER BY sub.received_at DESC 
    LIMIT 5
),
most_used_browser AS (
    SELECT 'most_used_browser' AS query_type, sub.browser_name AS name, COUNT(sub.browser_name) AS total
    FROM events sub
    WHERE sub.browser_name IS NOT NULL AND sub.browser_name <> ''
    AND sub.user_id = $2 AND sub.project_id = $1
    GROUP BY sub.browser_name 
    ORDER BY COUNT(*) DESC 
    LIMIT 5
),
most_event_type AS (
    SELECT 'most_event_type' AS query_type, sub.event_type AS name, COUNT(sub.event_type) AS total
    FROM events sub
    WHERE sub.event_type IS NOT NULL AND sub.event_type <> ''
    AND sub.user_id = $2 AND sub.project_id = $1
    GROUP BY sub.event_type 
    ORDER BY COUNT(*) DESC 
    LIMIT 5
),
most_event_label AS (
    SELECT 'most_event_label' AS query_type, sub.event_label AS name, COUNT(sub.event_label) AS total
    FROM events sub
    WHERE sub.event_label IS NOT NULL AND sub.event_label <> ''
    AND sub.user_id = $2 AND sub.project_id = $1
    GROUP BY sub.event_label 
    ORDER BY COUNT(*) DESC 
    LIMIT 5
)

SELECT query_type, name, CAST(total AS text) AS total -- Why cast total as text? so it can be used to also hold the timestamp
FROM (
    SELECT * FROM most_visited_url
    UNION ALL SELECT * FROM most_visited_country
    UNION ALL SELECT * FROM most_visited_city
    UNION ALL SELECT * FROM most_hit_element
    UNION ALL SELECT * FROM most_used_browser
    UNION ALL SELECT * FROM most_event_type
    UNION ALL SELECT * FROM most_event_label
) count_queries
UNION ALL

SELECT query_type, name, CAST(timestamp AS text) AS total
FROM last_visited_user;

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
