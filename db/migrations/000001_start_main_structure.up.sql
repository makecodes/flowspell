CREATE TYPE flow_definition_status AS ENUM ('active', 'inactive');
CREATE TABLE flow_definition (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    name TEXT NOT NULL,
    description TEXT,
    status flow_definition_status NOT NULL,
    metadata JSONB
);

CREATE TYPE flow_instance_status AS ENUM ('not_started', 'running', 'completed', 'failed', 'stopped');
CREATE TABLE flow_instance (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    flow_definition_id INTEGER NOT NULL REFERENCES flow_definition(id),
    status flow_instance_status NOT NULL,
    metadata JSONB
);
