package storage

import (
	"errors"
	"time"

	"github.com/Korpenter/club/internal/config"
	"github.com/Korpenter/club/internal/models"
	"github.com/Korpenter/club/internal/storage/queue"
)

var (
	ErrTableOccupied    = errors.New("table already occupied")
	ErrTableEmpty       = errors.New("table is empty")
	ErrClientNotInQueue = errors.New("client not in the queue")
	ErrClientExists     = errors.New("client already exists")
	ErrClientUnknown    = errors.New("unknown client")
)

type InMemRepo struct {
	tables  map[int]*models.Table
	queue   Queue
	clients map[string]*models.Client
	events  []*models.Event
}

type Queue interface {
	Enqueue(*models.Client) error
	Dequeue() *models.Client
	Remove(client *models.Client)
	Clear()
}

func NewInMemRepo(cfg *config.Config) *InMemRepo {
	queue := queue.NewQueue(cfg.NumberOfTables)
	tables := make(map[int]*models.Table, cfg.NumberOfTables)
	for i := 1; i <= cfg.NumberOfTables; i++ {
		tables[i] = &models.Table{Id: i}
	}
	return &InMemRepo{
		tables:  tables,
		queue:   queue,
		clients: make(map[string]*models.Client),
		events:  make([]*models.Event, 0),
	}
}

func (r *InMemRepo) AddClient(name string) error {
	if _, exists := r.clients[name]; exists {
		return ErrClientExists
	}
	r.clients[name] = &models.Client{Name: name}
	return nil
}

func (r *InMemRepo) ClientExists(name string) bool {
	_, ok := r.clients[name]
	return ok
}

func (r *InMemRepo) CheckFreeTables() bool {
	for _, v := range r.tables {
		if v.Client == nil {
			return true
		}
	}
	return false
}

func (r *InMemRepo) EnqueueClient(name string) error {
	return r.queue.Enqueue(r.clients[name])
}

func (r *InMemRepo) DequeueClient() *models.Client {
	return r.queue.Dequeue()
}

func (r *InMemRepo) SetClientTable(name string, tableID int, timeSat time.Time) error {
	if r.tables[tableID].Client != nil {
		return ErrTableOccupied
	}
	r.tables[tableID].Client = r.clients[name]
	r.tables[tableID].ClientSat = timeSat
	return nil
}

func (r *InMemRepo) FreedTableByClient(name string, timeSat time.Time) int {
	for i, v := range r.tables {
		if v.Client != nil && v.Client.Name == name {
			timeSpent := timeSat.Sub(v.ClientSat)
			r.tables[i].TotalTime = r.tables[i].TotalTime.Add(timeSpent)
			r.tables[i].Client = nil
			r.tables[i].ClientSat = time.Time{}
			return i
		}
	}
	return 0
}

func (r *InMemRepo) RemoveClient(name string) {
	r.queue.Remove(r.clients[name])
	delete(r.clients, name)
}

func (r *InMemRepo) KickAllClientsAndClearTables(kickTime time.Time) {
	for i, v := range r.tables {
		if v.Client != nil {
			timeSpent := kickTime.Sub(v.ClientSat)
			r.tables[i].TotalTime = r.tables[i].TotalTime.Add(timeSpent)
			r.tables[i].Client = nil
			r.tables[i].ClientSat = time.Time{}
		}
	}
}

func (r *InMemRepo) ClearAllClients() []*models.Client {
	var clients []*models.Client
	for _, v := range r.clients {
		clients = append(clients, v)
		delete(r.clients, v.Name)
	}
	r.queue.Clear()
	return clients
}

func (r *InMemRepo) GetAllTables() map[int]*models.Table {
	return r.tables
}
