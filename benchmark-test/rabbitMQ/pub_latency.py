import pika
import time

def pub():
    channel.basic_publish(exchange='',routing_key='hello',body=str(time.time()), properties=pika.BasicProperties(timestamp=int(time.time())))
    print(" [x] Sent 'Hello World!'")

if __name__=='__main__':
    connection = pika.BlockingConnection(pika.ConnectionParameters('localhost'))
    channel = connection.channel()
    channel.queue_declare(queue='hello')  #queue 생성
    for _ in range(1000):
        pub()
    connection.close()