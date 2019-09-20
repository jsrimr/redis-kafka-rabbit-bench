import time
from kafka import KafkaConsumer, KafkaProducer
from multiprocessing import Process

# from concurrent.futures import ProcessPoolExecutor
from multiprocessing import Process


def pub(n_msg):
    producer = KafkaProducer()
    key = bytes('key', encoding='utf-8')
    value = bytes(str(time.time()), encoding='utf-8')
    for n in range(n_msg):
        producer.send("test", value=value, key=key)

    print("published")


def sub(name):
    print("sub started")
    consumer = KafkaConsumer('test', consumer_timeout_ms=5000)
    latency_list = []
    Append = latency_list.append

    for msg in consumer:
        # print(f"Reader {name} arrived time {time.time()} and departure time {msg.timestamp}")
        # print("msg arrived")
        try:
            latency = time.time() - float(msg.value)
            print(latency)
            Append(latency)
        except:
            print(msg.value)

    print("latency average:", sum(latency_list) / len(latency_list))


if __name__ == '__main__':


    proc1 = Process(target=sub, kwargs={'name': 'reader1'})
    proc1.start()

    time.sleep(2)
    proc2 = Process(target=pub, kwargs={"n_msg":100})
    proc2.start()

    procs = [proc1, proc2]
    for proc in procs:
        proc.join()


