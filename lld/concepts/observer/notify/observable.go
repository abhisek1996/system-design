package notify

import (
	"fmt"
	"sync"
)

type Observable interface {
	AddObserver(o Observer)
	RemoveObserver(o Observer)
	Broadcast(e Event)
}

// --- Base Implementation ---

type BaseObservable struct {
	observers []Observer
	mu        sync.RWMutex
}

func (b *BaseObservable) AddObserver(o Observer) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.observers = append(b.observers, o)
}

func (b *BaseObservable) RemoveObserver(o Observer) {
	b.mu.Lock()
	defer b.mu.Unlock()
	newObservers := make([]Observer, 0, len(b.observers))
	for _, observer := range b.observers {
		if observer != o {
			newObservers = append(newObservers, observer)
		}
	}
	b.observers = newObservers
}

func (b *BaseObservable) Broadcast(e Event) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for _, observer := range b.observers {
		observer.Update(e)
	}
}

// --- Concrete Items ---

// ElectronicsItem
type ElectronicsItem struct {
	*BaseObservable
	model      string
	price      int
	stockCount int
}

func NewElectronicsItem(model string, price int) *ElectronicsItem {
	return &ElectronicsItem{
		BaseObservable: &BaseObservable{},
		model:          model,
		price:          price,
		stockCount:     0,
	}
}

func (i *ElectronicsItem) AddItem(count int) {

	if i.stockCount == 0 {
		i.stockCount += count
		i.Broadcast(Event{
			Type:    EventStockUpdate,
			Message: fmt.Sprintf("%s is back in stock!", i.model),
			Data: map[string]interface{}{
				"model": i.model,
				"stock": i.stockCount,
			},
		})
		return
	}
	i.stockCount += count
}

func (i *ElectronicsItem) RemoveItem(count int) {
	i.stockCount -= count
}

func (i *ElectronicsItem) GetStockCount() int {
	return i.stockCount
}

// GroceryItem
type GroceryItem struct {
	*BaseObservable
	name       string
	stockCount int
}

func NewGroceryItem(name string) *GroceryItem {
	return &GroceryItem{
		BaseObservable: &BaseObservable{},
		name:           name,
	}
}

func (g *GroceryItem) AddItem(count int) {
	if g.stockCount == 0 {
		g.stockCount += count
		g.Broadcast(Event{
			Type:    EventStockUpdate,
			Message: fmt.Sprintf("%s is refilled!", g.name),
			Data: map[string]interface{}{
				"name":  g.name,
				"stock": g.stockCount,
			},
		})
		return
	}
	g.stockCount += count
}

func (g *GroceryItem) RemoveItem(count int) {
	g.stockCount -= count
}

func (g *GroceryItem) GetStockCount() int {
	return g.stockCount
}
