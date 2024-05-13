package filehandler

import (
	"bufio"
	"errors"
	"github.com/vaberof/yadro-test-task/internal/app/entrypoint/event/eventhandler"
	"github.com/vaberof/yadro-test-task/internal/domain/computerclub"
	"github.com/vaberof/yadro-test-task/pkg/xtime"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	ErrInvalidFormatTablesCount   = errors.New("invalid format of tables count")
	ErrInvalidFormatOpeningHours  = errors.New("invalid format of opening hours")
	ErrInvalidFormatPricePerHour  = errors.New("invalid format of price per hour")
	ErrInvalidFormatClientName    = errors.New("invalid format of client name")
	ErrInvalidFormatTableNumber   = errors.New("invalid format of table number")
	ErrInvalidFormatEvent         = errors.New("invalid format of event")
	ErrInvalidFormatEventSequence = errors.New("invalid format of event sequence")
	ErrInvalidFormatFile          = errors.New("invalid format of file")
)

const (
	minFileLinesCount    = 4
	configLinesCount     = 3
	openingHoursSplitLen = 2

	minSplitEventLineLen = 3
	maxSplitEventLineLen = 4
)

const minTablesCount = 1

type InvalidLine string

type Handler struct {
	eventHandler eventhandler.Handler
}

func ProcessComputerClubConfig(filename string, config *computerclub.Config) (*InvalidLine, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	if fileInfo.Size() < minFileLinesCount {
		return nil, ErrInvalidFormatFile
	}

	scanner := bufio.NewScanner(file)

	return processConfigLines(scanner, config)
}

func processConfigLines(scanner *bufio.Scanner, config *computerclub.Config) (*InvalidLine, error) {
	invalidLine, err := scanTablesCountLine(scanner, config)
	if err != nil {
		return invalidLine, err
	}

	invalidLine, err = scanOpeningHoursLine(scanner, config)
	if err != nil {
		return invalidLine, err
	}

	invalidLine, err = scanPricePerHourLine(scanner, config)
	if err != nil {
		return invalidLine, err
	}

	return nil, nil
}

func scanTablesCountLine(scanner *bufio.Scanner, config *computerclub.Config) (*InvalidLine, error) {
	var invalidLine InvalidLine

	scanner.Scan()
	tablesCountLine := scanner.Text()

	tablesCount, err := parseTablesCountLine(tablesCountLine)
	if err != nil {
		invalidLine = InvalidLine(tablesCountLine)
		return &invalidLine, ErrInvalidFormatTablesCount
	}

	config.TablesCount = tablesCount

	return nil, nil
}

func scanOpeningHoursLine(scanner *bufio.Scanner, config *computerclub.Config) (*InvalidLine, error) {
	var invalidLine InvalidLine

	scanner.Scan()
	openingHoursLine := scanner.Text()

	openingTime, closingTime, err := parseOpeningHoursLine(openingHoursLine)
	if err != nil {
		invalidLine = InvalidLine(openingHoursLine)
		return &invalidLine, ErrInvalidFormatOpeningHours
	}

	config.OpeningTime = openingTime
	config.ClosingTime = closingTime

	return nil, nil
}

func scanPricePerHourLine(scanner *bufio.Scanner, config *computerclub.Config) (*InvalidLine, error) {
	var invalidLine InvalidLine

	scanner.Scan()
	pricePerHourLine := scanner.Text()

	pricePerHour, err := strconv.Atoi(pricePerHourLine)
	if err != nil {
		invalidLine = InvalidLine(pricePerHourLine)
		return &invalidLine, ErrInvalidFormatPricePerHour
	}

	if pricePerHour < 0 {
		invalidLine = InvalidLine(pricePerHourLine)
		return &invalidLine, ErrInvalidFormatPricePerHour
	}

	config.PricePerHour = pricePerHour

	return nil, nil
}

func parseTablesCountLine(tablesCountLine string) (int, error) {
	tablesCount, err := strconv.Atoi(tablesCountLine)
	if err != nil {
		return 0, ErrInvalidFormatTablesCount
	}

	if tablesCount < minTablesCount {
		return 0, ErrInvalidFormatTablesCount
	}

	return tablesCount, nil
}

func parseOpeningHoursLine(openingHoursLine string) (time.Time, time.Time, error) {
	splitOpeningHours := strings.Split(openingHoursLine, " ")
	if len(splitOpeningHours) != openingHoursSplitLen {
		return time.Time{}, time.Time{}, ErrInvalidFormatOpeningHours
	}

	openingTime, err := parseTime(splitOpeningHours[0])
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	closingTime, err := parseTime(splitOpeningHours[1])
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	if closingTime.Before(openingTime) {
		return time.Time{}, time.Time{}, ErrInvalidFormatOpeningHours
	}

	return openingTime, closingTime, nil
}

func parseTime(strTime string) (time.Time, error) {
	splitHours := strings.Split(strTime, ":")
	if len(splitHours) != 2 {
		return time.Time{}, ErrInvalidFormatOpeningHours
	}
	if len(splitHours[0]) != 2 {
		return time.Time{}, ErrInvalidFormatOpeningHours
	}
	if len(splitHours[1]) != 2 {
		return time.Time{}, ErrInvalidFormatOpeningHours
	}
	t, err := xtime.ParseHoursMinutesFromString(strTime)
	if err != nil {
		return time.Time{}, ErrInvalidFormatOpeningHours
	}
	return t, nil
}

func NewHandler(eventHandler eventhandler.Handler) *Handler {
	return &Handler{
		eventHandler: eventHandler,
	}
}

func (h *Handler) GetWorkingDayReport(filename string, computerClubTablesCount int) (WorkingDayReport, *InvalidLine, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", nil, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return "", nil, err
	}

	if fileInfo.Size() < minFileLinesCount {
		return "", nil, ErrInvalidFormatFile
	}

	return h.readFileByLine(file, computerClubTablesCount)
}

func (h *Handler) readFileByLine(file *os.File, computerClubTablesCount int) (WorkingDayReport, *InvalidLine, error) {
	scanner := bufio.NewScanner(file)

	h.moveScannerToFirstEventLine(scanner)

	return h.processEventLines(scanner, computerClubTablesCount)
}

func (h *Handler) moveScannerToFirstEventLine(scanner *bufio.Scanner) {
	// skip config lines
	for i := 1; i <= configLinesCount; i++ {
		scanner.Scan()
	}
}

func (h *Handler) processEventLines(scanner *bufio.Scanner, computerClubTablesCount int) (WorkingDayReport, *InvalidLine, error) {
	h.eventHandler.OpenComputerClub()

	var invalidLine InvalidLine
	var lastEventTime time.Time

	for scanner.Scan() {
		eventLine := scanner.Text()
		err := h.validateEventLine(eventLine, computerClubTablesCount, &lastEventTime)
		if err != nil {
			invalidLine = InvalidLine(eventLine)
			return "", &invalidLine, err
		}

		event, err := eventhandler.FromEventLine(eventLine)
		if err != nil {
			return "", nil, err
		}

		err = h.eventHandler.HandleEvent(event)
		if err != nil {
			return "", nil, err
		}
	}

	if err := scanner.Err(); err != nil {
		return "", nil, err
	}

	h.eventHandler.CloseComputerClub()

	workingDayReport := h.eventHandler.GetWorkingDayReport()

	return WorkingDayReport(workingDayReport), nil, nil
}

func (h *Handler) validateEventLine(eventLine string, tablesCount int, lastEventTime *time.Time) error {
	splitEventLine := strings.Split(eventLine, " ")
	if len(splitEventLine) < minSplitEventLineLen || len(splitEventLine) > maxSplitEventLineLen {
		return ErrInvalidFormatEvent
	}

	incomingEvent, err := strconv.Atoi(splitEventLine[1])
	if err != nil {
		return ErrInvalidFormatEvent
	}

	switch uint8(incomingEvent) {
	case computerclub.IncomingEventClientArrived, computerclub.IncomingEventClientWaiting, computerclub.IncomingEventClientLeft:
		return h.validateThreeArgsEvent(splitEventLine, lastEventTime)
	case computerclub.IncomingEventClientTookPlace:
		return h.validateFourArgsEvent(splitEventLine, tablesCount, lastEventTime)
	default:
		return ErrInvalidFormatEvent
	}
}

func (h *Handler) validateThreeArgsEvent(splitEventLine []string, lastEventTime *time.Time) error {
	if len(splitEventLine) != minSplitEventLineLen {
		return ErrInvalidFormatEvent
	}
	strEventTime := splitEventLine[0]

	eventTime, err := parseTime(strEventTime)
	if err != nil {
		return err
	}

	err = h.validateEventSequence(eventTime, lastEventTime)
	if err != nil {
		return err
	}

	*lastEventTime = eventTime

	clientName := splitEventLine[2]

	err = h.validateClientName(clientName)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) validateFourArgsEvent(splitEventLine []string, tablesCount int, lastEventTime *time.Time) error {
	if len(splitEventLine) != maxSplitEventLineLen {
		return ErrInvalidFormatEvent
	}

	strEventTime := splitEventLine[0]

	eventTime, err := parseTime(strEventTime)
	if err != nil {
		return err
	}

	err = h.validateEventSequence(eventTime, lastEventTime)
	if err != nil {
		return err
	}

	*lastEventTime = eventTime

	clientName := splitEventLine[2]

	err = h.validateClientName(clientName)
	if err != nil {
		return err
	}

	strTableNumber := splitEventLine[3]
	err = h.validateTableNumber(strTableNumber, tablesCount)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) validateEventSequence(currentEventTime time.Time, lastEventTime *time.Time) error {
	if !(*lastEventTime).IsZero() && currentEventTime.Before(*lastEventTime) {
		return ErrInvalidFormatEventSequence
	}

	return nil
}

func (h *Handler) validateClientName(clientName string) error {
	for i := 0; i < len(clientName); i++ {
		if !h.isLowerCaseLetter(clientName[i]) && !h.isDigit(clientName[i]) && !h.isSpecialSymbol(clientName[i]) {
			return ErrInvalidFormatClientName
		}
	}
	return nil
}

func (h *Handler) isLowerCaseLetter(char byte) bool {
	return char >= 'a' && char <= 'z'
}

func (h *Handler) isDigit(char byte) bool {
	return char >= '0' && char <= '9'
}

func (h *Handler) isSpecialSymbol(char byte) bool {
	return char == '_' || char == '-'
}

func (h *Handler) validateTableNumber(strTableNumber string, tablesCount int) error {
	tableNumber, err := strconv.Atoi(strTableNumber)
	if err != nil {
		return ErrInvalidFormatTableNumber
	}
	if tableNumber < minTablesCount || tableNumber > tablesCount {
		return ErrInvalidFormatTableNumber
	}
	return nil
}
