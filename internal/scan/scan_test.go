package scan

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func gitRun(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	cmd.Env = append(os.Environ(),
		"GIT_AUTHOR_NAME=test", "GIT_AUTHOR_EMAIL=test@example.com",
		"GIT_COMMITTER_NAME=test", "GIT_COMMITTER_EMAIL=test@example.com",
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
}

func gitRepo(t *testing.T, dir string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	gitRun(t, dir, "init", "-b", "main")
}

func TestScanRootsMixed(t *testing.T) {
	root := t.TempDir()
	gitRepo(t, root+"/alpha")
	gitRepo(t, root+"/beta")
	_ = os.Mkdir(root+"/plain", 0o755)

	report, err := ScanRoots([]string{root})
	if err != nil {
		t.Fatal(err)
	}
	if report.Summary.Repos != 2 {
		t.Fatalf("repos=%d", report.Summary.Repos)
	}
	// both are empty repos; still reported with an error marker, not dropped
	for _, r := range report.Repos {
		if r.Error == nil || *r.Error != "empty repository" {
			t.Fatalf("%s error=%v", r.Name, r.Error)
		}
	}
	if report.SchemaVersion != 1 {
		t.Fatalf("schema_version=%d", report.SchemaVersion)
	}
}

func TestScanRootsBareAsRoot(t *testing.T) {
	bare := t.TempDir() + "/b.git"
	if err := os.MkdirAll(bare, 0o755); err != nil {
		t.Fatal(err)
	}
	gitRun(t, bare, "init", "--bare", "-b", "main")

	report, err := ScanRoots([]string{bare})
	if err != nil {
		t.Fatal(err)
	}
	if report.Summary.Repos != 1 {
		t.Fatalf("repos=%d", report.Summary.Repos)
	}
	if !report.Repos[0].Bare {
		t.Fatal("expected bare")
	}
}

func TestScanRootsNotADirectory(t *testing.T) {
	if _, err := ScanRoots([]string{t.TempDir() + "/missing"}); err == nil {
		t.Fatal("expected error for missing root")
	}
}

func TestCheckRepoOnNonRepo(t *testing.T) {
	if _, err := CheckRepo(t.TempDir()); err == nil {
		t.Fatal("expected error for non-repo")
	}
}
