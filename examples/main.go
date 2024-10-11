package main

import (
	"fmt"

	"github.com/genekkion/squirrel/seeds"
	"github.com/rs/zerolog/log"
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

		// Lets create 2 V7 and print their info
		for range 2 {
			log.Info().
				Any("seed", seed.GetSeedBytes()).
				Str("V7UUID", seed.GenerateV7().String()).
				Msg("Displaying information about the seed and V7UUID")
			fmt.Println()
		}
	}
	store.LogDebug()
	fmt.Println()

	// Finally let's return the seeds to the store
	invalidSeeds := store.ReturnSeeds(seeds...)
	log.Info().
		Int("length", len(invalidSeeds)).
		Msg("Displaying information for invalid seeds returned")
	fmt.Println()

	store.LogDebug()
}
