ALTER TABLE devices ADD COLUMN api_key TEXT UNIQUE DEFAULT uuid_generate_v4()::text;
UPDATE devices SET api_key = uuid_generate_v4()::text WHERE api_key IS NULL;
ALTER TABLE devices ALTER COLUMN api_key SET NOT NULL;
ALTER TABLE devices ALTER COLUMN api_key DROP DEFAULT;
