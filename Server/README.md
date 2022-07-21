## Talje 
> Socks - [ws|tcp] server/client

## Usage
```bash
Usage:
  talje [command]

Available Commands:
  client      Run Client Mode
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  server      Run server Mode

Flags:
  -h, --help   help for talje

Use "talje [command] --help" for more information about a command.
```

## Install 

```bash
go build 
# get coffee

```

## License
Released under MIT. More information about the used dependencies can be found in the go.mod/go.sum file.

**Apache License 2.0**
- https://github.com/spf13/cobra
- https://github.com/OJ/gobuster (copied parts for the cmd/)
- https://github.com/inconshreveable/mousetrap
**MIT license**
- https://github.com/armon/go-socks5 

**BSD-3-Clause license**
- https://github.com/google/uuid
- https://github.com/gorilla/mux

**BSD-2-Clause license**
- https://github.com/gorilla/websocket

