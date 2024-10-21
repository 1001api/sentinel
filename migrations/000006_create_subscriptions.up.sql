CREATE TYPE billing_interval AS ENUM('monthly', 'annually');
CREATE TYPE billing_status AS ENUM('success', 'pending', 'cancelled', 'failed', 'other');

CREATE TABLE IF NOT EXISTS subscriptions (
    id SERIAL PRIMARY KEY,
    started_at TIMESTAMPTZ NOT NULL,
    next_billing_date TIMESTAMPTZ NOT NULL,
    interval billing_interval NOT NULL,
    currency TEXT NOT NULL,

    payment_id TEXT NOT NULL,
    payment_gateway TEXT NOT NULL,
    payment_status billing_status NOT NULL,
    
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    cancelled_at TIMESTAMPTZ,

    plan_id INT NOT NULL,
    user_id UUID NOT NULL,
    
    FOREIGN KEY(plan_id) REFERENCES plans(id),
    FOREIGN KEY(user_id) REFERENCES users(id)
);

CREATE INDEX idx_subs_plan_id ON subscriptions(plan_id);
CREATE INDEX idx_subs_user_id ON subscriptions(user_id);
