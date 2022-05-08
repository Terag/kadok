// The security package contains all the features related to security and permission management.
// The structure of the permissions is based on the RBAC model.
package security

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

// Properties information for the Security package.
// Yaml structure to add to the properties' hierarchy :
//
//  roles: <roles_file_path> (file path to the roles hierarchy file, see Kadok's wiki for more)
type Properties struct {
	RolesTreeFile  string
	RolesHierarchy RolesTree
}

// UnmarshalYAML implementation for the package Properties.
func (properties *Properties) UnmarshalYAML(unmarshal func(interface{}) error) error {

	type PropertiesYAML struct {
		RolesTreeFile string `yaml:"roles"`
	}

	var propertiesYAML PropertiesYAML
	err := unmarshal(&propertiesYAML)
	if err != nil {
		return err
	}

	rolesTree, err := MakeRolesTreeFromFile(propertiesYAML.RolesTreeFile)
	if err != nil {
		return err
	}
	properties.RolesTreeFile = propertiesYAML.RolesTreeFile
	properties.RolesHierarchy = rolesTree

	return nil
}

// RolesTree handles the hierarchy of roles and permissions using RBAC structures.
type RolesTree struct {
	Roles  map[string]*Role
	Buffer []Role
}

// MakeEmptyRolesTree returns an empty RolesTree structure.
func MakeEmptyRolesTree() RolesTree {
	var rolesTree RolesTree
	rolesTree.Roles = make(map[string]*Role)
	rolesTree.Buffer = make([]Role, 0)
	return rolesTree
}

// MakeRolesTreeFromFile returns a RolesTree generated from a file.
// See Kadok's wiki for more information regarding the file format
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

	err = jsonFile.Close()
	if err != nil {
		fmt.Println(err)
	}
	return rolesTree, nil
}

// UnmarshalJSON implementation for type RolesTree.
func (rolesTree *RolesTree) UnmarshalJSON(b []byte) error {

	// Defined the specific json structure for the context.
	type jsonRoles struct {
		Roles []struct {
			Name        string				`json:"name"`
			Parent      string        		`json:"parent"`
			Permissions []Permission  		`json:"permissions"`
			Type		RoleType		  	`json:"type"`
			Description	string				`json:"description"`
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
		// Roles of type Clan or Group must have an emoji defined and not used by another Role
		rolesTree.Buffer = append(rolesTree.Buffer, Role{jRole.Name, nil, jRole.Permissions, jRole.Type, jRole.Description})
		rolesTree.Roles[jRole.Name] = &rolesTree.Buffer[len(rolesTree.Buffer)-1]
	}
	// Bind the roles with their parent if present.
	for _, jRole := range jRoles.Roles {
		if len(jRole.Parent) > 0 {
			rolesTree.Roles[jRole.Name].Parent = rolesTree.Roles[jRole.Parent]
		}
	}
	return nil
}

// IsGranted must return true if the requested permission is granted.
type IsGranted func(entity PermissionRequester) bool

// MakeIsGranted return a function IsGranted that uses a RolesTree structure for the roles hierarchy
// and playerDiscodRoles representing a user's roles to return if a permission is granted.
func MakeIsGranted(rolesTree RolesTree, playerDiscordRoles []string) IsGranted {
	return func(entity PermissionRequester) bool {
		// In the case an entity does not require a permission
		if entity.GetPermission() == EmptyPermission {
			return true
		}
		// Check that the permission is granted
		for _, role := range playerDiscordRoles {
			if treeRole, ok := rolesTree.Roles[role]; ok {
				if treeRole.IsGranted(entity.GetPermission()) {
					return true
				}
			}
		}
		return false
	}
}

func (rolesTree *RolesTree) GetRolesByType(t RoleType) []Role {
	var roles []Role
	for _, role := range rolesTree.Buffer {
		if role.Type == t {
			roles = append(roles, role)
		}
	}
	return roles
}

type AddRole func(role Role) error

type RemoveRole func(role Role) error

type ErrorRoleNotFound string

func (error ErrorRoleNotFound) Error() string {
	return string(error)
}

func (rolesTree *RolesTree) GetGroups() []Role {
	return rolesTree.GetRolesByType(RoleGroup)
}

func (rolesTree *RolesTree) JoinGroup(addGroup AddRole, groupReference string) (Role, error) {
	if group, ok := rolesTree.Roles[groupReference]; ok && group.Type == RoleGroup {
		return *group, group.Join(addGroup)
	}

	groupIndex, err := strconv.Atoi(groupReference)
	if err == nil {
		groups := rolesTree.GetGroups()
		if -1 < groupIndex && groupIndex < len(groups) {
			return groups[groupIndex], groups[groupIndex].Join(addGroup)
		}
	}

	return Role{}, errors.New("role not configured")
}

func (rolesTree *RolesTree) LeaveGroup(removeRole RemoveRole, groupReference string) (Role, error) {
	if group, ok := rolesTree.Roles[groupReference]; ok && group.Type == RoleGroup {
		return *group, group.Leave(removeRole)
	}

	groupIndex, err := strconv.Atoi(groupReference)
	if err == nil {
		groups := rolesTree.GetGroups()
		if -1 < groupIndex && groupIndex < len(groups) {
			return groups[groupIndex], groups[groupIndex].Leave(removeRole)
		}
	}

	return Role{}, errors.New("role not configured")
}

func (rolesTree *RolesTree) GetClans() []Role {
	return rolesTree.GetRolesByType(RoleClan)
}

func (rolesTree *RolesTree) JoinClan(addClan AddRole, removeClan RemoveRole, clanReference string) (Role, error) {
	clans := rolesTree.GetClans()

	clanIndex, err := strconv.Atoi(clanReference)
	if err == nil && -1 < clanIndex && clanIndex < len(clans) {
		clanReference = clans[clanIndex].Name
	}

	if clan, ok := rolesTree.Roles[clanReference]; ok || clan.Type == RoleClan {
		err = nil
		for _, clan := range clans {
			if clan.Name == clanReference {
				err = clan.Join(addClan)
			} else {
				err = clan.Leave(removeClan)
			}
			if err != nil {
				return Role{}, err
			}
		}
		return *clan, nil
	}

	return Role{}, errors.New("role not configured or is not a clan")
}

func (rolesTree *RolesTree) LeaveClan(removeClan RemoveRole, clanReference string) (Role, error) {
	if clan, ok := rolesTree.Roles[clanReference]; ok && clan.Type == RoleClan {
		return *clan, clan.Leave(removeClan)
	}

	clanIndex, err := strconv.Atoi(clanReference)
	if err == nil {
		clans := rolesTree.GetClans()
		if -1 < clanIndex && clanIndex < len(clans) {
			return clans[clanIndex], clans[clanIndex].Leave(removeClan)
		}
	}

	return Role{}, errors.New("role not configured")
}
