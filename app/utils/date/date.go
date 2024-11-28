package date

import (
	"fmt"
	"time"
)

type DateOnly struct {
	time.Time
}

func (date *DateOnly) UnmarshalJSON(input []byte) error {
	formattedTime := fmt.Sprintf(`"%s"`, time.DateOnly)
	parsedDate, err := time.Parse(formattedTime, string(input))
	if err != nil {
		return err
	}

	date.Time = parsedDate
	return nil
}

func NewDateOnlyFromTime(fullTime time.Time) DateOnly {
	dateOnly := DateOnly{}
	dateOnly.Time = fullTime
	return dateOnly
}
