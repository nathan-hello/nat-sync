CREATE TABLE rooms (
        id TEXT PRIMARY KEY NOT NULL,
        name TEXT UNIQUE NOT NULL
)

CREATE TABLE users (
        id TEXT PRIMARY KEY NOT NULL,
        username TEXT UNIQUE NOT NULL
)


