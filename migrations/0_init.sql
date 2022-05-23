CREATE TABLE ufos
(
    id         uuid PRIMARY KEY,
    model      text,
    license    text,
    plates     text,
    tank       int,
    fuel       text,
    created_at timestamptz NOT NULL,
    updated_at timestamptz
);
