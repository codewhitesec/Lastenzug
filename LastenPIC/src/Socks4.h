#include "windows.h"
#include "windns.h"

#include "Defines.h"

typedef struct  {
	uint8_t version;
	uint8_t command;
	uint16_t port;
	uint32_t ip;
	uint8_t end;
	uint8_t buf[256];
} SocksConnectFrame;

typedef struct{
	uint8_t vn;
	uint8_t rep;
	uint16_t dstport;
	uint32_t dstip;
} SocksConnectReplyFrame;

DWORD newSocksConn(LastenPIC*, SocksConnectFrame* socks_frame, SOCKET* ptr_socket, SocksConnectReplyFrame*);

DWORD getIp(LastenPIC*, SocksConnectFrame*, uint32_t*);
