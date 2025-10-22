CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    display_name TEXT NOT NULL,
    password TEXT NOT NULL,
    created_at INTEGER NOT NULL DEFAULT (unixepoch()),
    updated_at INTEGER NOT NULL DEFAULT (unixepoch())
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

CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

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
