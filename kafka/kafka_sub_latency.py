import time
from kafka import KafkaConsumer, KafkaProducer

def sub(name):
    print("sub started")
    consumer = KafkaConsumer('test', consumer_timeout_ms=7000)
    latency_list = []
    Append = latency_list.append

    for msg in consumer:
        # print(f"Reader {name} arrived time {time.time()} and departure time {msg.timestamp}")
        # print("msg arrived")

        latency = time.time() - float(msg.value)
        print(latency)
        Append(latency)
        print(msg.value)

    print("latency average:", sum(latency_list) / len(latency_list))

if __name__ == '__main__':

    sub("reader1")