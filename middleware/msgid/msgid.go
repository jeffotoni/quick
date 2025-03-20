// Package msgid provides middleware for automatically generating and assigning
// a unique Message ID (MsgID) to incoming HTTP requests.
//
// This middleware ensures that every request has a unique identifier, which is useful
// for request tracking, debugging, logging, and tracing in distributed systems.
//
// Features:
// - Automatically generates a MsgID if the request does not already have one.
// - Allows customization of the MsgID format using a user-defined algorithm.
// - Adds the MsgID to both the request and response headers for tracking.
// - Supports configuring the MsgID name, range, and generation strategy.
package msgid

import (
	"crypto/rand"
	"log"
	"math/big"
	"net/http"
	"strconv"
)

// Default values for the generated message ID (MsgID) range
const (
	DefaultStartConfig = 900000000 // Default starting value for the MsgID range
	DefaultEndConfig   = 100000000 // Default ending value for the MsgID range
	KeyMsgID           = "Msgid"   // Header key name for the MsgID
)

// Config struct defines the settings for the MsgID middleware
type Config struct {
	Start int           // Start range for generating random MsgIDs
	End   int           // End range for generating random MsgIDs
	Name  string        // Name of the header where the MsgID will be stored
	Algo  func() string // Custom algorithm function to generate the MsgID (optional)
}

// Default configuration for the middleware
var (
	ConfigDefault = Config{
		Name:  KeyMsgID,           // Use "Msgid" as the default header key
		Start: DefaultStartConfig, // Set default starting value
		End:   DefaultEndConfig,   // Set default ending value
	}
)

// New creates a middleware that assigns a unique MsgID to requests.
//
// If a request does not contain a MsgID header, this middleware generates one based on the specified range.
//
// Parameters:
//   - config ...Config (optional): Custom configuration. If not provided, DefaultConfig is used.
//
// Returns:
//   - A middleware function that wraps an http.Handler, ensuring each request has a MsgID in the header.
func New(config ...Config) func(http.Handler) http.Handler {
	cfd := ConfigDefault // Use default configuration
	if len(config) > 0 {
		cfd = config[0] // Override with custom configuration if provided
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if the request already has a MsgID header
			msgId := r.Header.Get(cfd.Name)

			// If MsgID is not present, generate one
			if len(msgId) == 0 {
				if cfd.Algo == nil {
					// Ensure default start and end values are set
					if cfd.Start == 0 {
						cfd.Start = DefaultStartConfig
					}
					if cfd.End == 0 {
						cfd.End = DefaultEndConfig
					}
					// Generate a new MsgID using the default algorithm
					algo := AlgoDefault(cfd.Start, cfd.End)
					// Set the generated MsgID in request and response headers
					r.Header.Set(cfd.Name, algo)
					w.Header().Set(cfd.Name, algo)
				} else {
					// If a custom algorithm is provided, use it to generate the MsgID
					algo := cfd.Algo()
					r.Header.Set(cfd.Name, algo)
					w.Header().Set(cfd.Name, algo)
				}
				// Pass the request to the next handler in the chain
				next.ServeHTTP(w, r)
			}
		})
	}
}

// AlgoDefault generates a random MsgID within the specified range
func AlgoDefault(Start, End int) string {
	max := big.NewInt(int64(End))              // Define the upper limit for the random number
	randInt, err := rand.Int(rand.Reader, max) // Generate a secure random number
	if err != nil {
		log.Printf("error: %v", err)
		return "" // Return an empty string if an error occurs
	}
	// Return the generated MsgID as a string
	return strconv.Itoa(Start + int(randInt.Int64()))
}

// func AlgoDefault(Start, End int) string {
// 	rand.Seed(time.Now().UnixNano())
// 	return strconv.Itoa(Start + int(rand.Intn(End)))
// }
