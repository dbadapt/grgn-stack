// ============================================
// Migration: core/tenant/001_tenant_schema
// Description: Create Tenant and Membership schema
// ============================================

// ----- TENANT CONSTRAINTS -----

CREATE CONSTRAINT tenant_id_unique IF NOT EXISTS
FOR (t:Tenant) REQUIRE t.id IS UNIQUE;

CREATE CONSTRAINT tenant_slug_unique IF NOT EXISTS
FOR (t:Tenant) REQUIRE t.slug IS UNIQUE;

// ----- TENANT INDEXES -----

CREATE INDEX tenant_status IF NOT EXISTS
FOR (t:Tenant) ON (t.status);

CREATE INDEX tenant_plan IF NOT EXISTS
FOR (t:Tenant) ON (t.plan);

CREATE INDEX tenant_created_at IF NOT EXISTS
FOR (t:Tenant) ON (t.createdAt);

// ----- MEMBERSHIP CONSTRAINTS -----

CREATE CONSTRAINT membership_id_unique IF NOT EXISTS
FOR (m:Membership) REQUIRE m.id IS UNIQUE;

// ----- MEMBERSHIP INDEXES -----

CREATE INDEX membership_role IF NOT EXISTS
FOR (m:Membership) ON (m.role);

CREATE INDEX membership_joined_at IF NOT EXISTS
FOR (m:Membership) ON (m.joinedAt);
