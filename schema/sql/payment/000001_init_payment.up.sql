CREATE TABLE transactions (
    id bigserial PRIMARY KEY,

    order_id varchar(64),
    customer_id varchar(64) NOT NULL,

    amount money NOT NULL,

    reversed_at timestamp,
    processed_at timestamp NOT NULL DEFAULT timezone('utc', now())
);

CREATE TABLE customer_balances (
    customer_id varchar(64) PRIMARY KEY,
    balance money NOT NULL
);

CREATE TABLE event_outbox (
    id bigserial PRIMARY KEY,

    aggregateid varchar(128) NOT NULL,
    aggregatetype varchar(128) NOT NULL,
    payload bytea NOT NULL
);