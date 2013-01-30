// 
// go-sonos
// ========
// 
// Copyright (c) 2012, Ian T. Richards <ianr@panix.com>
// All rights reserved.
// 
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
// 
//   * Redistributions of source code must retain the above copyright notice,
//     this list of conditions and the following disclaimer.
//   * Redistributions in binary form must reproduce the above copyright
//     notice, this list of conditions and the following disclaimer in the
//     documentation and/or other materials provided with the distribution.
// 
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED
// TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR
// PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
// LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
// NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
// 

// csweb := (c)ontrol (s)onos from a (web) browser

package main

import (
	"encoding/json"
	"fmt"
	"github.com/ianr0bkny/go-sonos"
	"github.com/ianr0bkny/go-sonos/config"
	_ "github.com/ianr0bkny/go-sonos/model"
	"log"
	"net/http"
	"strings"
)

/*
func getCurrentQueue(sonos *sonos.Sonos, writer http.ResponseWriter, request *http.Request) {
	if result, err := sonos.GetQueueContents(); nil != err {
		log.Fatal(err)
	} else {
		encoder := json.NewEncoder(writer)
		encoder.Encode(model.ObjectMessageStream(result))
	}
}
*/

const (
	CSWEB_CONFIG        = "/home/ianr/.go-sonos"
	CSWEB_DEVICE        = "kitchen"
	CSWEB_DISCOVER_PORT = "13104"
	CSWEB_EVENTING_PORT = "13105"
	CSWEB_NETWORK       = "eth0"
	CSWEB_HTTP_PORT     = 8080
)

func initSonos(config *config.Config) *sonos.Sonos {
	var s *sonos.Sonos
	if dev := config.Lookup(CSWEB_DEVICE); nil != dev {
		reactor := sonos.MakeReactor(CSWEB_NETWORK, CSWEB_EVENTING_PORT)
		s = sonos.Connect(dev, reactor, sonos.SVC_CONTENT_DIRECTORY|sonos.SVC_AV_TRANSPORT)
	} else {
		log.Fatal("Could not create Sonos instance")
	}
	return s
}

func replyOk(w http.ResponseWriter) {
	encoder := json.NewEncoder(w)
	encoder.Encode(true)
}

func handleControl(s *sonos.Sonos, w http.ResponseWriter, r *http.Request) {
	f := r.FormValue("func")
	switch f {
	case "play": s.Play(0, "1")
	case "stop": s.Stop(0)
	case "seek-start":
	}
	replyOk(w)
}

func setupHttp(s *sonos.Sonos) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, strings.Join([]string{"res", r.RequestURI}, "/"))
	})

	http.HandleFunc("/control", func(w http.ResponseWriter, r *http.Request) {
		handleControl(s, w, r)
	})
}

func main() {
	log.SetFlags(log.Ltime | log.Lshortfile)
	config := config.MakeConfig(CSWEB_CONFIG)
	config.Init()

	s := initSonos(config)
	if nil != s {
		setupHttp(s)
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", CSWEB_HTTP_PORT), nil))
}