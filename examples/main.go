package main

import (
	"fmt"
	"log"

	"github.com/genekkion/squirrel/seeds"
	"github.com/google/uuid"
)

func main() {
	seedCount := 16

	// Create a new store for the seeds
	store := seeds.NewSeedStorage()
	// Should display an empty store
	store.LogDebug()
	fmt.Println()

	// Generate some seeds!
	store.GenerateNewSeeds(seedCount)
	// Should show that there are some seeds in the store
	store.LogDebug()
	fmt.Println()

	// Lets take a look at two of the seeds
	seeds := store.BorrowSeeds(2)

	for _, seed := range seeds {
		seedBytes := seed.GetSeedBytes()            // As [16]byte
		seedUUID, _ := uuid.FromBytes(seedBytes[:]) // As UUID

		// Lets create 2 V7 and print their info
		for range 2 {

			log.Printf("Displaying information about the seed and V7UUID\n { seed: %s, V7UUID: %s }\n",
				seedUUID.String(),
				seed.GenerateV7().String(),
			)
			fmt.Println()
		}
	}
	store.LogDebug()
	fmt.Println()

	// Finally let's return the seeds to the store
	invalidSeeds := store.ReturnSeeds(seeds...)
	log.Printf("Displaying information for the invalid seeds returned\n { count: %d }\n",
		len(invalidSeeds),
	)

	store.LogDebug()
}
