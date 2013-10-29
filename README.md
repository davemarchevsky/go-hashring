go-hashring
===========

Hashring implementation in Go. Based off of [this](https://pypi.python.org/pypi/hash_ring/1.3.1) python implementation and intended to function identically.

To use HashRing, instantiate via the `New` or `NewWithWeights` constructors and assign keys to nodes via 
the `GetNode` function.

Example:

```go
hosts := []string{"10.0.0.1", "10.0.0.2", "10.0.0.3", "10.0.0.4", "10.0.0.5"}
ring := hashring.New(hosts)
host := ring.getNode("key to hash on")
```
