package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/yourusername/grgn-stack/pkg/errors"
	shared "github.com/yourusername/grgn-stack/services/core/shared/controller"
	"github.com/yourusername/grgn-stack/services/core/shared/generated/graphql/model"
)

// TenantRepository implements ITenantRepository using Neo4j.
type TenantRepository struct {
	db shared.IDatabase
}

// NewTenantRepository creates a new TenantRepository.
func NewTenantRepository(db shared.IDatabase) *TenantRepository {
	return &TenantRepository{db: db}
}

// FindByID retrieves a tenant by their unique ID.
func (r *TenantRepository) FindByID(ctx context.Context, id string) (*model.Tenant, error) {
	result, err := r.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			MATCH (t:Tenant {id: $id})
			WHERE t.status <> 'DELETED'
			OPTIONAL MATCH (m:Membership)-[:IN_TENANT]->(t)
			RETURN t, count(m) as memberCount
		`, map[string]any{"id": id})
		if err != nil {
			return nil, err
		}

		record, err := result.Single(ctx)
		if err != nil {
			return nil, errors.ErrTenantNotFound
		}

		return r.mapRecordToTenant(record)
	})
	if err != nil {
		return nil, err
	}
	return result.(*model.Tenant), nil
}

// FindBySlug retrieves a tenant by their unique slug.
func (r *TenantRepository) FindBySlug(ctx context.Context, slug string) (*model.Tenant, error) {
	result, err := r.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			MATCH (t:Tenant {slug: $slug})
			WHERE t.status <> 'DELETED'
			OPTIONAL MATCH (m:Membership)-[:IN_TENANT]->(t)
			RETURN t, count(m) as memberCount
		`, map[string]any{"slug": slug})
		if err != nil {
			return nil, err
		}

		record, err := result.Single(ctx)
		if err != nil {
			return nil, errors.ErrTenantNotFound
		}

		return r.mapRecordToTenant(record)
	})
	if err != nil {
		return nil, err
	}
	return result.(*model.Tenant), nil
}

// FindByUserID retrieves all tenants a user is a member of.
func (r *TenantRepository) FindByUserID(ctx context.Context, userID string) ([]*model.Tenant, error) {
	result, err := r.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			MATCH (u:User {id: $userID})-[:HAS_MEMBERSHIP]->(m:Membership)-[:IN_TENANT]->(t:Tenant)
			WHERE t.status <> 'DELETED'
			WITH t
			OPTIONAL MATCH (m2:Membership)-[:IN_TENANT]->(t)
			RETURN t, count(m2) as memberCount
			ORDER BY t.createdAt DESC
		`, map[string]any{"userID": userID})
		if err != nil {
			return nil, err
		}

		var tenants []*model.Tenant
		for result.Next(ctx) {
			tenant, err := r.mapRecordToTenant(result.Record())
			if err != nil {
				return nil, err
			}
			tenants = append(tenants, tenant)
		}

		return tenants, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]*model.Tenant), nil
}

// Create creates a new tenant in the database.
func (r *TenantRepository) Create(ctx context.Context, tenant *model.Tenant) (*model.Tenant, error) {
	// Generate ID if not provided
	if tenant.ID == "" {
		tenant.ID = uuid.New().String()
	}

	now := time.Now()
	tenant.CreatedAt = now
	tenant.UpdatedAt = now

	// Set defaults
	if tenant.Status == "" {
		tenant.Status = model.TenantStatusActive
	}
	if tenant.Plan == "" {
		tenant.Plan = model.TenantPlanFree
	}
	if tenant.IsolationMode == "" {
		tenant.IsolationMode = model.TenantIsolationModeShared
	}

	result, err := r.db.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		// Check if slug already exists
		checkResult, err := tx.Run(ctx, `
			MATCH (t:Tenant {slug: $slug})
			WHERE t.status <> 'DELETED'
			RETURN count(t) > 0 as exists
		`, map[string]any{"slug": tenant.Slug})
		if err != nil {
			return nil, err
		}

		checkRecord, err := checkResult.Single(ctx)
		if err != nil {
			return nil, err
		}

		if exists, _ := checkRecord.Get("exists"); exists.(bool) {
			return nil, errors.ErrSlugTaken
		}

		// Create the tenant
		params := map[string]any{
			"id":            tenant.ID,
			"name":          tenant.Name,
			"slug":          tenant.Slug,
			"plan":          string(tenant.Plan),
			"isolationMode": string(tenant.IsolationMode),
			"status":        string(tenant.Status),
		}

		result, err := tx.Run(ctx, `
			CREATE (t:Tenant {
				id: $id,
				name: $name,
				slug: $slug,
				plan: $plan,
				isolationMode: $isolationMode,
				status: $status,
				createdAt: datetime(),
				updatedAt: datetime()
			})
			RETURN t, 0 as memberCount
		`, params)
		if err != nil {
			return nil, err
		}

		record, err := result.Single(ctx)
		if err != nil {
			return nil, err
		}

		return r.mapRecordToTenant(record)
	})
	if err != nil {
		return nil, err
	}
	return result.(*model.Tenant), nil
}

// Update updates an existing tenant.
func (r *TenantRepository) Update(ctx context.Context, id string, input model.UpdateTenantInput) (*model.Tenant, error) {
	result, err := r.db.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		params := map[string]any{
			"id":        id,
			"updatedAt": time.Now(),
		}

		// Build SET clause dynamically
		setClause := "t.updatedAt = datetime($updatedAt)"
		if input.Name != nil {
			params["name"] = *input.Name
			setClause += ", t.name = $name"
		}
		if input.Plan != nil {
			params["plan"] = string(*input.Plan)
			setClause += ", t.plan = $plan"
		}
		if input.Status != nil {
			params["status"] = string(*input.Status)
			setClause += ", t.status = $status"
		}

		query := `
			MATCH (t:Tenant {id: $id})
			WHERE t.status <> 'DELETED'
			SET ` + setClause + `
			WITH t
			OPTIONAL MATCH (m:Membership)-[:IN_TENANT]->(t)
			RETURN t, count(m) as memberCount
		`

		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		record, err := result.Single(ctx)
		if err != nil {
			return nil, errors.ErrTenantNotFound
		}

		return r.mapRecordToTenant(record)
	})
	if err != nil {
		return nil, err
	}
	return result.(*model.Tenant), nil
}

// Delete soft-deletes a tenant.
func (r *TenantRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			MATCH (t:Tenant {id: $id})
			WHERE t.status <> 'DELETED'
			SET t.status = 'DELETED', t.updatedAt = datetime()
			RETURN t
		`, map[string]any{"id": id})
		if err != nil {
			return nil, err
		}

		_, err = result.Single(ctx)
		if err != nil {
			return nil, errors.ErrTenantNotFound
		}

		return nil, nil
	})
	return err
}

// ExistsBySlug checks if a tenant with the given slug exists.
func (r *TenantRepository) ExistsBySlug(ctx context.Context, slug string) (bool, error) {
	result, err := r.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			MATCH (t:Tenant {slug: $slug})
			WHERE t.status <> 'DELETED'
			RETURN count(t) > 0 as exists
		`, map[string]any{"slug": slug})
		if err != nil {
			return nil, err
		}

		record, err := result.Single(ctx)
		if err != nil {
			return false, nil
		}

		exists, _ := record.Get("exists")
		return exists.(bool), nil
	})
	if err != nil {
		return false, err
	}
	return result.(bool), nil
}

// GetMemberCount returns the number of members in a tenant.
func (r *TenantRepository) GetMemberCount(ctx context.Context, tenantID string) (int, error) {
	result, err := r.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			MATCH (m:Membership)-[:IN_TENANT]->(t:Tenant {id: $tenantID})
			RETURN count(m) as count
		`, map[string]any{"tenantID": tenantID})
		if err != nil {
			return nil, err
		}

		record, err := result.Single(ctx)
		if err != nil {
			return 0, nil
		}

		count, _ := record.Get("count")
		return int(count.(int64)), nil
	})
	if err != nil {
		return 0, err
	}
	return result.(int), nil
}

// mapRecordToTenant converts a Neo4j record to a Tenant model.
func (r *TenantRepository) mapRecordToTenant(record *neo4j.Record) (*model.Tenant, error) {
	nodeVal, ok := record.Get("t")
	if !ok {
		return nil, errors.ErrTenantNotFound
	}

	node := nodeVal.(neo4j.Node)
	props := node.Props

	tenant := &model.Tenant{
		ID:            props["id"].(string),
		Name:          props["name"].(string),
		Slug:          props["slug"].(string),
		Plan:          model.TenantPlan(props["plan"].(string)),
		IsolationMode: model.TenantIsolationMode(props["isolationMode"].(string)),
		Status:        model.TenantStatus(props["status"].(string)),
	}

	if createdAt, ok := props["createdAt"]; ok {
		tenant.CreatedAt = createdAt.(time.Time)
	}

	if updatedAt, ok := props["updatedAt"]; ok {
		tenant.UpdatedAt = updatedAt.(time.Time)
	}

	// Get member count from the query result
	if memberCount, ok := record.Get("memberCount"); ok {
		tenant.MemberCount = int(memberCount.(int64))
	}

	return tenant, nil
}

// Ensure TenantRepository implements ITenantRepository
var _ ITenantRepository = (*TenantRepository)(nil)
