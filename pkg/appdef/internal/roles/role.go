/*
 * Copyright (c) 2024-present Sigma-Soft, Ltd.
 * @author: Nikolay Nikitin
 */

package roles

import (
	"errors"
	"iter"
	"slices"

	"github.com/voedger/voedger/pkg/appdef"
	"github.com/voedger/voedger/pkg/appdef/internal/acl"
	"github.com/voedger/voedger/pkg/appdef/internal/types"
)

// # Supports:
//   - appdef.IRole
type Role struct {
	types.Typ
	acl.WithACL
}

func NewRole(ws appdef.IWorkspace, name appdef.QName) *Role {
	r := &Role{
		Typ:     types.MakeType(ws.App(), ws, name, appdef.TypeKind_Role),
		WithACL: acl.MakeWithACL(),
	}
	types.Propagate(r)
	return r
}

func (r Role) Ancestors() iter.Seq[appdef.QName] {
	roles := appdef.QNames{}
	for rule := range r.WithACL.ACL() {
		if rule.Op(appdef.OperationKind_Inherits) {
			switch rule.Filter().Kind() {
			case appdef.FilterKind_QNames:
				for q := range rule.Filter().QNames() {
					roles.Add(q)
				}
			default:
				// complex filter
				for role := range appdef.Roles(appdef.FilterMatches(rule.Filter(), r.Workspace().Types())) {
					roles.Add(role.QName())
				}
			}
		}
	}
	return slices.Values(roles)
}

func (r *Role) grant(ops []appdef.OperationKind, flt appdef.IFilter, fields []appdef.FieldName, comment ...string) {
	acl.NewGrant(r.Workspace(), ops, flt, fields, r, comment...)
}

func (r *Role) grantAll(flt appdef.IFilter, comment ...string) {
	acl.NewGrantAll(r.Workspace(), flt, r, comment...)
}

func (r *Role) revoke(ops []appdef.OperationKind, flt appdef.IFilter, fields []appdef.FieldName, comment ...string) {
	acl.NewRevoke(r.Workspace(), ops, flt, fields, r, comment...)
}

func (r *Role) revokeAll(flt appdef.IFilter, comment ...string) {
	acl.NewRevokeAll(r.Workspace(), flt, r, comment...)
}

// Validates role.
//
// # Error if:
//   - ACL rule is not valid
func (r Role) Validate() (err error) {
	for acl := range r.ACL() {
		if acl, ok := acl.(interface{ Validate() error }); ok {
			err = errors.Join(err, acl.Validate())
		}
	}
	return err
}

// # Supports:
//   - appdef.IRoleBuilder
type RoleBuilder struct {
	types.TypeBuilder
	r *Role
}

func NewRoleBuilder(r *Role) *RoleBuilder {
	return &RoleBuilder{
		TypeBuilder: types.MakeTypeBuilder(&r.Typ),
		r:           r,
	}
}

func (rb *RoleBuilder) Grant(ops []appdef.OperationKind, flt appdef.IFilter, fields []appdef.FieldName, comment ...string) appdef.IRoleBuilder {
	rb.r.grant(ops, flt, fields, comment...)
	return rb
}

func (rb *RoleBuilder) GrantAll(flt appdef.IFilter, comment ...string) appdef.IRoleBuilder {
	rb.r.grantAll(flt, comment...)
	return rb
}

func (rb *RoleBuilder) Revoke(ops []appdef.OperationKind, flt appdef.IFilter, fields []appdef.FieldName, comment ...string) appdef.IRoleBuilder {
	rb.r.revoke(ops, flt, fields, comment...)
	return rb
}

func (rb *RoleBuilder) RevokeAll(flt appdef.IFilter, comment ...string) appdef.IRoleBuilder {
	rb.r.revokeAll(flt, comment...)
	return rb
}
