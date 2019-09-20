import time
from kafka import KafkaProducer

def pub(n_msg):
    producer = KafkaProducer()
    key = bytes('key', encoding='utf-8')
    for n in range(n_msg):
        value = bytes(str(time.time()), encoding='utf-8')
        producer.send("test", value=value, key=key)

    print("published")

if __name__ == '__main__':

    pub(1000)
    