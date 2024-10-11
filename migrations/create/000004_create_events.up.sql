CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    event_type TEXT NOT NULL,
    event_label TEXT,
    page_url TEXT,
    element_path TEXT,
    element_type TEXT,

    ip_addr INET,
    user_agent TEXT,
    browser_name VARCHAR(100),
    country VARCHAR(100),
    region VARCHAR(100),
    city VARCHAR(100),

    session_id VARCHAR(100),
    device_type VARCHAR(100),
    time_on_page INTEGER,
    screen_resolution VARCHAR(100),
    fired_at TIMESTAMPTZ NOT NULL,

    user_id UUID NOT NULL,
    project_id INTEGER NOT NULL,

    FOREIGN KEY(user_id) REFERENCES users(id),
    FOREIGN KEY(project_id) REFERENCES projects(id)
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_events_session_id ON events(session_id);
CREATE INDEX IF NOT EXISTS idx_events_fired_at ON events(fired_at);
CREATE INDEX IF NOT EXISTS idx_events_user_id ON events(user_id);
CREATE INDEX IF NOT EXISTS idx_events_project_id ON events(project_id);
