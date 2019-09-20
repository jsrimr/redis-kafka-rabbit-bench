import time
import datetime
import redis

def pub(myredis,):
    start = datetime.datetime.now()
    cnt = 0
    while True:
        myredis.publish('channel',time.time())
        cnt+=1
        if datetime.datetime.now() > start + datetime.timedelta(seconds=10):
            break
    print(f"publisher throughput {cnt / 10}")



if __name__ == '__main__':
    myredis = redis.StrictRedis()
    pub(myredis)
    # pub_with_pipe(myredis)