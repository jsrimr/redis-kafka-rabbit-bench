import argparse


def argparser():
    parser = argparse.ArgumentParser()

    parser.add_argument('--n_msg', type=int, default=1000,
                        help='how many messages to interact')

    parser.add_argument('--n_seconds', type=int, default=6,
                        help='how much time to run throughput')

    parser.add_argument('--n_threads', type=int, default=5,
                        help='how many pub/sub to run')


    args = parser.parse_args()
    return args