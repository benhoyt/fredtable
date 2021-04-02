// Hash table implementation by Frederico Bittencourt (here for code review only)
//
// Taken from: https://blog.fredrb.com/2021/04/01/hashtable-go/
// (c) 2021 Frederico Bittencourt

package fredtable

import (
	"fmt"
	"io"
)

const (
	FNV_OFFSET_BASIS      uint64 = 14695981039346656037
	FNV_PRIME             uint64 = 1099511628211
	INITIAL_M                    = 16
	LOAD_FACTOR_THRESHOLD        = 0.5
)

func hashValue(v PreHashable, mod int) int {
	hash := FNV_OFFSET_BASIS
	vBytes := v.Prehash()
	for _, b := range vBytes {
		hash = hash ^ uint64(b)
		hash = hash * FNV_PRIME
	}
	return int(hash % uint64(mod))
}

type IntKey int
type StringKey string

type PreHashable interface {
	Prehash() []byte
}

func (i IntKey) Prehash() []byte {
	return []byte{byte(i)}
}

func (str StringKey) Prehash() []byte {
	return []byte(str)
}

type HashTable struct {
	Length       int
	directAccess [][]HashTableEntry
	m            int
}

type HashTableEntry struct {
	key   PreHashable
	value interface{}
}

func NewSizedHashTable(initial int) HashTable {
	return HashTable{
		Length:       0,
		directAccess: make([][]HashTableEntry, initial),
		m:            initial,
	}
}

func NewHashTable() HashTable {
	return NewSizedHashTable(INITIAL_M)
}

func (ht *HashTable) Get(key PreHashable) *interface{} {
	hash := hashValue(key, ht.m)
	for _, v := range ht.directAccess[hash] {
		if v.key == key {
			return &v.value
		}
	}
	return nil
}

func (ht *HashTable) loadFactor() float32 {
	return float32(ht.Length) / float32(ht.m)
}

func (ht *HashTable) Set(key PreHashable, value interface{}) {
	hash := hashValue(key, ht.m)
	if ht.directAccess[hash] == nil {
		ht.directAccess[hash] = []HashTableEntry{}
	}
	ht.directAccess[hash] = append(ht.directAccess[hash], HashTableEntry{key, value})
	ht.Length += 1

	if ht.loadFactor() > LOAD_FACTOR_THRESHOLD {
		ht.doubleTable()
	}
}

func (ht *HashTable) doubleTable() {
	ht.m = ht.m * 2
	newTable := make([][]HashTableEntry, ht.m)
	var newHash int
	for _, bucket := range ht.directAccess {
		for _, v := range bucket {
			newHash = hashValue(v.key, ht.m)
			if newTable[newHash] == nil {
				newTable[newHash] = []HashTableEntry{}
			}
			newTable[newHash] = append(newTable[newHash], HashTableEntry{v.key, v.value})
		}
	}
	ht.directAccess = newTable
}

func (ht *HashTable) deleteEntry(hash int, key PreHashable) (interface{}, error) {
	index := -1
	for i, v := range ht.directAccess[hash] {
		if v.key == key {
			index = i
		}
	}
	if index == -1 {
		return nil, fmt.Errorf("Key error")
	} else {
		// Remove from array
		object := ht.directAccess[hash][index]
		ht.directAccess[hash][index] = ht.directAccess[hash][len(ht.directAccess[hash])-1]
		ht.directAccess[hash] = ht.directAccess[hash][:len(ht.directAccess[hash])-1]
		ht.Length -= 1
		return object.value, nil
	}
}

func (ht *HashTable) Delete(key PreHashable) (interface{}, error) {
	hash := hashValue(key, ht.m)
	if ht.directAccess[hash] == nil {
		return nil, fmt.Errorf("Key error")
	}
	return ht.deleteEntry(hash, key)
}

func (ht *HashTable) Dump(w io.Writer) {
	fmt.Fprintf(w, "length = %d\n", ht.Length)
	for i, entries := range ht.directAccess {
		fmt.Fprintf(w, "bucket %3d: ", i)
		for j, entry := range entries {
			fmt.Fprintf(w, "%s:%v", entry.key, entry.value)
			if j < len(entries)-1 {
				fmt.Fprintf(w, ", ")
			}
		}
		fmt.Fprintln(w)
	}
}
