package security

import "encoding/json"

// Permission is a int defined by its name. If present in a role, it means it is granted
type Permission int

// List of available permission
const (
	CallCharacter Permission = iota
	GetCharacterList
	GetHelp
)

func (p Permission) String() string {
	return [...]string{"CallCharacter", "GetCharacterList", "GetHelp"}[p]
}

// StringToPermission Convert a string to a Permission type
func StringToPermission(s string) Permission {
	return map[string]Permission{
		"CallCharacter":    CallCharacter,
		"GetCharacterList": GetCharacterList,
		"GetHelp":          GetHelp,
	}[s]
}

// UnmarshalJSON for custom Permission type
func (p *Permission) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*p = StringToPermission(s)
	return nil
}

// Role is defined by a name and can ihnerit permissions from a parent. It also contains a list of permissions
type Role struct {
	Name        string
	Parent      *Role
	Permissions []Permission
}

// Find the index of the permission in the role, returns -1 if not present
func (r Role) Find(permission Permission) int {
	for i, v := range r.Permissions {
		if v == permission {
			return i
		}
	}
	return -1
}

// IsGranted check if the permission is granted to the role
func (r Role) IsGranted(permission Permission) bool {
	result := r.Find(permission) > -1
	if !result && r.Parent != nil {
		result = r.Parent.IsGranted(permission)
	}
	return result
}

// Grant a permission to the role
func (r *Role) Grant(permission Permission) {
	if !r.IsGranted(permission) {
		r.Permissions = append(r.Permissions, permission)
	}
}

// Deny a permission from the role
func (r *Role) Deny(permission Permission) {
	i := r.Find(permission)
	if i > -1 {
		r.Permissions = append(r.Permissions[:i], r.Permissions[i+1:]...)
	}
}
