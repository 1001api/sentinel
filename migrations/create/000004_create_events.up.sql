CREATE TABLE IF NOT EXISTS events (
    id SERIAL PRIMARY KEY,
    event_type VARCHAR(100) NOT NULL,
    page_url VARCHAR(255) NOT NULL,
    element_path VARCHAR(100),
    element_type VARCHAR(100),
    ip_addr INET NOT NULL,
    user_agent VARCHAR(100) NOT NULL,
    country VARCHAR(100),
    region VARCHAR(100),
    city VARCHAR(100),

    user_id UUID,
    project_id INTEGER,
    fired_at TIMESTAMPTZ NOT NULL,

    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY(project_id) REFERENCES projects(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_events_user_id ON events(user_id);
CREATE INDEX IF NOT EXISTS idx_events_project_id ON events(project_id);
CREATE INDEX IF NOT EXISTS idx_events_fired_at ON events(fired_at);
