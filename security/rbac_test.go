package security

import "testing"

func TestCreateRoleAndPermission(t *testing.T) {
	member := Role{"member", nil, []Permission{CallCharacter},RoleDefault, ""}
	if member.Name != "member" || !member.IsGranted(CallCharacter) {
		t.Errorf("Error creating role member (%v), got => name: \"%s\" want \"member\", And \"CallCharacter\", IsGranted: %t want true", member, member.Name, member.IsGranted(CallCharacter))
	}
}

func TestGrantPermission(t *testing.T) {
	moderator := Role{"moderator", nil, []Permission{}, RoleDefault, ""}
	moderator.Grant(GetHelp)
	if !moderator.IsGranted(GetHelp) {
		t.Errorf("Error granting GetHelp to moderator (%v), IsGranted: %t want true,", moderator, moderator.IsGranted(GetHelp))
	}
}

func TestDenyPermission(t *testing.T) {
	administrator := Role{"administrator", nil, []Permission{GetCharacterList}, RoleDefault, ""}
	administrator.Deny(GetCharacterList)
	if administrator.IsGranted(GetCharacterList) {
		t.Errorf("Error denying GetHelp to moderator (%v), IsGranted: %t want false", administrator, administrator.IsGranted(GetCharacterList))
	}
}

func TestRoleHierarchy(t *testing.T) {
	member := Role{"member", nil, []Permission{CallCharacter}, RoleDefault, ""}
	moderator := Role{"moderator", &member, []Permission{}, RoleDefault, ""}
	administrator := Role{"administrator", &moderator, []Permission{GetCharacterList}, RoleDefault, ""}

	if !moderator.IsGranted(CallCharacter) {
		t.Errorf("Error checking CallCharacter on moderator (%v), IsGranted: %t want true", moderator, moderator.IsGranted(CallCharacter))
	}
	if administrator.IsGranted(GetHelp) {
		t.Errorf("Error checking permission GetHelp on administrator (%v), IsGranted: %t want false", administrator, administrator.IsGranted(GetHelp))
	}
}
