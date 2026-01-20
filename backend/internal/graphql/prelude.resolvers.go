package graphql

import (
	"context"

	"github.com/99designs/gqlgen/graphql/introspection"
)

type introspectionInputValueResolver struct{ *Resolver }

func (r *Resolver) __InputValue() __InputValueResolver {
	return &introspectionInputValueResolver{r}
}

func (r *introspectionInputValueResolver) IsDeprecated(ctx context.Context, obj *introspection.InputValue) (bool, error) {
	return false, nil
}

func (r *introspectionInputValueResolver) DeprecationReason(ctx context.Context, obj *introspection.InputValue) (*string, error) {
	return nil, nil
}

type introspectionTypeResolver struct{ *Resolver }

func (r *Resolver) __Type() __TypeResolver {
	return &introspectionTypeResolver{r}
}

func (r *introspectionTypeResolver) IsOneOf(ctx context.Context, obj *introspection.Type) (*bool, error) {
	return nil, nil
}
