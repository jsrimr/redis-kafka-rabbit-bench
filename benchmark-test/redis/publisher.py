import time
import redis

def pub(myredis,s=1000):
    # start = time.time()
    for n in range(s):
        myredis.publish('channel',time.time())


    # latency = time.time() - start
    # print( 1/ (latency / s))

# def pub_with_pipe(myredis,s=100000):
#     pipe = myredis.pipeline(transaction=False)
#     start = time.time()
#     for n in range(s):
#         pipe.publish('channel', time.time())
#
#     pipe.execute()
#     duration = time.time() - start
#     print(1 / (duration / s))
#


if __name__ == '__main__':
    myredis = redis.StrictRedis()
    pub(myredis)
    # pub_with_pipe(myredis)