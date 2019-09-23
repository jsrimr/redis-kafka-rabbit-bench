import pika
import time

cnt = 0

def sub(n_msg=1000):
    latency_list = []
    Append = latency_list.append

    def callback(ch, method, properties, body):
        global cnt
        cnt += 1
        latency = time.time() - float(body)
        Append(latency)
        if cnt == n_msg:
            channel.stop_consuming()

    connection = pika.BlockingConnection(pika.ConnectionParameters('localhost'))
    channel = connection.channel()
    channel.queue_declare(queue='hello')  # queue 생성
    channel.basic_consume(queue='hello', auto_ack=True, on_message_callback=callback)
    print(' [*] Waiting for messages. To exit press CTRL+C')
    channel.start_consuming()
    return sum(latency_list) / len(latency_list)


if __name__ == '__main__':
    avg_latency = sub(1000)
    print(avg_latency)