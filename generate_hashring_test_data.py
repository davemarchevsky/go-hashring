#!/usr/bin/env python2

from __future__ import with_statement
import sys
import random

DICT_PATH="/usr/share/dict/words"
NUM_HOSTS = 100
NUM_KEYS = 1000000
MAX_WORDS_PER_WORDSTRING = 3

def main(nodes_file, keys_file):
    with open(DICT_PATH) as f:
        lines = f.readlines()
    lines = map(lambda str: str.rstrip(), lines)

    with open(nodes_file, 'w') as f:
        for i in xrange(NUM_HOSTS):
            f.write(gen_random_wordstring(lines))

    with open(keys_file, 'w') as f:
        for i in xrange(NUM_KEYS):
            f.write(gen_random_wordstring(lines))

def gen_random_wordstring(lines):
    words_per_wordstring = random.randint(1, MAX_WORDS_PER_WORDSTRING)
    wordstring = ''.join(random.sample(lines, words_per_wordstring)) + "\n"
    return wordstring

if __name__ == "__main__":
    if len(sys.argv) < 3:
        print "Usage: generate_hashring_data.py nodes_output_filename keys_output_filename"
        sys.exit(1)
    main(sys.argv[1], sys.argv[2])
