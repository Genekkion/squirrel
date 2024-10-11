# Squirrel 🐿️

Squirrel is a library meant to help you generate unique V7 UUIDs across multiple threads or nodes concurrently.
It relies on the fact that V7 UUIDs have a set number of bits dedicated for the monotonic time-based component.

By using seeds as a form of bitmasking, we are able to generate V7 UUIDs with the randomness component fixed
by the seed, whilst ensuring uniquness via the time component.

## Installation

```bash
go get github.com/genekkion/squirrel/pkg/seeds
```

## Usage

```go
import "github.com/genekkion/squirrel/pkg/seeds"

// Create a new store for the seeds
store := seeds.NewSeedStorage()

// Generate some seeds!
store.GenerateNewSeeds(seedCount)

// Now we can borrow some seeds
seeds := store.BorrowSeeds(2)

// And then generate some unique V7 UUIDs!
seed1 := seeds[0]
v7UUID := seed1.GenerateV7()
```


## Attribution
This library adapts code from the library uuid by Google,
which is licensed under the BSD 3-Clause License.

You can view the original package here: [https://github.com/google/uuid](https://github.com/google/uuid).
