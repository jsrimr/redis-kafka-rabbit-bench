import os
from os.path import dirname
import sys
sys.path.append((dirname(sys.path[0])))
from arguments import argparser
import time
import pika
import multiprocessing
from multiprocessing import Process #, Queue

def sub(n_msg=1000, ):
    latency_list = []
    Append = latency_list.append
    cnt = 0

    def callback(ch, method, properties, body):
        nonlocal cnt
        cnt += 1
        latency = time.time_ns() - float(body)
        Append(latency)
        if cnt == n_msg:
            channel.stop_consuming()

    connection = pika.BlockingConnection(
        pika.ConnectionParameters('localhost'))
    channel = connection.channel()
    channel.queue_declare(queue='hello')  # queue 생성
    channel.basic_consume(queue='hello', auto_ack=True,
                          on_message_callback=callback)
    # print(' [*] Waiting for messages. To exit press CTRL+C')
    channel.start_consuming()
    channel.close()
    connection.close()
    print(f'Latency Average for {n_msg}: {sum(latency_list) / len(latency_list)}')
    # output.put(sum(latency_list) / len(latency_list))

def pub(n_msg):
    connection = pika.BlockingConnection(pika.ConnectionParameters('localhost'))
    channel = connection.channel()
    channel.queue_declare(queue='hello')  #queue 생성
    def pubish():
        channel.basic_publish(exchange='',routing_key='hello',body=str(time.time_ns()), \
                                properties=pika.BasicProperties(timestamp=int(time.time_ns())))
        # print(" [x] Sent 'Hello World!'")
    for _ in range(args.n_msg):
        pubish()
    channel.close() # for flushing
    connection.close()

if __name__ == '__main__':
    args = argparser()

    # output = Queue()

    sub_proc = Process(target=sub, args = [args.n_msg])
    sub_proc.start()

    pub_proc = Process(target=pub, args = [args.n_msg])
    pub_proc.start()

    procs = [sub_proc, pub_proc]
    for proc in procs:
        proc.join()
    # print(output.get())
    # output.close()
    
