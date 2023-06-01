package web

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

type CmdArgs struct {
	https   *bool
	ui      *bool
	address *string
	port    *uint
	domain  *string
	root    string
}

func GetCmd() *cobra.Command {

	var cmdArgs CmdArgs

	command := &cobra.Command{
		Use:   "start [flags] [root]",
		Short: "Serve static files",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				cmdArgs.root = args[len(args)-1]
			}

			server := Server{
				Https:   *cmdArgs.https,
				Address: *cmdArgs.address,
				Port:    *cmdArgs.port,
				Domain:  *cmdArgs.domain,
				Root:    cmdArgs.root,
				UI:      *cmdArgs.ui,
			}

			if err := server.Run(); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}

	cmdArgs = CmdArgs{
		https:   command.Flags().BoolP("https", "s", false, "enable https"),
		ui:      command.Flags().Bool("ui", false, "list files with ui"),
		address: command.Flags().StringP("address", "a", "0.0.0.0", "listen address"),
		port:    command.Flags().UintP("port", "p", 0, "listen port (default 80/443)"),
		domain:  command.Flags().String("domain", "localhost", "generate cert for domain"),
	}

	return command
}
