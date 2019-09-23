import sys
import time
import datetime
import pika


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





if __name__ == '__main__':

    pub(10)