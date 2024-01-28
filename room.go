package tghost

import (
	"fmt"
	"slices"
	"sync"
	"tghost/pkg/logger"
	"time"

	"github.com/gorilla/websocket"
)

type Room struct {
	Id       int       `json:"id"`
	Name     string    `json:"name"`
	GameName string    `json:"gameName"`
	Players  []*Player `json:"players"`
	Running  bool      `json:"running"`
	Hidden   bool      `json:"hidden"`

	script *Script
	game   *Game
}

var (
	rooms      = make(map[int]*Room)
	roomsMutex = sync.RWMutex{}
)

func NewRoom(script *Script, hidden bool, botCount int) *Room {
	g := NewGame()
	g.SetCode(script.Code)
	roomsMutex.Lock()
	r := Room{
		Id:       len(rooms),
		Name:     fmt.Sprintf("房间#%d", len(rooms)),
		GameName: script.Name,
		Hidden:   hidden,
		script:   script,
		game:     g,
	}
	rooms[r.Id] = &r
	roomsMutex.Unlock()

	// TODO: maybe implement player manager here (support join/leave)
	for i := 0; i < r.script.PlayerNum; i++ {
		r.Players = append(r.Players, NewPlayer(fmt.Sprintf("Player#%d", i), i))
	}
	for i := 0; i < botCount; i++ {
		r.Players[i].IsBot = true
	}
	r.game.Players = r.Players

	return &r
}

func GetRoom(id int) (*Room, error) {
	roomsMutex.RLock()
	defer roomsMutex.RUnlock()
	if _, ok := rooms[id]; !ok {
		return nil, fmt.Errorf("room not found")
	}
	return rooms[id], nil
}

func ListRooms(showHidden bool) []*Room {
	roomsMutex.RLock()
	defer roomsMutex.RUnlock()
	var result []*Room
	for _, room := range rooms {
		if !showHidden && !room.Hidden {
			result = append(result, room)
		}
	}
	slices.SortFunc(result, func(i, j *Room) int {
		return -(i.Id - j.Id)
	})
	return result
}

func (r *Room) Start() {
	r.Running = true
	if err := r.game.Run(); err != nil {
		logger.Error("Code failed with error: ", logger.Err(err))
	}
	r.Running = false
}

func (r *Room) GetPlayer(id int) (*Player, error) {
	for _, player := range r.Players {
		if player.Id == id {
			return player, nil
		}
	}
	return nil, fmt.Errorf("player not found")
}

type PlayerStatus struct {
	GameName    string `json:"gameName"`
	GameStatus  string `json:"gameStatus"`
	GameRunning bool   `json:"gameRunning"`
	InputDone   bool   `json:"inputDone"`
	InputDDL    int    `json:"inputDDL"`
	InputId     string `json:"inputId"`
	InputMsg    string `json:"inputMsg"`
	InputType   string `json:"inputType"`
}

func (r *Room) GetPlayerStatus(player *Player) PlayerStatus {
	return PlayerStatus{
		GameName:    r.script.Name,
		GameStatus:  player.status,
		GameRunning: r.game.Running,
		InputDone:   !r.game.Running || player.inputDone || player.inputTimeout,
		InputDDL:    int(time.Until(player.inputDDL).Seconds()),
		InputId:     player.inputId,
		InputMsg:    player.inputMsg,
		InputType:   player.inputType,
	}
}

func (r *Room) Echo(player *Player, c *websocket.Conn) {
	channel := make(chan string, 1024)
	player.wsMsgChans.Store(channel, 1)

	go func() {
		for {
			_, message, err := c.ReadMessage()
			_ = message
			if err != nil {
				player.wsMsgChans.Delete(channel)
				break
			}
		}
	}()
	for {
		msg, ok := <-channel
		if !ok {
			c.Close()
			break
		}
		err := c.WriteJSON(msg)
		if err != nil {
			logger.Error("ws WriteJSON msg failed", logger.Err(err))
			break
		}
	}
}
