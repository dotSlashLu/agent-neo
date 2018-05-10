# FnNameLen   uint32
# ParamsLen   uint32
# FnName      string
# Params      []byte

# UUID    [68]byte    // vm uuid
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
    "fnLen": len("volume.createVolume"),
    "paramLen": 111,
    "fnName": "volume.createVolume",
    "params": struct.pack("<68s32s3s4si", params['uuid'], params['name'],
        params['target'], params['slot'], params['size'])
}
sock = socket.socket(); sock.connect(("localhost", 18103));
header = struct.pack("<i", 0x53b);
sock.send(header)
body = struct.pack("<ii19s111s", args["fnLen"], args["paramLen"],
    args["fnName"], args["params"])
sock.send(body)
