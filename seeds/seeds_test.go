package seeds

import "testing"

func Test6(test *testing.T) {
	store := NewSeedStorage()
	store.GenerateNewSeeds(1)
	seeds := store.BorrowSeeds(1)
	if len(seeds) != 1 {
		test.Errorf("Test6 - something went wrong when borrowing seeds")
	}

	seed := seeds[0]

	for range 1_000 {
		id1 := seed.GenerateV7()
		bytes1 := id1[:]
		id2 := seed.GenerateV7()
		bytes2 := id2[:]

		flag := false
		for i := range id1 {
			if bytes1[i] == bytes2[i] {
				flag = true
				break
			}
		}
		if !flag {
			test.Errorf("Test6 - Repated uuid found, id1: %s, id2: %s", id1.String(), id2.String())
		}
	}
}

func Test7(test *testing.T) {
	store := NewSeedStorage()
	store.GenerateNewSeeds(1)
	seeds := store.BorrowSeeds(1)
	if len(seeds) != 1 {
		test.Errorf("Test7 - something went wrong when borrowing seeds")
	}

	seed := seeds[0]
	seedMap := map[[16]byte]struct{}{}

	for range 1_000_000 {
		id := seed.GenerateV7()
		_, exists := seedMap[id]
		if exists {
			test.Errorf("Test6 - Repated uuid found, id: %s", id.String())
		}
		seedMap[id] = struct{}{}
	}
}
