package models

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
	"github.com/jinzhu/now"
)

var (
	DailySalary = envy.Get("DAILY_SALARY", "0")
)

// Absence model struct.
type Absence struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Date      time.Time `json:"date" db:"date"`
	Reason    string    `json:"reason" db:"reason"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Absences array model struct of Absence.
type Absences []Absence

// String converts the struct into a string value.
func (a Absence) String() string {
	ja, err := json.Marshal(a)
	if err != nil {
		return ""
	}

	return string(ja)
}

// Validate checks for valid values in DeviceModel fields
func (ab *Absence) Validate(tx *pop.Connection) *validate.Errors {
	return validate.Validate(
		&validators.StringIsPresent{Field: ab.Reason, Name: "Reason", Message: "Hey, y la razón papi?"},

		&validators.FuncValidator{
			Name:    "Date",
			Message: "%v Esta fecha ya la agregaste",
			Fn: func() bool {
				exist, _ := tx.Where("date = ?", ab.Date.Format("01-02-2006")).Exists(&Absence{})

				return !exist
			},
		},
	)
}

func CalculateSalary(tx *pop.Connection) (int, error) {
	var total int
	today := time.Now()
	firstOfMonth := now.BeginningOfMonth()
	absenceRangeStart := today
	dailySalary, _ := strconv.Atoi(DailySalary)

	if today.Day() <= 15 {
		for i := firstOfMonth.Day(); i <= today.Day(); i++ {
			if IsSunday(i) {
				continue
			}

			total += dailySalary
		}

		// This is to check that first day of the fortnight is not a Sunday from 1 to 15
		if IsSunday(firstOfMonth.Day()) {
			firstOfMonth = firstOfMonth.Add(24 * time.Hour)
		}

		absenceRangeStart = firstOfMonth
	}

	if today.Day() > 15 && today.Day() <= 31 {
		for i := 16; i <= today.Day(); i++ {
			if IsSunday(i) {
				continue
			}

			total += dailySalary
		}

		// This is to check that first day of the fortnight is not a Sunday from 16 to 30/31
		absenceRangeStartDay := 16
		if IsSunday(16) {
			absenceRangeStartDay = 17
		}

		absenceRangeStart = time.Date(today.Year(), today.Month(), absenceRangeStartDay, 0, 0, 0, 0, today.Location())
	}

	absences := Absences{}
	if err := tx.Where("date >= ?", absenceRangeStart.Format("01-02-2006")).Where("date <= ?", time.Now().Format("01-02-2006")).All(&absences); err != nil {
		return 0, err
	}

	total = total - (len(absences) * dailySalary)
	return total, nil
}

func IsSunday(monthDay int) bool {
	today := time.Now()
	currentYear, currentMonth, _ := today.Date()
	currentLocation := today.Location()

	loopDay := time.Date(currentYear, currentMonth, monthDay, 0, 0, 0, 0, currentLocation)

	return loopDay.Weekday() == time.Sunday
}

func CurrentFortnightDetails(tx *pop.Connection) ([]time.Time, error) {
	today := time.Now()
	firstOfMonth := now.BeginningOfMonth()
	workedDays := []time.Time{}

	if today.Day() <= 15 {
		for i := firstOfMonth.Day(); i <= today.Day(); i++ {
			if IsSunday(i) {
				continue
			}

			currentDay := time.Date(today.Year(), today.Month(), i, 0, 0, 0, 0, today.Location())
			isAbsence, err := tx.Where("date = ?", currentDay.Format("01-02-2006")).Exists(&Absence{})

			if err != nil {
				return []time.Time{}, err
			}

			if !isAbsence {
				workedDays = append(workedDays, currentDay)
			}
		}
	}

	if today.Day() > 15 && today.Day() <= 31 {
		for i := 16; i <= today.Day(); i++ {
			if IsSunday(i) {
				continue
			}

			currentDay := time.Date(today.Year(), today.Month(), i, 0, 0, 0, 0, today.Location())
			isAbsence, err := tx.Where("date = ?", currentDay.Format("01-02-2006")).Exists(&Absence{})

			if err != nil {
				return []time.Time{}, err
			}

			if !isAbsence {
				workedDays = append(workedDays, currentDay)
			}
		}
	}

	return workedDays, nil
}
