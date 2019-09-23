import time
# from concurrent.futures import ProcessPoolExecutor
from multiprocessing import Process
import redis


def pub(myredis):
    for n in range(1000):
        myredis.publish('channel', time.time())
        # time.sleep(5)
    print("published")


def sub(myredis, name):
    pubsub = myredis.pubsub()
    pubsub.subscribe(['channel'])

    latency_list = []
    Append = latency_list.append

    try:
        for item in pubsub.listen():
            # print(f"Reader {name} arrived time {time.time()} and departure time {item['data']}")
            latency = time.time() - float(item['data'])
            Append(latency)
    except:
        latency_list = latency_list[1:]
        print("latency average:" , sum(latency_list) / len(latency_list))

if __name__ == '__main__':
    myredis = redis.StrictRedis(socket_timeout=10)


    proc2 = Process(target=sub, kwargs={'myredis': myredis, 'name': 'reader1'})
    proc2.start()
    # proc3 = Process(target=sub, kwargs={'myredis': myredis, 'name': 'reader2'})
    # proc3.start()
    # proc4 = Process(target=sub, kwargs={'myredis': myredis, 'name': 'reader3'})
    # proc4.start()
    proc1 = Process(target=pub, args=[myredis, ])
    proc1.start()

    procs = [proc1,proc2]
    for proc in procs:
        proc.join()
    