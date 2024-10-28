package seeds

import (
	"sync"
	"testing"
	"time"
)

func Test1(test *testing.T) {
	waitgroup := sync.WaitGroup{}
	store := NewSeedStore()
	count := 100
	for range count {
		waitgroup.Add(1)

		go func() {
			store.GenerateNewSeeds(5)
			waitgroup.Done()
		}()
	}
	waitgroup.Wait()
	result := store.Size()
	if result != count*5 {
		test.Errorf("Test1 - want: %d, have: %d", count, result)
	}
}

func Test2(test *testing.T) {
	waitgroup := sync.WaitGroup{}
	store := NewSeedStore()
	count := 1_000_000
	for range count {
		waitgroup.Add(1)

		go func() {
			store.GenerateNewSeeds(1)
			waitgroup.Done()
		}()
	}
	waitgroup.Wait()
	result := store.Size()
	if result != count {
		test.Errorf("Test2 - want: %d, have: %d", count, result)
	}
}

func Test3(test *testing.T) {
	waitgroup := sync.WaitGroup{}
	store := NewSeedStore()
	count := 100
	for range count {
		waitgroup.Add(1)
		go func() {
			store.GenerateNewSeeds(2)
			waitgroup.Done()
		}()
	}

	for range count / 2 {
		waitgroup.Add(1)

		go func() {
			for len(store.BorrowSeeds(2)) == 0 {
				time.Sleep(time.Millisecond)
			}
			waitgroup.Done()
		}()
	}
	waitgroup.Wait()
	result := store.borrowCount
	expected := count / 2 * 2
	if result != expected {
		test.Errorf("Test3 - want: %d, have: %d", expected, result)
	}
}

func Test4(test *testing.T) {
	waitgroup := sync.WaitGroup{}
	store := NewSeedStore()
	count := 1_000_000
	for range count {
		waitgroup.Add(1)
		go func() {
			store.GenerateNewSeeds(1)
			waitgroup.Done()
		}()
	}

	for range count / 2 {
		waitgroup.Add(1)

		go func() {
			store.BorrowSeeds(1)
			waitgroup.Done()
		}()
	}
	waitgroup.Wait()
	result := store.borrowCount
	if result != count/2 {
		test.Errorf("Test4 - want: %d, have: %d", count, result)
	}
}

func Test5(test *testing.T) {
	waitgroup := sync.WaitGroup{}
	store := NewSeedStore()
	count := 1_000_000
	store.GenerateNewSeeds(count)

	seedChan := make(chan *UUIDSeed, 128)

	for range count / 2 {
		waitgroup.Add(1)
		go func() {
			seeds := store.BorrowSeeds(1)
			for len(seeds) == 0 {
				seeds = store.BorrowSeeds(1)
			}
			seedChan <- seeds[0]
			waitgroup.Done()
		}()
	}

	for range count / 2 {
		waitgroup.Add(1)
		go func() {
			store.ReturnSeeds(<-seedChan)
			waitgroup.Done()
		}()
	}

	waitgroup.Wait()
	result := store.borrowCount
	expected := 0
	if result != expected {
		test.Errorf("Test5 - want: %d, have: %d", expected, result)
	}
}
