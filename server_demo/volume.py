from protlib import *

from api import call


def volume_create(host, uuid, name, size):
    # UUID    [36]byte    // vm uuid
    # Name    [32]byte    // random str
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
    ret = call(host, "volume.create", params)
    print(ret)


def volume_attach(host, uuid, name, target):
    class params_proto(CStruct):
        # UUID   llib.UUID // vm uuid
        # Name   [32]byte // random str
        # Target [3]byte  // vdb? vdc?
        uuid    = CString(length=36)
        name    = CString(length=32)
        target  = CString(length=3)
    params = {
        "uuid": uuid,
        "name": name,
        "target": target
    }
    ret = call(host, "volume.attach", params_proto(**params))
    print(ret)


def volume_detach(host, uuid, name, target):
    class params_proto(CStruct):
        # UUID   llib.UUID // vm uuid
        # Name   [32]byte // random str
        # Target [3]byte  // vdb? vdc?
        uuid    = CString(length=36)
        name    = CString(length=32)
        target  = CString(length=3)
    params = {
        "uuid": uuid,
        "name": name,
        "target": target
    }
    ret = call(host, "volume.detach", params_proto(**params))
    print(ret)


def volume_delete(host, uuid, name):
    class params_proto(CStruct):
        # UUID   llib.UUID // vm uuid
        # Name   [32]byte // random str
        uuid    = CString(length=36)
        name    = CString(length=32)
    params = {
        "uuid": uuid,
        "name": name,
    }
    ret = call(host, "volume.delete", params_proto(**params))
    print(ret)


if __name__ == "__main__":
    host = "localhost"
    uuid = "a5d95464-8e9a-7949-088d-99f889bf630c"
    name = "new-disk1"
    target = "vdd"
    volume_create(host, uuid, name, 1024)
    volume_attach(host, uuid, name, target)
    volume_detach(host, uuid, name, target)
    volume_delete(host, uuid, name)
