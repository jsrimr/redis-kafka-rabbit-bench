import sys
import time
import redis
import datetime

# def timer():
#
#     yield datetime.datetime.now()

def sub(myredis, name):
    pubsub = myredis.pubsub()
    pubsub.subscribe(['channel'])

    latency_list = []
    Append = latency_list.append
# 15초간만 listen 한다

    try:
        for item in pubsub.listen():
            latency = time.time()- float(item['data'])
            print(f"arrived time {time.time()} departure time {item['data']}, latency : {latency}")
            Append(latency)
    except:
        latency_list= latency_list[1:]
        print("latency average:" , sum(latency_list) / len(latency_list))
        # print(latency_list)

if __name__ == '__main__':
    # timer()
    myredis = redis.StrictRedis(socket_timeout=10)
    sub(myredis, 'reader1')

    # endTime = datetime.datetime.now() + datetime.timedelta(seconds=15)
    # if datetime.datetime.now() >= endTime:
    #     sys.exit()
    # print(endTime)