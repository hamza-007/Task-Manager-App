CREATE TABLE tasks (
    id VARCHAR(255) PRIMARY KEY UNIQUE,
    description TEXT NOT NULL,
    completed boolean,
    created_at TEXT NOT NULL,
    completed_at TEXT NOT NULL,
    userid VARCHAR(255) NOT NULL,
    FOREIGN KEY(userid) REFERENCES users(userid)
);