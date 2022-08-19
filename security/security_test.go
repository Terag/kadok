package security

import (
	"fmt"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestSecurityPropertiesUnmarshalYAML(t *testing.T) {
	propertiesYAML := []byte("roles: \"../config/roles.yaml\"")
	var properties Properties
	err := yaml.Unmarshal(propertiesYAML, &properties)
	if err != nil {
		t.Errorf("Error loading security module properties: %v", err)
	}
	if len(properties.RolesHierarchy.Buffer) == 0 {
		t.Errorf("Error, no roles were loaded, got: %v want more", len(properties.RolesHierarchy.Buffer))
	}
}

func TestMakeEmptyRolesTree(t *testing.T) {
	roles := MakeEmptyRolesTree()
	if len(roles.Buffer) != 0 {
		t.Errorf("Error, expected no roles in Buffer, got: %v want 0", len(roles.Buffer))
	}
	if len(roles.Roles) != 0 {
		t.Errorf("Error, expected no roles in Map, got: %v want 0", len(roles.Roles))
	}
}

func TestLoadRolesTree(t *testing.T) {
	rolesTree, err := MakeRolesTreeFromFile("../config/roles.yaml")
	if err != nil {
		t.Errorf("Error loading all roles, error: %v", err)
	}
	if len(rolesTree.Buffer) != 26 {
		t.Errorf("Error loading roles, got: %v want 11", len(rolesTree.Buffer))
	}
	if !rolesTree.Roles["Admin - Conseil de Guerre"].IsGranted(CallCharacter) {
		t.Errorf("Error checking generated permission hierarchy")
	}
}

func TestRolesTree_IsGranted(t *testing.T) {
	rolesTree, err := MakeRolesTreeFromFile("../config/roles.yaml")
	if err != nil {
		t.Errorf("Error loading all roles, error: %v", err)
	}
	isGranted := MakeIsGranted(rolesTree, []string{"Admin - Conseil de Guerre"})
	if !isGranted(CallCharacter) {
		t.Errorf("Error checking generated permission hierarchy")
	}
}

func TestRolesTree_GetGroups(t *testing.T) {
	rolesTree, err := MakeRolesTreeFromFile("../config/roles.yaml")
	if err != nil {
		t.Errorf("Error loading all roles, error: %v", err)
	}

	groups := rolesTree.GetGroups()
	if len(groups) != 15 {
		t.Errorf("Error retrieving groups, expected 15 and got %v", len(groups))
	}
	isCubisteFound := false
	for _, group := range groups {
		if group.Name == "Cubiste" {
			isCubisteFound = true
			break
		}
	}
	if !isCubisteFound {
		t.Errorf("Error, Cubiste group was expected in the list of groups")
	}
}

func TestRolesTree_GetClans(t *testing.T) {
	rolesTree, err := MakeRolesTreeFromFile("../config/roles.yaml")
	if err != nil {
		t.Errorf("Error loading all roles, error: %v", err)
	}

	clans := rolesTree.GetClans()
	if len(clans) != 4 {
		t.Errorf("Error retrieving clans, expected 4 and got %v", len(clans))
	}
	isPedestreSeniorsFound := false
	for _, clan := range clans {
		if clan.Name == "Pédèstres seniors" {
			isPedestreSeniorsFound = true
			break
		}
	}
	if !isPedestreSeniorsFound {
		t.Errorf("Error, Pédèstres seniors group was expected in the list of groups")
	}
}

func TestRolesTree_JoinClanByReference(t *testing.T) {
	rolesTree, err := MakeRolesTreeFromFile("../config/roles.yaml")
	if err != nil {
		t.Errorf("Error loading all roles, error: %v", err)
	}

	clans := rolesTree.GetClans()
	reference := clans[0].Name
	clansTested := make(map[string]bool, len(clans))
	rolesTree.JoinClan(
		func(role Role) error {
			clansTested[role.Name] = true
			return nil
		},
		func(role Role) error {
			clansTested[role.Name] = false
			if role.Name == reference {
				t.Errorf("Error bad leave role, didn't expect: " + reference + " and got: " + role.Name)
			}
			return nil
		},
		reference,
	)

	for _, clan := range clans {
		if value, ok := clansTested[clan.Name]; ok {
			if clan.Name == reference && !value {
				t.Errorf("Error bad join role, expected to be in : " + clan.Name)
			} else if clan.Name != reference && value {
				t.Errorf("Error bad leave role, didn't expect to be in : " + clan.Name)
			}
		} else {
			fmt.Println("Error, did not iterate on all the available clans, missing: " + clan.Name)
		}
	}
}

func TestRolesTree_JoinClanByIndex(t *testing.T) {
	rolesTree, err := MakeRolesTreeFromFile("../config/roles.yaml")
	if err != nil {
		t.Errorf("Error loading all roles, error: %v", err)
	}

	clans := rolesTree.GetClans()
	reference := clans[0].Name
	clansTested := make(map[string]bool, len(clans))
	rolesTree.JoinClan(
		func(role Role) error {
			clansTested[role.Name] = true
			return nil
		},
		func(role Role) error {
			clansTested[role.Name] = false
			if role.Name == reference {
				t.Errorf("Error bad leave role, didn't expect: " + reference + " and got: " + role.Name)
			}
			return nil
		},
		"0",
	)

	for _, clan := range clans {
		if value, ok := clansTested[clan.Name]; ok {
			if clan.Name == reference && !value {
				t.Errorf("Error bad join role, expected to be in : " + clan.Name)
			} else if clan.Name != reference && value {
				t.Errorf("Error bad leave role, didn't expect to be in : " + clan.Name)
			}
		} else {
			fmt.Println("Error, did not iterate on all the available clans, missing: " + clan.Name)
		}
	}
}

func TestRolesTree_LeaveClanByReference(t *testing.T) {
	rolesTree, err := MakeRolesTreeFromFile("../config/roles.yaml")
	if err != nil {
		t.Errorf("Error loading all roles, error: %v", err)
	}

	clans := rolesTree.GetClans()
	reference := clans[0].Name
	rolesTree.LeaveClan(
		func(role Role) error {
			if role.Name != reference {
				t.Errorf("Error, didn't expect to leave: " + role.Name + " wanted to leave: " + reference)
			}
			return nil
		},
		reference,
	)
}

func TestRolesTree_LeaveClanByIndex(t *testing.T) {
	rolesTree, err := MakeRolesTreeFromFile("../config/roles.yaml")
	if err != nil {
		t.Errorf("Error loading all roles, error: %v", err)
	}

	clans := rolesTree.GetClans()
	reference := clans[0].Name
	rolesTree.LeaveClan(
		func(role Role) error {
			if role.Name != reference {
				t.Errorf("Error, didn't expect to leave: " + role.Name + " wanted to leave: " + reference)
			}
			return nil
		},
		"0",
	)
}

func TestRolesTree_JoinGroupByReference(t *testing.T) {
	rolesTree, err := MakeRolesTreeFromFile("../config/roles.yaml")
	if err != nil {
		t.Errorf("Error loading all roles, error: %v", err)
	}

	clans := rolesTree.GetGroups()
	reference := clans[0].Name
	rolesTree.JoinGroup(
		func(role Role) error {
			if role.Name != reference {
				t.Errorf("Error, didn't expect to join: " + role.Name + " wanted to join: " + reference)
			}
			return nil
		},
		reference,
	)
}

func TestRolesTree_JoinGroupByIndex(t *testing.T) {
	rolesTree, err := MakeRolesTreeFromFile("../config/roles.yaml")
	if err != nil {
		t.Errorf("Error loading all roles, error: %v", err)
	}

	clans := rolesTree.GetGroups()
	reference := clans[0].Name
	rolesTree.JoinGroup(
		func(role Role) error {
			if role.Name != reference {
				t.Errorf("Error, didn't expect to join: " + role.Name + " wanted to join: " + reference)
			}
			return nil
		},
		"0",
	)
}

func TestRolesTree_LeaveGroupByReference(t *testing.T) {
	rolesTree, err := MakeRolesTreeFromFile("../config/roles.yaml")
	if err != nil {
		t.Errorf("Error loading all roles, error: %v", err)
	}

	clans := rolesTree.GetGroups()
	reference := clans[0].Name
	rolesTree.LeaveGroup(
		func(role Role) error {
			if role.Name != reference {
				t.Errorf("Error, didn't expect to leave: " + role.Name + " wanted to leave: " + reference)
			}
			return nil
		},
		reference,
	)
}

func TestRolesTree_LeaveGroupByIndex(t *testing.T) {
	rolesTree, err := MakeRolesTreeFromFile("../config/roles.yaml")
	if err != nil {
		t.Errorf("Error loading all roles, error: %v", err)
	}

	clans := rolesTree.GetGroups()
	reference := clans[0].Name
	rolesTree.LeaveGroup(
		func(role Role) error {
			if role.Name != reference {
				t.Errorf("Error, didn't expect to leave: " + role.Name + " wanted to leave: " + reference)
			}
			return nil
		},
		"0",
	)
}
