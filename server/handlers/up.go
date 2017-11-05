package handlers

import (
	api "../docker-api"
)

func Up(
	lang []byte,
	body []byte,
	connection *websocket.Conn,
	messageType int,
) {
	// go func() {
	var guid = xid.New().String()
	var dir = guid
	var name = guid

	port, err := freeport.GetFreePort()
	if err != nil {
		panic(err)
	}

	initWorkDirectory(string(lang), dir)
	notify(connection, messageType, "work dir initialized")
	dispatchFuncion(string(lang), body, dir)
	notify(connection, messageType, "function dispatched")

	api.Build(name, dir)

	notify(connection, messageType, "function built")
	api.Deploy(name, dir, strconv.Itoa(port))
	notify(connection, messageType, "function deployed: http://localhost:"+strconv.Itoa(port))

	closeConnection(connection)
	// }()
}

