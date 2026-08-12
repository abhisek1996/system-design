// ============================================================================
// PARKING LOT - IDIOMATIC GO STYLE
// ============================================================================
// This version demonstrates idiomatic Go patterns:
// - Simple structs over inheritance hierarchies
// - Small, focused interfaces
// - Composition over inheritance
// - Dependency injection over singletons
// - Functional options pattern
// ============================================================================

package design_parkinglot

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// ===============================
// TYPES (Simple enums)
// ===============================

type VehicleSize int

const (
	SizeSmall  VehicleSize = iota // Motorcycle
	SizeMedium                    // Car
	SizeLarge                     // Truck/Bus
)

func (s VehicleSize) String() string {
	return [...]string{"Small", "Medium", "Large"}[s]
}

// ===============================
// VEHICLE - Simple struct, no inheritance
// ===============================

// SimpleVehicle - just data, no behavior hierarchy
type SimpleVehicle struct {
	LicensePlate string
	Size         VehicleSize
	VehicleName  string
}

// Factory functions - cleaner than constructors
func NewMotorcycleSimple(plate string) SimpleVehicle {
	return SimpleVehicle{LicensePlate: plate, Size: SizeSmall, VehicleName: "Motorcycle"}
}

func NewCarSimple(plate string) SimpleVehicle {
	return SimpleVehicle{LicensePlate: plate, Size: SizeMedium, VehicleName: "Car"}
}

func NewTruckSimple(plate string) SimpleVehicle {
	return SimpleVehicle{LicensePlate: plate, Size: SizeLarge, VehicleName: "Truck"}
}

// ===============================
// PARKING SPOT - Simple struct
// ===============================

type SimpleSpot struct {
	ID      string
	Floor   int
	Size    VehicleSize
	Vehicle *SimpleVehicle // nil if empty
	mu      sync.Mutex
}

func (s *SimpleSpot) IsEmpty() bool {
	return s.Vehicle == nil
}

func (s *SimpleSpot) CanFit(v SimpleVehicle) bool {
	return s.IsEmpty() && s.Size >= v.Size
}

func (s *SimpleSpot) Park(v SimpleVehicle) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.IsEmpty() {
		return errors.New("spot occupied")
	}
	s.Vehicle = &v
	return nil
}

func (s *SimpleSpot) Unpark() (*SimpleVehicle, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.IsEmpty() {
		return nil, errors.New("spot empty")
	}
	v := s.Vehicle
	s.Vehicle = nil
	return v, nil
}

// ===============================
// PARKING FLOOR - Simple container
// ===============================

type SimpleFloor struct {
	Number int
	Spots  []*SimpleSpot
}

func NewSimpleFloor(number, small, medium, large int) *SimpleFloor {
	floor := &SimpleFloor{Number: number}
	id := 1

	for i := 0; i < small; i++ {
		floor.Spots = append(floor.Spots, &SimpleSpot{
			ID: fmt.Sprintf("F%d-%d", number, id), Floor: number, Size: SizeSmall,
		})
		id++
	}
	for i := 0; i < medium; i++ {
		floor.Spots = append(floor.Spots, &SimpleSpot{
			ID: fmt.Sprintf("F%d-%d", number, id), Floor: number, Size: SizeMedium,
		})
		id++
	}
	for i := 0; i < large; i++ {
		floor.Spots = append(floor.Spots, &SimpleSpot{
			ID: fmt.Sprintf("F%d-%d", number, id), Floor: number, Size: SizeLarge,
		})
		id++
	}
	return floor
}

func (f *SimpleFloor) FindSpot(v SimpleVehicle) *SimpleSpot {
	for _, spot := range f.Spots {
		if spot.CanFit(v) {
			return spot
		}
	}
	return nil
}

func (f *SimpleFloor) Available() map[VehicleSize]int {
	counts := make(map[VehicleSize]int)
	for _, spot := range f.Spots {
		if spot.IsEmpty() {
			counts[spot.Size]++
		}
	}
	return counts
}

// ===============================
// TICKET
// ===============================

type SimpleTicket struct {
	ID        string
	Vehicle   SimpleVehicle
	Spot      *SimpleSpot
	EntryTime time.Time
	ExitTime  time.Time
	Paid      bool
	Amount    float64
}

func (t *SimpleTicket) Duration() time.Duration {
	if t.Paid {
		return t.ExitTime.Sub(t.EntryTime)
	}
	return time.Since(t.EntryTime)
}

// ===============================
// PRICER - Interface (Strategy pattern is idiomatic)
// ===============================

type Pricer interface {
	Price(ticket *SimpleTicket) float64
}

// HourlyPricer implements Pricer
type HourlyPricer struct {
	Rates map[VehicleSize]float64
}

func NewHourlyPricer(small, medium, large float64) HourlyPricer {
	return HourlyPricer{
		Rates: map[VehicleSize]float64{
			SizeSmall:  small,
			SizeMedium: medium,
			SizeLarge:  large,
		},
	}
}

func (p HourlyPricer) Price(t *SimpleTicket) float64 {
	hours := t.Duration().Hours()
	if hours < 1 {
		hours = 1
	}
	return p.Rates[t.Spot.Size] * hours
}

// FlatPricer implements Pricer
type FlatPricer struct {
	Rate float64
}

func (p FlatPricer) Price(t *SimpleTicket) float64 {
	return p.Rate
}

// ===============================
// PARKING LOT - Dependency Injection (not Singleton)
// ===============================

type SimpleParkingLot struct {
	Name    string
	Address string
	Floors  []*SimpleFloor
	Pricer  Pricer

	tickets   map[string]*SimpleTicket // ticketID -> ticket
	vehicles  map[string]string        // plate -> ticketID
	ticketSeq int
	mu        sync.Mutex
}

// NewSimpleParkingLot - dependency injection instead of singleton
func NewSimpleParkingLot(name, address string, pricer Pricer) *SimpleParkingLot {
	return &SimpleParkingLot{
		Name:     name,
		Address:  address,
		Pricer:   pricer,
		tickets:  make(map[string]*SimpleTicket),
		vehicles: make(map[string]string),
	}
}

func (pl *SimpleParkingLot) AddFloor(f *SimpleFloor) {
	pl.Floors = append(pl.Floors, f)
}

func (pl *SimpleParkingLot) Park(v SimpleVehicle) (*SimpleTicket, error) {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	// Check if already parked
	if _, exists := pl.vehicles[v.LicensePlate]; exists {
		return nil, errors.New("vehicle already parked")
	}

	// Find spot
	var spot *SimpleSpot
	for _, floor := range pl.Floors {
		if s := floor.FindSpot(v); s != nil {
			spot = s
			break
		}
	}
	if spot == nil {
		return nil, errors.New("no available spot")
	}

	// Park
	if err := spot.Park(v); err != nil {
		return nil, err
	}

	// Create ticket
	pl.ticketSeq++
	ticket := &SimpleTicket{
		ID:        fmt.Sprintf("T-%d", pl.ticketSeq),
		Vehicle:   v,
		Spot:      spot,
		EntryTime: time.Now(),
	}

	pl.tickets[ticket.ID] = ticket
	pl.vehicles[v.LicensePlate] = ticket.ID

	return ticket, nil
}

func (pl *SimpleParkingLot) Exit(ticketID string) (float64, error) {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	ticket, exists := pl.tickets[ticketID]
	if !exists {
		return 0, errors.New("invalid ticket")
	}

	if ticket.Paid {
		return 0, errors.New("already paid")
	}

	// Calculate price
	amount := pl.Pricer.Price(ticket)

	// Unpark
	if _, err := ticket.Spot.Unpark(); err != nil {
		return 0, err
	}

	// Update ticket
	ticket.ExitTime = time.Now()
	ticket.Paid = true
	ticket.Amount = amount

	// Cleanup
	delete(pl.vehicles, ticket.Vehicle.LicensePlate)
	delete(pl.tickets, ticketID)

	return amount, nil
}

func (pl *SimpleParkingLot) Availability() map[int]map[VehicleSize]int {
	result := make(map[int]map[VehicleSize]int)
	for _, floor := range pl.Floors {
		result[floor.Number] = floor.Available()
	}
	return result
}

func (pl *SimpleParkingLot) PrintStatus() {
	fmt.Printf("\n===== %s =====\n", pl.Name)
	fmt.Printf("Address: %s\n\n", pl.Address)

	for floor, counts := range pl.Availability() {
		fmt.Printf("Floor %d: Small=%d, Medium=%d, Large=%d\n",
			floor, counts[SizeSmall], counts[SizeMedium], counts[SizeLarge])
	}
	fmt.Println("========================")
}
