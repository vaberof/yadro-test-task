package eventhandler

import (
	"errors"
	"github.com/vaberof/yadro-test-task/internal/domain/computerclub"
)

type Handler interface {
	HandleEvent(event *Event) error
	OpenComputerClub()
	CloseComputerClub()
	GetWorkingDayReport() computerclub.WorkingDayReport
}

type handlerImpl struct {
	computerClubService computerclub.ComputerClubService
}

func NewHandler(computerClubService computerclub.ComputerClubService) Handler {
	return &handlerImpl{computerClubService: computerClubService}
}

func (h *handlerImpl) HandleEvent(event *Event) error {
	switch event.Type {
	case computerclub.IncomingEventClientArrived:
		return h.handleEventClientArrived(event)
	case computerclub.IncomingEventClientTookPlace:
		return h.handleEventClientTookPlace(event)
	case computerclub.IncomingEventClientWaiting:
		return h.handleEventClientWaiting(event)
	case computerclub.IncomingEventClientLeft:
		return h.handleEventClientLeft(event)
	default:
		return errors.New("invalid event type")
	}
}

func (h *handlerImpl) OpenComputerClub() {
	h.computerClubService.Open()
}

func (h *handlerImpl) CloseComputerClub() {
	h.computerClubService.Close()
}

func (h *handlerImpl) GetWorkingDayReport() computerclub.WorkingDayReport {
	return h.computerClubService.GetWorkingDayReport()
}

func (h *handlerImpl) handleEventClientArrived(event *Event) error {
	err := h.computerClubService.ProcessEventClientArrived(event.Time, computerclub.ClientName(event.ClientName))
	if err != nil {
		if !errors.Is(err, computerclub.ErrYouShallNotPass) && !errors.Is(err, computerclub.ErrNotOpenYet) {
			return err
		}
	}
	return nil
}

func (h *handlerImpl) handleEventClientTookPlace(event *Event) error {
	err := h.computerClubService.ProcessEventClientTookPlace(event.Time, computerclub.ClientName(event.ClientName), computerclub.TableId(event.TableId))
	if err != nil {
		if !errors.Is(err, computerclub.ErrPlaceIsBusy) && !errors.Is(err, computerclub.ErrClientUnknown) {
			return err
		}
	}
	return nil
}

func (h *handlerImpl) handleEventClientWaiting(event *Event) error {
	err := h.computerClubService.ProcessEventClientWaiting(event.Time, computerclub.ClientName(event.ClientName))
	if err != nil {
		if !errors.Is(err, computerclub.ErrICanWaitNoLonger) && !errors.Is(err, computerclub.ErrQueueIsFull) {
			return err
		}
	}
	return nil
}

func (h *handlerImpl) handleEventClientLeft(event *Event) error {
	err := h.computerClubService.ProcessEventClientLeft(event.Time, computerclub.ClientName(event.ClientName))
	if err != nil {
		if !errors.Is(err, computerclub.ErrClientUnknown) {
			return err
		}
	}
	return nil
}
