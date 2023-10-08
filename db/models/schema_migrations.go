package models

//go:generate reform

// SchemaMigrations represents a row in schema_migrations table.
//
//reform:schema_migrations
type SchemaMigrations struct {
	Version int64 `reform:"version,pk"`
	Dirty   bool  `reform:"dirty"`
}
