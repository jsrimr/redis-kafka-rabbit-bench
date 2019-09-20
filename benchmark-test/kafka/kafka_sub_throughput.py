import sys
import time
from kafka import KafkaProducer, KafkaConsumer
import datetime
from multiprocessing import Process
import sys

def pub(seconds = 10):
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

def sub(seconds, name):
    consumer = KafkaConsumer('test', consumer_timeout_ms=15000)

    start = datetime.datetime.now()
    cnt = 0

    for _ in consumer:
        cnt+=1
        if datetime.datetime.now() > start + datetime.timedelta(seconds=seconds):
            break

    print(f"sub throughput {cnt / seconds}")

if __name__ == '__main__':
    # timer()

    sub = Process(target=sub, kwargs={'seconds': 10, 'name': 'reader1'})
    sub.start()

    pub = Process(target=pub , args=[10, ])
    pub.start()

    procs = [sub,pub]

    for proc in procs:
        proc.join()