from protlib import *

from api import call

def volume_create(uuid, name, size):
    # UUID    [36]byte    // vm uuid
    # Name    [32]byte    // random str
    # Target  [3]byte     // vdb? vdc?
    # Slot    [4]byte     // 0x007++
    # Size    int32       // size in MB
    class params_proto(CStruct):
        uuid    = CString(length=36)
        name    = CString(length=32)
        size    = CInt()
    params = {
        "uuid": uuid,
        "name": name,
        "size": size
    }
    params = params_proto(**params)
    #  params = params_proto(uuid="b628579d-ae3d-41f0-887e-895204190c70",
    #          name="new-disk1", target="vcd", size=1024)
    ret = call("localhost", "volume.create", params)
    print(ret)

if __name__ == "__main__":
    volume_create("b628579d-ae3d-41f0-887e-895204190c70", "new-disk1", 1024)
