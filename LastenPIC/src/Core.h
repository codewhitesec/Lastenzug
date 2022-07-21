#pragma once

#include "windows.h"
#include "winhttp.h"

#include "Defines.h"
#include "Socks4.h"
#include "WS.h"

DWORD inLoop(LastenPIC*);
void outLoop(LastenPIC*);

DWORD getSocksConn(LastenPIC*, uint64_t, SOCKET*);
void handleNewConn(LastenPIC*, LastenFrame* ptr_frame_in, LastenFrame*);
void handleLastenFrame(LastenPIC*, LastenFrame*, LastenFrame*);
void handleSend(LastenPIC*, LastenFrame*, LastenFrame*);
void handleStopConn(LastenPIC*, uint64_t);
