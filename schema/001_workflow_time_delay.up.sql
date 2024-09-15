-- workflow and time delay schema

CREATE TABLE IF NOT EXISTS workflows (
    id          varchar(36) PRIMARY KEY,
    name        varchar(255) NOT NULL,
    description varchar(255),
    status      smallint default 0, -- 0: inactive, 1: active
    created_at  timestamp default current_timestamp,
    updated_at  timestamp default current_timestamp
);

create index workflow_name_index on workflows(name);

CREATE TABLE IF NOT EXISTS tasks (
    id          SERIAL PRIMARY KEY,
    workflow_id varchar(36) NOT NULL REFERENCES workflows(id),
    name        varchar(255) NOT NULL,
    message     text NOT NULL, -- format message
    is_head_task boolean default false, -- identify head task
    job_type    varchar(20) NOT NULL, -- period, day_of_week, specific_day, specific_time
    job_time_value   varchar(255) NOT NULL, -- ex: 5m, monday, 2024-10-01, etc.
    next_task_id bigint,
    status      smallint default 0, -- 0: inactive, 1: active
    created_at  timestamp default current_timestamp,
    updated_at  timestamp default current_timestamp
);

create index task_workflow_id_index on tasks(workflow_id);

CREATE TABLE IF NOT EXISTS schedulers (
    id          SERIAL PRIMARY KEY,
    task_id     bigint NOT NULL REFERENCES tasks(id),
    run_time    timestamp not null, -- time to run task
    partition_id smallint default 0, -- partition for scheduler job to run
    status      smallint default 0, -- 0: inactive, 1: active, 2: completed, 3: failed
    created_at  timestamp default current_timestamp,
    updated_at  timestamp default current_timestamp
);

create index scheduler_task_id_index on schedulers(task_id);
create index scheduler_task_status_index on schedulers(partition_id, status);