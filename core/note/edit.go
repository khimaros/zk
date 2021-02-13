package note

import (
	"fmt"
	"os"
	"strings"

	"github.com/kballard/go-shellquote"
	"github.com/mickael-menu/zk/core/zk"
	"github.com/mickael-menu/zk/util/errors"
	executil "github.com/mickael-menu/zk/util/exec"
	"github.com/mickael-menu/zk/util/opt"
	osutil "github.com/mickael-menu/zk/util/os"
)

// Edit starts the editor with the notes at given paths.
func Edit(zk *zk.Zk, paths ...string) error {
	editor := editor(zk)
	if editor.IsNull() {
		return fmt.Errorf("no editor set in config")
	}

	cmd := executil.CommandFromString(editor.String() + " " + shellquote.Join(paths...))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return errors.Wrapf(cmd.Run(), "failed to launch editor: %s %s", editor, strings.Join(paths, " "))
}

// editor returns the editor command to use to edit a note.
func editor(zk *zk.Zk) opt.String {
	return zk.Config.Tool.Editor.
		Or(osutil.GetOptEnv("VISUAL")).
		Or(osutil.GetOptEnv("EDITOR"))
}
