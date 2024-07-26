CREATE TYPE flow_definitions_status AS ENUM ('active', 'inactive');
CREATE TABLE flow_definitions (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    name TEXT NOT NULL,
    description TEXT,
    status flow_definitions_status NOT NULL,
    metadata JSONB
);

CREATE TYPE flow_instances_status AS ENUM ('not_started', 'running', 'completed', 'failed', 'stopped');
CREATE TABLE flow_instances (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    flow_definition_id INTEGER NOT NULL REFERENCES flow_definitions(id),
    status flow_instances_status NOT NULL,
    metadata JSONB
);
