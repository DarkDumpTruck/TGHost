package tghost

import (
	"sync"
	"tghost/pkg/logger"

	"github.com/dop251/goja"
)

type Game struct {
	Players []*Player

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
	g.vm.Set("getInputAll", func(msg string, succMsg string, timeout int, defaultValue string, inputTypes []string) []string {
		wg := sync.WaitGroup{}
		wg.Add(len(g.Players))
		result := make([]string, len(g.Players))
		for i, player := range g.Players {
			index := i
			player := player
			go func() {
				defer wg.Done()
				value, err := player.GetInput(msg, timeout, defaultValue, inputTypes[index])
				if err != nil {
					logger.Info("player.GetInput error:", logger.Err(err))
					return
				}
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

	_, err := g.vm.RunString(g.code)
	if err != nil {
		return err
	}
	return nil
}
