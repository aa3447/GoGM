-- +goose up
CREATE TABLE IF NOT EXISTS campaigns_players (
    campaign_id INT NOT NULL,
    player_id INT NOT NULL,
    role TEXT,
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (campaign_id, player_id),
    FOREIGN KEY (campaign_id) REFERENCES campaigns(id) ON DELETE CASCADE,
    FOREIGN KEY (player_id) REFERENCES players(id) ON DELETE CASCADE
);

-- +goose down
DROP TABLE IF EXISTS campaigns_players CASCADE;