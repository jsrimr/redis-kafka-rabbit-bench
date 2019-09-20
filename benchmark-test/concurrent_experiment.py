from concurrent import futures
from concurrent.futures import ProcessPoolExecutor
import time

def proc1(cnt):

    for i in range(1000000):
        cnt+=1
    return cnt
def proc2(cnt):

    for i in range(1000000):
        cnt += 1
    return cnt
def proc3(cnt):

    for i in range(1000000):
        cnt += 1
    return cnt

if __name__=="__main__":
    start = time.time()

    # cnt1 = proc1(0)
    # cnt2 = proc2(0)
    # cnt3 =proc3(0)
    # print(cnt1+cnt2+cnt3)

    with ProcessPoolExecutor() as ex:
        t1 = ex.submit(proc1, 0)
        t2 = ex.submit(proc2, 0)
        t3 = ex.submit(proc3, 0)
    # print(t1.result()+t2.result()+t3.result())



    print(time.time() - start)