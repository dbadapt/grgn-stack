// ============================================
// Migration: core/feature/001_feature_flag_schema
// Description: Create FeatureFlag, OverrideSegment, and TargetingRule schema
//              for hierarchical feature flag system with A/B testing support
// Created: 2026-01-25 18:59:22
// ============================================

// ----- FEATURE FLAG CONSTRAINTS -----

// FeatureFlag must have unique ID
CREATE CONSTRAINT feature_flag_id_unique IF NOT EXISTS
FOR (f:FeatureFlag) REQUIRE f.id IS UNIQUE;

// FeatureFlag key must be unique (used for lookups in code)
CREATE CONSTRAINT feature_flag_key_unique IF NOT EXISTS
FOR (f:FeatureFlag) REQUIRE f.key IS UNIQUE;

// ----- FEATURE FLAG INDEXES -----

// Index for filtering by enabled status (common query)
CREATE INDEX feature_flag_is_enabled IF NOT EXISTS
FOR (f:FeatureFlag) ON (f.isEnabled);

// Index for filtering by flag type
CREATE INDEX feature_flag_type IF NOT EXISTS
FOR (f:FeatureFlag) ON (f.flagType);

// Index for timestamp-based queries
CREATE INDEX feature_flag_created_at IF NOT EXISTS
FOR (f:FeatureFlag) ON (f.createdAt);

CREATE INDEX feature_flag_updated_at IF NOT EXISTS
FOR (f:FeatureFlag) ON (f.updatedAt);

// ----- OVERRIDE SEGMENT CONSTRAINTS -----

// OverrideSegment must have unique ID
CREATE CONSTRAINT override_segment_id_unique IF NOT EXISTS
FOR (o:OverrideSegment) REQUIRE o.id IS UNIQUE;

// ----- OVERRIDE SEGMENT INDEXES -----

// CRITICAL: Composite index for efficient multi-tenant queries
// This prevents full table scans when fetching overrides for specific scopes
CREATE INDEX override_segment_scope IF NOT EXISTS
FOR (o:OverrideSegment) ON (o.scopeType, o.scopeId);

// Index for priority-based ordering
CREATE INDEX override_segment_priority IF NOT EXISTS
FOR (o:OverrideSegment) ON (o.priority);

// Index for filtering active/inactive overrides
CREATE INDEX override_segment_is_active IF NOT EXISTS
FOR (o:OverrideSegment) ON (o.isActive);

// ----- TARGETING RULE CONSTRAINTS -----

// TargetingRule must have unique ID
CREATE CONSTRAINT targeting_rule_id_unique IF NOT EXISTS
FOR (r:TargetingRule) REQUIRE r.id IS UNIQUE;

// ----- TARGETING RULE INDEXES -----

// Index for attribute-based queries (useful for impact analysis)
CREATE INDEX targeting_rule_attribute IF NOT EXISTS
FOR (r:TargetingRule) ON (r.attribute);

// ----- RELATIONSHIP INDEXES -----

// Index for efficient traversal from FeatureFlag to OverrideSegments
CREATE INDEX rel_flag_has_override IF NOT EXISTS
FOR ()-[r:HAS_OVERRIDE]->() ON (r.priority);

// Index for efficient traversal from OverrideSegment to TargetingRules
CREATE INDEX rel_segment_has_rule IF NOT EXISTS
FOR ()-[r:HAS_RULE]->() ON (r.order);
