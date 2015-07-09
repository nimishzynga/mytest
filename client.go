package main

type Auth_pkt struct {
	Device_id string   `json:"device"`
	Username  string   `json:"user"`
	Password  string   `json:"password"`
	Rid       int      `json:"rid"`
}

type Message struct {
	M        string //message for recruiter
	Rid      int    //recruiter id
	Clientid string
}

type wsHandler struct {
    RtoConn map[int]*connection
	CtoConn map[string]*connection
	db      *DbConn
}

func NewwsHandler() *wsHandler {
    return &wsHandler{
        RtoConn :make(map[int]*connection),
        CtoConn : make(map[string]*connection),
        db : NewDbConn("jobnotification", "root", "root"),
    }
}
