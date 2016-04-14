//	Sample use :
//
//	./pinger --ip http://dgraph.io/query --numuser 3
//
//
//

package main

import (
	"bytes"
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/dgraph-io/dgraph/x"
)

var (
	numUser    = flag.Int("numuser", 1, "number of users hitting simultaneously")
	numReq     = flag.Int("numreq", 10, "number of request per user")
	serverAddr = flag.String("ip", ":8081", "IP addr of server")
	avg        chan float64
	glog       = x.Log("Pinger")
	wg         sync.WaitGroup
)

func runUser() {
	var ti time.Duration
	var query = `{
		  me(_xid_: m.0f4vbz) {
			    type.object.name.en
			    film.actor.film {
				      film.performance.film {
					        type.object.name.en
				      }
			    }
		  }
		}`
	client := &http.Client{}
	for i := 0; i < *numReq; i++ {
		r, _ := http.NewRequest("POST", *serverAddr, bytes.NewBufferString(query))
		r.Header.Add("Content-Length", strconv.Itoa(len(query)))
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		t0 := time.Now()
		fmt.Println(i)
		resp, _ := client.Do(r)
		if resp.Status != "200 OK" {
			glog.WithField("Err", resp.Status).Fatalf("Error in query")
		}
		fmt.Println("user", i)
		t1 := time.Now()
		ti += t1.Sub(t0)
	}
	fmt.Println(ti.Seconds())
	avg <- ti.Seconds()
	fmt.Println("Done")
	wg.Done()
}

func main() {
	flag.Parse()
	var totTime float64
	avg = make(chan float64, *numUser)
	fmt.Println("user")

	for i := 0; i < *numUser; i++ {
		wg.Add(1)
		fmt.Println("user")
		go runUser()
		fmt.Println("qw")
	}
	wg.Wait()
	for it := range avg {
		totTime += it
	}

	fmt.Println(totTime / float64(*numUser*(*numReq)))
}
