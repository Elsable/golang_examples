package main

import (
	"fmt"
	"time"
	"math/rand"
)


var (
    Web1   = fakeSearch("web1")
    Web2   = fakeSearch("web2")
    Image1 = fakeSearch("image1")
    Image2 = fakeSearch("image2")
    Video1 = fakeSearch("video1")
    Video2 = fakeSearch("video2")
)

type Result string

type Search func(query string) Result

func fakeSearch(kind string) Search{
	return func(query string) Result{
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		return Result(fmt.Sprintf("%s result for %q \n", kind, query))
	}
}

func First(query string, replicas ...Search) Result{
	c := make(chan Result)
	searchReplica := func(i int) { c <- replicas[i](query) }
	for i := range replicas{
		go searchReplica(i)
	}
	return <- c
}

func Google(query string)(results []Result){
	c := make(chan Result)
	go func (){ c <- First(query, Web1, Web2) }()
	go func (){ c <- First(query, Video1, Video2) }()
	go func (){ c <- First(query, Image1, Image2) }()
	timeout := time.After(80 * time.Millisecond)
	for i := 0; i < 3; i++{
		select{
		case result := <- c: 
			results = append(results, result)
		case <- timeout:
			fmt.Println("server timed out")
			return
		}
	}
	
	// go func (){ c <- Web(query) }()
	// go func (){ c <- Videos(query) }()
	// go func (){ c <- Images(query) }()
	
	// results = append(results, Web(query))
	// results = append(results, Videos(query))
	// results = append(results, Images(query))
	return
}

func main(){
	rand.Seed(time.Now().UnixNano())
	start := time.Now()
	results := Google("golang")
	// results := First("golang", fakeSearch("rep 1"), fakeSearch("rep 2"))
	elaspe := time.Since(start)
	fmt.Println(results)
	fmt.Println(elaspe)
}