import os
from os.path import dirname
import sys
sys.path.append((dirname(sys.path[0])))
from arguments import argparser
import redis
from multiprocessing import Process
import time


def pub(myredis, n_msg):
    for n in range(n_msg):
        myredis.publish('channel', time.time())
    print("published")


def sub(myredis, name, n_msg):
    pubsub = myredis.pubsub()
    pubsub.subscribe(['channel'])

    latency_list = []
    Append = latency_list.append

    cnt = 0
    for item in pubsub.listen():
        # print(f"Reader {name} arrived time {time.time()} and departure time {item['data']}")
        if cnt==0: # 첫 메시지는 엄청 오래 걸리는거로 나와서 예외처리. 왜인지는 불분명
            cnt+=1
            continue 
        latency = time.time() - float(item['data'])
        print(latency)
        Append(latency)
        cnt += 1
        if cnt == n_msg:
            break
    print(
        f'Latency Average for {n_msg}: {sum(latency_list) / len(latency_list)}')


if __name__ == '__main__':
    args = argparser()
    myredis = redis.StrictRedis()
    proc2 = Process(target=sub, kwargs={
                    'myredis': myredis, 'name': 'reader1', 'n_msg': args.n_msg})
    proc2.start()
    # proc3 = Process(target=sub, kwargs={'myredis': myredis, 'name': 'reader2'})
    # proc3.start()
    # proc4 = Process(target=sub, kwargs={'myredis': myredis, 'name': 'reader3'})
    # proc4.start()
    proc1 = Process(target=pub, args=[myredis, args.n_msg])
    proc1.start()

    procs = [proc1, proc2]
    for proc in procs:
        proc.join()
