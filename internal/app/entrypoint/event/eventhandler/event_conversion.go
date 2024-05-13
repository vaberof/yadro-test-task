package eventhandler

import (
	"errors"
	"fmt"
	"github.com/vaberof/yadro-test-task/pkg/xtime"
	"strconv"
	"strings"
)

var ErrInvalidEventLine = errors.New("invalid event line")

func FromEventLine(eventLine string) (*Event, error) {
	splitEventLine := strings.Split(eventLine, " ")
	if len(splitEventLine) != 3 && len(splitEventLine) != 4 {
		return nil, ErrInvalidEventLine
	}

	eventTime, err := xtime.ParseHoursMinutesFromString(splitEventLine[0])
	if err != nil {
		return nil, fmt.Errorf("failed to parse event time: %w", err)
	}

	eventType, err := strconv.Atoi(splitEventLine[1])
	if err != nil {
		return nil, fmt.Errorf("failed to convert event type: %w", err)
	}

	tableId := 0
	if len(splitEventLine) == 4 {
		tableId, err = strconv.Atoi(splitEventLine[3])
		if err != nil {
			return nil, fmt.Errorf("failed to convert tableId: %w", err)
		}
	}

	clientName := splitEventLine[2]

	event := &Event{
		Time:       eventTime,
		Type:       uint8(eventType),
		ClientName: clientName,
		TableId:    tableId,
	}

	return event, nil
}
