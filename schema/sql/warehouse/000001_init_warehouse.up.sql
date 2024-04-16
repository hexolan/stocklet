CREATE TABLE product_stock (
    product_id varchar(64) PRIMARY KEY,
    quantity integer NOT NULL
);

CREATE TABLE reservations (
    id bigserial PRIMARY KEY,
    
    order_id varchar(64) NOT NULL,

    created_at timestamp NOT NULL DEFAULT timezone('utc', now())
);

CREATE TABLE reservation_items (
    reservation_id bigserial,
    product_id varchar(64),

    quantity integer NOT NULL,

    PRIMARY KEY (reservation_id, product_id),
    FOREIGN KEY (reservation_id) REFERENCES reservations (id) ON DELETE CASCADE
);

CREATE TABLE event_outbox (
    id bigserial PRIMARY KEY,

    aggregateid varchar(128) NOT NULL,
    aggregatetype varchar(128) NOT NULL,
    payload bytea NOT NULL
);