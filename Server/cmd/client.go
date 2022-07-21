package cmd

import (
	"log"
	"net/url"

	"github.com/c-f/talje/lib"
	"github.com/spf13/cobra"
)

var (
	clientFlags struct {
		SocksAddr   string
		WebSockAddr string

		Protocol string
	}
)

// runClient runs a local socks5 server, which transfer connections over the websocket conn
func runClient(cmd *cobra.Command, args []string) error {
	setProtocol(clientFlags.Protocol)

	log.Println("o--o=o=o")

	connectURI, err := url.Parse(clientFlags.WebSockAddr)
	if err != nil {
		return err
	}

	client := lib.WebsocketClient()
	// Start Connection
	return client.Send(connectURI.String())
}

// nolint:gochecknoinits
func init() {
	cmdSrv := &cobra.Command{
		Use:   "client",
		Short: "Run Client Mode",
		RunE:  runClient,
	}
	cmdSrv.Flags().StringVar(&clientFlags.WebSockAddr, "addr", "ws://127.0.0.1:1339/yolo", "Specify the remote websocket Addr")
	cmdSrv.Flags().StringVar(&clientFlags.Protocol, "protocol", "binary", "Define if binary or json protocol should be used")

	// --[Add Commands]--
	rootCmd.AddCommand(cmdSrv)
}
