package msguuid

import (
	"log"
	"net/http"

	"github.com/google/uuid"
)

const (
	UUID_VERSION_1 = iota + 1
	UUID_VERSION_2
	UUID_VERSION_3
	UUID_VERSION_4
)

const (
	KeyMsgUUID = "MsgUUID"
)

type Config struct {
	Version   int // this define the version you desire of your UUID (you can choose between 1 and 4)
	Name      string
	KeyString string // this is a key string to parse value of its bytes to a uuid value
}

var DefaultConfig = Config{
	Version: UUID_VERSION_4,
	Name:    KeyMsgUUID,
	// KeyString will be come as default "" to generate a new uuid, but always use something with this length
	// 00000000-0000-0000-0000-000000000000 to work properly
	KeyString: "",
}

func New(config ...Config) func(http.Handler) http.Handler {
	cfd := DefaultConfig
	if len(config) > 0 {
		cfd = config[0]
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			msgUuId := r.Header.Get(cfd.Name)
			if len(msgUuId) == 0 {
				uuid := generateDefaultUUID(cfd)
				r.Header.Set(cfd.Name, uuid)
				w.Header().Set(cfd.Name, uuid)
				next.ServeHTTP(w, r)
			}
		})
	}
}

func generateDefaultUUID(uidCfg Config) string {
	switch uidCfg.Version {
	case UUID_VERSION_1:
		if len(uidCfg.KeyString) != 0 {
			u, err := uuid.Parse(uidCfg.KeyString)
			if err != nil {
				log.Printf("error to generate UUID: %v", err)
			}
			return u.String()
		} else {
			u, err := uuid.NewUUID()
			if err != nil {
				log.Printf("error to generate UUID: %v", err)
			}
			return u.String()
		}
	case UUID_VERSION_2:
		if len(uidCfg.KeyString) != 0 {
			u, err := uuid.Parse(uidCfg.KeyString)
			if err != nil {
				log.Printf("error to generate UUID: %v", err)
			}
			return u.String()
		} else {
			u := uuid.New()
			return u.String()
		}
	case UUID_VERSION_3:
		if len(uidCfg.KeyString) != 0 {
			u, err := uuid.Parse(uidCfg.KeyString)
			if err != nil {
				log.Printf("error to generate UUID: %v", err)
			}
			return u.String()
		} else {
			u := uuid.NewMD5(uuid.New(), []byte(uidCfg.KeyString))
			return u.String()
		}
	default: // making uuid version 4 as default

		if len(uidCfg.KeyString) != 0 {
			u, err := uuid.Parse(uidCfg.KeyString)
			if err != nil {
				log.Printf("error to generate UUID: %v", err)
			}
			return u.String()
		} else {
			u, err := uuid.NewRandom()
			if err != nil {
				log.Printf("error to generate UUID: %v", err)
			}
			return u.String()
		}
	}
}
