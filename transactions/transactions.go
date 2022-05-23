package transactions

import (
	"fmt"
	"strconv"
	"sync"
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
	// Allows multiple reads to occur at the same time (blocked by a single update)
	lock         sync.RWMutex
	cardRegistry map[int]*Card
}

func NewClient() *Client {
	return &Client{
		cardRegistry: make(map[int]*Card),
	}
}

func (c *Client) GetActiveCards() []*Card {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.getActiveCards()
}

// Should be called in a synced method
func (c *Client) getActiveCards() []*Card {
	cards := make([]*Card, 0, ACTIVE_CARDS_LIMIT)
	for _, card := range c.cardRegistry {
		if card.isActive {
			cards = append(cards, card)
		}
	}
	return cards
}

func (c *Client) AddCard(card *Card) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	if _, ok := c.cardRegistry[card.id]; ok {
		return fmt.Errorf("card with id - %v, already exists", card.id)
	}

	c.cardRegistry[card.id] = card
	return nil
}

// Immediately returns and delegates processing to another goroutine
func (c *Client) ReceiveTransaction(id int) error {
	c.lock.Lock()
	if _, ok := c.cardRegistry[id]; !ok {
		return fmt.Errorf("Unkown card id - %v", id)
	}
	card := c.cardRegistry[id]
	card.lastUsed = time.Now()

	go func() {
		defer c.lock.Unlock()
		c.deactivateCard()
		card.isActive = true
	}()
	return nil
}

func (c *Client) deactivateCard() {
	cards := c.getActiveCards()
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
