import sys
import time
import datetime
import pika
from multiprocessing import Process

def pub(n_sec):
    start = datetime.datetime.now()
    cnt = 0
    connection = pika.BlockingConnection(pika.ConnectionParameters('localhost'))
    channel = connection.channel()
    channel.queue_declare(queue='hello')  # queue 생성

    def pub_():
        channel.basic_publish(exchange='', routing_key='hello', body=str(time.time()),
                              properties=pika.BasicProperties(timestamp=int(time.time())))

    while True:
        pub_()
        if datetime.datetime.now() > start + datetime.timedelta(seconds=n_sec):
            break
        cnt += 1
    print(f"published {cnt / n_sec} msgs")

def sub(n_sec, name):

    start = datetime.datetime.now()
    sub_cnt = 0

    def callback(ch, method, properties, body):
        nonlocal sub_cnt
        sub_cnt += 1
        if datetime.datetime.now() > start + datetime.timedelta(seconds=n_sec):
            channel.stop_consuming()

    connection = pika.BlockingConnection(pika.ConnectionParameters('localhost'))
    channel = connection.channel()
    channel.queue_declare(queue='hello')  # queue 생성
    channel.basic_consume(queue='hello', auto_ack=True, on_message_callback=callback)
    print(' [*] Waiting for messages. To exit press CTRL+C')
    channel.start_consuming()



    print(f"sub throughput {sub_cnt / n_sec}")



if __name__ == '__main__':
    sub_proc = Process(target=sub, args=[10, 'reader1'])
    sub_proc.start()

    pub_proc = Process(target=pub, args=[10,])
    pub_proc.start()

    procs = [sub_proc, pub_proc]
    for proc in procs:
        proc.join()