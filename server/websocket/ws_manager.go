package websocket

import (
	"context"
	"errors"
	"io"
	"net"
	"sync"
	"test/settings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/websocket"
)

type connection struct {
	ws       *websocket.Conn
	shutdown chan struct{}
}

type eww struct {
	events    []any
	timestamp time.Time
}

type waiter struct {
	wsID      string
	timestamp time.Time
}

type WebsocketManager struct {
	sync.Mutex

	connections         map[string]connection
	waiters             map[string]waiter
	eventsWithoutWaiter map[string]eww
}

func (wsm *WebsocketManager) WS(ctx echo.Context) error {
	s := websocket.Server{Handler: websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		id, shutdown := wsm.registerWS(ws)

	LOOP:
		for {
			select {
			case <-shutdown:
				break LOOP
			default:
				if err := ws.SetReadDeadline(time.Now().Add(settings.ReadDeadline)); err != nil {
					log.Error().Str("function", "WS").Err(err).Msg("failed to set read deadline")
				}

				var events []string
				err := websocket.JSON.Receive(ws, &events)
				if err != nil {
					if errors.Is(err, io.EOF) {
						break LOOP
					} else if err, ok := err.(net.Error); ok && err.Timeout() {
						continue
					}

					log.Error().Str("function", "WS").Err(err).Msg("failed to receive events from websocket")
					continue
				}

				wsm.registerWaiter(ws, id, events)
			}
		}

		wsm.removeWS(id)
	})}

	s.ServeHTTP(ctx.Response(), ctx.Request())
	return nil
}

func (wsm *WebsocketManager) Send(eventID string, data any) error {
	if eventID == "" {
		return errors.New("empty event ID")
	}

	wsm.Lock()
	defer wsm.Unlock()

	if waiter, ok := wsm.waiters[eventID]; ok {
		if conn, ok := wsm.connections[waiter.wsID]; ok {
			err := websocket.JSON.Send(conn.ws, map[string]any{
				"eventID": eventID,
				"data":    data,
			})
			if err != nil {
				return err
			}

			delete(wsm.waiters, eventID)
		} else {
			delete(wsm.waiters, eventID)
			return errors.New("connection not found")
		}
	} else {
		wsm.eventsWithoutWaiter[eventID] = eww{
			events:    append(wsm.eventsWithoutWaiter[eventID].events, data),
			timestamp: time.Now(),
		}
	}

	return nil
}

func (wsm *WebsocketManager) registerWS(ws *websocket.Conn) (string, chan struct{}) {
	wsm.Lock()
	defer wsm.Unlock()

	id := uuid.NewString()
	shutdown := make(chan struct{})
	wsm.connections[id] = connection{ws: ws, shutdown: shutdown}
	return id, shutdown
}

func (wsm *WebsocketManager) registerWaiter(ws *websocket.Conn, wsID string, subscrs []string) {
	wsm.Lock()
	defer wsm.Unlock()

	for _, eventID := range subscrs {
		err := websocket.JSON.Send(ws, map[string]any{
			"eventID": eventID,
			"data":    "Subscribed",
		})
		if err != nil {
			log.Error().Str("function", "ttl").Err(err).Str("eventID", eventID).Msg("failed to send message to websocket")
		}

		if data, ok := wsm.eventsWithoutWaiter[eventID]; ok {
			for _, event := range data.events {
				err := websocket.JSON.Send(ws, map[string]any{
					"eventID": eventID,
					"data":    event,
				})
				if err != nil {
					log.Error().Str("function", "registerWaiter").Err(err).Interface("event", event).Msg("failed to send websocket message")
				}
			}

			delete(wsm.eventsWithoutWaiter, eventID)
			continue
		}

		wsm.waiters[eventID] = waiter{
			wsID:      wsID,
			timestamp: time.Now(),
		}
	}
}

func (wsm *WebsocketManager) removeWS(id string) {
	wsm.Lock()
	defer wsm.Unlock()

	for event, waiter := range wsm.waiters {
		if waiter.wsID == id {
			delete(wsm.waiters, event)
		}
	}

	delete(wsm.connections, id)
}

func (wsm *WebsocketManager) ttl() {
	for now := range time.Tick(settings.WSTicker) {
		wsm.Lock()
		for eventID, waiter := range wsm.waiters {
			if now.Sub(waiter.timestamp) > settings.WaiterTTL {
				if conn, ok := wsm.connections[waiter.wsID]; ok {
					err := websocket.JSON.Send(conn.ws, map[string]any{
						"eventID": eventID,
						"data":    "Event timeout",
					})
					if err != nil {
						log.Error().Str("function", "ttl").Err(err).Str("eventID", eventID).Msg("failed to send message to websocket")
					}
				}

				delete(wsm.waiters, eventID)
			}
		}

		for k, v := range wsm.eventsWithoutWaiter {
			if now.Sub(v.timestamp) > settings.EventsWithoutWaiterTTL {
				delete(wsm.eventsWithoutWaiter, k)
			}
		}
		wsm.Unlock()
	}
}

func (wsm *WebsocketManager) Stop(ctx context.Context) error {
	var wg sync.WaitGroup

	wg.Add(len(wsm.connections))
	for _, conn := range wsm.connections {
		go func(shutdown chan struct{}, wg *sync.WaitGroup) {
			shutdown <- struct{}{}
			wg.Done()
		}(conn.shutdown, &wg)
	}

	wg.Wait()
	return nil
}

func NewWebsocketManager() *WebsocketManager {
	wsm := &WebsocketManager{
		connections:         make(map[string]connection),
		waiters:             make(map[string]waiter),
		eventsWithoutWaiter: make(map[string]eww),
	}

	go wsm.ttl()
	return wsm
}
