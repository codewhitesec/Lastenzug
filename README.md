# LastenZug

This project implements a Socka4a proxy based on websockets.    

The client component is implemented in C compiling down to fully position independent code (PIC).

During the compilation process, obfuscation is applied on assembly level by leveraging a second tool: **SpiderPIC** located in ```LastenPIC/SpiderPIC```

## SpiderPIC

The obfuscation includes:
- Instruction substitution
- Adding trash and a jump over the trash
- Adding useless instructions

This is meant to break static signatures, however you need to keep in mind that API hashes, strings and other constants are not obfuscated during this process.

## Usage

### Client
 
The makefile produces both: the PIC socks client and a sample loader for the shellcode.
You can call the shellcode using the following prototype:

```C
DWORD lastenzug(wchar_t* wServerName, PWSTR wPath, DWORD port, PWSTR proxy, PWSTR pUserName, PWSTR pPassword);
```

The sample loader embeds the shellcode in its ***.text*** segment and can be called as follows:
```bash
.\LastenLoader.exe --server [host] --path [path used by server] --port [port]
```

### Server

```bash
cd Server && go build -o LastenServer
./LastenServer server --addr ws://0.0.0.0:8080/lastenzug
```

# Credits
- Our [@invist](https://twitter.com/invist) for implementing the backend
- Our [@thefLinkk](https://twitter.com/thefLinkk) for implementing the client
