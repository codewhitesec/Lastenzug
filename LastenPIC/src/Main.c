#include "windows.h"

#include "ApiResolve.h"
#include "Core.h"
#include "Defines.h"
#include "Misc.h"
#include "WS.h"

void wrapOutLoop(LastenPIC *ptr_lastenPIC);

DWORD
lastenzug(wchar_t* wServerName, PWSTR wPath, DWORD port, PWSTR proxy, PWSTR pUserName, PWSTR pPassword) {

	DWORD dwSuccess = FAIL;
	HANDLE hThreadOutLoop = NULL;

	LastenPIC lastenPIC = { 0x00 };

	lastenPIC.bContinue = TRUE;
	FD_ZERO(&lastenPIC.master_set);

	dwSuccess = resolveFPointers(&lastenPIC.fPointers);
	if (dwSuccess == FAIL) 
		goto exit;
  
	lastenPIC.wServerName = wServerName;
	lastenPIC.port = port;
	lastenPIC.wPath = wPath;

	dwSuccess = initWS(&lastenPIC, proxy, pUserName, pPassword);
	if (dwSuccess == FAIL) 
		goto exit;

	hThreadOutLoop = lastenPIC.fPointers._CreateThread(NULL, 0, (LPTHREAD_START_ROUTINE)wrapOutLoop, &lastenPIC, 0, NULL);
	if (hThreadOutLoop == NULL)
		goto exit;

	dwSuccess = inLoop(&lastenPIC);
	if (dwSuccess == FAIL)
		goto exit;

	dwSuccess = SUCCESS;

exit:

	lastenPIC.bContinue = FALSE;
	wsClose(&lastenPIC);

	for (int i = 0; i < MAX_CONNECTIONS; i++) {
		if (lastenPIC.sockets[i])
			lastenPIC.fPointers._Closesocket(lastenPIC.sockets[i]);
	}

	if(hThreadOutLoop)
		lastenPIC.fPointers._WaitForSingleObject(hThreadOutLoop, INFINITE);

	if(hThreadOutLoop)
		lastenPIC.fPointers._CloseHandle(hThreadOutLoop);

	return dwSuccess;

}

void wrapOutLoop(LastenPIC *ptr_lastenPIC) {
  outLoop(ptr_lastenPIC);
}
