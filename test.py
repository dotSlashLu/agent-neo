# FnName      [32]byte
# ParamsLen   uint32
# Params      []byte

# UUID    [32]byte    // vm uuid
# Name    [32]byte    // random str
# Target  [3]byte     // vdb? vdc?
# Slot    [4]byte     // 0x007++
# Size    int32       // size in MB
import socket
import struct

params = {
    "uuid": "b628579d-ae3d-41f0-887e-895204190c70",
    "name": "new-disk1",
    "target": "vcd",
    "slot": "0x12",
    "size": 1024
}
args = {
    "paramLen": 68,
    "fnName": "volume.create",
    "params": struct.pack("<32s32si", params['uuid'], params['name'],
        params['size'])
}
sock = socket.socket(); sock.connect(("localhost", 18103));
header = struct.pack("<i", 0x53b)
sock.send(header)
body = struct.pack("<32s i 68s", args["fnName"], args["paramLen"], args["params"])
sock.send(body)
print(sock.recvfrom(1024))
sock.close()