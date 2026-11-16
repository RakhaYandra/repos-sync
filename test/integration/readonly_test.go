package integration

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func repoAt(t *testing.T, dir string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	gitRun(t, dir, "init", "-b", "main")
	if err := os.WriteFile(dir+"/seed.txt", []byte("seed"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitRun(t, dir, "add", "seed.txt")
	gitRun(t, dir, "commit", "-m", "seed")
}

func gitRun(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	cmd.Env = append(os.Environ(),
		"GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@example.com",
		"GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@example.com")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("%s: %v\n%s", strings.Join(args, " "), err, out)
	}
}

func TestScanDoesNotModifyRepositories(t *testing.T) {
	root := t.TempDir()
	repoAt(t, root+"/a")
	repoAt(t, root+"/b")

	mtimeOf := func(p string) int64 {
		fi, err := os.Stat(p)
		if err != nil {
			t.Fatal(err)
		}
		return fi.ModTime().UnixNano()
	}
	indexA := root + "/a/.git/index"
	headA := root + "/a/.git/HEAD"
	before := map[string]int64{"index": mtimeOf(indexA), "HEAD": mtimeOf(headA)}

	cmd := exec.Command("go", "run", "../../cmd/repos-sync", "--root", root, "--json")
	var out strings.Builder
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		t.Fatalf("run: %v\n%s", err, out.String())
	}
	if !strings.Contains(out.String(), `"schema_version": 1`) {
		t.Fatalf("unexpected output: %s", out.String())
	}
	if mtimeOf(indexA) != before["index"] || mtimeOf(headA) != before["HEAD"] {
		t.Fatal("repository metadata changed during scan")
	}
	entries, err := os.ReadDir(root + "/a/.git")
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), "index.lock") {
			t.Fatal("scan created a lock file")
		}
	}
}
