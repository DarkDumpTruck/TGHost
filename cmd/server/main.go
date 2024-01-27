package main

import (
	"flag"
	"fmt"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"
	"tghost"
	"tghost/pkg/logger"

	"github.com/gin-contrib/pprof"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var cfgPath string

func setupHTTPRoutes() (*gin.Engine, error) {
	app := gin.New()
	app.Use(func(c *gin.Context) {
		logger.Debug("http request",
			logger.String("host", c.Request.Host),
			logger.String("path", c.Request.URL.Path),
			logger.String("method", c.Request.Method),
			logger.String("client", c.ClientIP()),
		)
		if !strings.Contains(c.Request.Host, "w.fun") &&
			!strings.Contains(c.Request.Host, "localhost:") {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.Next()
		logger.Debug("http response",
			logger.Int("status", c.Writer.Status()),
		)
	})

	pprof.Register(app, "debug/pprof")
	app.Group("metrics").Use(func(c *gin.Context) {
		promhttp.Handler().ServeHTTP(c.Writer, c.Request)
	})
	app.Use(static.Serve("/", static.LocalFile("/frontend", false)))
	app.NoRoute(func(c *gin.Context) {
		accepted := c.NegotiateFormat(gin.MIMEJSON, gin.MIMEHTML)
		if c.Request.Method == http.MethodGet && accepted == gin.MIMEHTML {
			c.File("/frontend/index.html")
			return
		}
		c.AbortWithStatus(http.StatusMethodNotAllowed)
	})

	v1 := app.Group("api/v1")
	room := v1.Group("room")
	room.GET("list", func(c *gin.Context) {
		showHidden := c.Query("showHidden") == "true"
		c.IndentedJSON(200, tghost.ListRooms(showHidden))
	})
	room.POST("create", func(c *gin.Context) {
		req := struct {
			Name      string `json:"name"`
			Code      string `json:"code"`
			PlayerNum int    `json:"playerNum"`
			BotNum    int    `json:"botNum"`
			JudgeNum  int    `json:"judgeNum"`
			Hidden    bool   `json:"hidden"`
		}{}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.IndentedJSON(200, gin.H{"error": err.Error()})
			return
		}

		script := &tghost.Script{
			Name:      req.Name,
			Code:      req.Code,
			PlayerNum: req.PlayerNum + req.BotNum + req.JudgeNum,
		}

		r := tghost.NewRoom(script, req.Hidden, req.BotNum)
		go func() {
			defer func() {
				if err := recover(); err != nil {
					logger.Error("room.Run() panic", logger.Err(err.(error)), logger.String("stacktrace", string(debug.Stack())))
				}
			}()
			r.Start()
		}()

		c.IndentedJSON(200, gin.H{"status": "ok", "roomId": r.Id})
	})

	game := v1.Group("game")
	game.POST("status", func(c *gin.Context) {
		req := struct {
			RoomId   int `json:"roomId"`
			PlayerId int `json:"playerId"`
		}{}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.IndentedJSON(200, gin.H{"error": err.Error()})
			return
		}
		room, err := tghost.GetRoom(req.RoomId)
		if err != nil {
			c.IndentedJSON(200, gin.H{"error": err.Error()})
			return
		}
		player, err := room.GetPlayer(req.PlayerId)
		if err != nil {
			c.IndentedJSON(200, gin.H{"error": err.Error()})
			return
		}
		c.IndentedJSON(200, room.GetPlayerStatus(player))
	})
	game.POST("input", func(c *gin.Context) {
		req := struct {
			RoomId   int    `json:"roomId"`
			PlayerId int    `json:"playerId"`
			Msg      string `json:"msg"`
			InputId  string `json:"inputId"`
		}{}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.IndentedJSON(200, gin.H{"error": err.Error()})
			return
		}
		room, err := tghost.GetRoom(req.RoomId)
		if err != nil {
			c.IndentedJSON(200, gin.H{"error": err.Error()})
			return
		}
		player, err := room.GetPlayer(req.PlayerId)
		if err != nil {
			c.IndentedJSON(200, gin.H{"error": err.Error()})
			return
		}

		err = player.Input(req.InputId, req.Msg)
		if err != nil {
			c.IndentedJSON(200, gin.H{"error": err.Error()})
			return
		}
		c.IndentedJSON(200, gin.H{"status": "ok"})
	})

	app.GET("ws/:roomId/:playerId", func(c *gin.Context) {
		roomId, err := strconv.Atoi(c.Param("roomId"))
		if err != nil {
			c.IndentedJSON(200, gin.H{"error": "roomId must be an integer"})
			return
		}
		playerId, err := strconv.Atoi(c.Param("playerId"))
		if err != nil {
			c.IndentedJSON(200, gin.H{"error": "playerId must be an integer"})
			return
		}
		room, err := tghost.GetRoom(roomId)
		if err != nil {
			c.IndentedJSON(200, gin.H{"error": err.Error()})
			return
		}
		player, err := room.GetPlayer(playerId)
		if err != nil {
			c.IndentedJSON(200, gin.H{"error": err.Error()})
			return
		}

		var upgrader = websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}
		defer ws.Close()

		room.Echo(player, ws)
	})

	return app, nil
}

func main() {
	flag.StringVar(&cfgPath, "c", "config/server.toml", "path to config file")
	flag.Parse()

	cfg, err := tghost.LoadConfigFromFile(cfgPath)
	if err != nil {
		panic(err)
	}

	app, err := setupHTTPRoutes()
	if err != nil {
		panic(err)
	}
	app.Run(fmt.Sprintf(":%d", cfg.HTTPPort))
}
