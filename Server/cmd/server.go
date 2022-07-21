package cmd

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/c-f/talje/lib"
	"github.com/spf13/cobra"

	"github.com/kabukky/httpscerts"
)

var (
	serverFlags struct {
		IsHttps   bool
		Addr      string
		Cert      string
		Key       string
		SocksAddr string

		Protocol string
	}
)

// runServer runs a HTTP(s) server handling websocket requests and read socks :)
func runServer(cmd *cobra.Command, args []string) (err error) {
	setProtocol(serverFlags.Protocol)

	log.Println("o--o=o=o")

	socksURI := serverFlags.SocksAddr
	wsURI, err := url.Parse(serverFlags.Addr)
	if err != nil {
		return
	}

	socksSrv, websockSrv := lib.WebsocketServer(socksURI, wsURI.Host, wsURI.Path)

	// configure TLS
	if wsURI.Scheme == "wss" || serverFlags.IsHttps {
		websockSrv.Cert = serverFlags.Cert
		websockSrv.Key = serverFlags.Key
	}

	// Start Everything
	go socksSrv.Start()
	err = websockSrv.Start()
	return
}

// runGenerateCerts generates certs based on the provided userflags
func runGenerateCerts(cmd *cobra.Command, args []string) (err error) {
	var hostnames []string
	var fname string
	flags := cmd.Flags()
	if hostnames, err = flags.GetStringSlice("hostnames"); err != nil {
		return
	}
	if fname, err = flags.GetString("fname"); err != nil {
		return
	}
	return generateTestCerts(hostnames, fname)
}

// nolint:gochecknoinits
func init() {
	cmdSrv := &cobra.Command{
		Use:   "server",
		Short: "Run server Mode",
		RunE:  runServer,
	}

	cmdSrv.Flags().StringVar(&serverFlags.Addr, "addr", "ws://127.0.0.1:1339/yolo", "Specify the listen Addr")
	cmdSrv.Flags().StringVar(&serverFlags.SocksAddr, "socks", "127.0.0.1:1080", "Specify the socks listen Addr")
	cmdSrv.Flags().StringVar(&serverFlags.Protocol, "protocol", "binary", "Define if binary or json protocol should be used")

	// HTTPS options
	cmdSrv.Flags().BoolVar(&serverFlags.IsHttps, "https", false, "Https enabled traffic (can also specified via addr")
	cmdSrv.Flags().StringVar(&serverFlags.Cert, "crt", "test.cert.pem", "Specify the certificate path ")
	cmdSrv.Flags().StringVar(&serverFlags.Key, "key", "test.cert.key.pem", "Specify the private key path ")
	cmdSrv.Flags().SortFlags = false

	cmdCrtGen := &cobra.Command{
		Use:   "generate",
		Short: "Generate test HTTPS certificate",
		RunE:  runGenerateCerts,
	}
	cmdCrtGen.Flags().String("fname", "test.cert", "Name of the Certificates")
	cmdCrtGen.Flags().StringSlice("hostnames", []string{}, "Name of the Certificates")

	// --[Add Commands ]--
	cmdSrv.AddCommand(cmdCrtGen)
	cmdSrv.AddCommand(&cobra.Command{ // redundant helper
		Use:  "run",
		RunE: runServer,
	})
	rootCmd.AddCommand(cmdSrv)
}

// generateTestCerts generates test Certificates for a https wss server
func generateTestCerts(hostnames []string, certFname string) (err error) {

	keyFile := fmt.Sprintf("%s.key.pem", certFname)
	crtFile := fmt.Sprintf("%s.pem", certFname)

	if err = httpscerts.Check(crtFile, keyFile); err != nil {
		log.Println("Generate CA and HTTPS certs")
		serverlist := strings.Join(hostnames, ",")

		err = httpscerts.Generate(crtFile, keyFile, serverlist)
	}
	return err
}
