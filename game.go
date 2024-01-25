package tghost

import (
	"fmt"
	"sync"
	"tghost/pkg/logger"

	"github.com/dop251/goja"
)

type Game struct {
	Players []*Player
	Running bool

	vm   *goja.Runtime
	code string
}

func NewGame() *Game {
	g := Game{
		vm:   goja.New(),
		code: "",
	}

	return &g
}

func (g *Game) SetCode(code string) {
	g.code = code
}

func (g *Game) Run() error {
	g.vm.Set("getPlayers", func() goja.Value {
		return g.vm.ToValue(g.Players)
	})
	g.vm.Set("alertAll", func(msg string) {
		for _, player := range g.Players {
			player := player
			go func() {
				player.wsMsgChans.Range(func(key, value interface{}) bool {
					channel := key.(chan string)
					channel <- "alert:" + msg
					return true
				})
			}()
		}
	})
	g.vm.Set("getInputs", func(msg string, succMsg string, timeout int, defaultValue string, indexes []int, inputTypes []string, checkers []func(string)bool) []string {
		count := len(indexes)
		wg := sync.WaitGroup{}
		wg.Add(count)
		result := make([]string, len(g.Players))
		for _, index := range indexes {
			index := index
			go func() {
				defer wg.Done()
				player := g.Players[index]
				value, err := player.GetInput(msg, timeout, defaultValue, inputTypes[index])
				if err != nil {
					logger.Info("player.GetInput error:", logger.Err(err))
					return
				}
				checkResult := checkers[index](value)
				fmt.Println("checkResult:", checkResult) // TODO: use check result to retry
				result[index] = value
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
