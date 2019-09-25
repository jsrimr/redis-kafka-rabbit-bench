import os
from os.path import dirname
import sys
sys.path.append((dirname(sys.path[0])))
from arguments import argparser
import time
import redis
import datetime
from concurrent import futures

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
    futures_list = []
    with futures.ThreadPoolExecutor() as executor:
        executor.submit(sub, myredis, 'reader1', args.n_seconds)
        for _ in range(args.n_threads):
            future = executor.submit(pub, myredis, args.n_seconds)
            futures_list.append(future)

    pub_sum = 0
    for future in futures.as_completed(futures_list):
            result = future.result()
            pub_sum += result
            print(f'Result : {result}')

    print(f"final sum : {pub_sum}")