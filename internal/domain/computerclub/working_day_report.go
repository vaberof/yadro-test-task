package computerclub

import (
	"fmt"
	"time"
)

type WorkingDayReport []byte

const layoutHoursMinutes = "15:04"

func (w *WorkingDayReport) writeEvent(eventTime time.Time, eventType uint8, clientName ClientName) {
	*w = append(*w, []byte(w.buildEvent(eventTime, eventType, clientName))...)
}

func (w *WorkingDayReport) writeEventWithTableId(eventTime time.Time, eventType uint8, clientName ClientName, tableId TableId) {
	*w = append(*w, []byte(w.buildEventWithTableId(eventTime, eventType, clientName, tableId))...)
}

func (w *WorkingDayReport) writeEventError(eventTime time.Time, err error) {
	*w = append(*w, []byte(w.buildEventError(eventTime, err))...)
}

func (w *WorkingDayReport) writeTime(time time.Time) {
	*w = append(*w, []byte(w.buildTime(time))...)
}

func (w *WorkingDayReport) writeTableReport(tableId TableId, profit int, usageTimeStr string) {
	*w = append(*w, []byte(w.buildTableReport(tableId, profit, usageTimeStr))...)
}

func (w *WorkingDayReport) buildEvent(eventTime time.Time, eventType uint8, clientName ClientName) string {
	return fmt.Sprintf("%s %d %s\n", eventTime.Format(layoutHoursMinutes), eventType, clientName.String())
}

func (w *WorkingDayReport) buildEventWithTableId(eventTime time.Time, eventType uint8, clientName ClientName, tableId TableId) string {
	return fmt.Sprintf("%s %d %s %d\n", eventTime.Format(layoutHoursMinutes), eventType, clientName.String(), tableId.Int())
}

func (w *WorkingDayReport) buildEventError(eventTime time.Time, err error) string {
	return fmt.Sprintf("%s %d %s\n", eventTime.Format(layoutHoursMinutes), OutgoingEventError, err.Error())
}

func (w *WorkingDayReport) buildTime(time time.Time) string {
	return fmt.Sprintf("%s\n", time.Format(layoutHoursMinutes))
}

func (w *WorkingDayReport) buildTableReport(tableId TableId, profit int, usageTime string) string {
	return fmt.Sprintf("%d %d %s\n", tableId.Int(), profit, usageTime)
}
