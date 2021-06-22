package ui

import(
  "os"
  "os/exec"
  "fmt"
  "log"

	"github.com/mitchellh/go-homedir"
)

type TlEditor struct {
  Editor string
  TlFilePath string
  PostEditionScript string
}

var Tle TlEditor

func (tle *TlEditor) Init(tinylogPath string, postScriptPath string) error {
  // Retrieve editor from env variable.
  env_editor := os.Getenv("EDITOR")
  path_editor, err := exec.LookPath(env_editor)
  if err != nil {
    return fmt.Errorf("Couldn't find the editor. os.Getenv('EDITOR') returns:", env_editor, "\n", err)
  }
  tle.Editor = path_editor

  // Make sure tinylog file exists.
	tlFilePath, e := homedir.Expand(tinylogPath)
	if e != nil {
    return fmt.Errorf("Couldn't find tinylog file\n", e)
	}
	_, e = os.Stat(tlFilePath)
	if e != nil {
    return fmt.Errorf("Couldn't find tinylog file\n", e)
	}
  tle.TlFilePath = tlFilePath

  // postScriptPath is optional.
  if postScriptPath != "" {
    psp, e := homedir.Expand(postScriptPath)
    if e != nil {
      return fmt.Errorf("Couldn't find post script file\n", e)
    }
    f, e := os.Stat(psp)
    if e != nil {
      return fmt.Errorf("Couldn't find post script file\n", e)
    }
    // TODO: Check if script is executable.
    if f.Mode()&0111 == 0111 {
      log.Println("IN")
      tle.PostEditionScript = psp
    } else {
      log.Println("OUT")
      return fmt.Errorf("Post script is not executable.")
    }

  }

  return nil
}

func (tle *TlEditor) Push() error {
  cmd_editor := exec.Command(tle.PostEditionScript)
  cmd_editor.Stdin = os.Stdin
  cmd_editor.Stdout = os.Stdout
  cmd_editor.Stderr = os.Stderr

  log.Println(">", tle.PostEditionScript, cmd_editor)

  if err := cmd_editor.Run(); err != nil {
    return fmt.Errorf("Unable to run post script")
  }

  return nil
}

func editTl() {
  cmd_editor := exec.Command(Tle.Editor, Tle.TlFilePath)
  cmd_editor.Stdin = os.Stdin
  cmd_editor.Stdout = os.Stdout
  cmd_editor.Stderr = os.Stderr

  cmd_editor.Run()
}
