import socket
import datetime
import time

def pub(ip="127.0.0.1", port=4242):
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    sock.connect((ip, port))
    cnt = 0
    start = time.time_ns()
    while True:
        msg = b"PUB 23482093 PP %d\r\n" % start
        sock.send(msg)
        data = sock.recv(35).decode('utf-8')
        if data[:2] == "OK":
            cnt += 1
        if time.time_ns() > start + 10000000000:
            break
    print(f"publisher throughput {cnt / 10}")

def sub(ip="127.0.0.1", port=4242):
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    sock.connect((ip, port))
    cnt = 0
    start = datetime.datetime.now()

    # subscribe a channel
    msg = b"SUB 23482093 PP\r\n"
    sock.send(msg)
    data = sock.recv(35).decode('utf-8')
    if data[:2] != "OK":
        print("subs failed")
        return
    while True:
        data = sock.recv(1024).decude('utf-8')
        # match timestamps
        cnt += 1
    print(f"subscriber throughput")


if __name__ == '__main__':
    # sub(ip="127.0.0.1", port=4242)
    pub(ip="127.0.0.1", port=4242)
