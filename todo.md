# todo

- golang/x/net/websocket

- server struct holds a map about connections

- newserver function returns server and creates map

- s.handlews needs websocket.conn and adds connection to map

- s.readloop needs websocket.conn. makes infinite loop to read buffer and reply.

- main defines route and inits server. 
