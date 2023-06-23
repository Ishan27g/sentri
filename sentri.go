package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/DavidGamba/dgtools/run"
	"github.com/Ishan27g/sentri/hook"
	"github.com/Ishan27g/sentri/internal"
	"github.com/fatih/color"
	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"github.com/Ishan27g/sentri/internal/transport"
	"github.com/gorilla/websocket"
	"github.com/nats-io/nats.go"
)

type app string

var (
	port = func() string {
		p := os.Getenv("SENTRI_UI_PORT")
		if p == "" {
			return ":4999"
		}
		return p
	}

	toStdOut  = true
	toWebSock = false
	timestamp = false

	s           *sentri
	colorsShell = []color.Attribute{color.FgRed, color.FgGreen, color.FgYellow, color.FgBlue, color.FgMagenta, color.FgCyan}
	colorsHtml  = []string{"#E18E2E", "#80C5CD", "#83F52C", "#DB9370", "#7CFC00"}

	appShell   = map[app]*color.Attribute{}
	appHtml    = map[app]*string{}
	cIndexSh   = 0
	cIndexHt   = 0
	webLogChan chan interface{}
)

type sentri struct {
	count    int
	outputCh chan internal.Output
	p        *internal.Pool
}

func getOrNextHtml(from string) string {
	if appHtml[app(from)] != nil {
		return *appHtml[app(from)]
	}
	if cIndexHt == len(colorsHtml) {
		cIndexHt = 0
	}
	cIndexHt = cIndexHt + 1
	appHtml[app(from)] = &colorsHtml[cIndexHt]
	return *appHtml[app(from)]
}

func getOrNextShColor(from string) *color.Attribute {
	if appShell[app(from)] != nil {
		return appShell[app(from)]
	}
	cIndexSh = cIndexSh + 1
	if cIndexSh == len(colorsShell) {
		cIndexSh = 0
	}
	appShell[app(from)] = &colorsShell[cIndexSh]
	return appShell[app(from)]
}

type UiLog struct {
	From  string `json:"prefix"`
	Color string `json:"color"`
	Text  string `json:"text"`
}

func writeFromNats() {
	transport.NatsSubscribe(internal.Channel, func(msg *nats.Msg) {
		s.count++
		var o = internal.Output{}
		err := json.Unmarshal(msg.Data, &o)
		if err != nil {
			fmt.Println("sentri json error on nats msg", err.Error())
		}
		s.outputCh <- o
	})
}

func upgrade(writer http.ResponseWriter, request *http.Request) (*websocket.Conn, context.Context, context.CancelFunc) {
	upgrade := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrade.Upgrade(writer, request, nil)
	if err != nil {
		log.Println(err)
		return nil, nil, nil
	}
	ctx, cancel := context.WithCancel(context.Background())
	conn.SetCloseHandler(func(code int, text string) error {
		fmt.Println("closing web-socket connection error - " + err.Error())
		if ctx != nil {
			cancel()
		}
		_ = conn.Close()
		fmt.Println("closed web-socket connection error - " + err.Error())
		return nil
	})
	return conn, ctx, cancel
}

func main() {

	if len(os.Args) >= 2 {
		hook.Cmd(os.Args[1], os.Stdout, run.CMD(os.Args[1:]...))
		// wait for publish
		<-time.After(2 * time.Second)
		return
	}

	appShell = make(map[app]*color.Attribute)
	appHtml = make(map[app]*string)

	s = &sentri{}
	s.outputCh = make(chan internal.Output)
	s.p = internal.GetPool()

	go func() {
		lastCount := 0
		for {
			<-time.After(3 * time.Second)
			if lastCount != s.count {
				//fmt.Println(fmt.Sprintf("%d NEW LOGS, %d TOTAL LOGS", s.count-lastCount, s.count))
				lastCount = s.count
			}
		}
	}()
	go writeFromNats()
	go func() {
		for o := range s.outputCh {
			if toStdOut {
				shClr := getOrNextShColor(o.From)
				c := color.New(*shClr, color.Bold)
				if timestamp {
					fmt.Println(c.Sprintf(o.From), o.Timestamp.Format(time.RFC3339), o.Log)
				} else {
					fmt.Println(c.Sprintf(o.From), o.Log)
				}
			}
			if toWebSock {
				htmlClr := getOrNextHtml(o.From)
				u := UiLog{
					From:  o.From,
					Color: htmlClr,
				}
				if timestamp {
					u.Text = o.Timestamp.Format(time.RFC3339) + " " + o.Log
				} else {
					u.Text = o.Log
				}
				webLogChan <- u
			}
		}
		close(s.outputCh)
	}()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})
	router := mux.NewRouter().StrictSlash(true)
	handler := c.Handler(router)

	router.HandleFunc("/logs", func(writer http.ResponseWriter, request *http.Request) {
		conn, ctx, cancel := upgrade(writer, request)
		if conn == nil || ctx == nil || cancel == nil {
			return
		}
		defer cancel()
		defer conn.Close()

		toWebSock = true
		webLogChan = make(chan any, 1)
		defer func() {
			toWebSock = false
			close(webLogChan)
		}()

		for {
			messageType, p, err := conn.ReadMessage()
			if err != nil {
				break
			}
			if messageType != websocket.TextMessage {
				fmt.Println("messageType is not TextMessage", messageType)
				break
			}
			fmt.Println("Received request to log")
			for {
				select {
				case l := <-webLogChan:
					_ = conn.WriteJSON(l)
				case <-ctx.Done():
					fmt.Println("Stopped stream for", string(p))
					return
				}
			}
		}
		return
	})

	fmt.Println("http started on", port())
	err := http.ListenAndServe(port(), handler)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

}
