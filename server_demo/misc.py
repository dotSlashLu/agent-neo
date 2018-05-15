from api import call
from protlib import *

# SrcUUID llib.UUID
# DstUUID llib.UUID
# DstIP 	[15]byte
# DstMac	[17]byte
# DstVLAN int16
def clone(host, srcUUID, dstUUID, dstMac, dstVLAN):
    class paramsProto(CStruct):
        srcUUID = CString(length=36)
        dstUUID = CString(length=36)
        dstMac	= CString(length=17)
        dstVLAN = CShort()
    params = paramsProto(srcUUID=srcUUID, dstUUID=dstUUID, dstMac=dstMac,
            dstVLAN=dstVLAN)
    ret = call(host, "misc.clone", params)
    print(ret)

if __name__ == "__main__":
    host = "localhost"
    srcUUID = "3bd0b56d-eddf-4a41-b328-184e0f148367"
    dstUUID = "6a9e1b70-3dfa-4028-b867-2f0f48927014"
    dstMac = "52:54:00:e0:36:d5"
    dstVLAN = 69
    clone("localhost", srcUUID, dstUUID, dstMac, dstVLAN)
