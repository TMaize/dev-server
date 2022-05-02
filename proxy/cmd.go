package proxy

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
	cors    *bool
	target  *string
}

func GetCmd() *cobra.Command {
	var cmdArgs CmdArgs

	command := &cobra.Command{
		Use:   "proxy",
		Short: "HTTP reverse proxy",
		Run: func(cmd *cobra.Command, args []string) {
			server := Server{
				Https:   *cmdArgs.https,
				Address: *cmdArgs.address,
				Port:    *cmdArgs.port,
				Domain:  *cmdArgs.domain,
				Cors:    *cmdArgs.cors,
				Target:  *cmdArgs.target,
			}

			if err := server.Run(); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

		},
	}

	cmdArgs = CmdArgs{
		https:   command.Flags().BoolP("https", "s", false, "enable https"),
		address: command.Flags().StringP("address", "a", "0.0.0.0", "listen address"),
		port:    command.Flags().UintP("port", "p", 0, "listen port (default 80/443)"),
		domain:  command.Flags().String("domain", "localhost", "generate cert for domain"),
		cors:    command.Flags().Bool("cors", false, "add cors header"),
		target:  command.Flags().StringP("target", "t", "", "proxy target, https://uptream.com:8443"),
	}

	_ = command.MarkFlagRequired("target")

	return command
}
