package design_parkinglot

import (
	"fmt"
)

// RunInterviewStyleDemo demonstrates the OOP-style version (for interviews)
func RunInterviewStyleDemo() {
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║     PARKING LOT - INTERVIEW STYLE (OOP with Inheritance)     ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")

	// Reset singleton for clean demo
	ResetParkingLotInstance()

	// 1. Create Parking Lot (Singleton Pattern)
	parkingLot := GetParkingLotInstance("City Center Parking", "123 Main Street")

	// 2. Add Floors (SmallSpot, MediumSpot, LargeSpot - different classes)
	parkingLot.AddFloor(NewParkingFloor(1, 3, 5, 2)) // 3 small, 5 medium, 2 large

	// 3. Set Strategy
	parkingLot.SetPricingStrategy(NewHourlyPricingStrategy(10, 20, 30))

	// 4. Add Gates
	entryGate := parkingLot.AddEntryGate("ENTRY-1")
	exitGate := parkingLot.AddExitGate("EXIT-1")

	// Display initial state
	displayBoard := NewDisplayBoard(parkingLot)
	displayBoard.ShowAvailability()

	// Create vehicles using inheritance hierarchy
	fmt.Println("\n🔷 INHERITANCE HIERARCHY:")
	fmt.Println("   Vehicle (interface)")
	fmt.Println("      └── BaseVehicle (embedded struct)")
	fmt.Println("            ├── Motorcycle")
	fmt.Println("            ├── Car")
	fmt.Println("            ├── Truck")
	fmt.Println("            └── Bus")

	// Park vehicles polymorphically
	fmt.Println("\n📥 PARKING (using polymorphism):")
	vehicles := []Vehicle{
		NewMotorcycle("MOTO-001"),
		NewCar("CAR-001"),
		NewTruck("TRUCK-001"),
	}

	tickets := make([]*ParkingTicket, 0)
	for _, v := range vehicles {
		ticket, err := entryGate.ProcessEntry(v)
		if err != nil {
			fmt.Printf("   ❌ %s: %v\n", v.GetLicensePlate(), err)
		} else {
			fmt.Printf("   ✓ %s (%s) → %s (%s)\n",
				v.GetLicensePlate(), v.GetVehicleType(),
				ticket.Spot.GetID(), ticket.Spot.GetSize())
			tickets = append(tickets, ticket)
		}
	}

	// Exit one vehicle
	fmt.Println("\n📤 EXIT:")
	if len(tickets) > 0 {
		amount, _ := exitGate.ProcessExit(tickets[0].TicketID)
		fmt.Printf("   ✓ %s exited, paid $%.2f\n",
			tickets[0].Vehicle.GetLicensePlate(), amount)
	}

	displayBoard.ShowAvailability()
}

// RunIdiomaticStyleDemo demonstrates the idiomatic Go version
func RunIdiomaticStyleDemo() {
	fmt.Println("\n╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║         PARKING LOT - IDIOMATIC GO STYLE                     ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")

	// Dependency injection - no singleton!
	pricer := NewHourlyPricer(10, 20, 30)
	lot := NewSimpleParkingLot("Downtown Parking", "456 Oak Ave", pricer)

	// Add floors
	lot.AddFloor(NewSimpleFloor(1, 3, 5, 2))

	lot.PrintStatus()

	// Simple vehicle creation - no inheritance
	fmt.Println("\n🔷 NO INHERITANCE - Just simple structs:")
	fmt.Println("   SimpleVehicle { LicensePlate, Size, VehicleName }")

	// Create vehicles
	vehicles := []SimpleVehicle{
		NewMotorcycleSimple("MOTO-100"),
		NewCarSimple("CAR-100"),
		NewTruckSimple("TRUCK-100"),
	}

	fmt.Println("\n📥 PARKING:")
	tickets := make([]*SimpleTicket, 0)
	for _, v := range vehicles {
		ticket, err := lot.Park(v)
		if err != nil {
			fmt.Printf("   ❌ %s: %v\n", v.LicensePlate, err)
		} else {
			fmt.Printf("   ✓ %s (%s) → %s\n",
				v.LicensePlate, v.VehicleName, ticket.Spot.ID)
			tickets = append(tickets, ticket)
		}
	}

	// Exit
	fmt.Println("\n📤 EXIT:")
	if len(tickets) > 0 {
		amount, _ := lot.Exit(tickets[0].ID)
		fmt.Printf("   ✓ %s exited, paid $%.2f\n",
			tickets[0].Vehicle.LicensePlate, amount)
	}

	lot.PrintStatus()
}

// RunBothDemos runs both demos for comparison
func RunBothDemos() {
	RunInterviewStyleDemo()
	fmt.Println("\n══════════════════════════════════════════════════════════════════════\n")
	RunIdiomaticStyleDemo()

	// Comparison summary
	fmt.Println("\n╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║                    COMPARISON SUMMARY                         ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	fmt.Println(`
┌─────────────────────┬──────────────────────┬──────────────────────┐
│ Aspect              │ Interview Style      │ Idiomatic Go         │
├─────────────────────┼──────────────────────┼──────────────────────┤
│ Vehicle Design      │ Interface + Base     │ Simple struct        │
│                     │ struct + embedding   │                      │
├─────────────────────┼──────────────────────┼──────────────────────┤
│ Spot Design         │ SmallSpot, MediumSpot│ Single SimpleSpot    │
│                     │ LargeSpot classes    │ with Size field      │
├─────────────────────┼──────────────────────┼──────────────────────┤
│ ParkingLot          │ Singleton pattern    │ Dependency injection │
├─────────────────────┼──────────────────────┼──────────────────────┤
│ Strategy Pattern    │ ✓ (same)             │ ✓ (same)             │
├─────────────────────┼──────────────────────┼──────────────────────┤
│ Best For            │ LLD Interviews       │ Production Go        │
└─────────────────────┴──────────────────────┴──────────────────────┘
`)
}
