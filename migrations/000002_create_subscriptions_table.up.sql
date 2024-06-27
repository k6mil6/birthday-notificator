CREATE TABLE IF NOT EXISTS subscriptions (
  id UUID PRIMARY KEY,
  user_id UUID REFERENCES users(id) NOT NULL,
  subscribed_at_user_id UUID REFERENCES users(id) NOT NULL,
  notification_date TIMESTAMPTZ NOT NULL,
  UNIQUE (user_id, subscribed_at_user_id)
);