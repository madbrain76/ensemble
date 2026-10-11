package sandbox

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

// SanitizedEnv deliberately allowlists mundane process settings. A suffix
// blacklist misses arbitrary credential names and executable startup hooks.
// Full-access commands retain ordinary noncredential settings for compatibility.
func SanitizedEnv(env []string, root string) []string {
	var out []string
	for _, item := range env {
		name, _, _ := strings.Cut(item, "=")
		upper := strings.ToUpper(name)
		if root != "" {
			switch name {
			case "LANG", "LC_ALL", "LC_CTYPE", "TZ":
				out = append(out, item)
			}
			continue
		}
		if upper == "SSH_AUTH_SOCK" || upper == "TERM" || strings.Contains(upper, "TOKEN") || strings.Contains(upper, "SECRET") || strings.Contains(upper, "PASSWORD") || strings.HasSuffix(upper, "_KEY") {
			continue
		}
		out = append(out, item)
	}
	out = append(out, "TERM=dumb")
	if root != "" {
		out = append(out, "PATH=/usr/bin:/bin:/usr/sbin:/sbin", "HOME="+root, "TMPDIR="+root,
			"GIT_CONFIG_GLOBAL="+filepath.Join(root, ".ensemble-gitconfig"), "GIT_CONFIG_SYSTEM=/dev/null", "GIT_PAGER=cat", "PAGER=cat")
	}
	return out
}

// Command confines the whole descendant tree. The trusted binary's model HTTP
// transport stays outside this profile. Unsupported platforms fail closed.
func (s *sandbox) Command(script, cwd string) (*exec.Cmd, error) {
	root, err := s.root()
	if err != nil {
		return nil, errors.New("sandbox root unavailable")
	}
	if cwd == "" {
		cwd = s.parent.Config().Workspace
	}
	cwd, err = s.Resolve(cwd)
	if err != nil {
		return nil, err
	}
	command := exec.Command("/bin/sh", "-c", script)
	if root != "" {
		if runtime.GOOS != "darwin" {
			return nil, errors.New("sandbox commands unsupported on this platform")
		}
		// Generate only a local identity. Never read or copy the user's gitconfig,
		// which may contain credential helpers, includes and token-bearing rewrites.
		f, err := s.OpenFile(".ensemble-gitconfig", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			return nil, err
		}
		_, writeErr := f.WriteString("[user]\n name = Ensemble\n email = ensemble@localhost\n")
		closeErr := f.Close()
		if err := errors.Join(writeErr, closeErr); err != nil {
			return nil, errors.New("sandbox git identity write failed")
		}
		profile := `(version 1)
(allow default)
(deny file-read*)
(deny file-write*)
(allow file-read* (literal "/") (subpath "/usr") (subpath "/bin") (subpath "/sbin") (subpath "/System") (subpath "/Library") (subpath "/dev") (subpath "/private/etc/ssl") (literal "/private/var/select/sh") (literal "/private/var/select/developer_dir") (literal "/private/var/db/xcode_select_link"))
(allow file-write* (literal "/dev/null") (literal "/dev/tty") (literal "/dev/stdout") (literal "/dev/stderr"))
(allow file-read-metadata (literal "/var") (literal "/tmp") (literal "/etc") (literal "/private") (literal "/private/var") (literal "/private/var/select"))
(deny network*)
`
		profile += fmt.Sprintf("(allow file-read* file-write* (subpath %s))\n", strconv.Quote(root))
		// Git resolves cwd component by component; metadata on ancestors permits
		// that traversal without disclosing their contents or sibling file data.
		for dir := filepath.Dir(root); dir != "/"; dir = filepath.Dir(dir) {
			profile += fmt.Sprintf("(allow file-read-metadata (literal %s))\n", strconv.Quote(dir))
		}
		// Xcode shims require their own cache, not blanket access to the shared
		// Darwin temp directory where unrelated processes keep private files.
		tmp, err := exec.Command("/usr/bin/getconf", "DARWIN_USER_TEMP_DIR").Output()
		if err != nil {
			return nil, errors.New("sandbox cannot locate toolchain cache")
		}
		temp, err := filepath.EvalSymlinks(strings.TrimSpace(string(tmp)))
		if err != nil {
			return nil, errors.New("sandbox cannot resolve toolchain cache")
		}
		// QuoteMeta is unnecessary for SBPL's literal prefix when using a subpath
		// filter plus a basename regex: both must match the same file.
		profile += fmt.Sprintf("(allow file-read* file-write* (require-all (subpath %s) (regex #\"/xcrun_db[^/]*$\")))\n", strconv.Quote(temp))
		command = exec.Command("/usr/bin/sandbox-exec", "-p", profile, "/bin/sh", "-c", script)
	}
	command.Dir = cwd
	command.Env = SanitizedEnv(os.Environ(), root)
	return command, nil
}
