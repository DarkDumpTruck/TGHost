package tghost

import (
	"context"
	"sync"
	"tghost/pkg/logger"
	"time"

	"github.com/dop251/goja"
)

type Game struct {
	Players []*Player
	Running bool

	vm   *goja.Runtime
	code string
	ctx  context.Context
}

func NewGame() *Game {
	g := Game{
		vm:   goja.New(),
		code: "",
		ctx:  context.Background(),
	}

	return &g
}

func (g *Game) SetCode(code string) {
	g.code = baseCode + code
}

func (g *Game) Run() error {
	g.vm.Set("getPlayers", func() goja.Value {
		return g.vm.ToValue(g.Players)
	})
	g.vm.Set("alertAll", func(msg string) {
		for _, player := range g.Players {
			player := player
			go func() {
				player.AlertMessage(msg, "info")
			}()
		}
	})
	g.vm.Set("getInputs", func(msg string, succMsg string, timeout int, _defaultValue goja.Value, indexes []int, _inputType goja.Value, _checkers goja.Value) []string {
		count := len(indexes)
		if len(g.Players) < count {
			logger.Error("getInputs: too many indexes", logger.Int("count", count), logger.Int("playerNum", len(g.Players)))
			return []string{}
		}
		result := make([]string, count)
		var defaultValues []string
		if err := g.vm.ExportTo(_defaultValue, &defaultValues); err != nil {
			var defaultValue string
			if err := g.vm.ExportTo(_defaultValue, &defaultValue); err != nil {
				logger.Error("got error type for defaultValue", logger.Err(err))
				return result
			}
			defaultValues = make([]string, count)
			for i := 0; i < count; i++ {
				defaultValues[i] = defaultValue
			}
		}
		var inputTypes []string
		if err := g.vm.ExportTo(_inputType, &inputTypes); err != nil {
			var inputType string
			if err := g.vm.ExportTo(_inputType, &inputType); err != nil {
				logger.Error("got error type for inputTypes", logger.Err(err))
				return result
			}
			inputTypes = make([]string, count)
			for i := 0; i < count; i++ {
				inputTypes[i] = inputType
			}
		}
		var checkers []func(string) string
		if err := g.vm.ExportTo(_checkers, &checkers); err != nil {
			var checker func(string) string
			if err := g.vm.ExportTo(_checkers, &checker); err != nil {
				logger.Error("got error type for checkers", logger.Err(err))
				return result
			}
			checkers = make([]func(string) string, count)
			for i := 0; i < count; i++ {
				checkers[i] = checker
			}
		}
		wg := sync.WaitGroup{}
		wg.Add(count)
		for i, index := range indexes {
			i := i
			index := index
			go func() {
				defer wg.Done()
				player := g.Players[index]
				player.inputMutex.Lock()
				player.inputDDL = time.Now().Add(time.Duration(timeout) * time.Second)
				player.inputMutex.Unlock()
				ctx, cancel := context.WithTimeout(g.ctx, time.Duration(timeout)*time.Second)
				defer cancel()
				for {
					value, err := player.GetInput(ctx, msg, defaultValues[i], inputTypes[i])
					if err != nil {
						logger.Info("player.GetInput error:", logger.Err(err))
						return
					}
					if value == defaultValues[i] {
						result[i] = defaultValues[i]
						break
					}
					checkResult := checkers[i](value)
					if checkResult == "" {
						result[i] = value
						player.AlertMessage("提交成功", "success")
						break
					}
					player.AlertMessage(checkResult, "warning")
				}
				player.inputMsg = succMsg
			}()
		}
		wg.Wait()
		return result
	})
	g.vm.Set("updateStatus", func(index int, status string) {
		player := g.Players[index]
		player.statusMutex.Lock()
		player.status = status
		player.statusMutex.Unlock()

		player.wsMsgChans.Range(func(key, value interface{}) bool {
			channel := key.(chan string)
			channel <- "update"
			return true
		})
	})
	g.vm.Set("appendStatusAll", func(status string) {
		for _, player := range g.Players {
			player := player
			go func() {
				player.statusMutex.Lock()
				player.status += status
				player.statusMutex.Unlock()

				player.wsMsgChans.Range(func(key, value interface{}) bool {
					channel := key.(chan string)
					channel <- "update"
					return true
				})
			}()
		}
	})
	g.vm.Set("log", func(msg_type string, msg string) {
		logger.Info("game log: "+msg_type, logger.String("message", msg))
	})

	g.Running = true
	defer func() {
		g.Running = false
	}()
	_, err := g.vm.RunString(g.code)
	if err != nil {
		return err
	}

	var main func() int
	err = g.vm.ExportTo(g.vm.Get("main"), &main)
	if err != nil {
		return err
	}

	code := main()
	logger.Info("Game finished with code", logger.Int("code", code))

	return nil
}
