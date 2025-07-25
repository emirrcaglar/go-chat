
type Server struct {
	conns map[*websocket.Conn]bool
	mutex sync.RWMutex
}

func NewServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}

func (s *Server) handleWS(ws *websocket.Conn) {
	fmt.Printf("new connection: %v\n", ws.RemoteAddr())

	defer func() {
		s.mutex.Lock()
		delete(s.conns, ws)
		s.mutex.Unlock()
		log.Printf("connection closed: %v\n", ws.RemoteAddr())
	}()

	s.mutex.Lock()
	s.conns[ws] = true
	s.mutex.Unlock()

	s.readLoop(ws)
}

func (s *Server) readLoop(ws *websocket.Conn) {
	buf := make([]byte, 1024) // slice of bytes with size of 1024

	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("read error:%v", err)
			continue
		}
		msg := buf[:n]
		fmt.Println(string(msg))

		if _, err := ws.Write([]byte("Echo: " + string(msg))); err != nil {
			log.Printf("write error: %v", err)
			break
		}

		s.broadcast([]byte(fmt.Sprintf("Broadcast: %s", msg)))
	}
}

func (s *Server) broadcast(msg []byte) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	fmt.Println("Starting broadcast...")

	for conn := range s.conns {
		if _, err := conn.Write(msg); err != nil {
			log.Printf("Broadcast error to %v: %v", conn.RemoteAddr(), err)
		}
		log.Printf("broadcasting to connection: %v\n", conn.RemoteAddr())
		log.Printf("the broadcasted message: %s", string(msg[:]))
	}
	fmt.Println("Finishing broadcast...")
}
