import os
from os.path import dirname
import sys
sys.path.append((dirname(sys.path[0])))
from arguments import argparser
import time
import datetime
import pika
from multiprocessing import Process

def pub(n_sec, topic):
    start = time.time_ns()
    cnt = 0
    connection = pika.BlockingConnection(pika.ConnectionParameters('localhost'))
    channel = connection.channel()
    channel.queue_declare(queue=topic)  # queue 생성
    n_ns = n_sec * 1000000000
    
    def pub_():
        channel.basic_publish(exchange='', routing_key=topic, body=str(time.time()),
                              properties=pika.BasicProperties(timestamp=int(time.time())))

    while True:
        pub_()
        if time.time_ns() > start + n_ns:
            break
        cnt += 1
    print(f"pub throughput {cnt / n_sec} msgs")

def sub(n_sec, topic):

    start = time.time_ns()
    sub_cnt = 0
    n_ns = n_sec * 1000000000

    def callback(ch, method, properties, body):
        nonlocal sub_cnt
        sub_cnt += 1
        if time.time_ns() > start + n_ns:
            channel.stop_consuming()

    connection = pika.BlockingConnection(pika.ConnectionParameters('localhost'))
    channel = connection.channel()
    channel.queue_declare(queue=topic)  # queue 생성
    channel.basic_consume(queue=topic, auto_ack=True, on_message_callback=callback)
    # print(' [*] Waiting for messages. To exit press CTRL+C')
    channel.start_consuming()

    print(f"sub throughput {sub_cnt / n_sec}")



if __name__ == '__main__':
    args = argparser()
    sub_proc = Process(target=sub, args = [args.n_seconds,"hello"])
    sub_proc.start()

    pub_proc = Process(target=pub, args = [args.n_seconds, "hello"])
    pub_proc.start()

    procs = [sub_proc, pub_proc]
    for proc in procs:
        proc.join()
