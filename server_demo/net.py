from protlib import *
from api import call

def detach(host, uuid, mac, vlan, type):
    class params_proto(CStruct):
        # UUID llib.UUID
        # MAC  [17]byte
        # Type netType
        # VLAN int16
        uuid =  CString(length=36)
        mac =   CString(length=17)
        vlan =  CShort()
        type =  CShort()

    params = params_proto(uuid=uuid, mac=mac, vlan=vlan, type=type)
    ret = call(host, "net.detach", params)
    print(ret)


if __name__ == '__main__':
    host = "localhost"
    uuid = "f427452a-0fd6-41ce-b31b-2872ff0bca9f"
    mac = "52:54:00:70:0d:f6"
    vlan = 69

    detach(host, uuid, mac, vlan, 0)