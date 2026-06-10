package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func newGendocsCommand() *cobra.Command {
	return &cobra.Command{
		Use:    "gendocs <dir>",
		Short:  "Generate CLI markdown documentation",
		Hidden: true,
		Args:   cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			out := args[0]
			if err := os.MkdirAll(out, 0755); err != nil {
				return err
			}
			root := newRootCommand()
			root.DisableAutoGenTag = true
			if err := doc.GenMarkdownTree(root, out); err != nil {
				return err
			}
			rootFile := filepath.Join(out, root.Name()+".md")
			readmeFile := filepath.Join(out, "README.md")
			if _, err := os.Stat(rootFile); err == nil {
				_ = os.Rename(rootFile, readmeFile)
			}
			fmt.Printf("Generated docs in %s\n", out)
			return nil
		},
	}
}
