import os
from os.path import dirname
import sys
sys.path.append((dirname(sys.path[0])))
from arguments import argparser
import time
import redis
import datetime
from multiprocessing import Pool #, Process , QUEUE

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
    print("sub started")
    pubsub = myredis.pubsub()
    pubsub.subscribe(['channel'])

    start = time.time_ns()
    cnt = 0
    n_ns = n_seconds * 1000000000

    for _ in pubsub.listen():
        cnt+=1
        if time.time_ns() > start + n_ns:
            break

    print(f"sub throughput {cnt / n_seconds}")

if __name__ == '__main__':
    args = argparser()
    myredis = redis.StrictRedis()
    res_list = []
    with Pool() as p:
        sub_res = p.apply_async(sub, myredis, 'reader1', args.n_seconds)
        for _ in range(args.n_threads):
            res_obj = p.apply_async(pub, myredis, args.n_seconds)
            res_list.append(res_obj)

    pub_sum =0 
    for obj in res_list:
        pub_sum+=obj.get()
    print(f"final sum : {pub_sum}")