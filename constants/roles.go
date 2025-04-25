package constants

type Role string

const (
	RoleEngineer Role = "engineer"
	RoleLogistik Role = "logistik"
	RoleManajer  Role = "manajer"
	RoleAdmin    Role = "admin" // Tambahkan jika diperlukan
)

var RolePermissions = map[Role]map[string][]string{
	RoleEngineer: {
		"/assets":      {"GET", "PUT"},
		"/maintenance": {"GET", "POST", "PUT"},
	},
	RoleLogistik: {
		"/assets":      {"GET", "POST", "PUT", "DELETE"},
		"/maintenance": {"GET"},
	},
	RoleManajer: {
		"/assets":      {"GET"},
		"/maintenance": {"GET"},
	},
	RoleAdmin: {
		"/users": {"GET", "POST", "PUT", "DELETE"},
	},
}
