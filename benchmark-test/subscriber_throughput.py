import sys
import time
import redis
import datetime
from multiprocessing import Process
import sys

def pub(myredis):
    start = datetime.datetime.now()
    cnt = 0
    while True:
        myredis.publish('channel', time.time())
        if datetime.datetime.now() > start + datetime.timedelta(seconds=10):
            break
        cnt+=1
    print(f"published {cnt / 10} msgs")

def sub(myredis, name):
    pubsub = myredis.pubsub()
    pubsub.subscribe(['channel'])

    start = datetime.datetime.now()
    cnt = 0

    for _ in pubsub.listen():
        cnt+=1
        if datetime.datetime.now() > start + datetime.timedelta(seconds=10):
            break

    print(f"sub throughput {cnt / 10}")

if __name__ == '__main__':
    # timer()
    myredis = redis.StrictRedis(socket_timeout=15)

    sub = Process(target=sub, kwargs={'myredis': myredis, 'name': 'reader1'})
    sub.start()

    pub = Process(target=pub , args=[myredis, ])
    pub.start()

    procs = [sub,pub]

    for proc in procs:
        proc.join()