package main

import (
	"github.com/spf13/cobra"
)

func RootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "envgen [flags] [file or glob pattern]",
		Short:        "Use .yaml configs to generate typed env variable getters",
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
	}

	pkg := cmd.Flags().StringP("package", "p", "config", "Generated package name.")
	output := cmd.Flags().StringP("output", "o", "env_gen.go", "Path to write the generated files.")
	doc := cmd.Flags().StringP("doc", "d", "", "Path to write generated documentation.")

	cmd.RunE = func(_ *cobra.Command, args []string) error {
		return do(rootArgs{
			pkgName: *pkg,
			outName: *output,
			docName: *doc,
			files:   args,
		})
	}

	return cmd
}
