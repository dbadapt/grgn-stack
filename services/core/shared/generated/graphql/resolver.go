package graphql

import (
	identitySvc "github.com/yourusername/grgn-stack/services/core/identity/service"
	tenantSvc "github.com/yourusername/grgn-stack/services/core/tenant/service"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

// Resolver is the root resolver with service dependencies.
type Resolver struct {
	UserService   identitySvc.IUserService
	TenantService tenantSvc.ITenantService
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}
