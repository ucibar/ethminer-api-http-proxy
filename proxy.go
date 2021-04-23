package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

type Request struct {
	ID      uint64           `json:"id"`
	Jsonprc string           `json:"jsonrpc"`
	Method  string           `json:"method"`
	Params  *json.RawMessage `json:"params,omitempty"`
}

type Response struct {
	ID     uint64           `json:"id"`
	Result *json.RawMessage `json:"result,omitempty"`
	Error  *json.RawMessage `json:"error,omitempty"`
}

type Call struct {
	Req  Request
	Res  Response
	Done chan bool
}

func NewCall(req Request) *Call {
	done := make(chan bool)
	return &Call{Req: req, Done: done}
}

type Proxy struct {
	mutex   sync.Mutex
	conn    *ETHMinerConn
	pending map[uint64]*Call
	counter uint64

	reader *bufio.Reader
}

func NewProxy() *Proxy {
	return &Proxy{
		pending: make(map[uint64]*Call, 1),
		counter: 1,
	}
}

func (p *Proxy) HTTPHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		_, err := fmt.Fprintf(w, "Method Not Allowed")
		if err != nil {
			log.Println(err)
		}
		return
	}

	req, _ := ioutil.ReadAll(r.Body)
	res, err := p.Request(req)
	if err != nil {
		w.WriteHeader(500)
		_, err = fmt.Fprintf(w, err.Error())
		if err != nil {
			log.Println(err)
		}
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(200)
	_, err = w.Write(res)
	if err != nil {
		log.Println(err)
	}
}

func (p *Proxy) Connect(address string) error {
	conn, err := ETHMinerDial(address)
	if err != nil {
		return err
	}
	p.conn = conn

	p.reader = bufio.NewReader(p.conn)
	go p.read()

	return nil
}

func (p *Proxy) Close() error {
	return p.conn.Close()
}

func (p *Proxy) read() {
	for {
		r, err := p.reader.ReadBytes('\n')
		if err != nil {
			log.Println(err)
			continue
		}

		res := Response{}
		err = json.Unmarshal(r, &res)
		if err != nil {
			log.Println(err)
			continue
		}

		p.mutex.Lock()
		call := p.pending[res.ID]
		delete(p.pending, res.ID)
		p.mutex.Unlock()

		if call == nil {
			log.Println("no pending request found")
			continue
		}

		call.Res = res
		call.Done <- true
	}
}

func (p *Proxy) Request(r []byte) ([]byte, error) {
	req := Request{}
	err := json.Unmarshal(r, &req)
	if err != nil {
		return nil, err
	}

	call := NewCall(req)

	p.mutex.Lock()
	req.ID = p.counter
	p.pending[req.ID] = call

	r, err = json.Marshal(req)
	if err != nil {
		delete(p.pending, req.ID)
		p.mutex.Unlock()
		return nil, err
	}

	_, err = p.conn.Write(r)
	if err != nil {
		delete(p.pending, req.ID)
		p.mutex.Unlock()
		return nil, err
	}

	p.counter++
	p.mutex.Unlock()

	select {
	case <-call.Done:
	case <-time.After(1 * time.Second):
		return nil, errors.New("request timeout")
	}

	res, err := json.Marshal(call.Res)
	return res, nil
}
