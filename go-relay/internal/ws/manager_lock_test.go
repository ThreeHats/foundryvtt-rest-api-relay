package ws

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestRemoveClientCallbackRunsOutsideManagerLock(t *testing.T) {
	assertCallbackCanRead := func(t *testing.T, remove func(*ClientManager, *Client)) {
		t.Helper()

		manager := NewClientManager(nil, "test")
		client := &Client{id: "client-1", apiKey: "test-key"}
		manager.clients[client.id] = client
		manager.tokenGroups[client.apiKey] = map[string]bool{client.id: true}

		callbackLockWasFree := false
		manager.OnClientRemoved = func(string, string) {
			callbackLockWasFree = manager.mu.TryRLock()
			if callbackLockWasFree {
				manager.mu.RUnlock()
			}
		}

		remove(manager, client)
		if !callbackLockWasFree {
			t.Fatal("OnClientRemoved ran while the client-manager lock was held")
		}
	}

	t.Run("RemoveClient", func(t *testing.T) {
		assertCallbackCanRead(t, func(manager *ClientManager, client *Client) {
			manager.RemoveClient(client.id)
		})
	})

	t.Run("RemoveClientIfMatch", func(t *testing.T) {
		assertCallbackCanRead(t, func(manager *ClientManager, client *Client) {
			manager.RemoveClientIfMatch(client.id, client)
		})
	})
}

func TestReconnectWaitsForDisconnectCleanup(t *testing.T) {
	manager := NewClientManager(nil, "test")
	oldClient := &Client{id: "client-1", apiKey: "old-key"}
	manager.clients[oldClient.id] = oldClient
	manager.tokenGroups[oldClient.apiKey] = map[string]bool{oldClient.id: true}

	callbackEntered := make(chan struct{})
	finishCallback := make(chan struct{})
	manager.OnClientRemoved = func(id, _ string) {
		if id == oldClient.id {
			close(callbackEntered)
			<-finishCallback
		}
	}

	removed := make(chan struct{})
	go func() {
		manager.RemoveClient(oldClient.id)
		close(removed)
	}()
	<-callbackEntered

	otherServerConn, otherClientConn := testWebSocketPair(t)
	defer otherClientConn.Close()
	otherClient, err := manager.AddClient(otherServerConn, "client-2", "other-key", "test", "wands", "Wands", "14", "dnd5e", "D&D5e", "6", "", "")
	if err != nil {
		t.Fatalf("unrelated client registration blocked by disconnect cleanup: %v", err)
	}
	manager.RemoveClient(otherClient.id)

	serverConn, clientConn := testWebSocketPair(t)
	defer clientConn.Close()
	added := make(chan *Client, 1)
	addErr := make(chan error, 1)
	go func() {
		client, err := manager.AddClient(serverConn, oldClient.id, "new-key", "test", "wands", "Wands", "14", "dnd5e", "D&D5e", "6", "", "")
		if err != nil {
			addErr <- err
			return
		}
		added <- client
	}()

	select {
	case client := <-added:
		t.Fatalf("same-ID reconnect registered before disconnect cleanup finished: %p", client)
	case err := <-addErr:
		t.Fatalf("same-ID reconnect failed unexpectedly: %v", err)
	case <-time.After(25 * time.Millisecond):
	}

	close(finishCallback)
	select {
	case <-removed:
	case <-time.After(time.Second):
		t.Fatal("disconnect cleanup did not finish")
	}

	var replacement *Client
	select {
	case replacement = <-added:
	case err := <-addErr:
		t.Fatalf("same-ID reconnect failed: %v", err)
	case <-time.After(time.Second):
		t.Fatal("same-ID reconnect did not proceed after cleanup")
	}
	if got := manager.GetClient(oldClient.id); got != replacement {
		t.Fatal("disconnect cleanup removed or replaced the newly registered client")
	}

	manager.OnClientRemoved = nil
	manager.RemoveClient(oldClient.id)
}

func testWebSocketPair(t *testing.T) (*websocket.Conn, *websocket.Conn) {
	t.Helper()
	accepted := make(chan *websocket.Conn, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := (&websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}).Upgrade(w, r, nil)
		if err != nil {
			t.Errorf("accept websocket: %v", err)
			return
		}
		accepted <- conn
	}))
	t.Cleanup(server.Close)

	client, _, err := websocket.DefaultDialer.Dial("ws"+server.URL[len("http"):], nil)
	if err != nil {
		t.Fatalf("dial websocket: %v", err)
	}
	select {
	case serverConn := <-accepted:
		return serverConn, client
	case <-time.After(time.Second):
		client.Close()
		t.Fatal("server did not accept websocket")
		return nil, nil
	}
}
