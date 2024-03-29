CREATE TABLE IF NOT EXISTS DEVICES (
  id SERIAL PRIMARY KEY,
  public_id VARCHAR UNIQUE NOT NULL,
  serial_no VARCHAR UNIQUE NOT NULL,
  status VARCHAR NOT NULL,
  class VARCHAR NOT NULL,
  name VARCHAR NOT NULL,
  config JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at timestamp with time zone NOT NULL,
  updated_at timestamp with time zone NOT NULL
);