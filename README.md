Bloom Filter
============

[![Master Build Status](https://secure.travis-ci.org/Echelon9/bloomfilter-go.png?branch=master)](https://travis-ci.org/Echelon9/bloomfilter-go?branch=master)
[![Coverage Status](https://coveralls.io/repos/github/Echelon9/bloomfilter-go/badge.svg?branch=master)](https://coveralls.io/github/Echelon9/bloomfilter-go?branch=master)

This Go implementation of a Bloom filter uses the non-cryptographic
[Fowler–Noll–Vo hash function][1] for speed.

A Bloom filter is a space-efficient probabilistic data structure, conceived by
Burton Howard Bloom in 1970, that is used to test whether an element is a
member of a set. False positive matches are possible, but false negatives are
not, thus a Bloom filter has a 100% recall rate [(Bloom, 1970)][2]

[1]: http://isthe.com/chongo/tech/comp/fnv/
[2]: https://dx.doi.org/10.1145%2F362686.362692

## Installation

```bash
go get -u github.com/Echelon9/bloomfilter-go
```

## Usage

This implementation accepts keys for adding and testing as `string`

```go
import (
	...
	"github.com/Echelon9/bloomfilter-go"
	...
)

set := bloomfilter.New(1000) // Size
set.Add("Hello, world!")
set.Test("Hello, world!") // => true
set.Test("String Not-appearing-in-this-set") // => false
```

## Optimizations

The Bloom filter implemented in this library is the "Standard" variant.

There are no current plans to implement future modifications for additional
variants, but patches are welcome!

However, this Bloom filter implementation does adopt an idea from
[(Kirsh and Mitzenmachner, 2006)][3] that two hash functions can simulate an
arbitrary number of additional hash functions without losing functionality. 

>   Less hashing, same performance: Building a better bloom filter (2006)
>   Adam Kirsch, Michael Mitzenmacher
>
>   "A standard technique from the hashing literature is to use two hash functions h1(x)
>   and h2(x) to simulate additional hash functions of the form gi(x) = h1(x) + ih2(x).
>   We demonstrate that this technique can be usefully applied to Bloom filters and
>   related data structures. Specifically, only two hash functions are necessary to
>   effectively implement a Bloom filter without any loss in the asymptotic false
>   positive  probability. This leads to less computation and potentially less need
>   for randomness in practice."

[3]: http://citeseerx.ist.psu.edu/viewdoc/summary?doi=10.1.1.152.579

## Running all tests

Before committing the code, please check if it passes all tests using:
```bash
go test -v ./...
```