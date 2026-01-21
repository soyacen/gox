package mutexx

import (
	"hash/fnv"
	"runtime"
	"sync"
)

// GroupMutex is a sharded mutex implementation that provides fine-grained locking
// based on string keys. It uses multiple underlying mutexes (sync.Mutex) to reduce
// lock contention for different keys.
//
// Fields:
//   - N: int - Specifies the number of mutex shards to create. If <= 0, defaults to runtime.NumCPU().
//   - once: sync.Once - Ensures thread-safe one-time initialization of mutexes.
//   - mutexes: []sync.Mutex - The slice of underlying mutexes used for synchronization.
//
// This struct is useful when you need to lock different resources identified by string keys,
// allowing concurrent operations on different keys while serializing operations on the same key.
type GroupMutex struct {
	N       int
	once    sync.Once
	mutexes []sync.Mutex
}

// Lock acquires the mutex lock for the given key.
//
// This method ensures the GroupMutex is properly initialized before attempting to lock.
// It uses the key to determine which internal mutex to lock, providing fine-grained synchronization.
//
// Parameters:
//
//	key: The string key used to identify which mutex to lock. Different keys may map to the same mutex
//	     based on the internal hashing mechanism.
func (gm *GroupMutex) Lock(key string) {
	// Ensure the GroupMutex is initialized before proceeding
	gm.ensureInit()

	// Lock the mutex corresponding to the hashed index of the key
	gm.mutexes[gm.getMutexIndex(key)].Lock()
}

// Unlock releases the lock associated with the given key.
//
// This method first ensures the internal mutexes are initialized (if not already),
// then calculates the appropriate mutex index for the given key and unlocks it.
//
// Parameters:
//   - key: The string key used to identify which mutex to unlock.
//     The same key will always map to the same underlying mutex.
func (gm *GroupMutex) Unlock(key string) {
	// Ensure mutexes are initialized before attempting to unlock
	gm.ensureInit()

	// Unlock the mutex corresponding to the hashed key index
	gm.mutexes[gm.getMutexIndex(key)].Unlock()
}

// ensureInit ensures the GroupMutex is initialized with a slice of mutexes.
// It uses sync.Once to guarantee thread-safe initialization only once.
// The number of mutexes is determined by gm.N, or defaults to the number of CPUs if gm.N <= 0.
func (gm *GroupMutex) ensureInit() {
	// Use sync.Once to perform initialization exactly once
	gm.once.Do(func() {
		// Determine the number of mutexes to create
		n := gm.N
		if n <= 0 {
			// Default to number of CPUs if invalid value
			n = runtime.NumCPU()
		}
		// Initialize the slice of mutexes
		gm.mutexes = make([]sync.Mutex, n)
	})
}

// getMutexIndex calculates the index of the mutex to use for the given key.
// It uses a hash function to distribute keys evenly across the available mutexes.
//
// Parameters:
//
//	key - the string key used to determine which mutex to select
//
// Returns:
//
//	uint32 - the index of the mutex in the mutexes slice
func (gm *GroupMutex) getMutexIndex(key string) uint32 {
	// Hash the key and use modulo operation to ensure the index is within bounds
	return gm.hash(key) % uint32(len(gm.mutexes))
}

// hash computes a 32-bit FNV-1a hash for the given key.
// This is used to determine which mutex in the group should be used for a particular key.
//
// Parameters:
//
//	key: the string to be hashed
//
// Returns:
//
//	uint32: the 32-bit hash value of the input key
func (gm *GroupMutex) hash(key string) uint32 {
	// Create new FNV-1a 32-bit hash
	h := fnv.New32a()

	// Write key bytes to hash and return the result
	h.Write([]byte(key))
	return h.Sum32()
}
