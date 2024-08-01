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
    UNIQUE(name, version)
);

CREATE TYPE flow_instances_status AS ENUM (
    'not_started',
    'waiting',
    'running',
    'completed',
    'failed'
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
    status flow_instances_status NOT NULL DEFAULT 'not_started',
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
    UNIQUE(flow_definition_ref_id, name, version),
    CONSTRAINT fk_parent_task
        FOREIGN KEY(parent_task_id)
        REFERENCES task_definitions(id)
);


CREATE TYPE task_instances_status AS ENUM (
    'not_started',
    'waiting',
    'running',
    'completed',
    'failed'
);
CREATE TABLE task_instances (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    started_at TIMESTAMP,
    ended_at TIMESTAMP,
    error_at TIMESTAMP,
    task_definition_id INTEGER NOT NULL REFERENCES task_definitions(id),
    task_definition_ref_id UUID NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    parent_task_id INTEGER,
    version INTEGER NOT NULL DEFAULT 1,
    status task_instances_status NOT NULL DEFAULT 'not_started',
    input_data JSONB,
    output_data JSONB,
    metadata JSONB
);
COMMIT;
