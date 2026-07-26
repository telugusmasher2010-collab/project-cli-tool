-- {{project_name}} SQLite Schema
-- Stores bot state, messages, and contacts.

CREATE TABLE IF NOT EXISTS messages (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    remote_jid  TEXT    NOT NULL,
    message_id  TEXT    NOT NULL UNIQUE,
    from_me      INTEGER NOT NULL DEFAULT 0,
    message_type TEXT    NOT NULL DEFAULT 'text',
    content      TEXT    NOT NULL DEFAULT '',
    timestamp    INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
    created_at   TEXT    NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_messages_remote_jid ON messages(remote_jid);
CREATE INDEX IF NOT EXISTS idx_messages_timestamp ON messages(timestamp);

CREATE TABLE IF NOT EXISTS contacts (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    jid         TEXT    NOT NULL UNIQUE,
    name        TEXT    NOT NULL DEFAULT '',
    phone       TEXT    NOT NULL DEFAULT '',
    is_group    INTEGER NOT NULL DEFAULT 0,
    updated_at  TEXT    NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_contacts_jid ON contacts(jid);

CREATE TABLE IF NOT EXISTS bot_state (
    key   TEXT PRIMARY KEY,
    value TEXT NOT NULL
);

-- Insert default bot state
INSERT OR IGNORE INTO bot_state (key, value) VALUES ('started_at', datetime('now'));
