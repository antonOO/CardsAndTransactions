package transactions

import (
	"runtime"
	"testing"
	"time"
)
import "github.com/stretchr/testify/assert"

func TestCreateCard(t *testing.T) {
	card, err := NewCard(1234560123456789)
	assert.Nil(t, err)
	assert.NotNil(t, card)
}

func TestCreateCard_Invalid(t *testing.T) {
	card, err := NewCard(1)
	assert.Nil(t, card)
	assert.NotNil(t, err)
}

func TestCreateClient(t *testing.T) {
	client := NewClient()
	assert.NotNil(t, client)
}

func TestClient_AddCard(t *testing.T) {
	client := NewClient()
	card, _ := NewCard(1234560123456789)
	err := client.AddCard(card)
	assert.Nil(t, err)
}

func TestClient_AddCard_Duplicate(t *testing.T) {
	client := NewClient()
	card, _ := NewCard(1234560123456789)
	dupl, _ := NewCard(1234560123456789)
	client.AddCard(card)
	err := client.AddCard(dupl)
	assert.NotNil(t, err)
}

func TestClient_ReceiveTransaction(t *testing.T) {
	client := NewClient()
	card, _ := NewCard(1234560123456789)
	client.AddCard(card)

	client.ReceiveTransaction(1234560123456789)

	// Delay check and force reschedule, so update is made
	runtime.Gosched()
	time.Sleep(100)

	assert.NotNil(t, card.lastUsed)
	assert.True(t, card.isActive)
}

func TestClient_GetActiveCards(t *testing.T) {
	client := NewClient()
	cards := client.GetActiveCards()
	assert.Equal(t, 0, len(cards))

	card, _ := NewCard(1234560123456789)
	client.AddCard(card)

	cards = client.GetActiveCards()
	assert.Equal(t, 0, len(cards))

	client.ReceiveTransaction(1234560123456789)

	// Delay check and force reschedule, so update is made
	runtime.Gosched()
	time.Sleep(100)

	cards = client.GetActiveCards()
	assert.Equal(t, 1, len(cards))
}

func TestClient_GetActiveCards_AboveLimit(t *testing.T) {
	ACTIVE_CARDS_LIMIT = 2
	client := NewClient()
	card1, _ := NewCard(1111111111111111)
	card2, _ := NewCard(2222222222222222)
	card3, _ := NewCard(3333333333333333)
	client.AddCard(card1)
	client.AddCard(card2)
	client.AddCard(card3)

	client.ReceiveTransaction(1111111111111111)
	time.Sleep(100)
	client.ReceiveTransaction(2222222222222222)
	time.Sleep(100)
	client.ReceiveTransaction(3333333333333333)

	// Delay check and force reschedule, so update is made
	runtime.Gosched()
	time.Sleep(100)

	cards := client.GetActiveCards()
	assert.Equal(t, ACTIVE_CARDS_LIMIT, len(cards))
	// card1 is last used
	assert.False(t, card1.isActive)
}

func TestClient_GetActiveCards_Concurrent(t *testing.T) {
	ACTIVE_CARDS_LIMIT = 5
	client := NewClient()

	for i := 0; i < 10; i++ {
		cardId := 1000000000000000 + i
		var card, _ = NewCard(int(cardId))
		go func() {
			client.AddCard(card)
			client.ReceiveTransaction(int(cardId))
		}()
	}

	time.Sleep(2 * time.Second)

	cards := client.GetActiveCards()
	assert.Equal(t, ACTIVE_CARDS_LIMIT, len(cards))
}
