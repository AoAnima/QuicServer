package main

import (
	"bytes"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"

	. "aoanima.ru/ConnQuic"
	. "aoanima.ru/Logger"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

var hub = newHub()

func WSСервер() {
	Инфо("  %+v \n", "WSСервер")

	go hub.run()
	конфигТлс, err := КлиентскийТлсКонфиг()
	if err != nil {
		Ошибка("  %+v \n", err)
	}
	конфигТлс.InsecureSkipVerify = true
	конфигТлс.ServerName = "localhost"

	srv := &http.Server{
		Addr:      ":444",
		Handler:   http.HandlerFunc(serveWs),
		TLSConfig: конфигТлс,
	}

	err = srv.ListenAndServeTLS("cert/server.crt", "cert/server.key")
	if err != nil {
		Ошибка(" %s ", err)
	}
}

/*
TextMessage = 1

	// BinaryMessage denotes a binary data message.
	BinaryMessage = 2

	// CloseMessage denotes a close control message. The optional message
	// payload contains a numeric code and text. Use the FormatCloseMessage
	// function to format a close message payload.
	CloseMessage = 8

	// PingMessage denotes a ping control message. The optional message payload
	// is UTF-8 encoded text.
	PingMessage = 9

	// PongMessage denotes a pong control message. The optional message payload
	// is UTF-8 encoded text.
	PongMessage = 10
*/
func (c *Client) readPump() {
	// defer func() {
	// 	c.hub.unregister <- c
	// 	c.conn.Close()
	// }()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		t, message, err := c.conn.ReadMessage()
		if err != nil {
			Ошибка("error: %+v message %+v %+v", err.Error(), message, t)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				Ошибка("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.hub.broadcast <- message
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	// defer func() {
	// 	// ticker.Stop()
	// 	// c.conn.Close()
	// }()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				Инфо(" websocket.CloseMessage %+v \n", websocket.CloseMessage)
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				Ошибка("  %+v \n", err)
				return
			}
		case <-ticker.C:
			// c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			Инфо(" нужно отправляем пинг  %+v \n", websocket.PingMessage)
			// w, err := c.conn.NextWriter(websocket.PingMessage)
			// Инфо("w  %+v \n", w)
			// if err != nil {
			// 	Ошибка("  %+v \n", err.Error())
			// }
			// b, err := w.Write([]byte("ping"))
			// Инфо("передано b  %+v \n", b)
			// if err != nil {
			// 	Ошибка("  %+v \n", err)
			// }
			// // if err := c.conn.WriteMessage(websocket.PingMessage, []byte("ping")); err != nil {

			// 	Ошибка("  %+v \n", err.Error())

			// 	return
			// }
		}
	}
}

// serveWs handles websocket requests from the peer.
func serveWs(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}

type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				Ошибка("  %+v \n", "close(client.send)")
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					Ошибка("  %+v \n", "close(client.send)")
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
