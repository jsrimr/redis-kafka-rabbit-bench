import os
from os.path import dirname
import sys
sys.path.append((dirname(sys.path[0])))
from arguments import argparser
import time
import datetime
import redis

def pub(myredis, seconds):
    start = datetime.datetime.now()
    cnt = 0
    while True:
        myredis.publish('channel',time.time())
        cnt+=1
        if datetime.datetime.now() > start + datetime.timedelta(seconds=n_seconds):
            break
    print(f"publisher throughput {cnt / n_seconds}")

if __name__ == '__main__':
    args = argparser()
    myredis = redis.StrictRedis()
    pub(myredis, args.n_seconds)
    # pub_with_pipe(myredis)