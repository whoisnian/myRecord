CREATE TABLE IF NOT EXISTS items (
    id bigserial NOT NULL PRIMARY KEY,
    type integer NOT NULL,
    state integer NOT NULL,
    content text NOT NULL,
    date date NOT NULL,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS index_items_on_date ON items (date DESC);
