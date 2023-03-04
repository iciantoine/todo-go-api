CREATE TABLE todo (
    id         UUID NOT NULL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL,
    is_done    BOOLEAN NOT NULL,
    message    TEXT NOT NULL
);

---- create above / drop below ----

DROP TABLE todo;