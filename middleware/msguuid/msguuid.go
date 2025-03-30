// Package msguuid provides a middleware for assigning a UUID to incoming HTTP requests.
//
// This middleware ensures that each request contains a unique identifier in the header.
// If a request already has a UUID in the specified header, it remains unchanged.
// Otherwise, a new UUID is generated and assigned based on the configured version (1 to 4).
//
// The UUID can be retrieved from the request and response headers using the configured key.
// This is useful for tracking requests across services, logging, and debugging.
package msguuid

import (
	"log"
	"net/http"

	"github.com/jeffotoni/quick/pkg/uuid"
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

// Config defines the configuration for the MsgUUID middleware.
//
// Fields:
//   - Version (int): Specifies the UUID version to use (1 to 4).
//   - Name (string): Header key name where the UUID will be stored.
//   - KeyString (string): Custom key string for UUID parsing. Defaults to an empty string.
type Config struct {
	Version   int
	Name      string
	KeyString string
}

// DefaultConfig provides a default configuration for the MsgUUID middleware.
//
// Defaults:
//   - Version: UUID v4
//   - Name: "MsgUUID"
//   - KeyString: ""
var DefaultConfig = Config{
	Version:   UUID_VERSION_4,
	Name:      KeyMsgUUID,
	KeyString: "",
}

// New creates a middleware that assigns a UUID to requests.
//
// If a request does not contain a UUID header, this middleware generates one based on the specified version.
//
// Parameters:
//   - config ...Config (optional): Custom configuration. If not provided, DefaultConfig is used.
//
// Returns:
//   - A middleware function that wraps an http.Handler, ensuring each request has a UUID in the header.
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
			}
			next.ServeHTTP(w, r)
		})
	}
}

// generateDefaultUUID generates a UUID based on the provided configuration.
//
// Parameters:
//   - uidCfg (Config): The configuration specifying the UUID version and optional KeyString.
//
// Returns:
//   - A UUID string based on the chosen version.
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
