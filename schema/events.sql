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

-- name: GetEventSummaryDetailTLDR :one
SELECT 
    COUNT(id) AS total_events,
    COUNT(DISTINCT event_type) AS total_event_type,
    COUNT(DISTINCT ip_addr) AS total_unique_users,
    COUNT(DISTINCT country) AS total_country_visited,
    COUNT(DISTINCT page_url) AS total_page_url
FROM events WHERE user_id = $2 AND project_id = $1;

-- name: GetEventMostVisitedURL :many
SELECT page_url AS name, COUNT(page_url) AS total
FROM events WHERE user_id = $2 AND project_id = $1
GROUP BY page_url ORDER BY COUNT(page_url) DESC LIMIT 5;

-- name: GetEventMostVisitedCountry :many
SELECT country AS name, COUNT(country) AS total
FROM events 
WHERE country IS NOT NULL AND country <> ''
AND user_id = $2 AND project_id = $1
GROUP BY country ORDER BY COUNT(country) DESC LIMIT 5;

-- name: GetEventMostVisitedCity :many
SELECT city AS name, COUNT(city) AS total
FROM events 
WHERE city IS NOT NULL AND city <> '' 
AND user_id = $2 AND project_id = $1
GROUP BY city ORDER BY COUNT(city) DESC LIMIT 5;

-- name: GetEventMostHitElement :many
SELECT element_path AS name, COUNT(element_path) AS total
FROM events 
WHERE element_path IS NOT NULL AND element_path <> '' 
AND user_id = $2 AND project_id = $1
GROUP BY element_path ORDER BY COUNT(element_path) DESC LIMIT 5;

-- name: GetEventLastVisitedUser :many
SELECT ip_addr AS ip, received_at AS timestamp
FROM events
WHERE user_id = $2 AND project_id = $1
GROUP BY ip_addr, received_at ORDER BY received_at DESC LIMIT 5;

-- name: GetEventMostUsedBrowser :many
SELECT browser_name AS name, COUNT(browser_name) AS total
FROM events
WHERE browser_name IS NOT NULL AND browser_name <> ''
AND user_id = $2 AND project_id = $1
GROUP BY browser_name ORDER BY COUNT(browser_name) DESC LIMIT 5;

-- name: GetEventMostEventType :many
SELECT event_type AS name, COUNT(event_type) AS total
FROM events
WHERE event_type IS NOT NULL AND event_type <> ''
AND user_id = $2 AND project_id = $1
GROUP BY event_type ORDER BY COUNT(event_type) DESC LIMIT 5;

-- name: GetEventMostEventLabel :many
SELECT event_label AS name, COUNT(event_label) AS total
FROM events
WHERE event_label IS NOT NULL AND event_label <> ''
AND user_id = $2 AND project_id = $1
GROUP BY event_label ORDER BY COUNT(event_label) DESC LIMIT 5;

-- name: GetWeeklyEvents :many
SELECT
  DATE_TRUNC('day', received_at) AS timestamp,
  COUNT(*) AS total
FROM events
WHERE user_id = $2 AND project_id = $1 AND received_at >= NOW() - INTERVAL '7 days'
GROUP BY timestamp ORDER BY timestamp ASC;

-- name: GetWeeklyEventsTotal :one
SELECT
COUNT(id) AS total
FROM events WHERE received_at >= NOW() - INTERVAL '7 days'
AND user_id = $2 AND project_id = $1;
