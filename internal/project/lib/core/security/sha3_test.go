package security

import "testing"

func TestSha3_256(t *testing.T) {
	hash := Sha3Variant256("hello")

	if hash != "3338be694f50c5f338814986cdf0686453a888b84f424d792af4b9202398f392" {
		t.Fatal("hash should be 3338be694f50c5f338814986cdf0686453a888b84f424d792af4b9202398f392")
	}
}

func TestSha3_512(t *testing.T) {
	hash := Sha3Variant512("hello")

	if hash != "75d527c368f2efe848ecf6b073a36767800805e9eef2b1857d5f984f036eb6df891d75f72d9b154518c1cd58835286d1da9a38deba3de98b5a53e5ed78a84976" {
		t.Fatal("hash should be 75d527c368f2efe848ecf6b073a36767800805e9eef2b1857d5f984f036eb6df891d75f72d9b154518c1cd58835286d1da9a38deba3de98b5a53e5ed78a84976")
	}
}
