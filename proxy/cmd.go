package proxy

import (
	"github.com/spf13/cobra"
)

type CmdArgs struct {
	https   *bool
	address *string
	port    *uint
	domain  *string
	cors    *bool
	target  string
}

func GetCmd() *cobra.Command {

	var cmdArgs CmdArgs

	command := &cobra.Command{
		Use:   "proxy [flags] [target]",
		Short: "HTTP reverse proxy",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				cmdArgs.target = args[len(args)-1]
			}

			// TODO check target

		},
	}

	cmdArgs = CmdArgs{
		https:   command.Flags().BoolP("https", "s", false, "Enable https"),
		address: command.Flags().StringP("address", "a", "0.0.0.0", "Listen address"),
		port:    command.Flags().UintP("port", "p", 0, "Listen port (default 80/443)"),
		domain:  command.Flags().String("domain", "localhost", "Generate cert for domain"),
		cors:    command.Flags().Bool("cors", false, "Add cors header"),
	}

	return command
}
