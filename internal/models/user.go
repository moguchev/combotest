package models

const (
	AdminRole   Role = "admin"
	AnalystRole Role = "analyst"
	Anonymous   Role = "anonymous"
)

type Role string

type User struct {
	ID        string `json:"id"        bson:"_id"`
	Role      Role   `json:"role"      bson:"role"`
	Confirmed bool   `json:"confirmed" bson:"confirmed"`
}

type CreateUser struct {
	User     `json:"user" bson:"user"`
	AuthInfo `json:"auth" bson:"auth"`
}

type UserFilter struct {
	Limit  *uint32
	Offset *uint32

	Confirmed *bool
	Role      *Role
}
