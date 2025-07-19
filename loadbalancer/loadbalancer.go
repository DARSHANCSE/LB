package loadbalancer

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"sync"
)


var baseURL="http://localhost:808"
var ind =0




type loadbalancer struct {
	RevProxy httputil.ReverseProxy
}


type Endpoints struct{
	List []*url.URL;

}

func MakeLoadBalancer (amount int , wt *sync.WaitGroup ){
	defer wt.Done()

	var lb loadbalancer
	var ep Endpoints


	router :=http.NewServeMux()

	server := http.Server{
		Addr:":6969",
		Handler: router,
	}
	
	for i:=0;i<amount;i++{
	ep.List=append(ep.List, createurl(baseURL,i))
	}
	router.HandleFunc("/",makeRequest(&lb,&ep))
	router.HandleFunc("/health", makeRequest(&lb,&ep))
	router.HandleFunc("/disable", makeRequest(&lb,&ep))
	log.Fatal(server.ListenAndServe())
}

func makeRequest(lb *loadbalancer, ep *Endpoints) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := ind
		ind = (ind + 1) % len(ep.List)
		for {
			target := ep.List[ind]
			resp, err := http.Get(target.String() + "/health")
			if err == nil && resp.StatusCode == http.StatusOK {
				resp.Body.Close()
				break
			}
			ind = (ind + 1) % len(ep.List)
			if ind == start {
				http.Error(w, "No healthy backend servers", http.StatusServiceUnavailable)
				return
			}
		}

		director := func(req *http.Request) {
			req.URL.Scheme = ep.List[ind].Scheme
			req.URL.Host = ep.List[ind].Host
			req.URL.Path = r.URL.Path // Keep original path
			req.URL.RawQuery = r.URL.RawQuery
			req.Host = ep.List[ind].Host
		}

		lb.RevProxy = httputil.ReverseProxy{Director: director}
		lb.RevProxy.ServeHTTP(w, r)
		fmt.Printf("Forwarded to: %s%s\n", ep.List[ind].Host, r.URL.Path)
	}
}



func createurl (baseURL string , ind int)*url.URL{
	link:=baseURL + strconv.Itoa(ind)
	url,_:=url.Parse(link)
	return url 
}