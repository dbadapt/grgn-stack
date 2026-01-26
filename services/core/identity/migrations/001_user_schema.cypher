// ============================================
// Migration: core/identity/001_user_schema
// Description: Create User schema with constraints and indexes
// ============================================

// ----- CONSTRAINTS -----

// User constraints
CREATE CONSTRAINT user_id_unique IF NOT EXISTS
FOR (u:User) REQUIRE u.id IS UNIQUE;

CREATE CONSTRAINT user_email_unique IF NOT EXISTS
FOR (u:User) REQUIRE u.email IS UNIQUE;

// ----- INDEXES -----

// User indexes
CREATE INDEX user_status IF NOT EXISTS
FOR (u:User) ON (u.status);

CREATE INDEX user_email_idx IF NOT EXISTS
FOR (u:User) ON (u.email);

CREATE INDEX user_created_at IF NOT EXISTS
FOR (u:User) ON (u.createdAt);
