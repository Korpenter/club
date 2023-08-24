package service

import (
	"errors"
	"time"

	"github.com/Korpenter/club/internal/config"
	"github.com/Korpenter/club/internal/models"
	"github.com/Korpenter/club/internal/storage"
	"github.com/Korpenter/club/internal/storage/queue"
)

var (
	ErrYouShallNotPass  = errors.New("YouShallNotPass")
	ErrNotOpenYet       = errors.New("NotOpenYet")
	ErrPlaceIsBusy      = errors.New("PlaceIsBusy")
	ErrClientUnknown    = errors.New("ClientUnknown")
	ErrICanWaitNoLonger = errors.New("ICanWaitNoLonger!")
	ErrQueueFull        = errors.New("queue full")
)

type Service struct {
	cfg  *config.Config
	repo Storage
}

type Storage interface {
	AddClient(name string) error
	CheckFreeTables() bool
	EnqueueClient(name string) error
	DequeueClient() *models.Client
	RemoveClient(name string)
	FreedTableByClient(name string, timeSat time.Time) int
	ClientExists(name string) bool
	SetClientTable(name string, tableID int, timeSat time.Time) error
	KickAllClientsAndClearTables(kickTime time.Time)
	ClearAllClients() []*models.Client
	GetAllTables() map[int]*models.Table
}

func New(cfg *config.Config, repo Storage) *Service {
	return &Service{
		cfg:  cfg,
		repo: repo,
	}
}

func (s *Service) ClientArrive(timestamp time.Time, name string) error {
	if !timestamp.After(s.cfg.OpeningTime) {
		return ErrNotOpenYet
	}
	err := s.repo.AddClient(name)
	if err != nil {
		if errors.Is(err, storage.ErrClientExists) {
			return ErrYouShallNotPass
		}
	}
	return nil
}

func (s *Service) ClientSit(timestamp time.Time, name string, tableID int) error {
	exists := s.repo.ClientExists(name)
	if !exists {
		return ErrClientUnknown
	}
	_ = s.repo.FreedTableByClient(name, timestamp)
	err := s.repo.SetClientTable(name, tableID, timestamp)
	if err != nil {
		if errors.Is(err, storage.ErrTableOccupied) {
			return ErrPlaceIsBusy
		}
	}
	return nil
}

func (s *Service) ClientWait(timestamp time.Time, name string) error {
	if s.repo.CheckFreeTables() {
		return ErrICanWaitNoLonger
	}
	exists := s.repo.ClientExists(name)
	if !exists {
		return nil
	}
	err := s.repo.EnqueueClient(name)
	if err != nil {
		if errors.Is(err, queue.ErrQueueFull) {
			return ErrQueueFull
		}
	}
	return nil
}

func (s *Service) ClientLeave(timestamp time.Time, name string) (*models.Client, int, error) {
	exists := s.repo.ClientExists(name)
	if !exists {
		return nil, 0, ErrClientUnknown
	}
	freeTable := s.repo.FreedTableByClient(name, timestamp)
	s.repo.RemoveClient(name)
	dequeued := s.repo.DequeueClient()
	if dequeued == nil {
		return nil, 0, nil
	}
	err := s.repo.SetClientTable(dequeued.Name, freeTable, timestamp)
	if err != nil {
		return nil, 0, err
	}
	return dequeued, freeTable, nil
}

func (s *Service) KickClients(kickTime time.Time) []*models.Client {
	s.repo.KickAllClientsAndClearTables(kickTime)
	kicked := s.repo.ClearAllClients()
	return kicked
}

func (s *Service) CalcProfits() []*models.Profit {
	tables := s.repo.GetAllTables()

	profits := make([]*models.Profit, 0, len(tables))
	for _, v := range tables {
		total := v.TotalTime.Hour()
		if v.TotalTime.Minute() != 0 {
			total += 1
		}
		p := &models.Profit{
			Table: v,
			Sum:   int(total) * s.cfg.HourlyRate,
		}
		profits = append(profits, p)
	}

	return profits
}
