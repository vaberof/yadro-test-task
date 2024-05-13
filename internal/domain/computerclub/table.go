package computerclub

import (
	"fmt"
	"time"
)

type TableId int

func (tableId *TableId) Int() int {
	return int(*tableId)
}

type Table struct {
	Id              TableId
	State           uint8
	Profit          int
	StartTime       time.Time
	EndTime         time.Time
	UsageTimePerDay time.Duration
}

func (t *Table) calculateProfit(pricePerHour int) {
	if t.EndTime.Before(t.StartTime) {
		t.EndTime = t.EndTime.Add(24 * time.Hour)
	}

	usageTime := t.EndTime.Sub(t.StartTime)

	hours := int(usageTime.Hours())
	minutes := int(usageTime.Minutes()) % 60

	t.Profit += hours * pricePerHour
	if minutes > 0 {
		t.Profit += pricePerHour
	}
}

func (t *Table) calculateUsageTime() {
	if t.EndTime.Before(t.StartTime) {
		t.EndTime = t.EndTime.Add(24 * time.Hour)
	}
	t.UsageTimePerDay += t.EndTime.Sub(t.StartTime)
}

func (t *Table) usageTimePerDayString() string {
	hours := int(t.UsageTimePerDay.Hours())
	minutes := int(t.UsageTimePerDay.Minutes()) % 60
	return fmt.Sprintf("%02d:%02d", hours, minutes)
}
