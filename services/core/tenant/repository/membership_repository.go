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

// MembershipRepository implements IMembershipRepository using Neo4j.
type MembershipRepository struct {
	db shared.IDatabase
}

// NewMembershipRepository creates a new MembershipRepository.
func NewMembershipRepository(db shared.IDatabase) *MembershipRepository {
	return &MembershipRepository{db: db}
}

// FindByID retrieves a membership by its unique ID.
func (r *MembershipRepository) FindByID(ctx context.Context, id string) (*model.Membership, error) {
	result, err := r.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			MATCH (u:User)-[:HAS_MEMBERSHIP]->(m:Membership {id: $id})-[:IN_TENANT]->(t:Tenant)
			OPTIONAL MATCH (inviter:User)-[:INVITED]->(m)
			RETURN m, u, t, inviter
		`, map[string]any{"id": id})
		if err != nil {
			return nil, err
		}

		record, err := result.Single(ctx)
		if err != nil {
			return nil, errors.ErrMembershipNotFound
		}

		return r.mapRecordToMembership(record)
	})
	if err != nil {
		return nil, err
	}
	return result.(*model.Membership), nil
}

// FindByTenantID retrieves all memberships for a tenant.
func (r *MembershipRepository) FindByTenantID(ctx context.Context, tenantID string) ([]*model.Membership, error) {
	result, err := r.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			MATCH (u:User)-[:HAS_MEMBERSHIP]->(m:Membership)-[:IN_TENANT]->(t:Tenant {id: $tenantID})
			WHERE u.status <> 'DELETED'
			OPTIONAL MATCH (inviter:User)-[:INVITED]->(m)
			RETURN m, u, t, inviter
			ORDER BY m.joinedAt DESC
		`, map[string]any{"tenantID": tenantID})
		if err != nil {
			return nil, err
		}

		var memberships []*model.Membership
		for result.Next(ctx) {
			membership, err := r.mapRecordToMembership(result.Record())
			if err != nil {
				return nil, err
			}
			memberships = append(memberships, membership)
		}

		return memberships, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]*model.Membership), nil
}

// FindByUserID retrieves all memberships for a user.
func (r *MembershipRepository) FindByUserID(ctx context.Context, userID string) ([]*model.Membership, error) {
	result, err := r.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			MATCH (u:User {id: $userID})-[:HAS_MEMBERSHIP]->(m:Membership)-[:IN_TENANT]->(t:Tenant)
			WHERE t.status <> 'DELETED'
			OPTIONAL MATCH (inviter:User)-[:INVITED]->(m)
			RETURN m, u, t, inviter
			ORDER BY m.joinedAt DESC
		`, map[string]any{"userID": userID})
		if err != nil {
			return nil, err
		}

		var memberships []*model.Membership
		for result.Next(ctx) {
			membership, err := r.mapRecordToMembership(result.Record())
			if err != nil {
				return nil, err
			}
			memberships = append(memberships, membership)
		}

		return memberships, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]*model.Membership), nil
}

// FindByUserAndTenant retrieves a membership by user and tenant.
func (r *MembershipRepository) FindByUserAndTenant(ctx context.Context, userID, tenantID string) (*model.Membership, error) {
	result, err := r.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			MATCH (u:User {id: $userID})-[:HAS_MEMBERSHIP]->(m:Membership)-[:IN_TENANT]->(t:Tenant {id: $tenantID})
			OPTIONAL MATCH (inviter:User)-[:INVITED]->(m)
			RETURN m, u, t, inviter
		`, map[string]any{"userID": userID, "tenantID": tenantID})
		if err != nil {
			return nil, err
		}

		record, err := result.Single(ctx)
		if err != nil {
			return nil, errors.ErrMembershipNotFound
		}

		return r.mapRecordToMembership(record)
	})
	if err != nil {
		return nil, err
	}
	return result.(*model.Membership), nil
}

// Create creates a new membership.
func (r *MembershipRepository) Create(ctx context.Context, userID, tenantID string, role model.MembershipRole, invitedByID *string) (*model.Membership, error) {
	membershipID := uuid.New().String()

	result, err := r.db.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		// Check if user is already a member
		checkResult, err := tx.Run(ctx, `
			MATCH (u:User {id: $userID})-[:HAS_MEMBERSHIP]->(m:Membership)-[:IN_TENANT]->(t:Tenant {id: $tenantID})
			RETURN count(m) > 0 as exists
		`, map[string]any{"userID": userID, "tenantID": tenantID})
		if err != nil {
			return nil, err
		}

		checkRecord, err := checkResult.Single(ctx)
		if err != nil {
			return nil, err
		}

		if exists, _ := checkRecord.Get("exists"); exists.(bool) {
			return nil, errors.ErrAlreadyMember
		}

		// Create the membership
		params := map[string]any{
			"membershipID": membershipID,
			"userID":       userID,
			"tenantID":     tenantID,
			"role":         string(role),
		}

		query := `
			MATCH (u:User {id: $userID}), (t:Tenant {id: $tenantID})
			CREATE (m:Membership {id: $membershipID, role: $role, joinedAt: datetime()})
			CREATE (u)-[:HAS_MEMBERSHIP]->(m)-[:IN_TENANT]->(t)
			RETURN m, u, t
		`

		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		record, err := result.Single(ctx)
		if err != nil {
			return nil, err
		}

		// If there's an inviter, create the INVITED relationship
		if invitedByID != nil && *invitedByID != "" {
			_, err = tx.Run(ctx, `
				MATCH (inviter:User {id: $inviterID}), (m:Membership {id: $membershipID})
				CREATE (inviter)-[:INVITED]->(m)
			`, map[string]any{"inviterID": *invitedByID, "membershipID": membershipID})
			if err != nil {
				return nil, err
			}
		}

		return r.mapRecordToMembershipBasic(record)
	})
	if err != nil {
		return nil, err
	}
	return result.(*model.Membership), nil
}

// UpdateRole updates a membership's role.
func (r *MembershipRepository) UpdateRole(ctx context.Context, id string, role model.MembershipRole) (*model.Membership, error) {
	result, err := r.db.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			MATCH (u:User)-[:HAS_MEMBERSHIP]->(m:Membership {id: $id})-[:IN_TENANT]->(t:Tenant)
			SET m.role = $role
			OPTIONAL MATCH (inviter:User)-[:INVITED]->(m)
			RETURN m, u, t, inviter
		`, map[string]any{"id": id, "role": string(role)})
		if err != nil {
			return nil, err
		}

		record, err := result.Single(ctx)
		if err != nil {
			return nil, errors.ErrMembershipNotFound
		}

		return r.mapRecordToMembership(record)
	})
	if err != nil {
		return nil, err
	}
	return result.(*model.Membership), nil
}

// Delete removes a membership.
func (r *MembershipRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			MATCH (u:User)-[:HAS_MEMBERSHIP]->(m:Membership {id: $id})-[:IN_TENANT]->(t:Tenant)
			DETACH DELETE m
			RETURN u.id as userId
		`, map[string]any{"id": id})
		if err != nil {
			return nil, err
		}

		_, err = result.Single(ctx)
		if err != nil {
			return nil, errors.ErrMembershipNotFound
		}

		return nil, nil
	})
	return err
}

// CountOwners returns the number of owners in a tenant.
func (r *MembershipRepository) CountOwners(ctx context.Context, tenantID string) (int, error) {
	result, err := r.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			MATCH (m:Membership {role: 'OWNER'})-[:IN_TENANT]->(t:Tenant {id: $tenantID})
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

// GetTenantIDByMembershipID returns the tenant ID for a membership.
func (r *MembershipRepository) GetTenantIDByMembershipID(ctx context.Context, membershipID string) (string, error) {
	result, err := r.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			MATCH (m:Membership {id: $id})-[:IN_TENANT]->(t:Tenant)
			RETURN t.id as tenantID
		`, map[string]any{"id": membershipID})
		if err != nil {
			return nil, err
		}

		record, err := result.Single(ctx)
		if err != nil {
			return "", errors.ErrMembershipNotFound
		}

		tenantID, _ := record.Get("tenantID")
		return tenantID.(string), nil
	})
	if err != nil {
		return "", err
	}
	return result.(string), nil
}

// GetUserIDByMembershipID returns the user ID for a membership.
func (r *MembershipRepository) GetUserIDByMembershipID(ctx context.Context, membershipID string) (string, error) {
	result, err := r.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			MATCH (u:User)-[:HAS_MEMBERSHIP]->(m:Membership {id: $id})
			RETURN u.id as userID
		`, map[string]any{"id": membershipID})
		if err != nil {
			return nil, err
		}

		record, err := result.Single(ctx)
		if err != nil {
			return "", errors.ErrMembershipNotFound
		}

		userID, _ := record.Get("userID")
		return userID.(string), nil
	})
	if err != nil {
		return "", err
	}
	return result.(string), nil
}

// mapRecordToMembership converts a Neo4j record to a Membership model.
func (r *MembershipRepository) mapRecordToMembership(record *neo4j.Record) (*model.Membership, error) {
	mVal, ok := record.Get("m")
	if !ok {
		return nil, errors.ErrMembershipNotFound
	}

	mNode := mVal.(neo4j.Node)
	mProps := mNode.Props

	membership := &model.Membership{
		ID:   mProps["id"].(string),
		Role: model.MembershipRole(mProps["role"].(string)),
	}

	if joinedAt, ok := mProps["joinedAt"]; ok {
		membership.JoinedAt = joinedAt.(time.Time)
	}

	// Map user
	if uVal, ok := record.Get("u"); ok && uVal != nil {
		uNode := uVal.(neo4j.Node)
		uProps := uNode.Props
		membership.User = &model.User{
			ID:     uProps["id"].(string),
			Email:  uProps["email"].(string),
			Status: model.UserStatus(uProps["status"].(string)),
		}
		if name, ok := uProps["name"]; ok && name != nil {
			nameStr := name.(string)
			membership.User.Name = &nameStr
		}
		if avatarURL, ok := uProps["avatarUrl"]; ok && avatarURL != nil {
			avatarStr := avatarURL.(string)
			membership.User.AvatarURL = &avatarStr
		}
	}

	// Map tenant
	if tVal, ok := record.Get("t"); ok && tVal != nil {
		tNode := tVal.(neo4j.Node)
		tProps := tNode.Props
		membership.Tenant = &model.Tenant{
			ID:            tProps["id"].(string),
			Name:          tProps["name"].(string),
			Slug:          tProps["slug"].(string),
			Plan:          model.TenantPlan(tProps["plan"].(string)),
			IsolationMode: model.TenantIsolationMode(tProps["isolationMode"].(string)),
			Status:        model.TenantStatus(tProps["status"].(string)),
		}
	}

	// Map inviter (optional)
	if inviterVal, ok := record.Get("inviter"); ok && inviterVal != nil {
		inviterNode := inviterVal.(neo4j.Node)
		inviterProps := inviterNode.Props
		membership.InvitedBy = &model.User{
			ID:     inviterProps["id"].(string),
			Email:  inviterProps["email"].(string),
			Status: model.UserStatus(inviterProps["status"].(string)),
		}
		if name, ok := inviterProps["name"]; ok && name != nil {
			nameStr := name.(string)
			membership.InvitedBy.Name = &nameStr
		}
	}

	return membership, nil
}

// mapRecordToMembershipBasic maps a record without inviter info.
func (r *MembershipRepository) mapRecordToMembershipBasic(record *neo4j.Record) (*model.Membership, error) {
	mVal, ok := record.Get("m")
	if !ok {
		return nil, errors.ErrMembershipNotFound
	}

	mNode := mVal.(neo4j.Node)
	mProps := mNode.Props

	membership := &model.Membership{
		ID:   mProps["id"].(string),
		Role: model.MembershipRole(mProps["role"].(string)),
	}

	if joinedAt, ok := mProps["joinedAt"]; ok {
		membership.JoinedAt = joinedAt.(time.Time)
	}

	// Map user
	if uVal, ok := record.Get("u"); ok && uVal != nil {
		uNode := uVal.(neo4j.Node)
		uProps := uNode.Props
		membership.User = &model.User{
			ID:     uProps["id"].(string),
			Email:  uProps["email"].(string),
			Status: model.UserStatus(uProps["status"].(string)),
		}
		if name, ok := uProps["name"]; ok && name != nil {
			nameStr := name.(string)
			membership.User.Name = &nameStr
		}
	}

	// Map tenant
	if tVal, ok := record.Get("t"); ok && tVal != nil {
		tNode := tVal.(neo4j.Node)
		tProps := tNode.Props
		membership.Tenant = &model.Tenant{
			ID:            tProps["id"].(string),
			Name:          tProps["name"].(string),
			Slug:          tProps["slug"].(string),
			Plan:          model.TenantPlan(tProps["plan"].(string)),
			IsolationMode: model.TenantIsolationMode(tProps["isolationMode"].(string)),
			Status:        model.TenantStatus(tProps["status"].(string)),
		}
	}

	return membership, nil
}

// Ensure MembershipRepository implements IMembershipRepository
var _ IMembershipRepository = (*MembershipRepository)(nil)
