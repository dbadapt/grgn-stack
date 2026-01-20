package migrations

// GetAllMigrations returns all registered migrations in order
func GetAllMigrations() []Migration {
	return []Migration{
		Migration001InitialSchema,
		// Add new migrations here as they are created
	}
}

// RegisterAll registers all migrations with a migrator
func RegisterAll(migrator *Migrator) {
	for _, migration := range GetAllMigrations() {
		migrator.Register(migration)
	}
}

// NewMigratorWithAll creates a new migrator with all migrations registered
func NewMigratorWithAll(db Neo4jDB) *Migrator {
	migrator := NewMigrator(db)
	RegisterAll(migrator)
	return migrator
}
