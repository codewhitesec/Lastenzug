#include <stdio.h>
#include "stdint.h"
#include "windows.h"

#include "loader.h"

void args(PWSTR*, PDWORD, PWSTR*, PWSTR*, PWSTR*, PWSTR*, int, wchar_t **);
void help(wchar_t**);

int wmain( int argc, wchar_t *argv[], wchar_t *envp[] ) {
    
    PWSTR server = NULL, path = NULL, proxy = NULL, proxyuser = NULL, proxypassword = NULL;
    DWORD dwPort = 0, dwSuccess = 0;

    args( &server, &dwPort, &path, &proxy, &proxyuser, &proxypassword, argc, argv );

    wprintf(L"* Lastenzug - PIC Socks4a proxy by @invist and @theflinkk of @codewhitesec\n");
    wprintf(L"! This is a sample loader for Lastenzug\n");
    wprintf(L"! Using unsigned exe is not recommended\n");

    wprintf(L"+ Args parsed:\n");
    wprintf(L"*\tServer: %S\n", server);
    wprintf(L"*\tPort: %d\n", dwPort);
    wprintf(L"*\tPath: %S\n", path);
    wprintf(L"*\tProxy: %S\n", proxy);
    wprintf(L"*\tProxyuser: %S\n", proxyuser);
    wprintf(L"*\tProxypassword: %S\n", proxypassword);
    wprintf(L"* Now executing embedded PIC ... \n");

    dwSuccess = ( ( LASTENZUG* )lastenzug )( server, path, dwPort, proxy, proxyuser, proxypassword );
  
    wprintf(L"* Lastenzug returned: %d\n", dwSuccess);

cleanup:

    return 0;

}

void 
args( PWSTR* pServer, PDWORD pdwPort, PWSTR* pPath, PWSTR* pProxy, PWSTR* pProxyUser, PWSTR* pProxyPassword, int argc, wchar_t **argv ){

    if (argc < 4)
        help(argv);

    for (int i = 1; i < argc; i++) {

        if (!lstrcmpW(argv[i], L"--server"))
            *pServer = argv[i + 1];

        if (!lstrcmpW(argv[i], L"--port"))
            *pdwPort = _wtoi(argv[i + 1]);

        if (!lstrcmpW(argv[i], L"--path"))
            *pPath = argv[i + 1];

        if (!lstrcmpW(argv[i], L"--proxy"))
            *pProxy = argv[i + 1];

        if (!lstrcmpW(argv[i], L"--proxyuser"))
            *pProxyUser = argv[i + 1];

        if (!lstrcmpW(argv[i], L"--proxypassword"))
            *pProxyPassword = argv[i + 1];

    }

}

void
help(wchar_t** argv) {

    wprintf(L"%S --server [domain/ip] --port [port] --path [path registered to handle socks] --proxy http://proxy:8080 --proxyuser domain\\user --proxypassword Sommer2022\n", argv[0]);
    exit(0);

}