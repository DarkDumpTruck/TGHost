package tghost

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"tghost/pkg/logger"
	"time"

	"github.com/google/uuid"
)

type Player struct {
	Index int    `json:"index"`
	Id    int    `json:"id"`
	Name  string `json:"name"`

	IsBot bool `json:"isBot"`

	statusMutex *sync.Mutex
	status      string

	inputMutex   *sync.Mutex
	inputMsg     string
	inputId      string
	inputType    string
	inputChan    chan string
	inputDDL     time.Time
	inputDone    bool
	inputTimeout bool

	wsMsgChans sync.Map
}

func NewPlayer(name string, index int) *Player {
	p := Player{
		Index: index,
		Id:    int(rand.Int31()),
		Name:  name,

		statusMutex: &sync.Mutex{},
		inputMutex:  &sync.Mutex{},
	}
	return &p
}

func (p *Player) AlertMessage(msg string, level string) {
	p.wsMsgChans.Range(func(key, value interface{}) bool {
		channel := key.(chan string)
		channel <- "alert:" + level + ":" + msg
		return true
	})
}

func (p *Player) GetInput(ctx context.Context, msg string, defaultValue string, inputType string) (string, error) {
	if p.IsBot {
		return defaultValue, nil
	}
	if p.inputChan != nil {
		return "", fmt.Errorf("already waiting for input")
	}
	logger.Info("Waiting for input", logger.String("player", p.Name))
	p.inputMutex.Lock()
	p.inputMsg = msg
	p.inputId = uuid.New().String()
	p.inputChan = make(chan string)
	p.inputDone = false
	p.inputTimeout = false
	p.inputType = inputType
	p.inputMutex.Unlock()
	defer func() {
		p.inputMutex.Lock()
		p.inputChan = nil
		p.inputMutex.Unlock()
	}()
	select {
	case <-ctx.Done():
		logger.Info("Input timeout, using default value.", logger.String("defaultValue", defaultValue), logger.String("player", p.Name))
		p.AlertMessage("超时，默认提交"+defaultValue, "warning")
		p.inputTimeout = true
		return defaultValue, nil
	case value := <-p.inputChan:
		logger.Info("Input received.", logger.String("value", value), logger.String("player", p.Name))
		p.inputDone = true
		return value, nil
	}
}

func (p *Player) Input(inputId string, value string) error {
	if p.inputDone {
		return fmt.Errorf("already input")
	}
	if p.inputTimeout {
		return fmt.Errorf("input timeout")
	}
	if p.inputChan == nil {
		return fmt.Errorf("not waiting for input")
	}
	if p.inputId != inputId {
		return fmt.Errorf("invalid input id")
	}
	p.inputChan <- value
	return nil
}
