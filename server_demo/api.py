import socket
from protlib import *

class header_proto(CStruct):
    magic       = CInt(default=0x53b)
    fn_name     = CString(length=32)
    param_len   = CInt()


def call(host, fn_name, params):
    params = params.serialize()
    header = header_proto(fn_name=fn_name, param_len=len(params))
    sock = socket.socket()
    sock.connect(("localhost", 18103))
    sock.send(header.serialize())
    sock.send(params)
    ret = sock.recvfrom(1024)
    sock.close()

