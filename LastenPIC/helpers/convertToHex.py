cHex = ""
for b in open("bin/LastenPIC.bin", "rb").read():
    cHex += "0x" + hex(b)[2:] + ","

cHex = cHex[:-1]

chex = ("{" + cHex+ "}")

print(chex)

