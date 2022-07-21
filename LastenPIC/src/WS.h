#pragma once

#include "windows.h"
#include "ApiResolve.h"
#include <winhttp.h>

#include "Defines.h"

DWORD initWS(LastenPIC*, PWSTR, PWSTR, PWSTR);
DWORD readLastenFrame(LastenPIC*, LastenFrame*, LPDWORD);
DWORD sendLastenFrame(LastenPIC*, LastenFrame*);
void wsClose(LastenPIC*);
