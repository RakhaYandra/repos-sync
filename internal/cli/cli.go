// Package cli wires the cobra commands.
package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/RakhaYandra/repos-sync/internal/model"
	"github.com/RakhaYandra/repos-sync/internal/output"
	"github.com/RakhaYandra/repos-sync/internal/scan"
)

func NewRootCmd() *cobra.Command {
	var roots []string
	var asJSON bool
	root := &cobra.Command{
		Version: Version,
		Use:     "repos-sync",
		Short:   "One answer to: which of my repositories need attention right now.",
		Long: `repos-sync scans the direct children of its roots for git repositories and
reports branch, upstream sync state, dirty changes and untracked files.

Read-only: no fetch, no locks, no writes inside the scanned repositories.
No network access.`,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(roots) == 0 {
				home, err := os.UserHomeDir()
				if err != nil {
					return err
				}
				roots = []string{home + "/Projects"}
			}
			roots = dedupeRoots(roots)
			report, err := scan.ScanRoots(roots)
			if err != nil {
				_, _ = fmt.Fprintln(cmd.ErrOrStderr(), "repos-sync:", err)
				os.Exit(2)
			}
			render(cmd, report, asJSON)
			return nil
		},
	}
	root.Flags().StringSliceVar(&roots, "root", nil, "root directory to scan (repeatable); defaults to ~/Projects")
	root.Flags().BoolVar(&asJSON, "json", false, "machine-readable output (schema_version 1)")
	root.AddCommand(newCheckCmd())
	return root
}

func newCheckCmd() *cobra.Command {
	var asJSON bool
	cmd := &cobra.Command{
		Use:   "check [PATH]",
		Short: "Inspect one repository",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := "."
			if len(args) == 1 {
				dir = args[0]
			}
			report, err := scan.CheckRepo(dir)
			if err != nil {
				_, _ = fmt.Fprintln(cmd.ErrOrStderr(), "repos-sync:", err)
				os.Exit(2)
			}
			render(cmd, report, asJSON)
			return nil
		},
	}
	cmd.Flags().BoolVar(&asJSON, "json", false, "machine-readable output")
	return cmd
}

func render(cmd *cobra.Command, report model.Report, asJSON bool) {
	if asJSON {
		_ = output.JSON(cmd.OutOrStdout(), report)
		return
	}
	output.Human(cmd.OutOrStdout(), report, isTTY(cmd.OutOrStdout()))
}

func isTTY(w io.Writer) bool {
	f, ok := w.(*os.File)
	if !ok {
		return false
	}
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), syscall.TCGETS, 0); errno != 0 {
		return false
	}
	return true
}

func dedupeRoots(in []string) []string {
	seen := map[string]struct{}{}
	out := []string{}
	for _, r := range in {
		r = filepath.Clean(r)
		if _, ok := seen[r]; ok {
			continue
		}
		seen[r] = struct{}{}
		out = append(out, r)
	}
	return out
}
