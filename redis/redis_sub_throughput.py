import os
from os.path import dirname
import sys
sys.path.append((dirname(sys.path[0])))
from arguments import argparser
import time
import redis
import datetime
from multiprocessing import Process

def pub(myredis,n_seconds):
    start = time.time_ns()
    cnt = 0
    n_ns = n_seconds * 1000000000
    while True:
        myredis.publish('channel', time.time())
        if time.time_ns() > start + n_ns:
            break
        cnt+=1
    print(f"pub throughput {cnt / n_seconds} msgs")

def sub(myredis, name,n_seconds):
    pubsub = myredis.pubsub()
    pubsub.subscribe(['channel'])

    start = time.time_ns()
    cnt = 0
    n_ns = n_seconds * 1000000000

    for _ in pubsub.listen():
        # print("sub started")
        cnt+=1
        if time.time_ns() > start + n_ns:
            break

    print(f"sub throughput {cnt / n_seconds}")

if __name__ == '__main__':
    args = argparser()
    myredis = redis.StrictRedis()

    sub = Process(target=sub, kwargs={'myredis': myredis, 'name': 'reader1', 'n_seconds' :args.n_seconds})
    sub.start()

    pub = Process(target=pub , args=[myredis, args.n_seconds])
    pub.start()

    procs = [sub,pub]

    for proc in procs:
        proc.join()