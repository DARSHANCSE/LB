package servers

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

type ServerList struct {
	mu      sync.Mutex
	servers []server
}

type server struct {
	httpServer *http.Server
	active     int
	port       int
}

func (s *ServerList) populate(amount int) {
	if amount >= 10 {
		log.Fatal("ports can't exceed 10")
	}

	s.servers = make([]server, amount)
	for i := 0; i < amount; i++ {
		s.servers[i] = server{
			port:   i,
			active: 0,
		}
	}
}

func (s *ServerList) pop() *server {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.servers) == 0 {
		log.Fatal("no servers to pop")
	}
	s.servers = append(s.servers[1:],s.servers[0])
	return &s.servers[len(s.servers)-1]
}

func RunServer(amount int, wt *sync.WaitGroup) {
	defer wt.Done()
	var myserverlist ServerList
	myserverlist.populate(amount)

	var wg sync.WaitGroup
	wg.Add(amount)

	for i := 0; i < amount; i++ {
		go makeServer(&myserverlist, &wg)
	}

	wg.Wait()
}

func makeServer(sl *ServerList, wg *sync.WaitGroup) {
	defer wg.Done()

	s := sl.pop()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Server %d\n", s.port)
	})

	server := &http.Server{
		Addr:    fmt.Sprintf(":808%d\n",s.port),
		Handler: mux,
	}
	sl.mu.Lock()
	s.httpServer=server
	 sl.mu.Unlock()

	fmt.Printf("Starting server on port 808%d\n", s.port)
	if err := server.ListenAndServe(); err != nil {
		log.Printf("Server on port 808%d failed: %v\n", s.port , err)
	}
}
