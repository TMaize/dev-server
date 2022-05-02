package main

import (
	"github.com/TMaize/dev-server/proxy"
	"github.com/TMaize/dev-server/web"
	"github.com/spf13/cobra"
)

var Version = "dev"

func main() {

	command := &cobra.Command{
		Use:     "dev-server",
		Short:   "A simple dev server.",
		Version: Version,
	}

	command.AddCommand(web.GetCmd())
	command.AddCommand(proxy.GetCmd())
	_ = command.Execute()

}
