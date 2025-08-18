package model

import "time"

type Roles string

const (
	CustomerRole   Roles = "customer"
	AdminRole      Roles = "admin"
	SuperAdminRole Roles = "super_admin"
)

type User struct {
	BaseModel
	Auth0Id   string     `db:"auth0_id" json:"auth0Id"`
	Name      string     `db:"name" json:"name"`
	Email     string     `db:"email" json:"email"`
	Phone     string     `db:"phone" json:"phone"`
	Role      Roles      `db:"role" json:"role"`
	DeletedAt *time.Time `db:"deleted_at" json:"deletedAt,omitempty"`
}
