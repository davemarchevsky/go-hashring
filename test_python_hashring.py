#!/usr/bin/env python2

from __future__ import with_statement
import hash_ring
import sys

def main(nodes_file, keys_file):
    with open(nodes_file, 'r') as f:
        nodes = f.readlines()
    nodes = map(lambda str: str.rstrip(), nodes)

    with open(keys_file, 'r') as f:
        keys = f.readlines()
    keys = map(lambda str: str.rstrip(), keys)

    ring = hash_ring.HashRing(nodes)
    for key in keys:
        print "%s -> %s" % (key, ring.get_node(key))

if __name__ == "__main__":
    if len(sys.argv) < 3:
        print "Usage: %s nodes_file keys_file" % sys.argv[0]
        sys.exit(1)
    main(sys.argv[1], sys.argv[2])
