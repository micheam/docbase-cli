package docbasecli

// The original implementation can be found below.
// https://samrapdev.com/capturing-sensitive-input-with-editor-in-golang-from-the-cli/

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/micheam/go-docbase"
)

// DefaultEditor is vim because we're adults ;)
const DefaultEditor = "vim"

// PreferredEditorResolver is a function that returns an editor that the user
// prefers to use, such as the configured `$EDITOR` environment variable.
type PreferredEditorResolver func() string

// GetPreferredEditorFromEnvironment returns the user's editor as defined by the
// `$EDITOR` environment variable, or the `DefaultEditor` if it is not set.
func GetPreferredEditorFromEnvironment() string {
	editor := os.Getenv("EDITOR")

	if editor == "" {
		return DefaultEditor
	}

	return editor
}

func resolveEditorArguments(executable string, filename string) []string {
	args := []string{filename}

	if strings.Contains(executable, "Visual Studio Code.app") {
		args = append([]string{"--wait"}, args...)
	}

	// Other common editors

	return args
}

// OpenFileInEditor opens filename in a text editor.
func OpenFileInEditor(filename string, resolveEditor PreferredEditorResolver) error {
	// Get the full executable path for the editor.
	executable, err := exec.LookPath(resolveEditor())
	if err != nil {
		return err
	}

	cmd := exec.Command(executable, resolveEditorArguments(executable, filename)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// CaptureInputFromEditor opens a temporary file in a text editor and returns
// the written bytes on success or an error on failure. It handles deletion
// of the temporary file behind the scenes.
//
// TODO(micheam): tempfile は外部から受け取るべき（内部でつくるとテストできない）
func CaptureInputFromEditor(resolveEditor PreferredEditorResolver, value []byte) ([]byte, error) {
	file, err := ioutil.TempFile(os.TempDir(), "*.md")
	if err != nil {
		return []byte{}, err
	}
	defer func() { _ = os.Remove(file.Name()) }()

	// Edit(micheam): enable to specify, default value
	if len(value) > 0 {
		i, err := file.Write(value)
		if err != nil {
			return []byte{}, err
		}
		log.Printf("write %d bytes of default value", i)
	}

	if err = file.Close(); err != nil {
		return []byte{}, err
	}
	if err = OpenFileInEditor(file.Name(), resolveEditor); err != nil {
		return []byte{}, err
	}
	bytes, err := ioutil.ReadFile(file.Name())
	if err != nil {
		return []byte{}, err
	}
	return bytes, nil
}

func OpenBrowser(_ context.Context, post docbase.Post) error {
	openbrowser(post.URL)
	return nil
}

func openbrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}
