package packageinstanceauth

import (
	"context"
	"errors"
	"fmt"

	"github.com/kyma-incubator/compass/components/director/internal/domain/package/mock"
	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/kyma-incubator/compass/components/director/internal/persistence"
	"github.com/kyma-incubator/compass/components/director/pkg/graphql"
)

//go:generate mockery -name=Service -output=automock -outpkg=automock -case=underscore
type Service interface {
	Get(ctx context.Context, id string) (*model.PackageInstanceAuth, error)
	Delete(ctx context.Context, id string) error
}

//go:generate mockery -name=Converter -output=automock -outpkg=automock -case=underscore
type Converter interface {
	ToGraphQL(in *model.PackageInstanceAuth) *graphql.PackageInstanceAuth
}

type Resolver struct {
	transact persistence.Transactioner
	svc      Service
	conv     Converter
}

func NewResolver(transact persistence.Transactioner, svc Service, conv Converter) *Resolver {
	return &Resolver{
		transact: transact,
		svc:      svc,
		conv:     conv,
	}
}

var mockRequestTypeKey = "type"
var mockPackageID = "db5d3b2a-cf30-498b-9a66-29e60247c66b"

// TODO: Remove mock
func (r *Resolver) SetPackageInstanceAuthMock(ctx context.Context, authID string, in graphql.PackageInstanceAuthSetInput) (*graphql.PackageInstanceAuth, error) {
	return mock.FixPackageInstanceAuth(mockPackageID, graphql.PackageInstanceAuthStatusConditionSucceeded), nil
}

// TODO: Remove mock
func (r *Resolver) DeletePackageInstanceAuthMock(ctx context.Context, authID string) (*graphql.PackageInstanceAuth, error) {
	return mock.FixPackageInstanceAuth(mockPackageID, graphql.PackageInstanceAuthStatusConditionUnused), nil
}

// TODO: Remove mock
func (r *Resolver) RequestPackageInstanceAuthCreationMock(ctx context.Context, packageID string, in graphql.PackageInstanceAuthRequestInput) (*graphql.PackageInstanceAuth, error) {
	id := "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
	if in.Context == nil {
		return mock.FixPackageInstanceAuth(id, graphql.PackageInstanceAuthStatusConditionPending), nil
	}

	data, ok := (*in.Context).(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid context type: expected map[string]interface{}, actual %T", *in.Context)
	}

	if _, exists := data[mockRequestTypeKey]; !exists {
		return mock.FixPackageInstanceAuth(id, graphql.PackageInstanceAuthStatusConditionPending), nil
	}

	reqType, ok := data[mockRequestTypeKey].(string)
	if !ok {
		return nil, errors.New("invalid mock request type: expected string value (`success` or `error`)")
	}

	switch reqType {
	case "success":
		id = "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"
	case "error":
		id = "cccccccc-cccc-cccc-cccc-cccccccccccc"
	}

	return mock.FixPackageInstanceAuth(id, graphql.PackageInstanceAuthStatusConditionPending), nil
}

// TODO: Remove mock
func (r *Resolver) RequestPackageInstanceAuthDeletionMock(ctx context.Context, authID string) (*graphql.PackageInstanceAuth, error) {
	return mock.FixPackageInstanceAuth(mockPackageID, graphql.PackageInstanceAuthStatusConditionUnused), nil
}

func (r *Resolver) DeletePackageInstanceAuth(ctx context.Context, authID string) (*graphql.PackageInstanceAuth, error) {
	tx, err := r.transact.Begin()
	if err != nil {
		return nil, err
	}

	// TODO: Validate if client has access to given packageID

	defer r.transact.RollbackUnlessCommited(tx)
	ctx = persistence.SaveToContext(ctx, tx)

	instanceAuth, err := r.svc.Get(ctx, authID)
	if err != nil {
		return nil, err
	}

	err = r.svc.Delete(ctx, authID)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return r.conv.ToGraphQL(instanceAuth), nil
}

// TODO: Replace with real implementation
func (r *Resolver) SetPackageInstanceAuth(ctx context.Context, packageID string, authID string, in graphql.PackageInstanceAuthSetInput) (*graphql.PackageInstanceAuth, error) {
	panic("not implemented")
}

// TODO: Replace with real implementation
func (r *Resolver) RequestPackageInstanceAuthCreation(ctx context.Context, packageID string, in graphql.PackageInstanceAuthRequestInput) (*graphql.PackageInstanceAuth, error) {
	panic("not implemented")
}

// TODO: Replace with real implementation
func (r *Resolver) RequestPackageInstanceAuthDeletion(ctx context.Context, packageID string, authID string) (*graphql.PackageInstanceAuth, error) {
	panic("not implemented")
}
