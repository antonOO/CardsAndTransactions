package transactions

import (
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type Card struct {
	id       int
	lastUsed time.Time
	isActive bool
}

func NewCard(id int) (*Card, error) {
	if isIdValid(id) {
		return &Card{
			id:       id,
			isActive: false,
		}, nil
	}
	return nil, fmt.Errorf("invalid card ID")
}

func isIdValid(id int) bool {
	return len(strconv.Itoa(id)) == 16
}

var ACTIVE_CARDS_LIMIT = 10

type Client struct {
	lock                       sync.Mutex
	atomicCardRegistryAccessor atomic.Value
}

func NewClient() *Client {
	var atomicCardRegistryAccessor atomic.Value
	cardRegistry := make(map[int]*Card)
	atomicCardRegistryAccessor.Store(cardRegistry)

	return &Client{
		atomicCardRegistryAccessor: atomicCardRegistryAccessor,
	}
}

func (c *Client) getCardRegistry() map[int]*Card {
	return c.atomicCardRegistryAccessor.Load().(map[int]*Card)
}

// Guaranteed eventual consistency. If concurrent activation requests occur, immediate gets might contain old data
func (c *Client) GetActiveCards() []*Card {
	cards := make([]*Card, 0, ACTIVE_CARDS_LIMIT)
	for _, card := range c.getCardRegistry() {
		if card.isActive {
			cards = append(cards, card)
		}
	}
	return cards
}

func (c *Client) AddCard(card *Card) error {
	c.lock.Lock()
	c.lock.Unlock()
	cardRegistry := c.getCardRegistry()
	if _, ok := cardRegistry[card.id]; ok {
		return fmt.Errorf("card with id - %v, already exists", card.id)
	}

	cardRegistry[card.id] = card
	c.atomicCardRegistryAccessor.Store(cardRegistry)
	return nil
}

// Immediately returns and delegates processing to another goroutine
func (c *Client) ReceiveTransaction(id int) error {
	c.lock.Lock()
	cardRegistry := c.getCardRegistry()
	if _, ok := cardRegistry[id]; !ok {
		return fmt.Errorf("Unkown card id - %v", id)
	}
	card := cardRegistry[id]
	card.lastUsed = time.Now()

	go func() {
		defer c.lock.Unlock()
		c.deactivateCard()
		card.isActive = true
	}()
	return nil
}

func (c *Client) deactivateCard() {
	cards := c.GetActiveCards()
	if len(cards) < ACTIVE_CARDS_LIMIT {
		return
	}

	lastUsedCard := cards[0]
	for _, card := range cards {
		if lastUsedCard.lastUsed.After(card.lastUsed) {
			lastUsedCard = card
		}
	}
	lastUsedCard.isActive = false
}
