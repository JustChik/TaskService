CREATE TABLE users (
    id UUID PRIMARY KEY NOT NULL,
    username VARCHAR(50) UNIQUE,
    password VARCHAR(128),
);