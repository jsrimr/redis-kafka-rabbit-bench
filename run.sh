#redis pub/sub throughput
/home/jeffrey/anaconda3/envs/benchmark/bin/python redis/redis_sub_throughput.py --n_msg=1000 > redis_pub_sub_throughput.txt
#kafka pub/sub throughput
/home/jeffrey/anaconda3/envs/benchmark/bin/python kafka/kafka_sub_throughput.py --n_msg=1000 > kafka_pub_sub_throughput.txt
#rabbit pub/sub throughput
/home/jeffrey/anaconda3/envs/benchmark/bin/python rabbitMQ/rabbit_sub_throughput.py --n_msg=1000 > rabbit_pub_sub_throughput.txt

#redis latency
/home/jeffrey/anaconda3/envs/benchmark/bin/python redis/redis_latency_onefile.py --n_seconds=5 > redis_avg_latency.txt
#kafka latency
/home/jeffrey/anaconda3/envs/benchmark/bin/python kafka/kafka_latency_onefile.py --n_seconds=5 > kafka_avg_latency.txt
#rabbit latency
/home/jeffrey/anaconda3/envs/benchmark/bin/python rabbitMQ/rabbit_latency_onefile.py --n_seconds=5 > rabbit_avg_latency.txt