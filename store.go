package evplanner

import (
	"sync"
)

type PlaceholderObject struct {
	contents string
}

func (p *PlaceholderObject) Serialize() string {
	return p.contents
}

func PlaceholderFromString(contents string) *PlaceholderObject {
	return &PlaceholderObject{contents}
}

// Manages the persistence of all our data
type Store interface {
	Read() *PlaceholderObject
	Write(*PlaceholderObject)
	LockForModification()
	UnlockForModification()
	IsLockedForModification() bool
}

type InMemoryStore struct {
	storeage PlaceholderObject
	modMutex sync.Mutex
	isLocked bool
}

// returns a copy of the placeholder
func (i *InMemoryStore) Read() *PlaceholderObject {
	return &PlaceholderObject{i.storeage.contents}
}

func (i *InMemoryStore) Write(p *PlaceholderObject) {
	i.storeage.contents = p.contents
}

func (i *InMemoryStore) LockForModification() {
	i.modMutex.Lock()
	i.isLocked = true
}

func (i *InMemoryStore) UnlockForModification() {
	i.isLocked = false
	i.modMutex.Unlock()
}

func (i *InMemoryStore) IsLockedForModification() bool {
	return i.isLocked
}
