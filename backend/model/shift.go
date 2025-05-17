package model

import (
	"errors"
	"time"
)

// Business Requirements
// These are basic business requirements:
// ● Workers cannot request shifts already assigned to someone else
// ● No overlapping shift requests allowed per worker
// ● Max 1 shift per day, max 5 shifts per week per worker
// ● Admin can override or reassign approved shifts
// ● Conflict checking must occur on both worker request and admin approval
// ● Shift times are stored and compared in UTC

type Shift struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Date        string    `gorm:"not null" json:"date"`       // Format: YYYY-MM-DD
	StartTime   string    `gorm:"not null" json:"start_time"` // Format: HH:MM
	EndTime     string    `gorm:"not null" json:"end_time"`   // Format: HH:MM
	Role        string    `gorm:"not null" json:"role"`       // cashier|delivery|driver
	Status      string    `gorm:"not null" json:"status"`     // pending|accepted|rejected
	WorkerID    int64     `gorm:"not null" json:"worker_id"`
	Assigned    bool      `gorm:"default:false" json:"assigned"`
	Overridden  bool      `gorm:"default:false" json:"overridden"`
	RequestedAt time.Time `json:"requested_at"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}

// Add new error
var (
	ErrValidation      = errors.New("invalid_request")
	ErrTimeConflict    = errors.New("time_conflict")
	ErrShiftLimit      = errors.New("shift_limit")
	ErrInvalidRole     = errors.New("invalid_role")
	ErrInvalidState    = errors.New("invalid_status")
	ErrAlreadyAssigned = errors.New("shift_already_assigned")
)

// Update error messages
var ErrMessages = map[error]string{
	ErrValidation:   "Invalid request format",
	ErrTimeConflict: "Shift times conflict with existing schedule",
	ErrShiftLimit:   "Maximum shift limit exceeded",
	ErrInvalidRole:  "Role must be cashier, delivery, or driver",
	ErrInvalidState: "Status must be pending, accepted, or rejected",
}

const (
	RoleCashier  = "cashier"
	RoleDelivery = "delivery"
	RoleDriver   = "driver"
)

const (
	StatusPending  = "pending"
	StatusApproved = "approved"
	StatusRejected = "rejected"
)

const DateTimeLayout = "2006-01-02 15:04"

func (s Shift) CreateShiftValidation() error {
	// Check required fields
	switch 0 {
	case
		len(s.Date),
		len(s.StartTime),
		len(s.EndTime),
		len(s.Role):

		return ErrValidation
	}

	// Combined datetime parsing and validation
	start, err := time.Parse(DateTimeLayout, s.Date+" "+s.StartTime)
	if err != nil {
		return ErrValidation
	}
	end, err := time.Parse(DateTimeLayout, s.Date+" "+s.EndTime)
	if err != nil {
		return ErrValidation
	}

	// Time sequence check
	startUTC, endUTC := start.UTC(), end.UTC()
	if !endUTC.After(startUTC) {
		return ErrTimeConflict
	}

	// Role validation
	switch s.Role {
	case
		RoleCashier,
		RoleDelivery,
		RoleDriver:
		// Valid role
	default:
		return ErrInvalidRole
	}

	return nil
}

func (s Shift) VerifyOverlap(existing []Shift) error {
	newStart, newEnd, _ := s.parseTimes()

	for _, shift := range existing {
		existingStart, existingEnd, _ := shift.parseTimes()

		if newStart.Before(existingEnd) && newEnd.After(existingStart) {
			return ErrShiftLimit
		}
	}
	return nil
}

const (
	LimitDaily  = 1
	LimitWeekly = 5
)

func (s Shift) VerifyShiftLimits(existing []Shift) error {
	newStart, _, _ := s.parseTimes()

	var dailyCount, weeklyCount int

	for _, shift := range existing {
		existingStart, _, _ := shift.parseTimes()

		if sameDay(newStart, existingStart) {
			dailyCount++
		}
		if sameWeek(newStart, existingStart) {
			weeklyCount++
		}
	}

	if dailyCount >= LimitDaily || weeklyCount >= LimitWeekly {
		return ErrShiftLimit
	}
	return nil
}

var TimeUTC = time.UTC

func (s Shift) parseTimes() (time.Time, time.Time, error) {
	start, err := time.ParseInLocation(DateTimeLayout, s.Date+" "+s.StartTime, TimeUTC)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	end, err := time.ParseInLocation(DateTimeLayout, s.Date+" "+s.EndTime, TimeUTC)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	return start, end, nil
}

func sameDay(a, b time.Time) bool {
	return a.Year() == b.Year() && a.YearDay() == b.YearDay()
}

func sameWeek(a, b time.Time) bool {
	y1, w1 := a.ISOWeek()
	y2, w2 := b.ISOWeek()
	return y1 == y2 && w1 == w2
}

func (s Shift) CheckAssignment() error {
	if s.Assigned && !s.Overridden {
		return ErrAlreadyAssigned
	}
	return nil
}
