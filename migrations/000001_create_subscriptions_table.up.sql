CREATE TABLE subscriptions(
    id UUID PRIMARY KEY,
    service_name VARCHAR(30) NOT NULL,
    price integer NOT NULL,
    user_id UUID NOT NULL,
    start_date  TIMESTAMPTZ NOT NULL
    end_date TIMESTAMPTZ DEFAULT NULL
);
CREATE INDEX idx_subscriptions_user_id ON subscriptions(user_id);