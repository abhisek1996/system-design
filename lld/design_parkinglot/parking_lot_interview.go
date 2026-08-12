// ============================================================================
// PARKING LOT - INTERVIEW STYLE (OOP with Inheritance)
// ============================================================================
// This version demonstrates OOP concepts using Go's struct embedding
// to simulate inheritance. Good for LLD interviews.
//
// Patterns used:
// - Inheritance (via struct embedding)
// - Interface polymorphism
// - Singleton pattern
// - Strategy pattern
// - Factory pattern
// ============================================================================

package design_parkinglot

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// ===============================
// ENUMS & CONSTANTS
// ===============================

type VehicleType int

const (
	VehicleTypeMotorcycle VehicleType = iota
	VehicleTypeCar
	VehicleTypeTruck
)

func (v VehicleType) String() string {
	return [...]string{"Motorcycle", "Car", "Truck"}[v]
}

type SpotSize int

const (
	SpotSizeSmall  SpotSize = iota // For motorcycles
	SpotSizeMedium                 // For cars
	SpotSizeLarge                  // For trucks/buses
)

func (s SpotSize) String() string {
	return [...]string{"Small", "Medium", "Large"}[s]
}

type TicketStatus int

const (
	TicketStatusActive TicketStatus = iota
	TicketStatusPaid
)

// ===============================
// VEHICLE - Interface + Inheritance
// ===============================

// Vehicle interface - base type for all vehicles (polymorphism)
type Vehicle interface {
	GetLicensePlate() string
	GetVehicleType() VehicleType
	GetRequiredSpotSize() SpotSize
}

// BaseVehicle - common fields (composition for inheritance)
type BaseVehicle struct {
	LicensePlate string
	Type         VehicleType
}

func (v *BaseVehicle) GetLicensePlate() string {
	return v.LicensePlate
}

func (v *BaseVehicle) GetVehicleType() VehicleType {
	return v.Type
}

// Motorcycle - inherits from BaseVehicle
type Motorcycle struct {
	BaseVehicle // Embedding = Inheritance in Go
}

func NewMotorcycle(licensePlate string) *Motorcycle {
	return &Motorcycle{
		BaseVehicle: BaseVehicle{
			LicensePlate: licensePlate,
			Type:         VehicleTypeMotorcycle,
		},
	}
}

func (m *Motorcycle) GetRequiredSpotSize() SpotSize {
	return SpotSizeSmall // Motorcycle only needs small spot
}

// Car - inherits from BaseVehicle
type Car struct {
	BaseVehicle
}

func NewCar(licensePlate string) *Car {
	return &Car{
		BaseVehicle: BaseVehicle{
			LicensePlate: licensePlate,
			Type:         VehicleTypeCar,
		},
	}
}

func (c *Car) GetRequiredSpotSize() SpotSize {
	return SpotSizeMedium // Car needs medium or larger
}

// Truck - inherits from BaseVehicle
type Truck struct {
	BaseVehicle
}

func NewTruck(licensePlate string) *Truck {
	return &Truck{
		BaseVehicle: BaseVehicle{
			LicensePlate: licensePlate,
			Type:         VehicleTypeTruck,
		},
	}
}

func (t *Truck) GetRequiredSpotSize() SpotSize {
	return SpotSizeLarge // Truck needs large spot
}

// Bus - inherits from BaseVehicle
type Bus struct {
	BaseVehicle
}

func NewBus(licensePlate string) *Bus {
	return &Bus{
		BaseVehicle: BaseVehicle{
			LicensePlate: licensePlate,
			Type:         VehicleTypeTruck, // Treated same as truck
		},
	}
}

func (b *Bus) GetRequiredSpotSize() SpotSize {
	return SpotSizeLarge
}

// ===============================
// PARKING SPOT - Interface + Inheritance
// ===============================

// ParkingSpot interface - base type for all spots (polymorphism)
type ParkingSpot interface {
	GetID() string
	GetFloorNum() int
	GetSpotNum() int
	GetSize() SpotSize
	IsAvailable() bool
	CanFitVehicle(vehicle Vehicle) bool
	Park(vehicle Vehicle) error
	Unpark() (Vehicle, error)
	GetVehicle() Vehicle
}

// BaseParkingSpot - common fields and behavior
type BaseParkingSpot struct {
	ID         string
	FloorNum   int
	SpotNum    int
	Size       SpotSize
	IsOccupied bool
	Vehicle    Vehicle
	mu         sync.Mutex
}

func (ps *BaseParkingSpot) GetID() string {
	return ps.ID
}

func (ps *BaseParkingSpot) GetFloorNum() int {
	return ps.FloorNum
}

func (ps *BaseParkingSpot) GetSpotNum() int {
	return ps.SpotNum
}

func (ps *BaseParkingSpot) GetSize() SpotSize {
	return ps.Size
}

func (ps *BaseParkingSpot) IsAvailable() bool {
	return !ps.IsOccupied
}

func (ps *BaseParkingSpot) GetVehicle() Vehicle {
	return ps.Vehicle
}

func (ps *BaseParkingSpot) CanFitVehicle(vehicle Vehicle) bool {
	if ps.IsOccupied {
		return false
	}
	requiredSize := vehicle.GetRequiredSpotSize()
	return ps.Size >= requiredSize
}

func (ps *BaseParkingSpot) Park(vehicle Vehicle) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if ps.IsOccupied {
		return errors.New("spot is already occupied")
	}
	ps.IsOccupied = true
	ps.Vehicle = vehicle
	return nil
}

func (ps *BaseParkingSpot) Unpark() (Vehicle, error) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if !ps.IsOccupied {
		return nil, errors.New("spot is already empty")
	}
	vehicle := ps.Vehicle
	ps.IsOccupied = false
	ps.Vehicle = nil
	return vehicle, nil
}

// SmallSpot - for motorcycles (inherits BaseParkingSpot)
type SmallSpot struct {
	BaseParkingSpot
}

func NewSmallSpot(floorNum, spotNum int) *SmallSpot {
	return &SmallSpot{
		BaseParkingSpot: BaseParkingSpot{
			ID:       fmt.Sprintf("F%d-S%d", floorNum, spotNum),
			FloorNum: floorNum,
			SpotNum:  spotNum,
			Size:     SpotSizeSmall,
		},
	}
}

// Override: Small spot can only fit motorcycles
func (s *SmallSpot) CanFitVehicle(vehicle Vehicle) bool {
	if s.IsOccupied {
		return false
	}
	return vehicle.GetRequiredSpotSize() == SpotSizeSmall
}

// MediumSpot - for cars (inherits BaseParkingSpot)
type MediumSpot struct {
	BaseParkingSpot
}

func NewMediumSpot(floorNum, spotNum int) *MediumSpot {
	return &MediumSpot{
		BaseParkingSpot: BaseParkingSpot{
			ID:       fmt.Sprintf("F%d-M%d", floorNum, spotNum),
			FloorNum: floorNum,
			SpotNum:  spotNum,
			Size:     SpotSizeMedium,
		},
	}
}

// Override: Medium spot can fit motorcycles and cars
func (m *MediumSpot) CanFitVehicle(vehicle Vehicle) bool {
	if m.IsOccupied {
		return false
	}
	requiredSize := vehicle.GetRequiredSpotSize()
	return requiredSize <= SpotSizeMedium
}

// LargeSpot - for trucks/buses (inherits BaseParkingSpot)
type LargeSpot struct {
	BaseParkingSpot
}

func NewLargeSpot(floorNum, spotNum int) *LargeSpot {
	return &LargeSpot{
		BaseParkingSpot: BaseParkingSpot{
			ID:       fmt.Sprintf("F%d-L%d", floorNum, spotNum),
			FloorNum: floorNum,
			SpotNum:  spotNum,
			Size:     SpotSizeLarge,
		},
	}
}

// Override: Large spot can fit any vehicle
func (l *LargeSpot) CanFitVehicle(vehicle Vehicle) bool {
	if l.IsOccupied {
		return false
	}
	return true // Can fit any vehicle type
}

// ===============================
// PARKING FLOOR
// ===============================

type ParkingFloor struct {
	FloorNumber int
	Spots       []ParkingSpot // Interface type - polymorphism!
	mu          sync.RWMutex
}

func NewParkingFloor(floorNumber int, smallSpots, mediumSpots, largeSpots int) *ParkingFloor {
	floor := &ParkingFloor{
		FloorNumber: floorNumber,
		Spots:       make([]ParkingSpot, 0),
	}

	spotNum := 1
	// Create small spots
	for i := 0; i < smallSpots; i++ {
		floor.Spots = append(floor.Spots, NewSmallSpot(floorNumber, spotNum))
		spotNum++
	}
	// Create medium spots
	for i := 0; i < mediumSpots; i++ {
		floor.Spots = append(floor.Spots, NewMediumSpot(floorNumber, spotNum))
		spotNum++
	}
	// Create large spots
	for i := 0; i < largeSpots; i++ {
		floor.Spots = append(floor.Spots, NewLargeSpot(floorNumber, spotNum))
		spotNum++
	}

	return floor
}

func (pf *ParkingFloor) FindAvailableSpot(vehicle Vehicle) ParkingSpot {
	pf.mu.RLock()
	defer pf.mu.RUnlock()

	for _, spot := range pf.Spots {
		if spot.CanFitVehicle(vehicle) {
			return spot
		}
	}
	return nil
}

func (pf *ParkingFloor) GetAvailableSpotCount() map[SpotSize]int {
	pf.mu.RLock()
	defer pf.mu.RUnlock()

	counts := make(map[SpotSize]int)
	for _, spot := range pf.Spots {
		if spot.IsAvailable() {
			counts[spot.GetSize()]++
		}
	}
	return counts
}

// ===============================
// PARKING TICKET
// ===============================

type ParkingTicket struct {
	TicketID   string
	Vehicle    Vehicle     // Interface type - polymorphism!
	Spot       ParkingSpot // Interface type - polymorphism!
	EntryTime  time.Time
	ExitTime   time.Time
	Status     TicketStatus
	AmountPaid float64
}

func NewParkingTicket(ticketID string, vehicle Vehicle, spot ParkingSpot) *ParkingTicket {
	return &ParkingTicket{
		TicketID:  ticketID,
		Vehicle:   vehicle,
		Spot:      spot,
		EntryTime: time.Now(),
		Status:    TicketStatusActive,
	}
}

func (pt *ParkingTicket) GetParkingDuration() time.Duration {
	if pt.Status == TicketStatusPaid {
		return pt.ExitTime.Sub(pt.EntryTime)
	}
	return time.Since(pt.EntryTime)
}

// ===============================
// PRICING STRATEGY (Strategy Pattern)
// ===============================

type PricingStrategy interface {
	CalculatePrice(ticket *ParkingTicket) float64
}

// HourlyPricingStrategy - implements PricingStrategy
type HourlyPricingStrategy struct {
	SmallSpotRate  float64
	MediumSpotRate float64
	LargeSpotRate  float64
}

func NewHourlyPricingStrategy(smallRate, mediumRate, largeRate float64) *HourlyPricingStrategy {
	return &HourlyPricingStrategy{
		SmallSpotRate:  smallRate,
		MediumSpotRate: mediumRate,
		LargeSpotRate:  largeRate,
	}
}

func (h *HourlyPricingStrategy) CalculatePrice(ticket *ParkingTicket) float64 {
	duration := ticket.GetParkingDuration()
	hours := duration.Hours()
	if hours < 1 {
		hours = 1
	}

	var rate float64
	switch ticket.Spot.GetSize() {
	case SpotSizeSmall:
		rate = h.SmallSpotRate
	case SpotSizeMedium:
		rate = h.MediumSpotRate
	case SpotSizeLarge:
		rate = h.LargeSpotRate
	}

	return rate * hours
}

// FlatRatePricingStrategy - another implementation (polymorphism)
type FlatRatePricingStrategy struct {
	FlatRate float64
}

func NewFlatRatePricingStrategy(rate float64) *FlatRatePricingStrategy {
	return &FlatRatePricingStrategy{FlatRate: rate}
}

func (f *FlatRatePricingStrategy) CalculatePrice(ticket *ParkingTicket) float64 {
	return f.FlatRate
}

// ===============================
// ENTRY & EXIT GATES
// ===============================

type EntryGate struct {
	ID         string
	ParkingLot *ParkingLot
}

func NewEntryGate(id string, parkingLot *ParkingLot) *EntryGate {
	return &EntryGate{
		ID:         id,
		ParkingLot: parkingLot,
	}
}

func (eg *EntryGate) ProcessEntry(vehicle Vehicle) (*ParkingTicket, error) {
	return eg.ParkingLot.ParkVehicle(vehicle)
}

type ExitGate struct {
	ID         string
	ParkingLot *ParkingLot
}

func NewExitGate(id string, parkingLot *ParkingLot) *ExitGate {
	return &ExitGate{
		ID:         id,
		ParkingLot: parkingLot,
	}
}

func (eg *ExitGate) ProcessExit(ticketID string) (float64, error) {
	return eg.ParkingLot.UnparkVehicle(ticketID)
}

// ===============================
// PARKING LOT (Singleton)
// ===============================

type ParkingLot struct {
	Name            string
	Address         string
	Floors          []*ParkingFloor
	EntryGates      []*EntryGate
	ExitGates       []*ExitGate
	ActiveTickets   map[string]*ParkingTicket
	VehicleTickets  map[string]string
	PricingStrategy PricingStrategy
	ticketCounter   int
	mu              sync.Mutex
}

var (
	parkingLotInstance *ParkingLot
	once               sync.Once
)

func GetParkingLotInstance(name, address string) *ParkingLot {
	once.Do(func() {
		parkingLotInstance = &ParkingLot{
			Name:           name,
			Address:        address,
			Floors:         make([]*ParkingFloor, 0),
			EntryGates:     make([]*EntryGate, 0),
			ExitGates:      make([]*ExitGate, 0),
			ActiveTickets:  make(map[string]*ParkingTicket),
			VehicleTickets: make(map[string]string),
		}
	})
	return parkingLotInstance
}

func ResetParkingLotInstance() {
	once = sync.Once{}
	parkingLotInstance = nil
}

func (pl *ParkingLot) AddFloor(floor *ParkingFloor) {
	pl.mu.Lock()
	defer pl.mu.Unlock()
	pl.Floors = append(pl.Floors, floor)
}

func (pl *ParkingLot) AddEntryGate(gateID string) *EntryGate {
	pl.mu.Lock()
	defer pl.mu.Unlock()
	gate := NewEntryGate(gateID, pl)
	pl.EntryGates = append(pl.EntryGates, gate)
	return gate
}

func (pl *ParkingLot) AddExitGate(gateID string) *ExitGate {
	pl.mu.Lock()
	defer pl.mu.Unlock()
	gate := NewExitGate(gateID, pl)
	pl.ExitGates = append(pl.ExitGates, gate)
	return gate
}

func (pl *ParkingLot) SetPricingStrategy(strategy PricingStrategy) {
	pl.mu.Lock()
	defer pl.mu.Unlock()
	pl.PricingStrategy = strategy
}

func (pl *ParkingLot) generateTicketID() string {
	pl.ticketCounter++
	return fmt.Sprintf("TKT-%d-%d", time.Now().Unix(), pl.ticketCounter)
}

func (pl *ParkingLot) findAvailableSpot(vehicle Vehicle) ParkingSpot {
	for _, floor := range pl.Floors {
		if spot := floor.FindAvailableSpot(vehicle); spot != nil {
			return spot
		}
	}
	return nil
}

func (pl *ParkingLot) ParkVehicle(vehicle Vehicle) (*ParkingTicket, error) {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	if _, exists := pl.VehicleTickets[vehicle.GetLicensePlate()]; exists {
		return nil, errors.New("vehicle is already parked")
	}

	spot := pl.findAvailableSpot(vehicle)
	if spot == nil {
		return nil, errors.New("no available parking spot for this vehicle")
	}

	if err := spot.Park(vehicle); err != nil {
		return nil, err
	}

	ticketID := pl.generateTicketID()
	ticket := NewParkingTicket(ticketID, vehicle, spot)
	pl.ActiveTickets[ticketID] = ticket
	pl.VehicleTickets[vehicle.GetLicensePlate()] = ticketID

	return ticket, nil
}

func (pl *ParkingLot) UnparkVehicle(ticketID string) (float64, error) {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	ticket, exists := pl.ActiveTickets[ticketID]
	if !exists {
		return 0, errors.New("invalid ticket")
	}

	if ticket.Status == TicketStatusPaid {
		return 0, errors.New("ticket already paid")
	}

	var amount float64
	if pl.PricingStrategy != nil {
		amount = pl.PricingStrategy.CalculatePrice(ticket)
	}

	_, err := ticket.Spot.Unpark()
	if err != nil {
		return 0, err
	}

	ticket.ExitTime = time.Now()
	ticket.Status = TicketStatusPaid
	ticket.AmountPaid = amount

	delete(pl.VehicleTickets, ticket.Vehicle.GetLicensePlate())
	delete(pl.ActiveTickets, ticketID)

	return amount, nil
}

func (pl *ParkingLot) GetAvailability() map[int]map[SpotSize]int {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	availability := make(map[int]map[SpotSize]int)
	for _, floor := range pl.Floors {
		availability[floor.FloorNumber] = floor.GetAvailableSpotCount()
	}
	return availability
}

func (pl *ParkingLot) IsFull() bool {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	for _, floor := range pl.Floors {
		counts := floor.GetAvailableSpotCount()
		for _, count := range counts {
			if count > 0 {
				return false
			}
		}
	}
	return true
}

// ===============================
// DISPLAY BOARD
// ===============================

type DisplayBoard struct {
	ParkingLot *ParkingLot
}

func NewDisplayBoard(parkingLot *ParkingLot) *DisplayBoard {
	return &DisplayBoard{ParkingLot: parkingLot}
}

func (db *DisplayBoard) ShowAvailability() {
	availability := db.ParkingLot.GetAvailability()
	fmt.Println("\n========== PARKING AVAILABILITY ==========")
	fmt.Printf("Parking Lot: %s\n", db.ParkingLot.Name)
	fmt.Printf("Address: %s\n\n", db.ParkingLot.Address)

	for floorNum, spots := range availability {
		fmt.Printf("Floor %d:\n", floorNum)
		fmt.Printf("  Small spots (Motorcycle):  %d\n", spots[SpotSizeSmall])
		fmt.Printf("  Medium spots (Car):        %d\n", spots[SpotSizeMedium])
		fmt.Printf("  Large spots (Truck/Bus):   %d\n", spots[SpotSizeLarge])
		fmt.Println()
	}
	fmt.Println("===========================================")
}
