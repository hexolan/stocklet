CREATE TABLE products (
    id bigserial PRIMARY KEY,

    name varchar(128) NOT NULL,
    description varchar(256) NOT NULL,
    price money NOT NULL,

    created_at timestamp NOT NULL DEFAULT timezone('utc', now()),
    updated_at timestamp
);

CREATE TABLE event_outbox (
    id bigserial PRIMARY KEY,

    aggregateid varchar(128) NOT NULL,
    aggregatetype varchar(128) NOT NULL,
    payload bytea NOT NULL
);