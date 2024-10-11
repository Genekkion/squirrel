package seeds

import (
	"sync"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

/*
The main struct used to generate, store and distribute UUIDSeeds safely
across multiple threads. It ensures that the seeds created are all
unique and that the V7 UUIDs created by the seeds will not collide.
*/
type SeedStore struct {
	/*
		Key is the V7 UUID generated by the seedBytes of the UUIDSeed,
		based on the reference timestamp within the store.

		A nil value is used to imply that the seed has been
		borrowed but not returned.
	*/
	store map[[16]byte]*UUIDSeed

	/*
		Counter for the number of UUIDSeeds which have been borrowed out
	*/
	borrowCount int

	/*
		The timestamp to be used as a reference for comparison. The theory is that
		so long as two seeds produce different V7 UUIDs given the same seed, any
		V7 UUID generated by the two seeds will never collide (at least until the
		limits of the timestamp value).
	*/
	referenceTimestamp int64

	/*
		The mutex required for multi threaded access.
	*/
	mutex sync.RWMutex
}

/*
The main function used to generate a store.
*/
func NewSeedStorage() *SeedStore {
	return &SeedStore{
		store:              map[[16]byte]*UUIDSeed{},
		borrowCount:        0,
		referenceTimestamp: 0,
		mutex:              sync.RWMutex{},
	}
}

/*
Returns the number of seeds held by this store.
*/
func (store *SeedStore) Size() int {
	store.mutex.RLock()
	size := len(store.store)
	store.mutex.RUnlock()
	return size
}

/*
Used to add seedBytes into this store. It will validate which
seeds have already been added in previously and return the
seeds which have successfully been inserted.

Note that this function should not be used in typical conditions,
mostly for restoring the state of the store from a backup.
*/
func (store *SeedStore) AddNewSeeds(seedBytesSlice ...[16]byte) [][16]byte {
	validSeeds := make([][16]byte, 0, len(seedBytesSlice))

	store.mutex.Lock()
	for _, seedBytes := range seedBytesSlice {
		V7seedBytes := GenerateV7WithTimestamp(seedBytes, store.referenceTimestamp)
		_, exists := store.store[V7seedBytes]

		if !exists {
			store.store[V7seedBytes] = FromSeedBytes(seedBytes)
			validSeeds = append(validSeeds, seedBytes)
		}
	}
	store.mutex.Unlock()

	return validSeeds
}

/*
Looks through the store and returns up to the number of seeds specified.

Note that it will not generate new seeds if there are not enough seeds.
Use GenerateNewSeeds(int) for that instead.
*/
func (store *SeedStore) BorrowSeeds(count int) []*UUIDSeed {
	seeds := make([]*UUIDSeed, 0, count)

	store.mutex.Lock()
	for seedSeed, seed := range store.store {
		if seed != nil {
			seeds = append(seeds, seed)

			// Use nil value to represent borrowed seed
			store.store[seedSeed] = nil
			store.borrowCount++
		}

		// Break out of the loop when enough seeds are found
		if len(seeds) == count {
			break
		}
	}
	store.mutex.Unlock()

	return seeds
}

/*
Returns the UUIDSeeds previously borrowed from the store. Will ignore
any seeds which did not belong to the store previously.

Returns the seeds which DID NOT belong to this store previously.

Note that after returning, the seeds should no longer be used by the
called of this function. The seeds should be gotten from the GetSeeds(int)
function instead.
*/
func (store *SeedStore) ReturnSeeds(seeds ...*UUIDSeed) []*UUIDSeed {
	invalidSeeds := []*UUIDSeed{}

	store.mutex.Lock()
	for _, seed := range seeds {
		V7seedBytes := GenerateV7WithTimestamp(seed.seedBytes,
			store.referenceTimestamp)
		_, exists := store.store[V7seedBytes]

		if exists {
			store.store[V7seedBytes] = seed
			store.borrowCount--
		} else {
			invalidSeeds = append(invalidSeeds, seed)
		}
	}

	store.mutex.Unlock()
	return invalidSeeds
}

/*
Generates the specified number of seeds. Is safe to call in most
cases unless you are saturating the limit of using the seed as a
bitmask for randomness.
*/
func (store *SeedStore) GenerateNewSeeds(count int) {
	seedsAdded := 0
	for seedsAdded < count {
		seedBatch := make([][16]byte, 0, count-seedsAdded)

		for range count {
			seedBatch = append(seedBatch, uuid.New())
		}

		// Mutex acquired during batch seed adding.
		// This potentially improves performance since we do not try 1 by 1.
		seedsAdded += len(store.AddNewSeeds(seedBatch...))
	}
}

/*
Simple wrapper around "github.com/rs/zerolog/log" for quick
debugging purposes.
*/
func (store *SeedStore) LogDebug() {
	store.mutex.RLock()
	log.Debug().
		Int("size", len(store.store)).
		Int("borrowCount", store.borrowCount).
		Int64("referenceTimestamp", store.referenceTimestamp).
		Msg("Displaying debug information of the SeedStore")
	store.mutex.RUnlock()
}
