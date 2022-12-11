package bot

import (
	"log"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/terag/kadok/internal/info"
	"github.com/terag/kadok/pkg/security"
)

const templatesFolderPath = "../assets/templates"

func TestResolveAndActionStatus(t *testing.T) {
	action, _ := ResolveAction(&RootAction, []string{"tatan"})
	if action != &StatusAction {
		t.Errorf("Tried to resolve status action, different action was found")
	}
}

func TestResolveActionNil(t *testing.T) {
	action, _ := ResolveAction(&RootAction, []string{"nothing_to_be_found"})
	if action != &RootAction {
		t.Errorf("RootAction should have been resolved")
	}
}

func TestResolveActionStatusHelp(t *testing.T) {
	packageDir, err := os.Getwd()
	if err != nil {
		log.Fatal("Error getting current dir")
	}
	err = os.Chdir("..")
	if err != nil {
		log.Fatal("Error moving the current working directory")
	}

	action, execute := ResolveAction(&RootAction, []string{"tatan", "aide"})
	if action != &StatusAction {
		t.Errorf("Tried to resolve status action, different action was found")
	}
	information, _ := execute(nil, nil, nil)
	if information != StatusAction.Information {
		t.Errorf("Did not retrieve the expected information about the action")
	}

	err = os.Chdir(packageDir)
	if err != nil {
		log.Fatal("Error moving the current working directory")
	}
}

func TestActionGetPermission(t *testing.T) {
	permission := StatusAction.GetPermission()
	if permission != security.EmptyPermission {
		t.Errorf("Unexpected permission for StatusAction, got: %s", permission.String())
	}
}

func TestExecuteActionStatus(t *testing.T) {
	packageDir, err := os.Getwd()
	if err != nil {
		log.Fatal("Error getting current dir")
	}
	err = os.Chdir(packageDir)
	if err != nil {
		log.Fatal("Error moving the current working directory")
	}

	Configuration.Templates = path.Join(path.Dir(packageDir), templatesFolderPath)
	statusFirstLine := "> " + info.GetInfo().About
	var message discordgo.MessageCreate = discordgo.MessageCreate{
		Message: &discordgo.Message{},
	}
	result, err := StatusAction.Execute(nil, nil, &message, nil)
	if err != nil {
		t.Errorf("Error executing status action")
	}
	if result[:strings.IndexRune(result, '\n')] != statusFirstLine {
		t.Errorf("Got wrong status message")
	}

	err = os.Chdir(packageDir)
	if err != nil {
		log.Fatal("Error moving the current working directory")
	}
}

func TestGetFuncMap(t *testing.T) {
	functions := getFuncMap([]string{"param1", "param2"})

	if withParams, check := functions["withParams"]; check {
		if !withParams.(func(param string) bool)("param1") {
			t.Errorf("Error, unexpected output, expected true")
		}
	} else {
		t.Errorf("Did not find withParams function")
	}
}
