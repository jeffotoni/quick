package msgid

import (
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type MsgID struct {
	Config Config
}

type Config struct {
	Name string
}

var ConfigDefault = Config{
	Name: "Msgid",
}

func (m *MsgID) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	msgId := r.Header.Get(m.Config.Name)
	if len(msgId) == 0 {
		r.Header.Set(m.Config.Name, RandAlgo1())
		w.Header().Set(m.Config.Name, RandAlgo1())
		next.ServeHTTP(w, r)
	}
}

func (m *MsgID) New(config ...Config) *MsgID {
	cfd := ConfigDefault
	if len(config) > 0 {
		cfd = config[0]
	}
	return &MsgID{
		Config: cfd,
	}

	// return func(next http.Handler) http.Handler {
	// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 		msgId := r.Header.Get(cfd.Name)
	// 		if len(msgId) == 0 {
	// 			r.Header.Set(cfd.Name, RandAlgo1())
	// 			w.Header().Set(cfd.Name, RandAlgo1())
	// 			next.ServeHTTP(w, r)
	// 		}
	// 	})
	// }
}

func RandAlgo1() string {
	rand.Seed(time.Now().UnixNano())
	return strconv.Itoa(int(rand.Intn(100000)))
}
