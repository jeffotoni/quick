package msgid

import (
	"crypto/rand"
	"log"
	"math/big"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

const (
	uuidVersion = iota
	UUID_VERSION_1
	UUID_VERSION_2
	UUID_VERSION_3
	UUID_VERSION_4
)

const (
	KeyMsgID   = "Msgid"
	KeyMsgUUID = "MsgUUID"
)

type Config struct {
	UUID       bool
	Name       string
	Start      int
	End        int
	Algo       func() string
	ConfigUUID ConfigUUID
}

type ConfigUUID struct {
	Version   int // this define the version you desire of your UUID (you can choose between 1 and 4)
	Name      string
	KeyString string // this is a key string to parse value of its bytes to a uuid value
}

func NewUUID(uuidCnfg ...ConfigUUID) ConfigUUID {
	config := ConfigUUID{
		Version:   UUID_VERSION_4,
		Name:      KeyMsgUUID,
		KeyString: "quick",
	}

	if len(uuidCnfg) > 0 {
		config = uuidCnfg[0]
	}

	return config
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
		if cfd.UUID {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				msgId := r.Header.Get(cfd.Name)
				if len(msgId) == 0 {
					uuid := generateDefaultUUID(NewUUID(cfd.ConfigUUID))
					r.Header.Set(cfd.ConfigUUID.Name, uuid)
					w.Header().Set(cfd.ConfigUUID.Name, uuid)
					next.ServeHTTP(w, r)
				}
			})
		}

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

func generateDefaultUUID(uidCfg ConfigUUID) string {
	switch uidCfg.Version {
	case UUID_VERSION_1:
		u, err := uuid.NewUUID()
		if err != nil {
			log.Printf("error to generate UUID: %v", err)
		}
		if len(uidCfg.KeyString) != 0 {
			u = uuid.MustParse(uidCfg.KeyString)
		}
		return u.String()
	case UUID_VERSION_2:
		u := uuid.New()
		if len(uidCfg.KeyString) != 0 {
			u = uuid.MustParse(uidCfg.KeyString)
		}
		return u.String()
	case UUID_VERSION_3:
		u := uuid.NewMD5(uuid.New(), []byte(uidCfg.KeyString))
		if len(uidCfg.KeyString) != 0 {
			u = uuid.MustParse(uidCfg.KeyString)
		}
		return u.String()
	default: // making uuid version 4 as default
		u, err := uuid.NewRandom()
		if err != nil {
			log.Printf("error to generate UUID: %v", err)
		}
		if len(uidCfg.KeyString) != 0 {
			u = uuid.MustParse(uidCfg.KeyString)
		}
		return u.String()
	}
}

// func AlgoDefault(Start, End int) string {
// 	rand.Seed(time.Now().UnixNano())
// 	return strconv.Itoa(Start + int(rand.Intn(End)))
// }

func (c Config) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
			// w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
