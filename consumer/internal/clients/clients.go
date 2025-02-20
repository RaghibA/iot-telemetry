package clients

import (
	"sync"

	"github.com/gorilla/websocket"
)

/**
 * Client Directory
 *
 * Stores the connection to the active ws conn
 * using the deviceId as the key
 */
type ClientDirectory struct {
	mu      sync.Mutex
	clients map[string]*websocket.Conn
}

/**
 * Init the client directory map & return ref to global var
 *
 * @output - *ClientDirectory: pointer to client dir
 */
func InitClientDirectory() *ClientDirectory {
	return &ClientDirectory{
		clients: make(map[string]*websocket.Conn),
	}
}

/**
 * Adds a client to the directory
 *
 * @params (deviceId string, conn *websocket.Conn): deviceId from header & active ws conn
 */
func (cd *ClientDirectory) AddClient(deviceID string, conn *websocket.Conn) {
	cd.mu.Lock()
	defer cd.mu.Unlock()

	cd.clients[deviceID] = conn
}

/**
 * Removes a client from the dir
 *
 * @params (deviceID string): user deviceID
 */
func (cd *ClientDirectory) RemoveClient(deviceID string) {
	cd.mu.Lock()
	defer cd.mu.Unlock()

	delete(cd.clients, deviceID)
}

/**
 * Get ws conn for a user with deviceID
 *
 * @params (deviceID string)
 * @output (*websocket.Conn, bool): active conn & true if value exists in dir map
 */
func (cd *ClientDirectory) GetClientConn(deviceID string) (*websocket.Conn, bool) {
	cd.mu.Lock()
	defer cd.mu.Unlock()

	conn, exists := cd.clients[deviceID]
	return conn, exists
}
