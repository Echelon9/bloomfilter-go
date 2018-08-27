/*
bloomfilter.go - A simple Bloom filter, written in Go
Copyright (C) 2018  Rhys Kidd

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package bloomfilter

import (
	"hash/fnv"
	"math"
)

// Interface provides a template for a Bloom filter implementation
type Interface interface {
	Add(item string)             // Adds item to the BloomFilter
	Test(item string) bool       // Test if item is maybe in the BloomFilter
	EstimatedFillRatio() float64 // Estimate fill ratio of the BloomFilter
}

// BloomFilter probabilistic data structure definition
type BloomFilter struct {
	bitset []uint32 // The Bloom filter bitset
	m      uint     // Number of bits in the Bloom filter bitset
	k      uint     // Number of effective hashing functions
	n      uint     // Number of elements in the filter
	b      uint     // Number of buckets
}

// New returns a new BloomFilter object
func New(size uint) *BloomFilter {
	return &BloomFilter{
		bitset: make([]uint32, size/32),
		m:      size,
		k:      4, // 4 effective hash functions
		n:      uint(0),
		b:      size / 32,
	}
}

// Add item to the Bloom filter
func (f *BloomFilter) Add(item string) {
	hash1 := hashFNV1a([]byte(item))
	hash2 := hashFNV1a([]byte(item))

	for i := uint(0); i < f.k; i++ {
		index := uint((hash1 + hash2*uint32(i))) % f.m
		bucket := (index / 32) % f.b
		f.bitset[bucket] |= (1 << uint(index%32))
	}

	f.n++
}

// Test returns true if the item is in the Bloom filter, false otherwise.
// If true, the result might be a false positive. If false, the data
// is definitely not in the set.
func (f *BloomFilter) Test(item string) (exists bool) {
	hash1 := hashFNV1a([]byte(item))
	hash2 := hashFNV1a([]byte(item))
	exists = false

	for i := uint(0); i < f.k; i++ {
		index := uint((hash1 + hash2*uint32(i))) % f.m
		bucket := (index / 32) % f.b
		if f.bitset[bucket]&(1<<uint(index%32)) != 0 {
			exists = true
			break
		}
	}

	return
}

// EstimatedFillRatio returns the estimated fill ratio of the Bloom filter
func (f *BloomFilter) EstimatedFillRatio() float64 {
	return 1 - math.Exp((-float64(f.n)*float64(f.k))/float64(f.m))
}

/*
 See http://willwhim.wordpress.com/2011/09/03/producing-n-hash-functions-by-hashing-only-once/
 and http://citeseerx.ist.psu.edu/viewdoc/summary?doi=10.1.1.152.579

    Less hashing, same performance: Building a better bloom filter (2006)
    Adam Kirsch, Michael Mitzenmacher

    "A standard technique from the hashing literature is to use two hash functions h1(x)
    and h2(x) to simulate additional hash functions of the form gi(x) = h1(x) + ih2(x).
    We demonstrate that this technique can be usefully applied to Bloom filters and
    related data structures. Specifically, only two hash functions are necessary to
    effectively implement a Bloom filter without any loss in the asymptotic false
    positive  probability. This leads to less computation and potentially less need
    for randomness in practice."
*/
func hashFNV1a(input []byte) uint32 {
	hash := fnv.New32a()

	hash.Write(input)
	val32 := hash.Sum32()

	return val32
}
