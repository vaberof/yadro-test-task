package computerclub

import (
	"errors"
	"slices"
	"time"
)

const (
	IncomingEventClientArrived uint8 = iota + 1
	IncomingEventClientTookPlace
	IncomingEventClientWaiting
	IncomingEventClientLeft
)

const (
	OutgoingEventClientLeft uint8 = iota + 11
	OutgoingEventClientTookPlace
	OutgoingEventError
)

var (
	ErrNotOpenYet       = errors.New("NotOpenYet")
	ErrYouShallNotPass  = errors.New("YouShallNotPass")
	ErrClientUnknown    = errors.New("ClientUnknown")
	ErrPlaceIsBusy      = errors.New("PlaceIsBusy")
	ErrICanWaitNoLonger = errors.New("ICanWaitNoLonger!")

	ErrQueueIsFull = errors.New("Queue is full")
)

type ComputerClubService interface {
	Open()
	ProcessEventClientArrived(eventTime time.Time, clientName ClientName) error
	ProcessEventClientTookPlace(eventTime time.Time, clientName ClientName, tableId TableId) error
	ProcessEventClientWaiting(eventTime time.Time, clientName ClientName) error
	ProcessEventClientLeft(eventTime time.Time, clientName ClientName) error
	Close()
	GetWorkingDayReport() WorkingDayReport
}

type Config struct {
	TablesCount  int
	OpeningTime  time.Time
	ClosingTime  time.Time
	PricePerHour int
}

type computerClubServiceImpl struct {
	tablesCount  int
	openingTime  time.Time
	closingTime  time.Time
	pricePerHour int

	clients map[ClientName]Client
	tables  map[TableId]Table

	clientQueue *ClientQueue

	// buf contains all output for the day
	buf WorkingDayReport
}

const startBufSize = 512

const minTablesCount = 1

func NewComputerClub(config *Config) ComputerClubService {
	tables := make(map[TableId]Table)

	for tableId := TableId(minTablesCount); tableId <= TableId(config.TablesCount); tableId++ {
		tables[tableId] = Table{
			Id:    tableId,
			State: StateTableIsFree,
		}
	}

	computerClub := &computerClubServiceImpl{
		tablesCount:  config.TablesCount,
		openingTime:  config.OpeningTime,
		closingTime:  config.ClosingTime,
		pricePerHour: config.PricePerHour,
		clients:      make(map[ClientName]Client),
		tables:       tables,
		clientQueue:  NewClientQueue(config.TablesCount + 1),
		buf:          make([]byte, 0, startBufSize),
	}

	return computerClub
}

func (c *computerClubServiceImpl) ProcessEventClientArrived(eventTime time.Time, clientName ClientName) error {
	c.buf.writeEvent(eventTime, IncomingEventClientArrived, clientName)

	if c.isClientInComputerClub(clientName) {
		c.buf.writeEventError(eventTime, ErrYouShallNotPass)

		return ErrYouShallNotPass
	}

	if c.isNonWorkingHours(eventTime) {
		c.buf.writeEventError(eventTime, ErrNotOpenYet)

		return ErrNotOpenYet
	}

	c.addClient(clientName)

	return nil
}

func (c *computerClubServiceImpl) ProcessEventClientTookPlace(eventTime time.Time, clientName ClientName, tableId TableId) error {
	c.buf.writeEventWithTableId(eventTime, IncomingEventClientTookPlace, clientName, tableId)

	if c.isBusyTable(tableId) {
		c.buf.writeEventError(eventTime, ErrPlaceIsBusy)

		return ErrPlaceIsBusy
	}

	if !c.isClientInComputerClub(clientName) {
		c.buf.writeEventError(eventTime, ErrClientUnknown)

		return ErrClientUnknown
	}

	client := c.clients[clientName]

	if client.State == StateClientTookPlace {
		busyTableId := client.BusyTableId
		c.freeTable(busyTableId, eventTime)

		if !c.clientQueue.IsEmpty() {
			clientFromQueue := c.clientQueue.Pop()
			c.takeTable(busyTableId, eventTime, clientFromQueue)
			c.buf.writeEventWithTableId(eventTime, OutgoingEventClientTookPlace, clientFromQueue.Name, busyTableId)
		}
	}

	c.takeTable(tableId, eventTime, &client)

	return nil
}

func (c *computerClubServiceImpl) ProcessEventClientWaiting(eventTime time.Time, clientName ClientName) error {
	c.buf.writeEvent(eventTime, IncomingEventClientWaiting, clientName)

	if c.isThereFreeTable() {
		c.buf.writeEventError(eventTime, ErrICanWaitNoLonger)

		return ErrICanWaitNoLonger
	}

	if c.clientQueue.IsFull() {
		c.deleteClient(clientName)

		c.buf.writeEvent(eventTime, OutgoingEventClientLeft, clientName)

		return ErrQueueIsFull
	}

	c.addClientToQueue(clientName)

	return nil
}

func (c *computerClubServiceImpl) ProcessEventClientLeft(eventTime time.Time, clientName ClientName) error {
	c.buf.writeEvent(eventTime, IncomingEventClientLeft, clientName)

	if !c.isClientInComputerClub(clientName) {
		c.buf.writeEventError(eventTime, ErrClientUnknown)

		return ErrClientUnknown
	}

	client := c.clients[clientName]
	busyTableId := client.BusyTableId

	c.freeTable(busyTableId, eventTime)

	if !c.clientQueue.IsEmpty() {
		clientFromQueue := c.clientQueue.Pop()
		c.takeTable(busyTableId, eventTime, clientFromQueue)

		c.buf.writeEventWithTableId(eventTime, OutgoingEventClientTookPlace, clientFromQueue.Name, busyTableId)
	}

	c.deleteClient(client.Name)

	return nil
}

func (c *computerClubServiceImpl) Open() {
	c.buf.writeTime(c.openingTime)
}

func (c *computerClubServiceImpl) Close() {
	clientNames := c.getRemainingClientNames()
	slices.Sort(clientNames)

	for _, clientName := range clientNames {
		c.buf.writeEvent(c.closingTime, OutgoingEventClientLeft, clientName)
	}

	c.buf.writeTime(c.closingTime)

	for tableId := TableId(minTablesCount); tableId <= TableId(c.tablesCount); tableId++ {
		table := c.tables[tableId]
		c.buf.writeTableReport(tableId, table.Profit, table.usageTimePerDayString())
	}
}

func (c *computerClubServiceImpl) GetWorkingDayReport() WorkingDayReport {
	return c.buf
}

func (c *computerClubServiceImpl) getRemainingClientNames() []ClientName {
	var clientNames []ClientName

	for _, client := range c.clients {
		clientNames = append(clientNames, client.Name)

		if client.State == StateClientTookPlace {
			busyTableId := client.BusyTableId
			c.freeTable(busyTableId, c.closingTime)
		}

		c.deleteClient(client.Name)
	}

	return clientNames
}

func (c *computerClubServiceImpl) takeTable(tableId TableId, startTime time.Time, client *Client) {
	table := c.tables[tableId]
	table.State = StateTableIsBusy
	table.StartTime = startTime
	c.tables[tableId] = table

	client.State = StateClientTookPlace
	client.BusyTableId = tableId
	c.clients[client.Name] = *client
}

func (c *computerClubServiceImpl) freeTable(tableId TableId, endTime time.Time) {
	table := c.tables[tableId]

	table.EndTime = endTime
	table.State = StateTableIsFree
	table.calculateProfit(c.pricePerHour)
	table.calculateUsageTime()

	c.tables[tableId] = table
}

func (c *computerClubServiceImpl) isBusyTable(tableId TableId) bool {
	table := c.tables[tableId]
	return table.State == StateTableIsBusy
}

func (c *computerClubServiceImpl) addClient(clientName ClientName) {
	client := Client{
		Name:  clientName,
		State: StateClientArrived,
	}
	c.clients[clientName] = client
}

func (c *computerClubServiceImpl) addClientToQueue(clientName ClientName) {
	client := c.clients[clientName]
	client.State = StateClientIsWaiting
	c.clients[clientName] = client
	c.clientQueue.Push(&client)
}

func (c *computerClubServiceImpl) deleteClient(clientName ClientName) {
	delete(c.clients, clientName)
}

func (c *computerClubServiceImpl) isClientInComputerClub(clientName ClientName) bool {
	_, ok := c.clients[clientName]
	return ok
}

func (c *computerClubServiceImpl) isNonWorkingHours(time time.Time) bool {
	return time.Before(c.openingTime) || time.After(c.closingTime)
}

func (c *computerClubServiceImpl) isThereFreeTable() bool {
	for _, table := range c.tables {
		if table.State == StateTableIsFree {
			return true
		}
	}
	return false
}
