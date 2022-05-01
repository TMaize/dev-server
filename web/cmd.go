package web

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

type CmdArgs struct {
	https   *bool
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

			server := StaticServer{
				Https:   *cmdArgs.https,
				Address: *cmdArgs.address,
				Port:    *cmdArgs.port,
				Domain:  *cmdArgs.domain,
				Root:    cmdArgs.root,
			}

			if err := server.Run(); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}

	cmdArgs = CmdArgs{
		https:   command.Flags().BoolP("https", "s", false, "Enable https"),
		address: command.Flags().StringP("address", "a", "0.0.0.0", "Listen address"),
		port:    command.Flags().UintP("port", "p", 0, "Listen port (default 80/443)"),
		domain:  command.Flags().String("domain", "localhost", "Generate cert for domain"),
	}

	return command
}
