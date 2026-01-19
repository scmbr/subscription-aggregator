package dto

import (
	"fmt"
	"strings"
	"time"
)

type MonthYear struct {
	time.Time
}

func (m *MonthYear) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), `"`)
	t, err := time.Parse("01-2006", str)
	if err != nil {
		return err
	}
	m.Time = t
	return nil
}

func (m MonthYear) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%02d-%d"`, m.Month(), m.Year())), nil
}
