package servers

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

type ServerList struct {
	Ports []int
	mu    sync.Mutex  
}

func (s *ServerList) populate(amount int) {
	if amount >= 10 {
		log.Fatal("ports can't exceed 10")
	}
	
	for i := 0; i < amount; i++ {
		s.Ports = append(s.Ports, i)
	}
}

func (s *ServerList) pop() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.Ports) == 0 {
		log.Fatal("no ports to pop")
	}
	port := s.Ports[0]
	s.Ports = s.Ports[1:]
	return port
}

func RunServer(amount int,wt *sync.WaitGroup) {
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
	port := sl.pop()  

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Server %d", port)
	})

	server := http.Server{
		Addr:    fmt.Sprintf(":808%d", port),
		Handler: mux,
	}
	
	fmt.Printf("Starting server on port 808%d\n", port)
	if err := server.ListenAndServe(); err != nil {
		log.Printf("Server %d failed: %v", port, err)
	}
}
