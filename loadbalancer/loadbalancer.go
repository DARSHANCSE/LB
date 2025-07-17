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
	router.HandleFunc("/loadbalancer",makeRequest(&lb,&ep))
	log.Fatal(server.ListenAndServe())
}

func makeRequest(lb *loadbalancer,ep *Endpoints) func (w http.ResponseWriter,r *http.Request){
	return func(w http.ResponseWriter, r *http.Request) {
		ind=(ind+1)%len(ep.List)
		lb.RevProxy = *httputil.NewSingleHostReverseProxy(ep.List[ind])
		lb.RevProxy.ServeHTTP(w,r)
		fmt.Printf("testdone at port :808%d",ind)

	}
}


func createurl (baseURL string , ind int)*url.URL{
	link:=baseURL + strconv.Itoa(ind)
	url,_:=url.Parse(link)
	return url 
}