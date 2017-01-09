package api

import (
	"encoding/json"
	"github.com/googollee/go-socket.io"
	"log"
	"net/http"
	"time"
)

type customServer struct {
	Server *socketio.Server
}

func (s *customServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	origin := r.Header.Get("Origin")
	w.Header().Set("Access-Control-Allow-Origin", origin)
	s.Server.ServeHTTP(w, r)
}

func Ready() {
	ioServer := SocketIoServer()
	wsServer := new(customServer)
	wsServer.Server = ioServer

	http.Handle("/socket.io/", wsServer)

	port := "7000"
	log.Println("[Main] Starting Server Port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

type Packet struct {
	Cmd string `json:"cmd"`
}

type Response struct {
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
}

func sendMessage(so socketio.Socket, channel string, message string) {
	response := Response{
		time.Now(),
		message,
	}
	bytes, err := json.Marshal(response)
	if err != nil {
		return
	}
	so.Emit(channel, string(bytes))
}

func processMessage(so socketio.Socket, packet Packet) {

	ticker := time.NewTicker(1 * time.Second)
	stop := make(chan bool)
	go func() {
	loop:
		for {
			select {
			case t := <-ticker.C:
				log.Println("Send", t)
				sendMessage(so, "hoge", t.String())
			case <-stop:
				break loop
			}
		}
		//fmt.Println("Reachable!")
	}()

	time.Sleep(10 * time.Second)
	ticker.Stop()
	close(stop)

}

func SocketIoServer() *socketio.Server {

	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	server.On("connection", func(so socketio.Socket) {

		println(so.Id() + " joined clients.")
		so.Join("chat")

		so.On("msg", func(msg string) {
			packet := Packet{}
			json.Unmarshal([]byte(msg), &packet)
			processMessage(so, packet)
			// log.Println("emit:", so.Emit("chat message", msg))
			// so.BroadcastTo("chat", "chat message", msg)
		})

		so.On("disconnection", func() {
			log.Println("on disconnect")
		})

	})
	server.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})

	return server
}
