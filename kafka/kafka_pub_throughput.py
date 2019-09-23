import os
from os.path import dirname
import sys
sys.path.append((dirname(sys.path[0])))
from arguments import argparser
import time
import datetime
from kafka import KafkaProducer

def pub():
    producer = KafkaProducer()
    key = bytes('key', encoding='utf-8')
    start = datetime.datetime.now()
    cnt = 0
    while True:
        value = bytes(str(time.time()), encoding='utf-8')
        producer.send("test", value=value, key=key)
        cnt+=1
        if datetime.datetime.now() > start + datetime.timedelta(seconds=10):
            break
    print(f"publisher throughput {cnt / 10}")



if __name__ == '__main__':
    args = argparser()
    pub()
    # pub_with_pipe(myredis)