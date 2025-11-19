-- +goose up
CREATE TABLE IF NOT EXISTS player_characters (
    id SERIAL PRIMARY KEY,
    player_id INT NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    class TEXT,
    level INT DEFAULT 1,
    experience INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (player_id) REFERENCES players(id) ON DELETE CASCADE
);

-- +goose down
DROP TABLE IF EXISTS player_characters CASCADE;