CREATE TABLE IF NOT EXISTS project_aggregations (
    id SERIAL PRIMARY KEY,
    project_id UUID NOT NULL,
    user_id UUID NOT NULL,

    total_events INTEGER NOT NULL,
    total_event_types INTEGER NOT NULL,
    total_unique_users INTEGER NOT NULL,
    total_locations INTEGER NOT NULL,
    total_unique_page_urls INTEGER NOT NULL,

    most_visited_urls JSONB,
    most_visited_countries JSONB,
    most_visited_cities JSONB,
    most_visited_regions JSONB,
    most_firing_elements JSONB,
    last_visited_users JSONB,
    most_used_browsers JSONB,
    most_fired_event_types JSONB,
    most_fired_event_labels JSONB,

    aggregated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    aggregated_at_str VARCHAR(255) NOT NULL,

    -- store start of the week 
    aggregated_time_bucket TIMESTAMPTZ NOT NULL,

    FOREIGN KEY (project_id) REFERENCES projects(id),
    FOREIGN KEY (user_id) REFERENCES users(id),

    -- Ensure only one aggregation per time bucket
    UNIQUE (project_id, aggregated_time_bucket)
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_project_aggr_project_id ON project_aggregations (project_id);
CREATE INDEX IF NOT EXISTS idx_project_aggr_user_id ON project_aggregations (user_id);
CREATE INDEX IF NOT EXISTS idx_project_aggr_time_bucket ON project_aggregations (aggregated_time_bucket);

-- Helper function to get the start of the week for a timestamp
CREATE OR REPLACE FUNCTION get_week_start(ts TIMESTAMPTZ)
RETURNS TIMESTAMPTZ AS $$
    -- week start in monday
    SELECT date_trunc('week', ts)::TIMESTAMPTZ;
$$ LANGUAGE SQL IMMUTABLE;
