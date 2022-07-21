#include "Socks4.h"

DWORD newSocksConn(LastenPIC* pLastenPIC, SocksConnectFrame* pSocksFrame, SOCKET* pSocket, SocksConnectReplyFrame* pReplyFrame) {

	DWORD dwSuccess = FAIL;
	int on = 1, rv = 0;
	uint32_t ip = 0;
	SOCKADDR_IN clientService = { 0x00 };

	if (pSocksFrame->version != 0x4)
		goto exit;

	if (pSocksFrame->command != 0x1)
		goto exit;

	*pSocket = pLastenPIC->fPointers._socket(AF_INET, SOCK_STREAM, IPPROTO_TCP);
	if (*pSocket == INVALID_SOCKET)
		goto exit;

	dwSuccess = getIp(pLastenPIC, pSocksFrame, &ip);
	if (dwSuccess == FAIL)
		goto exit;

	clientService.sin_family = AF_INET;
	clientService.sin_addr.s_addr = ip;
	clientService.sin_port = pSocksFrame->port;

	rv = pLastenPIC->fPointers._Connect(*pSocket, (const SOCKADDR_IN*)&clientService, sizeof(SOCKADDR_IN));
	if (rv == SOCKET_ERROR)
		goto exit;

	rv = pLastenPIC->fPointers._ioctlsocket(*pSocket, FIONBIO, (u_long*)&on);
	if (rv != 0) {
		pLastenPIC->fPointers._Closesocket(*pSocket);
		goto exit;
	}

	dwSuccess = SUCCESS;

exit:

	pReplyFrame->rep = (dwSuccess == FAIL) ? 0x5b : 0x5a;

	return dwSuccess;

}

DWORD getIp(LastenPIC* pLastenPIC, SocksConnectFrame* pSocksFrame, uint32_t* pResolvedIp) {

	DWORD dwSuccess = FAIL;
	DNS_RECORD dnsRecord = { 0x00 };
	PDNS_RECORD pDnsRecord = &dnsRecord;
	DNS_STATUS dns_status = { 0x00 };

	if ( 
		 ( (pSocksFrame->ip >> (8 * 0)) & 0xff ) == 0x00 && 
	     ( (pSocksFrame->ip >> (8 * 1)) & 0xff) == 0x00 && 
		 ( (pSocksFrame->ip >> (8 * 2)) & 0xff) == 0x00 &&
		 ( (pSocksFrame->ip >> (8 * 3)) & 0xff) == 0x1
	)  {

		dns_status = pLastenPIC->fPointers._DnsQuery_A((PCSTR)&pSocksFrame->buf, DNS_TYPE_A, DNS_QUERY_BYPASS_CACHE, NULL, &pDnsRecord, NULL);
		if (dns_status != ERROR_SUCCESS || pDnsRecord == NULL)
			goto exit;

		*pResolvedIp = pDnsRecord->Data.A.IpAddress;
		
	}
	else {
		*pResolvedIp = pSocksFrame->ip;
	}

	dwSuccess = SUCCESS;
exit:

	return dwSuccess;

}
