import sys
import time
import datetime
import pika

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
    # timer()
    sub(10, 'reader1')
