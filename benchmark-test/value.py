from multiprocessing import Value, Pool
def count1(i):
    val.value += 1

for run in range(3):
    val = Value('i', 0)
    with Pool(processes=4) as pool:
        pool.map(count1, range(1000))

    print(val.value)