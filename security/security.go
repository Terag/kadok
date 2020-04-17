package security

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// RolesTree handles the hierarchy of roles and permissions using RBAC structures
type RolesTree struct {
	Roles  map[string]*Role
	Buffer []Role
}

// MakeEmptyRolesTree returns an empty SecurityRoles structure
func MakeEmptyRolesTree() RolesTree {
	var rolesTree RolesTree
	rolesTree.Roles = make(map[string]*Role)
	rolesTree.Buffer = make([]Role, 0)
	return rolesTree
}

// MakeRolesTreeFromFile returns a RolesTree generated from a file
func MakeRolesTreeFromFile(path string) (RolesTree, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return MakeEmptyRolesTree(), err
	}

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return MakeEmptyRolesTree(), err
	}

	var rolesTree RolesTree
	err = json.Unmarshal(byteValue, &rolesTree)
	if err != nil {
		return MakeEmptyRolesTree(), err
	}

	jsonFile.Close()
	return rolesTree, nil
}

// UnmarshalJSON implementation for type Role
func (rolesTree *RolesTree) UnmarshalJSON(b []byte) error {

	// Defined the specific json structure for the context
	type jsonRoles struct {
		Roles []struct {
			Name        string       `json:"name"`
			Parent      string       `json:"parent"`
			Permissions []Permission `json:"permissions"`
		} `json:"roles"`
	}

	// Unmarshal it
	var jRoles jsonRoles
	err := json.Unmarshal(b, &jRoles)
	if err != nil {
		return err
	}

	// Populate the context
	rolesTree.Roles = make(map[string]*Role)
	for _, jRole := range jRoles.Roles {
		rolesTree.Buffer = append(rolesTree.Buffer, Role{jRole.Name, nil, jRole.Permissions})
		rolesTree.Roles[jRole.Name] = &rolesTree.Buffer[len(rolesTree.Buffer)-1]
	}
	// Bind the roles with their parent if present
	for _, jRole := range jRoles.Roles {
		if len(jRole.Parent) > 0 {
			rolesTree.Roles[jRole.Name].Parent = rolesTree.Roles[jRole.Parent]
		}
	}
	return nil
}

// IsGranted must return true if the requested permission is granted
type IsGranted func(permission Permission) bool

// MakeIsGranted return a function is Granted that use the domains RolesTree and playerDiscodRoles to return if a permission is granted
func MakeIsGranted(rolesTree RolesTree, playerDiscordRoles []string) func(permission Permission) bool {
	return func(permission Permission) bool {
		for _, role := range playerDiscordRoles {
			if rolesTree.Roles[role].IsGranted(permission) {
				return true
			}
		}
		return false
	}
}
