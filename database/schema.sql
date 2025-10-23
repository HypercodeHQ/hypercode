CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    display_name TEXT NOT NULL,
    password TEXT,
    github_user_id TEXT UNIQUE,
    created_at INTEGER NOT NULL DEFAULT (unixepoch()),
    updated_at INTEGER NOT NULL DEFAULT (unixepoch())
);

CREATE TABLE IF NOT EXISTS access_tokens (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    token_hash TEXT NOT NULL UNIQUE,
    last_used_at INTEGER,
    created_at INTEGER NOT NULL DEFAULT (unixepoch()),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS organizations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    display_name TEXT NOT NULL,
    created_at INTEGER NOT NULL DEFAULT (unixepoch()),
    updated_at INTEGER NOT NULL DEFAULT (unixepoch())
);

CREATE TABLE IF NOT EXISTS repositories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT,
    default_branch TEXT NOT NULL DEFAULT 'main',
    visibility TEXT NOT NULL CHECK(visibility IN ('public', 'private')),
    owner_user_id INTEGER,
    owner_org_id INTEGER,
    created_at INTEGER NOT NULL DEFAULT (unixepoch()),
    updated_at INTEGER NOT NULL DEFAULT (unixepoch()),
    FOREIGN KEY (owner_user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (owner_org_id) REFERENCES organizations(id) ON DELETE CASCADE,
    CHECK ((owner_user_id IS NOT NULL AND owner_org_id IS NULL) OR (owner_user_id IS NULL AND owner_org_id IS NOT NULL))
);

CREATE TABLE IF NOT EXISTS contributors (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    repository_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    role TEXT NOT NULL CHECK(role IN ('admin', 'write', 'read')),
    created_at INTEGER NOT NULL DEFAULT (unixepoch()),
    FOREIGN KEY (repository_id) REFERENCES repositories(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(repository_id, user_id)
);

CREATE TABLE IF NOT EXISTS stars (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    repository_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    created_at INTEGER NOT NULL DEFAULT (unixepoch()),
    FOREIGN KEY (repository_id) REFERENCES repositories(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(repository_id, user_id)
);

-- Tickets (similar to GitHub Issues)
CREATE TABLE IF NOT EXISTS tickets (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    repository_id INTEGER NOT NULL,
    number INTEGER NOT NULL,
    title TEXT NOT NULL,
    body TEXT,
    status TEXT NOT NULL DEFAULT 'open' CHECK(status IN ('open', 'closed')),
    author_id INTEGER NOT NULL,
    closed_at INTEGER,
    closed_by_id INTEGER,
    created_at INTEGER NOT NULL DEFAULT (unixepoch()),
    updated_at INTEGER NOT NULL DEFAULT (unixepoch()),
    FOREIGN KEY (repository_id) REFERENCES repositories(id) ON DELETE CASCADE,
    FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (closed_by_id) REFERENCES users(id) ON DELETE SET NULL,
    UNIQUE(repository_id, number)
);

-- Ticket comments
CREATE TABLE IF NOT EXISTS ticket_comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    ticket_id INTEGER NOT NULL,
    author_id INTEGER NOT NULL,
    body TEXT NOT NULL,
    created_at INTEGER NOT NULL DEFAULT (unixepoch()),
    updated_at INTEGER NOT NULL DEFAULT (unixepoch()),
    FOREIGN KEY (ticket_id) REFERENCES tickets(id) ON DELETE CASCADE,
    FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Ticket labels (defined per repository)
CREATE TABLE IF NOT EXISTS ticket_labels (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    repository_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    color TEXT NOT NULL,
    description TEXT,
    created_at INTEGER NOT NULL DEFAULT (unixepoch()),
    FOREIGN KEY (repository_id) REFERENCES repositories(id) ON DELETE CASCADE,
    UNIQUE(repository_id, name)
);

-- Ticket label assignments (many-to-many)
CREATE TABLE IF NOT EXISTS ticket_label_assignments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    ticket_id INTEGER NOT NULL,
    label_id INTEGER NOT NULL,
    created_at INTEGER NOT NULL DEFAULT (unixepoch()),
    FOREIGN KEY (ticket_id) REFERENCES tickets(id) ON DELETE CASCADE,
    FOREIGN KEY (label_id) REFERENCES ticket_labels(id) ON DELETE CASCADE,
    UNIQUE(ticket_id, label_id)
);

-- Ticket assignees (many-to-many, tickets can have multiple assignees)
CREATE TABLE IF NOT EXISTS ticket_assignees (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    ticket_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    assigned_by_id INTEGER NOT NULL,
    created_at INTEGER NOT NULL DEFAULT (unixepoch()),
    FOREIGN KEY (ticket_id) REFERENCES tickets(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (assigned_by_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(ticket_id, user_id)
);

-- Reactions on tickets and comments (emoji reactions)
CREATE TABLE IF NOT EXISTS ticket_reactions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    ticket_id INTEGER,
    comment_id INTEGER,
    user_id INTEGER NOT NULL,
    emoji TEXT NOT NULL,
    created_at INTEGER NOT NULL DEFAULT (unixepoch()),
    FOREIGN KEY (ticket_id) REFERENCES tickets(id) ON DELETE CASCADE,
    FOREIGN KEY (comment_id) REFERENCES ticket_comments(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CHECK ((ticket_id IS NOT NULL AND comment_id IS NULL) OR (ticket_id IS NULL AND comment_id IS NOT NULL)),
    UNIQUE(ticket_id, user_id, emoji),
    UNIQUE(comment_id, user_id, emoji)
);

CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

CREATE INDEX IF NOT EXISTS idx_access_tokens_user ON access_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_access_tokens_token_hash ON access_tokens(token_hash);

CREATE INDEX IF NOT EXISTS idx_organizations_username ON organizations(username);

CREATE INDEX IF NOT EXISTS idx_repositories_owner_user ON repositories(owner_user_id);
CREATE INDEX IF NOT EXISTS idx_repositories_owner_org ON repositories(owner_org_id);
CREATE INDEX IF NOT EXISTS idx_repositories_name ON repositories(name);
CREATE INDEX IF NOT EXISTS idx_repositories_visibility ON repositories(visibility);

CREATE UNIQUE INDEX IF NOT EXISTS idx_repositories_unique_name_user ON repositories(owner_user_id, name) WHERE owner_user_id IS NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS idx_repositories_unique_name_org ON repositories(owner_org_id, name) WHERE owner_org_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_contributors_repository ON contributors(repository_id);
CREATE INDEX IF NOT EXISTS idx_contributors_user ON contributors(user_id);
CREATE INDEX IF NOT EXISTS idx_contributors_role ON contributors(role);

CREATE INDEX IF NOT EXISTS idx_stars_repository ON stars(repository_id);
CREATE INDEX IF NOT EXISTS idx_stars_user ON stars(user_id);

-- Ticket indexes
CREATE INDEX IF NOT EXISTS idx_tickets_repository ON tickets(repository_id);
CREATE INDEX IF NOT EXISTS idx_tickets_author ON tickets(author_id);
CREATE INDEX IF NOT EXISTS idx_tickets_status ON tickets(status);
CREATE INDEX IF NOT EXISTS idx_tickets_closed_by ON tickets(closed_by_id);
CREATE INDEX IF NOT EXISTS idx_tickets_created_at ON tickets(created_at);
CREATE INDEX IF NOT EXISTS idx_tickets_updated_at ON tickets(updated_at);

CREATE INDEX IF NOT EXISTS idx_ticket_comments_ticket ON ticket_comments(ticket_id);
CREATE INDEX IF NOT EXISTS idx_ticket_comments_author ON ticket_comments(author_id);
CREATE INDEX IF NOT EXISTS idx_ticket_comments_created_at ON ticket_comments(created_at);

CREATE INDEX IF NOT EXISTS idx_ticket_labels_repository ON ticket_labels(repository_id);

CREATE INDEX IF NOT EXISTS idx_ticket_label_assignments_ticket ON ticket_label_assignments(ticket_id);
CREATE INDEX IF NOT EXISTS idx_ticket_label_assignments_label ON ticket_label_assignments(label_id);

CREATE INDEX IF NOT EXISTS idx_ticket_assignees_ticket ON ticket_assignees(ticket_id);
CREATE INDEX IF NOT EXISTS idx_ticket_assignees_user ON ticket_assignees(user_id);

CREATE INDEX IF NOT EXISTS idx_ticket_reactions_ticket ON ticket_reactions(ticket_id);
CREATE INDEX IF NOT EXISTS idx_ticket_reactions_comment ON ticket_reactions(comment_id);
CREATE INDEX IF NOT EXISTS idx_ticket_reactions_user ON ticket_reactions(user_id);
CREATE INDEX IF NOT EXISTS idx_ticket_reactions_emoji ON ticket_reactions(emoji);

CREATE TRIGGER IF NOT EXISTS update_users_timestamp
AFTER UPDATE ON users
BEGIN
    UPDATE users SET updated_at = unixepoch() WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS update_organizations_timestamp
AFTER UPDATE ON organizations
BEGIN
    UPDATE organizations SET updated_at = unixepoch() WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS update_repositories_timestamp
AFTER UPDATE ON repositories
BEGIN
    UPDATE repositories SET updated_at = unixepoch() WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS update_tickets_timestamp
AFTER UPDATE ON tickets
BEGIN
    UPDATE tickets SET updated_at = unixepoch() WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS update_ticket_comments_timestamp
AFTER UPDATE ON ticket_comments
BEGIN
    UPDATE ticket_comments SET updated_at = unixepoch() WHERE id = NEW.id;
END;

-- Trigger to auto-increment ticket numbers per repository
CREATE TRIGGER IF NOT EXISTS tickets_auto_number
BEFORE INSERT ON tickets
WHEN NEW.number IS NULL OR NEW.number = 0
BEGIN
    SELECT RAISE(FAIL, 'Ticket number must be set explicitly')
    WHERE NEW.number IS NULL OR NEW.number = 0;
END;
