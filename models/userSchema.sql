CREATE TABLE users (
    userid VARCHAR(255) PRIMARY KEY UNIQUE,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    passwrd LONGBLOB
);