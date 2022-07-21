CC := x86_64-w64-mingw32-gcc
CXX := x86_64-w64-mingw32-g++
LD := x86_64-w64-mingw32-ld
CFLAGS := -Wall -m64 -ffunction-sections -fno-asynchronous-unwind-tables -nostdlib -fno-ident -O2 -masm=intel
CFLAGS_ASM := -Wall -m64 -ffunction-sections -fno-asynchronous-unwind-tables -nostdlib -fno-ident -O2 -S -masm=intel
CCLDFLAGS := -Wl,-Tsrc/linker.ld,--no-seh -DC2

S_SRCS := src/entry.asm src/chkstk_ms.asm
C_SRCS := src/ApiResolve.c src/Core.c src/Main.c src/Socks4.c src/WS.c

S_OBJS := $(patsubst src/%.asm,%.o,$(S_SRCS)) $(patsubst src/%.c,%.s,$(C_SRCS)) 
OBJS := $(patsubst src/%.asm,%.o,$(S_SRCS)) $(patsubst src/%.c,%.o,$(C_SRCS)) 

all: SpiderPIC/SpiderPIC bin/LastenPIC.exe bin/LastenPIC.bin bin/LastenLoader.exe

SpiderPIC/SpiderPIC:
	/bin/bash -c "cd SpiderPIC && go build -o SpiderPIC"

bin/LastenPIC.exe: $(S_OBJS) $(OBJS)
	mkdir -p $(@D)
	$(LD) -s $(OBJS) -o $@

bin/LastenPIC.bin: bin/LastenPIC.exe
	objcopy -j .text -O binary $< $@
	rm bin/LastenPIC.exe

loader/loader.h: loader/loader.h.prefix bin/LastenPIC.bin
		(                                                        \
				set -e;                                          \
				cat loader/loader.h.prefix; 				     \
				python3 helpers/convertToHex.py | xargs -i echo '{}';\
				echo ';'                                         \
		) > $@.t
		mv $@.t $@

bin/LastenLoader.exe: loader/loader.cpp loader/loader.h
		mkdir -p $(@D)
		$(CXX) -s -o $@ $< -municode

%.o: src/%.asm
	nasm -f win64 $< -o $@

%.s: src/%.c
	$(CC) $< $(CFLAGS_ASM) -c -o $@ $(CCLDFLAGS)
	SpiderPIC/SpiderPIC -asm $@ -o $@

%.o: %.s
	$(CC) $< $(CFLAGS) -c -o $@ $(CCLDFLAGS)

.PHONY: clean
clean:
	rm -rf $(OBJS) $(S_OBJS)\
		bin/LastenPIC.exe bin/LastenPIC.bin loader/loader.h bin/LastenLoader.exe

