package queue

import (
	"errors"

	"github.com/Korpenter/club/internal/models"
)

var (
	ErrQueueFull     = errors.New("queue full")
	ErrAlreadyExists = errors.New("alreay in queue")
)

type Node struct {
	value *models.Client
	prev  *Node
	next  *Node
}

type Queue struct {
	head      *Node
	tail      *Node
	maxLength int
	set       map[string]*Node
}

func NewQueue(maxLength int) *Queue {
	return &Queue{
		maxLength: maxLength,
		set:       make(map[string]*Node, maxLength),
	}
}

func (q *Queue) Enqueue(client *models.Client) error {
	if _, exists := q.set[client.Name]; exists {
		return ErrAlreadyExists
	}
	if len(q.set) == q.maxLength {
		return ErrQueueFull
	}
	newNode := &Node{value: client}
	if q.tail == nil {
		q.head = newNode
		q.tail = newNode
	} else {
		q.tail.next = newNode
		newNode.prev = q.tail
		q.tail = newNode
	}
	q.set[client.Name] = newNode
	return nil
}

func (q *Queue) Dequeue() *models.Client {
	if len(q.set) == 0 {
		return nil
	}

	oldHead := q.head
	q.head = q.head.next
	if q.head != nil {
		q.head.prev = nil
	} else {
		q.tail = nil
	}

	delete(q.set, oldHead.value.Name)
	return oldHead.value
}

func (q *Queue) Remove(client *models.Client) {
	if node, exists := q.set[client.Name]; exists {
		if node.prev != nil {
			node.prev.next = node.next
		} else {
			q.head = node.next
		}
		if node.next != nil {
			node.next.prev = node.prev
		} else {
			q.tail = node.prev
		}
		delete(q.set, client.Name)
	}
}

func (q *Queue) Clear() {
	q.head = nil
	q.tail = nil
	q.set = make(map[string]*Node, q.maxLength)
}
