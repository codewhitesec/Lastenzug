package lib

import (
	"log"

	"github.com/c-f/talje/lib/component/server"
	"github.com/c-f/talje/lib/config"
	"github.com/c-f/talje/lib/ladung"
	"github.com/c-f/talje/lib/laeufer/http"
	"github.com/c-f/talje/lib/laeufer/tcp"
	"github.com/c-f/talje/lib/laeufer/websocket"
)

/*
	Helper structs for dev
*/

func TCP2WebsocketServer() {

}

func TcpAgentServer() {


	incoming := config.Config{
		CommType:  config.Direct,
		Direction: config.Receiving,

		InC:  make(chan ladung.Packet),
		OutC: make(chan ladung.Packet),
	}

	outgoing := config.Config{
		CommType:  config.Agent2Agent,
		Direction: config.Receiving,

		InC:  incoming.OutC,
		OutC: incoming.InC,
	}

	inLaeufer := tcp.New(incoming)
	rs := server.NewSocks5("0.0.0.0:1080")
	rs.Bind(inLaeufer.Receive)

	outLaeufer := tcp.New(outgoing)
	sec := server.NewSocks5("0.0.0.0:1337")
	sec.Bind(outLaeufer.Receive)

	go rs.Start()
	sec.Start()
}

func ServerHTTPSend() {
	incoming := config.Config{
		CommType:  config.Direct,
		Direction: config.Receiving,

		InC:  make(chan ladung.Packet),
		OutC: make(chan ladung.Packet),
	}

	outgoing := config.Config{
		CommType:  config.Agent2Agent,
		Direction: config.Sending,

		InC:  incoming.OutC,
		OutC: incoming.InC,
	}

	inLaeufer := tcp.New(incoming)
	rs := server.NewSocks5("0.0.0.0:1080")
	rs.Bind(inLaeufer.Receive)

	outLaeufer := http.New(outgoing)

	// outLaeufer := tcp.New(outgoing)
	//sec := server.NewHttp("0.0.0.0:1337")
	// sec := server.NewSocks5("0.0.0.0:1337")
	//sec.Bind("/yolo", outLaeufer.Receive)

	go rs.Start()
	outLaeufer.Send("http://127.0.0.1:1337/yolo")
}

func ServerHTTPReceive() {
	incoming := config.Config{
		CommType:  config.Direct,
		Direction: config.Receiving,

		InC:  make(chan ladung.Packet),
		OutC: make(chan ladung.Packet),
	}

	outgoing := config.Config{
		CommType:  config.Agent2Agent,
		Direction: config.Receiving,

		InC:  incoming.OutC,
		OutC: incoming.InC,
	}

	inLaeufer := tcp.New(incoming)
	rs := server.NewSocks5("0.0.0.0:1080")
	rs.Bind(inLaeufer.Receive)

	outLaeufer := http.New(outgoing)
	sec := server.NewHttp("0.0.0:1337")
	sec.Bind("/yolo", outLaeufer.Receive)

	go rs.Start()
	sec.Start()
}

func ClientHTTPSend() {

	incoming := config.Config{
		CommType:  config.Agent2Agent,
		Direction: config.Sending,

		InC:  make(chan ladung.Packet),
		OutC: make(chan ladung.Packet),
	}
	outgoing := config.Config{
		CommType:  config.Direct,
		Direction: config.Sending,

		InC:  incoming.OutC,
		OutC: incoming.InC,
	}

	inLaeufer := http.New(incoming)
	//inLaeufer := tcp.New(incoming)

	outLaeufer := tcp.New(outgoing)
	log.Println("o", outLaeufer)

	// Start the rolls
	inLaeufer.Send("http://127.0.0.1:1337/yolo")
}

func WebsocketServer(socksAddr string, httpAddr string, ctxPath string) (*server.Socks5, *server.Http) {
	incoming := config.NewServer()

	outgoing := config.NewAgent(config.Receiving)
	incoming.Bind(outgoing)

	// Socks Server
	inLaeufer := tcp.New(*incoming)
	rs := server.NewSocks5(socksAddr)
	rs.Bind(inLaeufer.Receive)

	// Websocket server
	outLaeufer := websocket.New(*outgoing)
	sec := server.NewHttp(httpAddr)

	// bind both together
	sec.Bind(ctxPath, outLaeufer.Receive)

	return rs, sec
}
func WebsocketClient() *websocket.Laeufer {

	incoming := config.NewAgent(config.Sending)

	outgoing := config.NewClient()
	incoming.Bind(outgoing)

	inLaeufer := websocket.New(*incoming)

	tcp.New(*outgoing)

	return inLaeufer

}

//
func TCPAgentClient() {

	incoming := config.Config{
		CommType:  config.Agent2Agent,
		Direction: config.Sending,

		InC:  make(chan ladung.Packet),
		OutC: make(chan ladung.Packet),
	}
	outgoing := config.Config{
		CommType:  config.Direct,
		Direction: config.Sending,

		InC:  incoming.OutC,
		OutC: incoming.InC,
	}

	inLaeufer := tcp.New(incoming)

	outLaeufer := tcp.New(outgoing)
	log.Println("o", outLaeufer)

	// Start the rolls
	inLaeufer.Send("127.0.0.1:1337")
}
