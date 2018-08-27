/*
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
	"encoding/binary"
	"testing"
)

func TestBasic(t *testing.T) {
	f := New(1000)
	n1 := "Bess"
	n2 := "Jane"

	f.Add(n1)

	n1b := f.Test(n1)
	n2b := f.Test(n2)

	if !n1b {
		t.Errorf("%s should be in.", n1)
	}
	if n2b {
		t.Errorf("%s should not be in.", n2)
	}
}

func TestJabberwocky(t *testing.T) {
	f := New(1000)

	jabberwocky := "'Twas brillig, and the slithy toves\n  Did gyre and gimble in the wabe:\n" +
		"All mimsy were the borogoves,\n  And the mome raths outgrabe.\n\n" +
		"'Beware the Jabberwock, my son!\n  The jaws that bite, the claws that catch!\n" +
		"Beware the Jubjub bird, and shun\n  The frumious Bandersnatch!'\n\n" +
		"He took his vorpal sword in hand:\n  Long time the manxome foe he sought --\n" +
		"So rested he by the Tumtum tree,\n  And stood awhile in thought.\n\n" +
		"And, as in uffish thought he stood,\n  The Jabberwock, with eyes of flame,\n" +
		"Came whiffling through the tulgey wood,\n  And burbled as it came!\n\n" +
		"One, two! One, two! And through and through\n  The vorpal blade went snicker-snack!\n" +
		"He left it dead, and with its head\n  He went galumphing back.\n\n" +
		"'And, has thou slain the Jabberwock?\n  Come to my arms, my beamish boy!\n" +
		"O frabjous day! Callooh! Callay!'\n  He chortled in his joy.\n\n" +
		"'Twas brillig, and the slithy toves\n  Did gyre and gimble in the wabe;"

	n1 := jabberwocky
	n2 := "'Twas brillig, and the slithy toves"

	f.Add(n1)

	n1b := f.Test(n1)
	n2b := f.Test(n2)

	if !n1b {
		t.Errorf("%s should be in.", n1)
	}
	if n2b {
		t.Errorf("%s should not be in.", n2)
	}
}

func TestBasicByte(t *testing.T) {
	f := New(1000)

	n1 := make([]byte, 4)
	n2 := make([]byte, 4)
	n3 := make([]byte, 4)
	n4 := make([]byte, 4)
	binary.BigEndian.PutUint32(n1, 100)
	binary.BigEndian.PutUint32(n2, 101)
	binary.BigEndian.PutUint32(n3, 102)
	binary.BigEndian.PutUint32(n4, 103)

	f.Add(string(n1))

	n1b := f.Test(string(n1))
	n2b := f.Test(string(n2))
	n3b := f.Test(string(n3))
	n4b := f.Test(string(n4))

	if !n1b {
		t.Errorf("%v should be in.", n1)
	}
	if n2b {
		t.Errorf("%v should not be in.", n2)
	}
	if n3b {
		t.Errorf("%v should not be in.", n3)
	}
	if n4b {
		t.Errorf("%v should not be in.", n4)
	}
}

func TestWtf(t *testing.T) {
	f := New(1000)

	f.Add("abc")

	n1b := f.Test("wtf")

	if n1b {
		t.Errorf("%s should not be in.", "wtf")
	}
}

func TestEstimatedFillRatio(t *testing.T) {
	f := New(1000)
	fr1Target := 0.015873
	fr2Target := 0.973110
	ep := 0.000005

	for i := uint(0); i < f.k; i++ {
		f.Add(string(i))
	}

	fr1 := f.EstimatedFillRatio()

	if (fr1 < fr1Target-ep) || (fr1 > fr1Target+ep) {
		t.Errorf("%f should be %f (+/- %f).", fr1, fr1Target, ep)
	}

	for i := 100; i < 1000; i++ {
		f.Add(string(i))
	}

	fr2 := f.EstimatedFillRatio()

	if (fr2 < fr2Target-ep) || (fr2 > fr2Target+ep) {
		t.Errorf("%f should be %f (+/- %f).", fr2, fr2Target, ep)
	}
}

func TestKnownFNV1aCollisions(t *testing.T) {
	f := New(1000)

	/*
	  These are known collisions with FNV-1a (32 bit) hash.
	  i.e. this unit test is implementation algorithm-specific, not Bloom filter-specific.
	  It will fail if the underlying hash function is changed.
	*/

	// c.f. https://programmers.stackexchange.com/questions/49550/which-hashing-algorithm-is-best-for-uniqueness-and-speed/145633#145633
	f.Add("costarring")
	f.Add("declinate")
	f.Add("altarage")
	f.Add("altarages")

	n1b := f.Test("liquid")
	n2b := f.Test("macallums")
	n4b := f.Test("zinkes")
	n3b := f.Test("zinke")

	if !n1b {
		t.Errorf("%s should be in.", "liquid")
	}
	if !n2b {
		t.Errorf("%s should be in.", "macallums")
	}
	if !n3b {
		t.Errorf("%s should be in.", "zinke")
	}
	if !n4b {
		t.Errorf("%s should be in.", "zinkes")
	}
}
