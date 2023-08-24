package handler

import (
	"bufio"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"slices"

	"github.com/Korpenter/club/internal/config"
	"github.com/Korpenter/club/internal/models"
	"github.com/Korpenter/club/internal/service"
	"github.com/Korpenter/club/internal/utils"
)

type FileHandler struct {
	Scanner *bufio.Scanner
	Service Service
	cfg     *config.Config
	ee      []*models.Event
}

type Service interface {
	ClientArrive(timestamp time.Time, name string) error
	ClientWait(timestamp time.Time, name string) error
	ClientLeave(timestamp time.Time, name string) (*models.Client, int, error)
	ClientSit(timestamp time.Time, name string, tableID int) error
	KickClients(kickTime time.Time) []*models.Client
	CalcProfits() []*models.Profit
}

func NewFileHandler(scanner *bufio.Scanner, service Service, cfg *config.Config) *FileHandler {
	return &FileHandler{
		Scanner: scanner,
		Service: service,
		cfg:     cfg,
	}
}

func (h *FileHandler) ProcessEvents() error {
	for h.Scanner.Scan() {
		eventString := h.Scanner.Text()
		eventSplit := strings.Split(h.Scanner.Text(), " ")
		if len(eventSplit) < 3 || len(eventSplit) > 5 {
			return errors.New(eventString)
		}

		eventTime, err := utils.Parse(eventSplit[0])
		if err != nil {
			return errors.New(eventString)
		}

		eventCode, err := strconv.Atoi(eventSplit[1])
		if err != nil {
			return errors.New(eventString)
		}

		if !models.ValidClientName.MatchString(eventSplit[2]) {
			return errors.New(eventString)
		}
		clientName := eventSplit[2]
		event := &models.Event{
			Code:       eventCode,
			Timestamp:  eventTime,
			ClientName: clientName,
		}
		switch eventCode {
		case models.ClientArrived:
			h.logEvent(event)
			err := h.Service.ClientArrive(eventTime, clientName)
			if err != nil {
				errEvent := &models.Event{
					Code:      models.EventError,
					Timestamp: eventTime,
					ErrorMsg:  err,
				}
				h.logEvent(errEvent)
			}
		case models.ClientSat:
			if len(eventSplit) < 4 {
				return errors.New(eventString)
			}
			tableID, err := strconv.Atoi(eventSplit[3])
			if err != nil {
				return errors.New(eventString)
			}
			event.TableID = tableID
			h.logEvent(event)
			err = h.Service.ClientSit(eventTime, clientName, tableID)
			if err != nil {
				errEvent := &models.Event{
					Code:      models.EventError,
					Timestamp: eventTime,
					ErrorMsg:  err,
				}
				h.logEvent(errEvent)
			}
		case models.ClientWaiting:
			h.logEvent(event)
			err := h.Service.ClientWait(eventTime, clientName)
			if err != nil {
				if errors.Is(err, service.ErrICanWaitNoLonger) {
					errEvent := &models.Event{
						Code:      models.EventError,
						Timestamp: eventTime,
						ErrorMsg:  err,
					}
					h.logEvent(errEvent)
				}
				if errors.Is(err, service.ErrQueueFull) {
					leftEvent := &models.Event{
						Code:       models.ClientForceLeft,
						Timestamp:  eventTime,
						ClientName: clientName,
					}
					h.logEvent(leftEvent)
				}
			}
		case models.ClientLeft:
			h.logEvent(event)
			dequeued, tableID, err := h.Service.ClientLeave(eventTime, clientName)
			if err != nil && errors.Is(err, service.ErrClientUnknown) {
				errEvent := &models.Event{
					Code:      models.EventError,
					Timestamp: eventTime,
					ErrorMsg:  err,
				}
				h.logEvent(errEvent)
			}
			if dequeued != nil {
				sitEvent := &models.Event{
					Code:       models.ClientSat,
					Timestamp:  eventTime,
					ClientName: dequeued.Name,
					TableID:    tableID,
				}
				h.logEvent(sitEvent)
			}
		default:
			return errors.New(eventString)
		}
	}
	return nil
}

func (h *FileHandler) EndDay() error {
	fmt.Println(utils.Format(h.cfg.OpeningTime))
	kicked := h.Service.KickClients(h.cfg.ClosingTime)
	cmp := func(a, b *models.Client) int {
		return strings.Compare(a.Name, b.Name)
	}
	slices.SortFunc(kicked, cmp)
	for _, v := range h.ee {
		fmt.Println(v)
	}
	for _, v := range kicked {
		kickEvent := &models.Event{
			Code:       models.ClientForceLeft,
			Timestamp:  h.cfg.ClosingTime,
			ClientName: v.Name,
		}
		fmt.Println(kickEvent)
	}
	profits := h.Service.CalcProfits()
	cmpInt := func(a, b *models.Profit) int {
		if a.Table.Id < b.Table.Id {
			return -1
		}
		if a.Table.Id > b.Table.Id {
			return 1
		}
		return 0
	}
	fmt.Println(utils.Format(h.cfg.ClosingTime))
	slices.SortFunc(profits, cmpInt)
	for _, v := range profits {
		fmt.Println(v)
	}
	return nil
}

func (h *FileHandler) logEvent(event *models.Event) {
	h.ee = append(h.ee, event)
}
