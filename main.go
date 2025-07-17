package main

import (
	"loadbalancer/loadbalancer"
	"loadbalancer/servers"
	"sync"
)


func main (){
	var wg sync.WaitGroup
	wg.Add(2)	
	go servers.RunServer(5,&wg)
	go loadbalancer.MakeLoadBalancer(5,&wg)
	wg.Wait()
}