package persistence

func schemaStatements() []string {
	return []string{
		`CREATE TABLE IF NOT EXISTS phrases (
			id TEXT PRIMARY KEY,
			text TEXT NOT NULL,
			mode TEXT NOT NULL,
			priority INTEGER NOT NULL,
			enabled INTEGER NOT NULL,
			created_at TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS sessions (
			id TEXT PRIMARY KEY,
			visitor_name TEXT NOT NULL,
			mode TEXT NOT NULL,
			started_at TEXT NOT NULL,
			ended_at TEXT NOT NULL,
			active INTEGER NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS gesture_events (
			id TEXT PRIMARY KEY,
			session_id TEXT NOT NULL,
			kind TEXT NOT NULL,
			strength INTEGER NOT NULL,
			sequence_no INTEGER NOT NULL,
			accepted INTEGER NOT NULL,
			created_at TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS display_states (
			session_id TEXT PRIMARY KEY,
			mode TEXT NOT NULL,
			phrase_id TEXT NOT NULL,
			phrase_text TEXT NOT NULL,
			particle_form TEXT NOT NULL,
			revision INTEGER NOT NULL,
			updated_at TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS audit_entries (
			id TEXT PRIMARY KEY,
			session_id TEXT NOT NULL,
			action TEXT NOT NULL,
			detail TEXT NOT NULL,
			created_at TEXT NOT NULL
		)`,
	}
}
