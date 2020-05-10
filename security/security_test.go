package security

import "testing"

func TestLoadRolesTree(t *testing.T) {
	rolesTree, err := MakeRolesTreeFromFile("roles.json")
	if err != nil {
		t.Errorf("Error loading all roles member, invalid number, error: %v", err)
	}
	if len(rolesTree.Buffer) != 11 {
		t.Errorf("Error loading roles, got: %v want 11", len(rolesTree.Buffer))
	}
	if !rolesTree.Roles["Admin - Conseil de Guerre"].IsGranted(CallCharacter) {
		t.Errorf("Error checking generated permission hierarchy")
	}
}

func TestRolesTreeIsGranted(t *testing.T) {
	rolesTree, err := MakeRolesTreeFromFile("roles.json")
	if err != nil {
		t.Errorf("Error loading all roles member, invalid number, error: %v", err)
	}
	isGranted := MakeIsGranted(rolesTree, []string{"Admin - Conseil de Guerre"})
	if !isGranted(CallCharacter) {
		t.Errorf("Error checking generated permission hierarchy")
	}
}
