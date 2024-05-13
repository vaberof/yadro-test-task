package computerclub

import (
	"errors"
	"fmt"
	"slices"
	"testing"
	"time"
)

func TestOpen(t *testing.T) {
	config, err := getConfig(3)
	if err != nil {
		t.Fatalf("TestOpen: %s", err.Error())
	}

	computerClubService := NewComputerClub(config)

	computerClubService.Open()

	var workingDayReport WorkingDayReport

	workingDayReport.writeTime(config.OpeningTime)

	expectedWorkingDayReport := computerClubService.GetWorkingDayReport()

	if !slices.Equal(workingDayReport, expectedWorkingDayReport) {
		err = fmt.Errorf("invalid wokring day report: expected: '%v', got: '%v'", string(expectedWorkingDayReport), string(workingDayReport))
		t.Fatalf("TestOpen: %v", err)
	}
}

func TestClose(t *testing.T) {
	config, err := getConfig(1)
	if err != nil {
		t.Fatalf("TestClose: %s", err.Error())
	}

	computerClubService := NewComputerClub(config)

	computerClubService.Close()

	tableId := TableId(1)
	profit := 0
	usageTimeStr := "00:00"

	var workingDayReport WorkingDayReport

	workingDayReport.writeTime(config.ClosingTime)
	workingDayReport.writeTableReport(tableId, profit, usageTimeStr)

	expectedWorkingDayReport := computerClubService.GetWorkingDayReport()

	if !slices.Equal(workingDayReport, expectedWorkingDayReport) {
		err = fmt.Errorf("invalid wokring day report: expected: '%v', got: '%v'", string(expectedWorkingDayReport), string(workingDayReport))
		t.Fatalf("TestClose: %v", err)
	}
}

func TestGetWorkingDayReport(t *testing.T) {
	config, err := getConfig(1)
	if err != nil {
		t.Fatalf("TestGetWorkingDayReport: %s", err.Error())
	}

	var workingDayReport WorkingDayReport

	computerClubService := NewComputerClub(config)

	computerClubService.Open()

	workingDayReport.writeTime(config.OpeningTime)

	eventTime := config.OpeningTime.Add(time.Minute)
	clientName := ClientName("client1")

	tableId := TableId(1)
	profit := 0
	usageTimeStr := "00:00"

	err = computerClubService.ProcessEventClientArrived(eventTime, clientName)
	if err != nil {
		t.Fatalf("TestProcessEventClientArrived: %s", err.Error())
	}

	workingDayReport.writeEvent(eventTime, IncomingEventClientArrived, clientName)

	computerClubService.Close()

	workingDayReport.writeEvent(config.ClosingTime, OutgoingEventClientLeft, clientName)
	workingDayReport.writeTime(config.ClosingTime)
	workingDayReport.writeTableReport(tableId, profit, usageTimeStr)

	expectedWorkingDayReport := computerClubService.GetWorkingDayReport()

	if !slices.Equal(workingDayReport, expectedWorkingDayReport) {
		err = fmt.Errorf("invalid wokring day report: expected: '%v', got: '%v'", string(expectedWorkingDayReport), string(workingDayReport))
		t.Fatalf("TestGetWorkingDayReport: %v", err)
	}
}

func TestProcessEventClientArrived(t *testing.T) {
	config, err := getConfig(3)
	if err != nil {
		t.Fatalf("TestProcessEventClientArrived: %s", err.Error())
	}

	computerClubService := NewComputerClub(config)

	eventTime := config.OpeningTime.Add(time.Minute)
	clientName := ClientName("client1")

	err = computerClubService.ProcessEventClientArrived(eventTime, clientName)
	if err != nil {
		t.Fatalf("TestProcessEventClientArrived: %s", err.Error())
	}

	var workingDayReport WorkingDayReport

	workingDayReport.writeEvent(eventTime, IncomingEventClientArrived, clientName)

	expectedWorkingDayReport := computerClubService.GetWorkingDayReport()

	if !slices.Equal(workingDayReport, expectedWorkingDayReport) {
		err = fmt.Errorf("invalid wokring day report: expected: '%v', got: '%v'", string(expectedWorkingDayReport), string(workingDayReport))
		t.Fatalf("TestProcessEventClientArrived: %v", err)
	}
}

func TestProcessEventClientArrivedError(t *testing.T) {
	config, err := getConfig(3)
	if err != nil {
		t.Fatalf("TestProcessEventClientArrivedError: %s", err.Error())
	}

	computerClubService := NewComputerClub(config)

	eventTimeBeforeOpeningHours := config.OpeningTime.Add(-time.Minute)
	clientName1 := ClientName("client1")

	clientName2 := ClientName("client2")

	eventTimeInOpeningHours := config.OpeningTime.Add(time.Minute)

	// need to get ErrYouShallNotPass in second test case
	err = computerClubService.ProcessEventClientArrived(eventTimeInOpeningHours, clientName2)
	if err != nil {
		newErr := fmt.Errorf("expected error: '%v', got: '%v'", nil, err)
		t.Fatalf("TestProcessEventClientArrivedError: %s", newErr.Error())
	}

	testCases := []struct {
		name        string
		eventTime   time.Time
		clientName  ClientName
		expectedErr error
	}{
		{
			name:        "err_not_open_yet",
			eventTime:   eventTimeBeforeOpeningHours,
			clientName:  clientName1,
			expectedErr: ErrNotOpenYet,
		},
		{
			name:        "err_you_shall_not_pass",
			eventTime:   eventTimeInOpeningHours,
			clientName:  clientName2,
			expectedErr: ErrYouShallNotPass,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err = computerClubService.ProcessEventClientArrived(testCase.eventTime, testCase.clientName)
			if !errors.Is(err, testCase.expectedErr) {
				newErr := fmt.Errorf("expected error: '%v', got: '%v'", testCase.expectedErr, err)
				t.Fatalf("TestProcessEventClientArrivedError: %s", newErr.Error())
			}
		})
	}
}

func TestProcessEventClientTookPlace(t *testing.T) {
	config, err := getConfig(3)
	if err != nil {
		t.Fatalf("TestProcessEventClientTookPlace: %s", err.Error())
	}

	computerClubService := NewComputerClub(config)

	eventTime := config.OpeningTime.Add(time.Minute)
	clientName := ClientName("client1")
	tableId := TableId(1)

	var workingDayReport WorkingDayReport

	err = computerClubService.ProcessEventClientArrived(eventTime, clientName)
	if err != nil {
		t.Fatalf("TestProcessEventClientTookPlace: %s", err.Error())
	}

	workingDayReport.writeEvent(eventTime, IncomingEventClientArrived, clientName)

	err = computerClubService.ProcessEventClientTookPlace(eventTime, clientName, tableId)
	if err != nil {
		t.Fatalf("TestProcessEventClientTookPlace: %s", err.Error())
	}

	workingDayReport.writeEventWithTableId(eventTime, IncomingEventClientTookPlace, clientName, tableId)

	expectedWorkingDayReport := computerClubService.GetWorkingDayReport()

	if !slices.Equal(workingDayReport, expectedWorkingDayReport) {
		err = fmt.Errorf("invalid wokring day report: expected: '%v', got: '%v'", string(expectedWorkingDayReport), string(workingDayReport))
		t.Fatalf("TestProcessEventClientTookPlace: %v", err)
	}
}

func TestProcessEventClientTookPlaceError(t *testing.T) {
	config, err := getConfig(3)
	if err != nil {
		t.Fatalf("TestProcessEventClientTookPlaceError: %s", err.Error())
	}

	computerClubService := NewComputerClub(config)

	clientName1 := ClientName("client1")
	clientName2 := ClientName("client2")
	eventTime := config.OpeningTime.Add(time.Minute)

	tableId1 := TableId(1)
	tableId2 := TableId(2)

	// need to get ErrPlaceIsBusy in first test case
	err = computerClubService.ProcessEventClientArrived(eventTime, clientName1)
	if err != nil {
		newErr := fmt.Errorf("expected error: '%v', got: '%v'", nil, err)
		t.Fatalf("TestProcessEventClientTookPlaceError: %s", newErr.Error())
	}

	err = computerClubService.ProcessEventClientTookPlace(eventTime, clientName1, tableId1)
	if err != nil {
		newErr := fmt.Errorf("expected error: '%v', got: '%v'", nil, err)
		t.Fatalf("TestProcessEventClientTookPlaceError: %s", newErr.Error())
	}

	testCases := []struct {
		name        string
		eventTime   time.Time
		clientName  ClientName
		tableId     TableId
		expectedErr error
	}{
		{
			name:        "err_place_is_busy",
			eventTime:   eventTime,
			clientName:  clientName1,
			expectedErr: ErrPlaceIsBusy,
			tableId:     tableId1,
		},
		{
			name:        "err_client_unknown",
			eventTime:   eventTime,
			clientName:  clientName2,
			expectedErr: ErrClientUnknown,
			tableId:     tableId2,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err = computerClubService.ProcessEventClientTookPlace(testCase.eventTime, testCase.clientName, testCase.tableId)
			if !errors.Is(err, testCase.expectedErr) {
				newErr := fmt.Errorf("expected error: '%v', got: '%v'", testCase.expectedErr, err)
				t.Fatalf("TestProcessEventClientTookPlaceError: %s", newErr.Error())
			}
		})
	}
}

func TestProcessEventClientWaiting(t *testing.T) {
	config, err := getConfig(1)
	if err != nil {
		t.Fatalf("TestProcessEventClientWaiting: %s", err.Error())
	}

	computerClubService := NewComputerClub(config)

	eventTime := config.OpeningTime.Add(time.Minute)
	clientName := ClientName("client1")
	tableId := TableId(1)

	var workingDayReport WorkingDayReport

	// need to set busy table to get no error

	err = computerClubService.ProcessEventClientArrived(eventTime, clientName)
	if err != nil {
		t.Fatalf("TestProcessEventClientWaiting: %s", err.Error())
	}

	workingDayReport.writeEvent(eventTime, IncomingEventClientArrived, clientName)

	err = computerClubService.ProcessEventClientTookPlace(eventTime, clientName, tableId)
	if err != nil {
		newErr := fmt.Errorf("expected error: '%v', got: '%v'", nil, err)
		t.Fatalf("TestProcessEventClientWaiting: %s", newErr.Error())
	}

	workingDayReport.writeEventWithTableId(eventTime, IncomingEventClientTookPlace, clientName, tableId)

	err = computerClubService.ProcessEventClientWaiting(eventTime, clientName)
	if err != nil {
		t.Fatalf("TestProcessEventClientWaiting: %s", err.Error())
	}

	workingDayReport.writeEvent(eventTime, IncomingEventClientWaiting, clientName)

	expectedWorkingDayReport := computerClubService.GetWorkingDayReport()

	if !slices.Equal(workingDayReport, expectedWorkingDayReport) {
		err = fmt.Errorf("invalid wokring day report: expected: '%v', got: '%v'", string(expectedWorkingDayReport), string(workingDayReport))
		t.Fatalf("TestProcessEventClientWaiting: %v", err)
	}
}

func TestProcessEventClientWaitingErrorICanWaitNoLonger(t *testing.T) {
	config, err := getConfig(1)
	if err != nil {
		t.Fatalf("TestProcessEventClientWaitingErrorICanWaitNoLonger: %s", err.Error())
	}

	computerClubService := NewComputerClub(config)

	eventTime := config.OpeningTime.Add(time.Minute)
	clientName := ClientName("client1")

	err = computerClubService.ProcessEventClientWaiting(eventTime, clientName)
	if !errors.Is(err, ErrICanWaitNoLonger) {
		err = fmt.Errorf("expected error: '%v', got: '%v'", ErrICanWaitNoLonger, err)
		t.Fatalf("TestProcessEventClientWaitingErrorICanWaitNoLonger: %s", err.Error())
	}
}

func TestProcessEventClientWaitingErrorQueueIsFull(t *testing.T) {
	config, err := getConfig(1)
	if err != nil {
		t.Fatalf("TestProcessEventClientWaitingErrorQueueIsFull: %s", err.Error())
	}

	computerClubService := NewComputerClub(config)

	eventTime := config.OpeningTime.Add(time.Minute)
	clientName1 := ClientName("client1")
	clientName2 := ClientName("client2")
	clientName3 := ClientName("client3")
	clientName4 := ClientName("client4")
	tableId := TableId(1)

	err = computerClubService.ProcessEventClientArrived(eventTime, clientName1)
	if err != nil {
		t.Fatalf("TestProcessEventClientWaitingErrorQueueIsFull: %s", err.Error())
	}

	err = computerClubService.ProcessEventClientTookPlace(eventTime, clientName1, tableId)
	if err != nil {
		newErr := fmt.Errorf("expected error: '%v', got: '%v'", nil, err)
		t.Fatalf("TestProcessEventClientWaitingErrorQueueIsFull: %s", newErr.Error())
	}

	err = computerClubService.ProcessEventClientWaiting(eventTime, clientName2)
	if err != nil {
		newErr := fmt.Errorf("expected error: '%v', got: '%v'", nil, err)
		t.Fatalf("TestProcessEventClientWaitingErrorQueueIsFull: %s", newErr.Error())
	}

	err = computerClubService.ProcessEventClientWaiting(eventTime, clientName3)
	if err != nil {
		newErr := fmt.Errorf("expected error: '%v', got: '%v'", nil, err)
		t.Fatalf("TestProcessEventClientWaitingErrorQueueIsFull: %s", newErr.Error())
	}

	err = computerClubService.ProcessEventClientWaiting(eventTime, clientName4)
	if !errors.Is(err, ErrQueueIsFull) {
		err = fmt.Errorf("expected error: '%v', got: '%v'", ErrQueueIsFull, err)
		t.Fatalf("TestProcessEventClientWaitingErrorQueueIsFull: %s", err.Error())
	}
}

func TestProcessEventClientLeft(t *testing.T) {
	config, err := getConfig(3)
	if err != nil {
		t.Fatalf("TestProcessEventClientLeft: %s", err.Error())
	}

	computerClubService := NewComputerClub(config)

	eventTime := config.OpeningTime.Add(time.Minute)
	clientName := ClientName("client1")

	var workingDayReport WorkingDayReport

	err = computerClubService.ProcessEventClientArrived(eventTime, clientName)
	if err != nil {
		t.Fatalf("TestProcessEventClientLeft: %s", err.Error())
	}

	workingDayReport.writeEvent(eventTime, IncomingEventClientArrived, clientName)

	err = computerClubService.ProcessEventClientLeft(eventTime, clientName)
	if err != nil {
		t.Fatalf("TestProcessEventClientLeft: %s", err.Error())
	}

	workingDayReport.writeEvent(eventTime, IncomingEventClientLeft, clientName)

	expectedWorkingDayReport := computerClubService.GetWorkingDayReport()

	if !slices.Equal(workingDayReport, expectedWorkingDayReport) {
		err = fmt.Errorf("invalid wokring day report: expected: '%v', got: '%v'", string(expectedWorkingDayReport), string(workingDayReport))
		t.Fatalf("TestProcessEventClientLeft: %v", err)
	}
}

func TestProcessEventClientLeftError(t *testing.T) {
	config, err := getConfig(1)
	if err != nil {
		t.Fatalf("TestProcessEventClientLeftError: %s", err.Error())
	}

	computerClubService := NewComputerClub(config)

	eventTime := config.OpeningTime.Add(time.Minute)
	clientName := ClientName("client1")

	err = computerClubService.ProcessEventClientLeft(eventTime, clientName)
	if !errors.Is(err, ErrClientUnknown) {
		err = fmt.Errorf("expected error: '%v', got: '%v'", ErrClientUnknown, err)
		t.Fatalf("TestProcessEventClientLeftError: %s", err.Error())
	}
}

func getConfig(tablesCount int) (*Config, error) {
	const layout = "15:04"

	openingTime, err := time.Parse(layout, "09:00")
	if err != nil {
		return nil, err
	}

	closingTime, err := time.Parse(layout, "19:00")
	if err != nil {
		return nil, err
	}

	config := Config{
		TablesCount:  tablesCount,
		OpeningTime:  openingTime,
		ClosingTime:  closingTime,
		PricePerHour: 10,
	}

	return &config, nil
}
