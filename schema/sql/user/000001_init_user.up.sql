CREATE TABLE users (
    id bigserial PRIMARY KEY,

    first_name varchar(64) NOT NULL,
    last_name varchar(64) NOT NULL,

    email varchar NOT NULL UNIQUE,

    created_at timestamp NOT NULL DEFAULT timezone('utc', now()),
    updated_at timestamp
);

CREATE TABLE event_outbox (
    id bigserial PRIMARY KEY,

    aggregateid varchar(128) NOT NULL,
    aggregatetype varchar(128) NOT NULL,
    payload bytea NOT NULL
);