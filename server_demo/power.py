from protlib import *
from api import call

def power(host, op, uuid):
    class params_proto(CStruct):
        uuid =  CString(length=36)
        op =    CString(length=10)

    params = params_proto(op=op, uuid=uuid)
    ret = call(host, "power." + op, params)
    print(ret)


if __name__ == '__main__':
    host = "localhost"
    uuid = "a5d95464-8e9a-7949-088d-99f889bf630c"

    power(host, "suspend", uuid)
    power(host, "resume", uuid)