import os
from os.path import dirname
import sys
sys.path.append((dirname(sys.path[0])))
from arguments import argparser
import time
from kafka import KafkaConsumer, KafkaProducer
from multiprocessing import Process
# from concurrent.futures import ProcessPoolExecutor


def pub(n_msg):
    producer = KafkaProducer()
    key = bytes('key', encoding='utf-8')
    value = bytes(str(time.time()), encoding='utf-8')
    for n in range(n_msg):
        producer.send("test", value=value, key=key)

    # print("published")


def sub(n_msg,consumer):
    # print("sub started")
    latency_list = []
    Append = latency_list.append

    cnt = 0
    for msg in consumer:
        latency = time.time() - float(msg.value)
        # print(latency , cnt)
        Append(latency)
        cnt+=1
        if cnt == n_msg:
            break
    print("latency average:", sum(latency_list) / len(latency_list))

if __name__ == '__main__':
    args = argparser()
    consumer = KafkaConsumer('test', )
    proc1 = Process(target=sub, kwargs={'n_msg': args.n_msg, "consumer": consumer})
    proc1.start()
    time.sleep(2) # 약간의 딜레이가 있어야 작동함

    proc2 = Process(target=pub, kwargs={"n_msg": 2*args.n_msg})
    proc2.start()

    procs = [proc1, proc2]
    for proc in procs:
        proc.join()


