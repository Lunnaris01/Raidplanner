-- +goose Up
CREATE TABLE characters (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    charname TEXT NOT NULL UNIQUE,
    charclass TEXT NOT NULL,
    supporter BOOLEAN NOT NULL default false,
    roaster_id INTEGER,
    FOREIGN KEY (roaster_id) REFERENCES roasters(id)  -- Reference to the raosters table
);

-- +goose Down
DROP TABLE characters;