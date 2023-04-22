package msgid

import (
	"crypto/rand"
	"math/big"
	"net/http"
	"strconv"
)

const (
	KeyMsgID = "Msgid"
)

type Config struct {
	UUID  bool
	Name  string
	Start int
	End   int
	Algo  func() string
}

var (
	ConfigDefault = Config{
		UUID:  false,
		Name:  KeyMsgID,
		Start: 900000000,
		End:   100000000,
	}
)

func New(config ...Config) func(http.Handler) http.Handler {
	cfd := ConfigDefault
	if len(config) > 0 {
		cfd = config[0]
	}

	return func(next http.Handler) http.Handler {
		// return default MsgID
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			msgId := r.Header.Get(cfd.Name)
			if len(msgId) == 0 {
				if cfd.Algo == nil {
					algo := AlgoDefault(cfd.Start, cfd.End)
					r.Header.Set(cfd.Name, algo)
					w.Header().Set(cfd.Name, algo)
				} else {
					algo := cfd.Algo()
					r.Header.Set(cfd.Name, algo)
					w.Header().Set(cfd.Name, algo)
				}
				next.ServeHTTP(w, r)
			}
		})
	}
}

func AlgoDefault(Start, End int) string {
	max := big.NewInt(int64(End))
	randInt, err := rand.Int(rand.Reader, max)
	if err != nil {
		panic(err)
	}
	return strconv.Itoa(Start + int(randInt.Int64()))
}

// func AlgoDefault(Start, End int) string {
// 	rand.Seed(time.Now().UnixNano())
// 	return strconv.Itoa(Start + int(rand.Intn(End)))
// }
