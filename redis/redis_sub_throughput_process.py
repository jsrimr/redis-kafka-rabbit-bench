import os
from os.path import dirname
import sys
sys.path.append((dirname(sys.path[0])))
from arguments import argparser
import time
import redis
import datetime
from multiprocessing import Process , Value

def pub(myredis,n_seconds, value):
    start = time.time_ns()
    cnt = 0
    n_ns = n_seconds * 1000000000
    due = start + n_ns
    while True:
        myredis.publish('channel', time.time())
        if time.time_ns() > due:
            break
        cnt+=1
    result = cnt / n_seconds
    with value.get_lock():
        value.value += result
    # print(f"pub throughput {result} msgs")
    # return result

def sub(myredis, name,n_seconds):
    # print("sub started")
    pubsub = myredis.pubsub()
    pubsub.subscribe(['channel'])

    start = time.time_ns()
    cnt = 0
    n_ns = n_seconds * 1000000000
    due = start + n_ns
    for _ in pubsub.listen():
        cnt+=1
        if time.time_ns() > due:
            break

    print(f"sub throughput {cnt / n_seconds}")

if __name__ == '__main__':
    args = argparser()
    output = Value('f',0)
    myredis = redis.StrictRedis()
    pub_list = []

    sub = Process(target=sub, kwargs={'myredis': myredis, 'name': 'reader1', 'n_seconds' :args.n_seconds})
    sub.start()

    for _ in range(args.n_threads):
        pub_proc = Process(target=pub , args=[myredis, args.n_seconds, output])
        pub_list.append(pub_proc)
        pub_proc.start()
        
    procs = [sub]+pub_list
    for proc in procs:
        proc.join()
    print(f'pub_sum : {output.value}') 