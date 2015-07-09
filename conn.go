package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"net/http"
    "database/sql"
    "fmt"
    "log"
)

const (
	conn_recruiter = 0
	conn_candidate = 1
)

type connection struct {
	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send      chan *Message
	conn_type int
	auth      bool
	//in recruiter connection it will be id recuiter sent
	// in client connection - the id in which client is interested
	rid int
	did string //device id
	// The hub.
	w *wsHandler
}

func (c *connection) verifyRecruiter(pkt *Auth_pkt) (int, bool) {
	var data int
    gl := c.w
	err := gl.db.c.QueryRow("select id from recuiterinfo where email = ? and password = ?", pkt.Username, pkt.Password).Scan(&data)
	switch {
	case err == sql.ErrNoRows:
		fmt.Printf("No user with that ID.")
		return data, false
	case err != nil:
		log.Fatal(err)
	default:
        gl.db.c.QueryRow("update recuiterinfo set online='YES' where id = ?", data)
		fmt.Printf("Id is %d\n", data)
	}
	return data, true
}

func (c *connection) getRid(pkt *Auth_pkt) (int, bool) {
	var data int
    gl := c.w
	err := gl.db.c.QueryRow("select recruiter_id from job where id = ?", pkt.Rid).Scan(&data)
	switch {
	case err == sql.ErrNoRows:
		fmt.Printf("No recruiter found with that job ID.")
		return data, false
	case err != nil:
		log.Fatal(err)
	default:
		fmt.Printf("recruiter Id is %d for job\n", data, pkt.Rid)
	}
	return data, true
}

func (c *connection) setOffline() {
    if c.conn_type == conn_recruiter {
        gl := c.w
        gl.db.c.QueryRow("update recuiterinfo set online='NO' where id = ?", c.rid)
    }
}

func (c *connection) handlePkt(m []byte) {
    gl := c.w
    fmt.Println("got data")
	if c.auth == false {
		pkt := &Auth_pkt{}
		err := json.Unmarshal(m, pkt)
		if err != nil {
			fmt.Println("error in unmarshalling")
		}
		if pkt.Device_id != "" {
            fmt.Println("received auth pkt for client")
			c.conn_type = conn_candidate
		    //rid, ok := c.getRid(pkt)
            {
                gl.CtoConn[pkt.Device_id] = c
			    c.did = pkt.Device_id
			    c.rid = pkt.Rid
            }
		} else if pkt.Rid > 0 {
            fmt.Println("received auth pkt for recruiter")
			c.conn_type = conn_recruiter
            rid := pkt.Rid
            /*
            rid, e := c.verifyRecruiter(pkt)
			if !e {
				return
			}
            */
			gl.RtoConn[rid] = c
			c.rid = rid
		}
		c.auth = true
	} else {
		pkt := &Message{}
		json.Unmarshal(m, pkt)
		if c.conn_type == conn_recruiter {
            fmt.Println("got data from recruiter")
			// sent by recruiter
			cc, ok := gl.CtoConn[pkt.Clientid]
			if ok && cc.rid == c.rid {
				pkt.Rid = c.rid
				cc.send <- pkt
				pkt.Clientid = ""
			}
		} else {
			// sent by client
            fmt.Println("got data from client")
			rc, ok := gl.RtoConn[c.rid]
			if ok {
                fmt.Println("recruiter connection found")
				pkt.Clientid = c.did
				rc.send <- pkt
			} else {
                fmt.Println(gl.RtoConn)
            }
		}
	}
    c.w = gl
}

func (c *connection) reader() {
    //defer c.setOffline()
    fmt.Println("in reader")
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			break
		}
		c.handlePkt(message)
		//c.h.broadcast <- message
	}
	c.ws.Close()
}

func (c *connection) writer() {
    fmt.Println("in writer")
	for m := range c.send {
		message, err := json.Marshal(m)
		if err != nil {
			fmt.Println("Error in marshalling")
			continue
		}
		err = c.ws.WriteMessage(websocket.TextMessage, message)
		if err != nil {
            fmt.Println("error in sending data")
			break
		}
	}
	c.ws.Close()
}

func MyHandler(w http.ResponseWriter, r *http.Request) {
    //fmt.Println("Received new connection")
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if err != nil {
        fmt.Fprintf(w, "%s", err.Error())
		return
	}
	c := &connection{send: make(chan *Message, 256), ws: ws, w: wsh}
	go c.writer()
	c.reader()
}
