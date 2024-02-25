CREATE TABLE IF NOT EXISTS DEVICES (
  id SERIAL PRIMARY KEY,
  public_id VARCHAR UNIQUE NOT NULL,
  status VARCHAR NOT NULL,
  amount NUMERIC(10, 2) NOT NULL,
  created_at timestamp with time zone NOT NULL,
  updated_at timestamp with time zone NOT NULL
);