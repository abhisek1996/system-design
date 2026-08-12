Requirement: 
- user
- location(city)
- movie
- theater
    - seat

**Objects:**
- user
- city
- movie
- theater
- screen(hall)
- shows
- seat
- 
- booking
- payment


top --> bottom approve
- movie -> theater --> screen --> show --> booking --> payment



# Movie Ticket Booking System — Low-Level Design (Go)

---

## Table of Contents

1. [Requirements](#requirements)
2. [Core Entities & Enums](#core-entities--enums)
3. [Data Models (Structs)](#data-models-structs)
4. [Interfaces](#interfaces)
5. [Service Layer](#service-layer)
6. [Repository Layer](#repository-layer)
7. [Concurrency & Seat Locking](#concurrency--seat-locking)
8. [Design Patterns Used](#design-patterns-used)
9. [Class Diagram (ASCII)](#class-diagram-ascii)
10. [Complete Go Code](#complete-go-code)
11. [Key Design Decisions & Trade-offs](#key-design-decisions--trade-offs)

---

## Requirements

### Functional Requirements
- Browse movies currently showing
- View cinemas and their screens/halls
- View shows (movie + screen + time slot)
- Search shows by movie, city, or date
- Select seats interactively from a seat map
- Temporarily lock seats during booking (prevent double-booking)
- Book tickets with payment processing
- Cancel bookings and trigger refund
- Apply discount/promo codes
- Generate booking confirmation with unique ID
- Support multiple seat types (REGULAR, PREMIUM, RECLINER)
- Admin: add movies, cinemas, screens, shows

### Non-Functional Requirements
- Thread-safe seat reservation (no two users book the same seat)
- Seat lock expires after a configurable TTL (e.g., 10 minutes)
- Horizontal scalability (stateless services)
- Idempotent booking to handle retries

---

## Core Entities & Enums

```go
// models/enums.go

package models

type SeatType string
const (
    SeatTypeRegular  SeatType = "REGULAR"
    SeatTypePremium  SeatType = "PREMIUM"
    SeatTypeRecliner SeatType = "RECLINER"
)

type SeatStatus string
const (
    SeatStatusAvailable SeatStatus = "AVAILABLE"
    SeatStatusLocked    SeatStatus = "LOCKED"
    SeatStatusBooked    SeatStatus = "BOOKED"
)

type BookingStatus string
const (
    BookingStatusPending   BookingStatus = "PENDING"
    BookingStatusConfirmed BookingStatus = "CONFIRMED"
    BookingStatusCancelled BookingStatus = "CANCELLED"
    BookingStatusFailed    BookingStatus = "FAILED"
)

type PaymentStatus string
const (
    PaymentStatusPending  PaymentStatus = "PENDING"
    PaymentStatusSuccess  PaymentStatus = "SUCCESS"
    PaymentStatusFailed   PaymentStatus = "FAILED"
    PaymentStatusRefunded PaymentStatus = "REFUNDED"
)

type PaymentMethod string
const (
    PaymentMethodCard   PaymentMethod = "CARD"
    PaymentMethodUPI    PaymentMethod = "UPI"
    PaymentMethodWallet PaymentMethod = "WALLET"
)

type Genre string
const (
    GenreAction  Genre = "ACTION"
    GenreDrama   Genre = "DRAMA"
    GenreComedy  Genre = "COMEDY"
    GenreHorror  Genre = "HORROR"
    GenreSci_Fi  Genre = "SCI_FI"
)
```

---

## Data Models (Structs)

```go
// models/models.go

package models

import "time"

// ─── User ────────────────────────────────────────────────────────────────────

type User struct {
    ID        string
    Name      string
    Email     string
    Phone     string
    CreatedAt time.Time
}

// ─── Movie ───────────────────────────────────────────────────────────────────

type Movie struct {
    ID          string
    Title       string
    Description string
    DurationMin int      // duration in minutes
    Language    string
    Genre       Genre
    Rating      float32  // e.g., 8.4
    ReleaseDate time.Time
    IsActive    bool
}

// ─── City / Cinema / Screen ──────────────────────────────────────────────────

type City struct {
    ID   string
    Name string
}

type Cinema struct {
    ID      string
    Name    string
    City    City
    Address string
    Screens []Screen
}

type Screen struct {
    ID       string
    Name     string          // "Screen 1", "IMAX Hall"
    CinemaID string
    Rows     int
    Cols     int
    Seats    [][]Seat        // 2D seat grid [row][col]
}

// ─── Seat ────────────────────────────────────────────────────────────────────

type Seat struct {
    ID       string
    ScreenID string
    Row      int
    Col      int
    Label    string     // e.g., "A1", "B12"
    Type     SeatType
}

// ─── Show ────────────────────────────────────────────────────────────────────

type Show struct {
    ID        string
    Movie     Movie
    Screen    Screen
    StartTime time.Time
    EndTime   time.Time
    Language  string
    Format    string     // "2D", "3D", "IMAX"
    IsActive  bool
}

// ─── ShowSeat (runtime seat state per show) ──────────────────────────────────

type ShowSeat struct {
    ID         string
    ShowID     string
    Seat       Seat
    Status     SeatStatus
    Price      float64
    LockedBy   string     // userID who locked it
    LockedAt   *time.Time
    LockExpiry *time.Time
}

// ─── Booking ─────────────────────────────────────────────────────────────────

type Booking struct {
    ID            string
    UserID        string
    ShowID        string
    ShowSeats     []ShowSeat
    TotalAmount   float64
    DiscountAmount float64
    FinalAmount   float64
    PromoCode     string
    Status        BookingStatus
    Payment       *Payment
    CreatedAt     time.Time
    UpdatedAt     time.Time
}

// ─── Payment ─────────────────────────────────────────────────────────────────

type Payment struct {
    ID            string
    BookingID     string
    Amount        float64
    Method        PaymentMethod
    Status        PaymentStatus
    TransactionID string
    CreatedAt     time.Time
}

// ─── Promo Code ──────────────────────────────────────────────────────────────

type PromoCode struct {
    Code           string
    DiscountPercent float64
    MaxDiscount    float64
    ValidFrom      time.Time
    ValidTo        time.Time
    IsActive       bool
    UsageLimit     int
    UsedCount      int
}

// ─── Notification ────────────────────────────────────────────────────────────

type Notification struct {
    ID        string
    UserID    string
    BookingID string
    Message   string
    SentAt    time.Time
}
```

---

## Interfaces

```go
// interfaces/repositories.go

package interfaces

import (
    "time"
    "moviebooking/models"
)

type MovieRepository interface {
    Save(movie models.Movie) error
    FindByID(id string) (*models.Movie, error)
    FindAll(active bool) ([]models.Movie, error)
    Update(movie models.Movie) error
}

type CinemaRepository interface {
    Save(cinema models.Cinema) error
    FindByID(id string) (*models.Cinema, error)
    FindByCity(cityID string) ([]models.Cinema, error)
}

type ShowRepository interface {
    Save(show models.Show) error
    FindByID(id string) (*models.Show, error)
    FindByMovieAndCity(movieID, cityID string, date time.Time) ([]models.Show, error)
    FindByCinema(cinemaID string) ([]models.Show, error)
}

type ShowSeatRepository interface {
    SaveAll(seats []models.ShowSeat) error
    FindByShow(showID string) ([]models.ShowSeat, error)
    FindByIDForUpdate(id string) (*models.ShowSeat, error) // row-level lock hint
    UpdateStatus(id string, status models.SeatStatus) error
    BulkUpdateStatus(ids []string, status models.SeatStatus) error
    LockSeats(ids []string, userID string, expiry time.Time) error
    ReleaseExpiredLocks() error
}

type BookingRepository interface {
    Save(booking models.Booking) error
    FindByID(id string) (*models.Booking, error)
    FindByUser(userID string) ([]models.Booking, error)
    UpdateStatus(id string, status models.BookingStatus) error
}

type PaymentRepository interface {
    Save(payment models.Payment) error
    FindByBooking(bookingID string) (*models.Payment, error)
    UpdateStatus(id string, status models.PaymentStatus, txnID string) error
}

type PromoRepository interface {
    FindByCode(code string) (*models.PromoCode, error)
    IncrementUsage(code string) error
}
```

```go
// interfaces/services.go

package interfaces

import (
    "time"
    "moviebooking/models"
)

type MovieService interface {
    AddMovie(movie models.Movie) (models.Movie, error)
    GetMovie(id string) (models.Movie, error)
    ListMovies() ([]models.Movie, error)
}

type ShowService interface {
    CreateShow(show models.Show) (models.Show, error)
    GetShow(id string) (models.Show, error)
    SearchShows(movieID, cityID string, date time.Time) ([]models.Show, error)
    GetSeatMap(showID string) ([]models.ShowSeat, error)
}

type BookingService interface {
    InitiateBooking(userID, showID string, seatIDs []string) (models.Booking, error)
    ConfirmBooking(bookingID string, paymentReq PaymentRequest) (models.Booking, error)
    CancelBooking(bookingID string) (models.Booking, error)
}

type PaymentService interface {
    ProcessPayment(req PaymentRequest) (models.Payment, error)
    RefundPayment(paymentID string) (models.Payment, error)
}

type NotificationService interface {
    SendBookingConfirmation(booking models.Booking) error
    SendCancellationNotice(booking models.Booking) error
}

// DTOs
type PaymentRequest struct {
    BookingID string
    Amount    float64
    Method    models.PaymentMethod
    Token     string // card token / UPI VPA
}
```

---

## Service Layer

```go
// services/booking_service.go

package services

import (
    "errors"
    "fmt"
    "time"

    "moviebooking/interfaces"
    "moviebooking/models"
    "github.com/google/uuid"
)

const SeatLockTTL = 10 * time.Minute

type bookingService struct {
    showSeatRepo  interfaces.ShowSeatRepository
    bookingRepo   interfaces.BookingRepository
    promoRepo     interfaces.PromoRepository
    paymentSvc    interfaces.PaymentService
    notifySvc     interfaces.NotificationService
    seatLocker    SeatLocker   // interface for distributed lock
}

func NewBookingService(
    showSeatRepo interfaces.ShowSeatRepository,
    bookingRepo interfaces.BookingRepository,
    promoRepo interfaces.PromoRepository,
    paymentSvc interfaces.PaymentService,
    notifySvc interfaces.NotificationService,
    locker SeatLocker,
) interfaces.BookingService {
    return &bookingService{
        showSeatRepo: showSeatRepo,
        bookingRepo:  bookingRepo,
        promoRepo:    promoRepo,
        paymentSvc:   paymentSvc,
        notifySvc:    notifySvc,
        seatLocker:   locker,
    }
}

// InitiateBooking — locks seats and creates a PENDING booking
func (b *bookingService) InitiateBooking(userID, showID string, seatIDs []string) (models.Booking, error) {
    // 1. Acquire distributed locks for each seat (sorted to prevent deadlock)
    locked, err := b.seatLocker.LockSeats(seatIDs, userID, SeatLockTTL)
    if err != nil || !locked {
        return models.Booking{}, errors.New("one or more seats are unavailable")
    }

    // 2. Fetch seat details & validate availability
    total := 0.0
    var showSeats []models.ShowSeat
    for _, id := range seatIDs {
        ss, err := b.showSeatRepo.FindByIDForUpdate(id)
        if err != nil {
            b.seatLocker.UnlockSeats(seatIDs)
            return models.Booking{}, fmt.Errorf("seat %s not found", id)
        }
        if ss.Status != models.SeatStatusAvailable {
            b.seatLocker.UnlockSeats(seatIDs)
            return models.Booking{}, fmt.Errorf("seat %s is not available", id)
        }
        total += ss.Price
        showSeats = append(showSeats, *ss)
    }

    // 3. Mark seats as LOCKED in DB
    expiry := time.Now().Add(SeatLockTTL)
    if err := b.showSeatRepo.LockSeats(seatIDs, userID, expiry); err != nil {
        b.seatLocker.UnlockSeats(seatIDs)
        return models.Booking{}, err
    }

    // 4. Create PENDING booking
    booking := models.Booking{
        ID:          uuid.NewString(),
        UserID:      userID,
        ShowID:      showID,
        ShowSeats:   showSeats,
        TotalAmount: total,
        FinalAmount: total,
        Status:      models.BookingStatusPending,
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }
    if err := b.bookingRepo.Save(booking); err != nil {
        b.seatLocker.UnlockSeats(seatIDs)
        return models.Booking{}, err
    }
    return booking, nil
}

// ConfirmBooking — applies promo, processes payment, marks seats BOOKED
func (b *bookingService) ConfirmBooking(bookingID string, payReq interfaces.PaymentRequest) (models.Booking, error) {
    booking, err := b.bookingRepo.FindByID(bookingID)
    if err != nil {
        return models.Booking{}, err
    }
    if booking.Status != models.BookingStatusPending {
        return models.Booking{}, errors.New("booking is not in PENDING state")
    }

    // Apply promo if provided
    if booking.PromoCode != "" {
        discount, err := b.applyPromo(booking.PromoCode, booking.TotalAmount)
        if err == nil {
            booking.DiscountAmount = discount
            booking.FinalAmount = booking.TotalAmount - discount
        }
    }
    payReq.Amount = booking.FinalAmount

    // Process payment
    payment, err := b.paymentSvc.ProcessPayment(payReq)
    if err != nil {
        _ = b.bookingRepo.UpdateStatus(bookingID, models.BookingStatusFailed)
        return models.Booking{}, fmt.Errorf("payment failed: %w", err)
    }
    booking.Payment = &payment

    // Mark seats as BOOKED
    var seatIDs []string
    for _, ss := range booking.ShowSeats {
        seatIDs = append(seatIDs, ss.ID)
    }
    if err := b.showSeatRepo.BulkUpdateStatus(seatIDs, models.SeatStatusBooked); err != nil {
        return models.Booking{}, err
    }

    // Confirm booking
    booking.Status = models.BookingStatusConfirmed
    booking.UpdatedAt = time.Now()
    _ = b.bookingRepo.UpdateStatus(bookingID, models.BookingStatusConfirmed)

    // Send notification (async)
    go b.notifySvc.SendBookingConfirmation(*booking)

    return *booking, nil
}

// CancelBooking — refunds payment, releases seats
func (b *bookingService) CancelBooking(bookingID string) (models.Booking, error) {
    booking, err := b.bookingRepo.FindByID(bookingID)
    if err != nil {
        return models.Booking{}, err
    }
    if booking.Status != models.BookingStatusConfirmed {
        return models.Booking{}, errors.New("only confirmed bookings can be cancelled")
    }

    // Refund
    if booking.Payment != nil {
        _, err := b.paymentSvc.RefundPayment(booking.Payment.ID)
        if err != nil {
            return models.Booking{}, fmt.Errorf("refund failed: %w", err)
        }
    }

    // Release seats
    var seatIDs []string
    for _, ss := range booking.ShowSeats {
        seatIDs = append(seatIDs, ss.ID)
    }
    _ = b.showSeatRepo.BulkUpdateStatus(seatIDs, models.SeatStatusAvailable)

    // Update booking status
    booking.Status = models.BookingStatusCancelled
    booking.UpdatedAt = time.Now()
    _ = b.bookingRepo.UpdateStatus(bookingID, models.BookingStatusCancelled)

    go b.notifySvc.SendCancellationNotice(*booking)

    return *booking, nil
}

func (b *bookingService) applyPromo(code string, amount float64) (float64, error) {
    promo, err := b.promoRepo.FindByCode(code)
    if err != nil || !promo.IsActive {
        return 0, errors.New("invalid promo code")
    }
    if time.Now().Before(promo.ValidFrom) || time.Now().After(promo.ValidTo) {
        return 0, errors.New("promo code expired")
    }
    if promo.UsedCount >= promo.UsageLimit {
        return 0, errors.New("promo usage limit reached")
    }
    discount := amount * promo.DiscountPercent / 100
    if discount > promo.MaxDiscount {
        discount = promo.MaxDiscount
    }
    _ = b.promoRepo.IncrementUsage(code)
    return discount, nil
}
```

```go
// services/show_service.go

package services

import (
    "time"
    "fmt"

    "moviebooking/interfaces"
    "moviebooking/models"
    "github.com/google/uuid"
)

type showService struct {
    showRepo     interfaces.ShowRepository
    showSeatRepo interfaces.ShowSeatRepository
}

func NewShowService(sr interfaces.ShowRepository, ssr interfaces.ShowSeatRepository) interfaces.ShowService {
    return &showService{showRepo: sr, showSeatRepo: ssr}
}

func (s *showService) CreateShow(show models.Show) (models.Show, error) {
    show.ID = uuid.NewString()
    if err := s.showRepo.Save(show); err != nil {
        return models.Show{}, err
    }
    // Auto-generate ShowSeats from screen's seat layout
    showSeats := generateShowSeats(show)
    return show, s.showSeatRepo.SaveAll(showSeats)
}

func (s *showService) GetShow(id string) (models.Show, error) {
    show, err := s.showRepo.FindByID(id)
    if err != nil {
        return models.Show{}, fmt.Errorf("show not found: %w", err)
    }
    return *show, nil
}

func (s *showService) SearchShows(movieID, cityID string, date time.Time) ([]models.Show, error) {
    return s.showRepo.FindByMovieAndCity(movieID, cityID, date)
}

func (s *showService) GetSeatMap(showID string) ([]models.ShowSeat, error) {
    return s.showSeatRepo.FindByShow(showID)
}

func generateShowSeats(show models.Show) []models.ShowSeat {
    var seats []models.ShowSeat
    for _, row := range show.Screen.Seats {
        for _, seat := range row {
            seats = append(seats, models.ShowSeat{
                ID:     uuid.NewString(),
                ShowID: show.ID,
                Seat:   seat,
                Status: models.SeatStatusAvailable,
                Price:  priceForSeatType(seat.Type),
            })
        }
    }
    return seats
}

func priceForSeatType(t models.SeatType) float64 {
    switch t {
    case models.SeatTypeRecliner:
        return 500
    case models.SeatTypePremium:
        return 350
    default:
        return 200
    }
}
```

```go
// services/payment_service.go

package services

import (
    "errors"
    "time"

    "moviebooking/interfaces"
    "moviebooking/models"
    "github.com/google/uuid"
)

type paymentService struct {
    paymentRepo interfaces.PaymentRepository
    gateway     PaymentGateway // external gateway abstraction
}

func NewPaymentService(repo interfaces.PaymentRepository, gw PaymentGateway) interfaces.PaymentService {
    return &paymentService{paymentRepo: repo, gateway: gw}
}

func (p *paymentService) ProcessPayment(req interfaces.PaymentRequest) (models.Payment, error) {
    txnID, err := p.gateway.Charge(req.Token, req.Amount, req.Method)
    if err != nil {
        return models.Payment{}, errors.New("gateway charge failed")
    }

    payment := models.Payment{
        ID:            uuid.NewString(),
        BookingID:     req.BookingID,
        Amount:        req.Amount,
        Method:        req.Method,
        Status:        models.PaymentStatusSuccess,
        TransactionID: txnID,
        CreatedAt:     time.Now(),
    }
    if err := p.paymentRepo.Save(payment); err != nil {
        return models.Payment{}, err
    }
    return payment, nil
}

func (p *paymentService) RefundPayment(paymentID string) (models.Payment, error) {
    payment, err := p.paymentRepo.FindByBooking(paymentID)
    if err != nil {
        return models.Payment{}, err
    }
    if err := p.gateway.Refund(payment.TransactionID); err != nil {
        return models.Payment{}, err
    }
    _ = p.paymentRepo.UpdateStatus(paymentID, models.PaymentStatusRefunded, "")
    payment.Status = models.PaymentStatusRefunded
    return *payment, nil
}
```

---

## Repository Layer

```go
// repositories/in_memory_show_seat_repo.go
// (In-memory implementation; swap with DB implementation for production)

package repositories

import (
    "errors"
    "sort"
    "sync"
    "time"

    "moviebooking/models"
)

type InMemoryShowSeatRepo struct {
    mu    sync.RWMutex
    seats map[string]*models.ShowSeat // key: ShowSeat.ID
}

func NewInMemoryShowSeatRepo() *InMemoryShowSeatRepo {
    return &InMemoryShowSeatRepo{seats: make(map[string]*models.ShowSeat)}
}

func (r *InMemoryShowSeatRepo) SaveAll(seats []models.ShowSeat) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    for i := range seats {
        r.seats[seats[i].ID] = &seats[i]
    }
    return nil
}

func (r *InMemoryShowSeatRepo) FindByShow(showID string) ([]models.ShowSeat, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    var result []models.ShowSeat
    for _, s := range r.seats {
        if s.ShowID == showID {
            result = append(result, *s)
        }
    }
    return result, nil
}

func (r *InMemoryShowSeatRepo) FindByIDForUpdate(id string) (*models.ShowSeat, error) {
    r.mu.Lock() // exclusive lock simulates SELECT FOR UPDATE
    defer r.mu.Unlock()
    s, ok := r.seats[id]
    if !ok {
        return nil, errors.New("show seat not found")
    }
    return s, nil
}

func (r *InMemoryShowSeatRepo) UpdateStatus(id string, status models.SeatStatus) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    s, ok := r.seats[id]
    if !ok {
        return errors.New("show seat not found")
    }
    s.Status = status
    return nil
}

func (r *InMemoryShowSeatRepo) BulkUpdateStatus(ids []string, status models.SeatStatus) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    for _, id := range ids {
        if s, ok := r.seats[id]; ok {
            s.Status = status
        }
    }
    return nil
}

func (r *InMemoryShowSeatRepo) LockSeats(ids []string, userID string, expiry time.Time) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    now := time.Now()
    for _, id := range ids {
        s, ok := r.seats[id]
        if !ok {
            return errors.New("seat not found: " + id)
        }
        s.Status = models.SeatStatusLocked
        s.LockedBy = userID
        s.LockedAt = &now
        s.LockExpiry = &expiry
    }
    return nil
}

func (r *InMemoryShowSeatRepo) ReleaseExpiredLocks() error {
    r.mu.Lock()
    defer r.mu.Unlock()
    now := time.Now()
    for _, s := range r.seats {
        if s.Status == models.SeatStatusLocked && s.LockExpiry != nil && now.After(*s.LockExpiry) {
            s.Status = models.SeatStatusAvailable
            s.LockedBy = ""
            s.LockedAt = nil
            s.LockExpiry = nil
        }
    }
    return nil
}

// Ensure seats are locked in sorted order to avoid deadlock
func sortedSeatIDs(ids []string) []string {
    sorted := make([]string, len(ids))
    copy(sorted, ids)
    sort.Strings(sorted)
    return sorted
}
```

---

## Concurrency & Seat Locking

```go
// services/seat_locker.go

package services

import (
    "context"
    "fmt"
    "sort"
    "sync"
    "time"
)

// SeatLocker abstracts distributed locking.
// In production: replace with Redis-based Redlock.
type SeatLocker interface {
    LockSeats(seatIDs []string, userID string, ttl time.Duration) (bool, error)
    UnlockSeats(seatIDs []string) error
}

// ─── In-Memory SeatLocker (for single-node / testing) ────────────────────────

type InMemorySeatLocker struct {
    mu    sync.Mutex
    locks map[string]lockEntry
}

type lockEntry struct {
    ownerID string
    expiry  time.Time
}

func NewInMemorySeatLocker() *InMemorySeatLocker {
    return &InMemorySeatLocker{locks: make(map[string]lockEntry)}
}

// LockSeats acquires locks on all seatIDs atomically (sorted to prevent deadlock)
func (l *InMemorySeatLocker) LockSeats(seatIDs []string, userID string, ttl time.Duration) (bool, error) {
    sorted := sorted(seatIDs)
    l.mu.Lock()
    defer l.mu.Unlock()

    now := time.Now()
    // Check all seats are free
    for _, id := range sorted {
        if entry, exists := l.locks[id]; exists {
            if now.Before(entry.expiry) {
                return false, fmt.Errorf("seat %s is already locked", id)
            }
        }
    }
    // Acquire all
    expiry := now.Add(ttl)
    for _, id := range sorted {
        l.locks[id] = lockEntry{ownerID: userID, expiry: expiry}
    }
    return true, nil
}

func (l *InMemorySeatLocker) UnlockSeats(seatIDs []string) error {
    l.mu.Lock()
    defer l.mu.Unlock()
    for _, id := range seatIDs {
        delete(l.locks, id)
    }
    return nil
}

func sorted(ids []string) []string {
    s := make([]string, len(ids))
    copy(s, ids)
    sort.Strings(s)
    return s
}

// ─── Redis SeatLocker (Production Sketch) ────────────────────────────────────

/*
type RedisSeatLocker struct {
    client *redis.Client
}

func (r *RedisSeatLocker) LockSeats(seatIDs []string, userID string, ttl time.Duration) (bool, error) {
    ctx := context.Background()
    sorted := sorted(seatIDs)
    acquired := []string{}

    for _, id := range sorted {
        key := "seat_lock:" + id
        ok, err := r.client.SetNX(ctx, key, userID, ttl).Result()
        if err != nil || !ok {
            // Rollback acquired locks
            for _, k := range acquired {
                r.client.Del(ctx, "seat_lock:"+k)
            }
            return false, fmt.Errorf("could not lock seat %s", id)
        }
        acquired = append(acquired, id)
    }
    return true, nil
}

func (r *RedisSeatLocker) UnlockSeats(seatIDs []string) error {
    ctx := context.Background()
    for _, id := range seatIDs {
        r.client.Del(ctx, "seat_lock:"+id)
    }
    return nil
}
*/

// ─── Lock Expiry Background Worker ───────────────────────────────────────────

func StartLockExpiryWorker(ctx context.Context, repo interface{ ReleaseExpiredLocks() error }, interval time.Duration) {
    go func() {
        ticker := time.NewTicker(interval)
        defer ticker.Stop()
        for {
            select {
            case <-ticker.C:
                _ = repo.ReleaseExpiredLocks()
            case <-ctx.Done():
                return
            }
        }
    }()
}
```

---

## Design Patterns Used

| Pattern | Where Applied |
|---|---|
| **Repository Pattern** | `MovieRepository`, `ShowRepository`, `BookingRepository` — separates data access from business logic |
| **Service Layer** | `BookingService`, `ShowService`, `PaymentService` — encapsulates business rules |
| **Strategy Pattern** | `PaymentGateway` interface — swap card/UPI/wallet processors at runtime |
| **Factory Method** | `generateShowSeats()` — creates ShowSeat instances from a Screen layout |
| **Observer / Event** | `NotificationService` called asynchronously (via goroutines) after state changes |
| **Mutex / Distributed Lock** | `SeatLocker` — prevents concurrent double-booking; Redis Redlock in production |
| **Singleton** | Repositories & services wired once via constructor injection |
| **Decorator (future)** | Logging/tracing middleware wrapped around service interfaces |

---

## Class Diagram (ASCII)

```
┌─────────────┐      books       ┌───────────────┐
│    User     │─────────────────>│    Booking     │
└─────────────┘                  └───────┬────────┘
                                         │ contains
                                         ▼
                                  ┌─────────────┐
                                  │  ShowSeat   │◄─── status: AVAILABLE/LOCKED/BOOKED
                                  └──────┬──────┘
                                         │ belongs to
                                         ▼
┌─────────────┐    has many      ┌──────────────┐
│    Movie    │────────────────> │     Show     │
└─────────────┘                  └──────┬───────┘
                                         │ runs on
                                         ▼
                               ┌─────────────────┐
                               │     Screen      │
                               └──────┬──────────┘
                                      │ belongs to
                                      ▼
                               ┌─────────────────┐
                               │     Cinema      │
                               └──────┬──────────┘
                                      │ located in
                                      ▼
                               ┌─────────────────┐
                               │      City       │
                               └─────────────────┘

┌───────────────┐     has one   ┌───────────────┐
│    Booking    │──────────────>│    Payment    │
└───────────────┘               └───────────────┘
        │ applies
        ▼
┌───────────────┐
│  PromoCode    │
└───────────────┘
```

---

## Complete Go Code — main.go (Wiring Example)

```go
// main.go

package main

import (
    "context"
    "fmt"
    "time"

    "moviebooking/models"
    "moviebooking/repositories"
    "moviebooking/services"
)

func main() {
    // ── Wire dependencies ────────────────────────────────────────────────────
    showSeatRepo  := repositories.NewInMemoryShowSeatRepo()
    bookingRepo   := repositories.NewInMemoryBookingRepo()
    movieRepo     := repositories.NewInMemoryMovieRepo()
    showRepo      := repositories.NewInMemoryShowRepo()
    promoRepo     := repositories.NewInMemoryPromoRepo()
    paymentRepo   := repositories.NewInMemoryPaymentRepo()

    gateway       := services.NewMockPaymentGateway()
    paymentSvc    := services.NewPaymentService(paymentRepo, gateway)
    notifySvc     := services.NewConsoleNotificationService()
    seatLocker    := services.NewInMemorySeatLocker()
    showSvc       := services.NewShowService(showRepo, showSeatRepo)
    bookingSvc    := services.NewBookingService(showSeatRepo, bookingRepo, promoRepo, paymentSvc, notifySvc, seatLocker)

    // Start expired-lock cleanup worker (every 30 seconds)
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    services.StartLockExpiryWorker(ctx, showSeatRepo, 30*time.Second)

    // ── Seed data ────────────────────────────────────────────────────────────
    movie := models.Movie{
        ID: "m1", Title: "Interstellar", DurationMin: 169,
        Language: "English", Genre: models.GenreSci_Fi, Rating: 8.6,
    }
    _ = movieRepo.Save(movie)

    screen := models.Screen{
        ID: "sc1", Name: "IMAX Hall", CinemaID: "c1",
        Rows: 2, Cols: 3,
        Seats: [][]models.Seat{
            {{ID: "s1", Label: "A1", Type: models.SeatTypeRegular},
             {ID: "s2", Label: "A2", Type: models.SeatTypeRegular},
             {ID: "s3", Label: "A3", Type: models.SeatTypePremium}},
            {{ID: "s4", Label: "B1", Type: models.SeatTypePremium},
             {ID: "s5", Label: "B2", Type: models.SeatTypeRecliner},
             {ID: "s6", Label: "B3", Type: models.SeatTypeRecliner}},
        },
    }

    show, _ := showSvc.CreateShow(models.Show{
        Movie: movie, Screen: screen,
        StartTime: time.Now().Add(2 * time.Hour),
        Format: "IMAX", Language: "English", IsActive: true,
    })

    // ── Booking flow ─────────────────────────────────────────────────────────
    booking, err := bookingSvc.InitiateBooking("user1", show.ID, []string{"s1", "s2"})
    if err != nil {
        fmt.Println("InitiateBooking error:", err)
        return
    }
    fmt.Println("Booking initiated:", booking.ID, "Status:", booking.Status)

    confirmed, err := bookingSvc.ConfirmBooking(booking.ID, interfaces.PaymentRequest{
        BookingID: booking.ID,
        Amount:    booking.FinalAmount,
        Method:    models.PaymentMethodUPI,
        Token:     "user1@upi",
    })
    if err != nil {
        fmt.Println("ConfirmBooking error:", err)
        return
    }
    fmt.Println("Booking confirmed:", confirmed.ID, "| Amount paid:", confirmed.FinalAmount)
}
```

---

## Key Design Decisions & Trade-offs

### 1. Optimistic vs Pessimistic Locking
- **Chosen**: Pessimistic (mutex + DB row lock) for seat reservation.
- **Why**: At high concurrency, optimistic retries cause poor UX for seat booking. A 10-minute TTL lock mimics real-world booking flows.
- **Trade-off**: Slightly lower throughput; mitigated by seat-level granularity (not table-level).

### 2. ShowSeat as a Join Table with State
- `ShowSeat` is NOT a static entity — it carries runtime `Status`, `LockedBy`, `LockExpiry`.
- This separates the physical seat (`Seat`) from its per-show availability state (`ShowSeat`).

### 3. Sorted Lock Acquisition
- Seats are always locked in sorted order to prevent **deadlocks** in concurrent transactions.

### 4. Async Notifications
- `NotificationService` is called inside a goroutine to avoid blocking the booking response.
- In production, push to a message queue (Kafka/SQS) instead.

### 5. Strategy Pattern for Payment Gateway
- `PaymentGateway` interface allows swapping between Razorpay, Stripe, PayU without changing business logic.

### 6. Lock Expiry Worker
- A background goroutine periodically calls `ReleaseExpiredLocks()` to free seats for abandoned bookings (user closed browser mid-flow).

### 7. Idempotency
- Booking IDs are UUID-based. Re-confirming a non-PENDING booking returns an error immediately, preventing duplicate charges.

### 8. Price Tiers
- Seat pricing is determined by `SeatType` at show-creation time and stored in `ShowSeat.Price`, allowing dynamic pricing per show without touching core `Seat` data.

---

*Designed for a LLD interview — covers entities, interfaces, services, repositories, concurrency, and design patterns in idiomatic Go.*



