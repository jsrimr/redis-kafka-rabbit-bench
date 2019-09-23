import os
from os.path import dirname
import sys
sys.path.append((dirname(sys.path[0])))
from arguments import argparser
import time
from kafka import KafkaProducer, KafkaConsumer
import datetime
from multiprocessing import Process

def pub(seconds):
    producer = KafkaProducer()
    key = bytes('key', encoding='utf-8')
    start = datetime.datetime.now()
    cnt = 0
    while True:
        value = bytes(str(time.time()), encoding='utf-8')
        producer.send("test", value=value, key=key)
        cnt+=1
        if datetime.datetime.now() > start + datetime.timedelta(seconds=seconds):
            break
    print(f"publisher throughput {cnt / seconds}")

def sub(seconds, ):
    consumer = KafkaConsumer('test')

    start = datetime.datetime.now()
    cnt = 0

    for _ in consumer:
        cnt+=1
        if datetime.datetime.now() > start + datetime.timedelta(seconds=seconds):
            break

    print(f"sub throughput {cnt / seconds}")

if __name__ == '__main__':
    # timer()
    args = argparser()
    sub = Process(target=sub, kwargs={'seconds': args.n_seconds,})
    sub.start()

    pub = Process(target=pub , args=[args.n_seconds, ])
    pub.start()

    procs = [sub,pub]

    for proc in procs:
        proc.join()