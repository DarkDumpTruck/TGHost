package tghost

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Player struct {
	Index int    `json:"index"`
	Id    int    `json:"id"`
	Name  string `json:"name"`

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

func (p *Player) GetInput(msg string, timeout int, defaultValue string, inputType string) (string, error) {
	if p.inputChan != nil {
		return "", fmt.Errorf("already waiting for input")
	}
	fmt.Println("Waiting for input from", p.Name)
	p.inputMutex.Lock()
	p.inputMsg = msg
	p.inputId = uuid.New().String()
	p.inputChan = make(chan string)
	p.inputDone = false
	p.inputTimeout = false
	p.inputDDL = time.Now().Add(time.Duration(timeout) * time.Second)
	p.inputType = inputType
	p.inputMutex.Unlock()
	defer func() {
		p.inputMutex.Lock()
		p.inputChan = nil
		p.inputMutex.Unlock()
	}()
	select {
	case <-time.After(time.Duration(timeout) * time.Second):
		p.inputTimeout = true
		return defaultValue, nil
	case value := <-p.inputChan:
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
