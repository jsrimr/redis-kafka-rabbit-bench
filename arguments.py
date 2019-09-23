import argparse


def argparser():
    parser = argparse.ArgumentParser()

    parser.add_argument('--n_msg', type=int, default=1000,
                        help='how many messages to interact')

    parser.add_argument('--n_seconds', type=int, default=10,
                        help='how much time to run throughput')


    args = parser.parse_args()
    return args