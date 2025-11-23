CREATE TABLE IF NOT EXISTS teams (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    team_id INTEGER REFERENCES teams(id),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS pull_requests (
    id SERIAL PRIMARY KEY,
    title VARCHAR(500) NOT NULL,
    author_id INTEGER REFERENCES users(id),
    status VARCHAR(50) DEFAULT 'OPEN',
    reviewers INTEGER[] DEFAULT ARRAY[]::INTEGER[],
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_users_team_active ON users(team_id, is_active);
CREATE INDEX IF NOT EXISTS idx_pr_author ON pull_requests(author_id);
CREATE INDEX IF NOT EXISTS idx_pr_status ON pull_requests(status);