package seeds

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

const (
	// Required for generating V7 UUIDs.
	NANOSECONDS_PER_MILLISECONDS = 1_000_000
)

// The main struct for generating unique UUIDs based on seed and timestamp,
// adapted from the V7 UUID standard. This struct should not be created on
// its own, but rather by FromSeedBytes([16]byte) or NewUUIDSeed().
//
// The better way to generate one is to use a SeedStore.
type UUIDSeed struct {

	// The seedBytes used for generating the UUIDs. Further implementations
	// may reduce the size of seedBytes since not all bits are required for
	// randomness.
	seedBytes [16]byte

	// Timestamp for monotonic increments of the timestamp embedded
	// in the generated V7 UUID
	lastV7Time int64

	// The mutex lock required for safely generating the V7 UUIDs across multiple
	// threads.
	mutex sync.Mutex
}

// This should not be used in most cases. This is only for when the seed bytes are
// received from another source, such as from a different machine, and sent as bytes.
//
// Note that the timestamp will be set to the current time.
func FromSeedBytes(seedBytes [16]byte) *UUIDSeed {
	return &UUIDSeed{
		seedBytes:  seedBytes,
		lastV7Time: time.Now().UnixNano(),
		mutex:      sync.Mutex{},
	}
}

// Returns a new random uuidSeed struct. May error due to
// the underlying call to uuid.NewRandom() from
// "github.com/google/uuid".
func NewUUIDSeed() (*UUIDSeed, error) {
	baseId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	return &UUIDSeed{
		seedBytes:  [16]byte(baseId[:]),
		lastV7Time: time.Now().UnixNano(),
	}, nil
}

// Returns a new copy of the seed bytes.
//
// Note that the built-in seed bytes is considered immutable after initialisation.
func (seed *UUIDSeed) GetSeedBytes() [16]byte {
	return seed.seedBytes
}

// Generates a V7 UUID based on the current time, with randomness based on the
// seed bytes specified within the seed struct.
func (seed *UUIDSeed) GenerateV7() uuid.UUID {
	return uuid.UUID(seed.GenerateV7Bytes())
}

// Similar to GenerateV7(), but returns as bytes instead
func (seed *UUIDSeed) GenerateV7Bytes() [16]byte {
	milliseconds, sequence := seed.getV7Time()
	seedCopy := seed.GetSeedBytes()
	makeV7(&seedCopy, milliseconds, sequence)
	return seedCopy
}

// Adapted from "github.com/google/uuid". Used for generating V7 UUIDs.
func nanosecondsToMillisecondsAndSequence(nanoseconds int64) (int64, int64) {
	milliseconds := nanoseconds / NANOSECONDS_PER_MILLISECONDS
	sequence := (nanoseconds - milliseconds*NANOSECONDS_PER_MILLISECONDS) >> 8
	return milliseconds, sequence
}

// Generates a V7 UUID with a specified timestamp. The timestamp can be gotten from
// a time.Time object via the UnixNano() function.
func GenerateV7WithTimestamp(seedBytes [16]byte, timestamp int64) [16]byte {
	milliseconds, sequence := nanosecondsToMillisecondsAndSequence(timestamp)
	makeV7(&seedBytes, milliseconds, sequence)
	return seedBytes
}

// Converts seed bytes into a V7 UUID in place.
//
// Adapted from "github.com/google/uuid".
func makeV7(seed *[16]byte, milliseconds int64, sequence int64) {
	seed[0] = byte(milliseconds >> 40)
	seed[1] = byte(milliseconds >> 32)
	seed[2] = byte(milliseconds >> 24)
	seed[3] = byte(milliseconds >> 16)
	seed[4] = byte(milliseconds >> 8)
	seed[5] = byte(milliseconds)

	seed[6] = 0x70 | (0x0F & byte(sequence>>8))
	seed[7] = byte(sequence)

	seed[8] = (seed[8] & 0x3f) | 0x80
}

// Gets the current time and outputs the associated milliseconds and
// sequence. Outputs unique values due to monotonic nature by comparing
// previous timestamp data. Mutex used to ensure safe thread access.
//
// Adapted from "github.com/google/uuid".
func (seed *UUIDSeed) getV7Time() (int64, int64) {
	seed.mutex.Lock()
	defer seed.mutex.Unlock()

	nanoseconds := time.Now().UnixNano()
	milliseconds := nanoseconds / NANOSECONDS_PER_MILLISECONDS
	sequence := (nanoseconds - milliseconds*NANOSECONDS_PER_MILLISECONDS) >> 8
	now := milliseconds<<12 + sequence

	if now <= seed.lastV7Time {
		now = seed.lastV7Time + 1
		milliseconds = now >> 12
		sequence = now & 0xfff
	}
	seed.lastV7Time = now

	return milliseconds, sequence
}
