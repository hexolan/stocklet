CREATE TABLE shipments (
    id bigserial PRIMARY KEY,
    order_id varchar(64) NOT NULL,

    dispatched boolean DEFAULT FALSE,

    created_at timestamp NOT NULL DEFAULT timezone('utc', now())
);

CREATE TABLE shipment_items (
    shipment_id bigserial,
    product_id varchar(64),

    quantity integer NOT NULL,

    PRIMARY KEY (shipment_id, product_id),
    FOREIGN KEY (shipment_id) REFERENCES shipments (id) ON DELETE CASCADE
);

CREATE TABLE event_outbox (
    id bigserial PRIMARY KEY,

    aggregateid varchar(128) NOT NULL,
    aggregatetype varchar(128) NOT NULL,
    payload bytea NOT NULL
);