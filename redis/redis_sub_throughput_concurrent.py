import os
from os.path import dirname
import sys
sys.path.append((dirname(sys.path[0])))
from arguments import argparser
import time
import redis
import datetime
from concurret.futures import ThreadPoolExecutor

def pub(myredis,n_seconds):
    start = datetime.datetime.now()
    cnt = 0
    while True:
        myredis.publish('channel', time.time())
        if datetime.datetime.now() > start + datetime.timedelta(seconds=n_seconds):
            break
        cnt+=1
    print(f"pub throughput {cnt / n_seconds} msgs")

def sub(myredis, name,n_seconds):
    pubsub = myredis.pubsub()
    pubsub.subscribe(['channel'])

    start = datetime.datetime.now()
    cnt = 0

    for _ in pubsub.listen():
        print("sub started")
        cnt+=1
        if datetime.datetime.now() > start + datetime.timedelta(seconds=n_seconds):
            break

    print(f"sub throughput {cnt / n_seconds}")

if __name__ == '__main__':
    args = argparser()
    myredis = redis.StrictRedis()
    executor.submit(sub, [myredis, 'reader1', args.n_seconds])
    future_list = []
    with futures.ThreadPoolExecutor() as executor:
        for _ in range(args.n_threads):
            future = executor.submit(pub, [myredis, args.n_seconds])
            future_list.append(future)

    for future in futures.as_completed(futures_list):
            result = future.result()
            print(f'Result : {result}')
            
    procs = [sub,pub]

    for proc in procs:
        proc.join()