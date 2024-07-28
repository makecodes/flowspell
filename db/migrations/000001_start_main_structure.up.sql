BEGIN;
CREATE TYPE flow_definitions_status AS ENUM ('active', 'inactive');
CREATE TABLE flow_definitions (
    id SERIAL PRIMARY KEY,
    reference_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    name TEXT NOT NULL,
    description TEXT,
    status flow_definitions_status NOT NULL DEFAULT 'inactive',
    version INTEGER NOT NULL DEFAULT 1,
    input_schema JSONB,
    output_schema JSONB,
    metadata JSONB,
    UNIQUE(reference_id, name, version)
);

CREATE TYPE flow_instances_status AS ENUM (
    'waiting',
    'running',
    'completed',
    'failed',
    'stopped'
);
CREATE TABLE flow_instances (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    started_at TIMESTAMP,
    ended_at TIMESTAMP,
    error_at TIMESTAMP,
    flow_definition_id INTEGER NOT NULL REFERENCES flow_definitions(id),
    flow_definition_ref_id UUID NOT NULL,
    status flow_instances_status NOT NULL DEFAULT 'waiting',
    version INTEGER NOT NULL DEFAULT 1,
    input_data JSONB,
    output_data JSONB,
    metadata JSONB
);


CREATE TABLE task_definitions (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    reference_id UUID NOT NULL,
    flow_definition_id INTEGER NOT NULL REFERENCES flow_definitions(id),
    flow_definition_ref_id UUID NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    parent_task_id INTEGER,
    input_schema JSONB,
    output_schema JSONB,
    version INTEGER NOT NULL DEFAULT 1,
    metadata JSONB,
    CONSTRAINT fk_parent_task
        FOREIGN KEY(parent_task_id)
        REFERENCES task_definitions(id)
);

COMMIT;
