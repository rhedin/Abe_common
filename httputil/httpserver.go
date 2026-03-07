/*
 * Public Domain Software
 *
 * I (Matthias Ladkau) am the author of the source code in this file.
 * I have placed the source code in this file in the public domain.
 *
 * For further information see: http://creativecommons.org/publicdomain/zero/1.0/
 */

/*
Package httputil contains a HTTP/HTTPS Server which can be stopped via signals
or a Shutdown() call.
*/
package httputil

import (
	"crypto/tls"
	"errors"
	"fmt"
	abelog "github.com/rhedin/Abe_common/abelogutil"
	"github.com/rhedin/Abe_eliasdb/config"
	"github.com/rs/cors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

/*
HTTPServer data structure
*/
type HTTPServer struct {
	signalling chan os.Signal    // Channel for receiving signals
	LastError  error             // Last recorded error
	Running    bool              // Flag if the server is running
	listener   signalTCPListener // TCP listener of the server
}

/*
Shutdown sends a shutdown signal.
*/
func (hs *HTTPServer) Shutdown() {
	if hs.signalling != nil {
		hs.signalling <- syscall.SIGINT
	}
}

/*
RunHTTPServer starts a HTTP Server which can be stopped via ^C (Control-C).
It is assumed that all routes have been added prior to this call.

laddr should be the local address which should be given to net.Listen.
wgStatus is an optional wait group which will be notified once the server is listening
and once the server has shutdown.

This function will not return unless the server is shutdown.
*/
func (hs *HTTPServer) RunHTTPServer(laddr string, wgStatus *sync.WaitGroup) error {

	hs.Running = false

	// Create normal TCP listener

	originalListener, err := net.Listen("tcp", laddr)
	if err != nil {
		hs.LastError = err

		if wgStatus != nil {
			wgStatus.Done()
		}

		return err
	}

	// Wrap listener in a signal aware listener

	sl := newSignalTCPListener(originalListener, originalListener.(*net.TCPListener), wgStatus)

	return hs.runServer(sl, wgStatus)
}

/*
RunHTTPSServer starts a HTTPS Server which can be stopped via ^C (Control-C).
It is assumed that all routes have been added prior to this call.

keypath should be set to a path containing the TLS certificate and key.
certFile should be the file containing the TLS certificate.
keyFile should be the file containing the private key for the TLS connection.
laddr should be the local address which should be given to net.Listen.
wgStatus is an optional wait group which will be notified once the server is listening
and once the server has shutdown.

This function will not return unless the server is shutdown.
*/
func (hs *HTTPServer) RunHTTPSServer(keypath string, certFile string, keyFile string,
	laddr string, wgStatus *sync.WaitGroup) error {

	// Check parameters

	if keypath != "" && !strings.HasSuffix(keypath, "/") {
		keypath += "/"
	}

	// Load key pair and create a TLS config

	cert, err := tls.LoadX509KeyPair(keypath+certFile, keypath+keyFile)
	if err != nil {
		hs.LastError = err

		if wgStatus != nil {
			wgStatus.Done()
		}

		return err
	}

	hs.Running = false

	// Create normal TCP listener

	originalListener, err := net.Listen("tcp", laddr)
	if err != nil {
		hs.LastError = err

		if wgStatus != nil {
			wgStatus.Done()
		}

		return err
	}

	// Wrap the listener in a TLS listener

	tlsConfig := tls.Config{Certificates: []tls.Certificate{cert}}

	originalTLSListener := tls.NewListener(originalListener, &tlsConfig)

	// Wrap listeners in a signal aware listener

	sl := newSignalTCPListener(originalTLSListener, originalListener.(*net.TCPListener), wgStatus)

	return hs.runServer(sl, wgStatus)
}

/*
runServer starts the actual server and notifies the wait group.
*/
func (hs *HTTPServer) runServer(sl *signalTCPListener, wgStatus *sync.WaitGroup) error {

	// I wanted to make a configuration entry like this:
	//
	// 	CORSAllowedOrigins:       ["*", "https://studio.apollographql.com"],
	//
	// But Grok insisted I had to do it like this:
	//
	// 	CORSAllowedOrigins:       "*,https://studio.apollographql.com",
	//
	// He said that if I were reading the information back in from the file,
	// Go wouldn't have enough information to reconstruct an array of strings.
	//
	// Not sure that's true.  Grok suggested that I could take the current
	// configuration stuff apart, and give everything explicit types.
	//
	// Used this complicated bit of code instead, to convert a comma-separated
	// string into a slice of strings.  Aim for good, not great.  That's the
	// new policy.  Have to get done.

	var allowedOrigins []string

	// 1. Read the string value from config
	stringWithCommas := config.Str(config.CORSAllowedOrigins)

	// 2. Split on commas
	sliceOfStrings := strings.Split(stringWithCommas, ",")

	// 3. Trim whitespace and filter empty entries
	allowedOrigins = make([]string, 0, len(sliceOfStrings))
	for _, indivString := range sliceOfStrings {
		trimmedString := strings.TrimSpace(indivString)
		if trimmedString != "" {
			allowedOrigins = append(allowedOrigins, trimmedString)
		}
	}

	// 4. Fallback if everything was empty after trimming
	if len(allowedOrigins) == 0 {
		fmt.Printf("CORSAllowedOrigins configuration is wrong.\n")
		// I hate to "fall back" to a possibly inappropriate value,
		// but at least we logged that there was a problem.
		allowedOrigins = []string{"*"}
	}

	abelog.UnderPrintf("allowedOrigins = %v\n", allowedOrigins)
	if abelog.UnderEnabled {
		for index, origin := range allowedOrigins {
			fmt.Printf("     %d. -->%s<--\n", index, origin)
		}
	}
	// I wasn't sure I was understanding how Go printed slices using %v.

	c := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Authorization"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
		Debug:            config.Bool(config.CORSDebug),
	})

	// Wrap the default mux
	corsHandler := c.Handler(http.DefaultServeMux)

	// Use the http server from the standard library
	// If we had left the braces empty, it would have used
	// the default multiplexed server.  We want it to use
	// our wrapper around DefaultServerMux.
	server := http.Server{
		Handler: corsHandler,
	}

	// Attach SIGINT handler - on unix and windows this is send
	// when the user presses ^C (Control-C).

	hs.signalling = make(chan os.Signal)
	signal.Notify(hs.signalling, syscall.SIGINT)

	// Put the serve call into a wait group so we can wait until shutdown
	// completed

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		hs.Running = true
		server.Serve(sl)
	}()

	for true {
		signal := <-hs.signalling

		if signal == syscall.SIGINT {

			// Shutdown the server

			sl.Shutdown()

			// Wait until the server has shut down

			wg.Wait()

			hs.Running = false

			break
		}
	}

	if wgStatus != nil {
		wgStatus.Done()
	}

	return nil
}

/*
signalTCPListener models a TCPListener which can receive signals.
*/
type signalTCPListener struct {
	net.Listener                  // Wrapped new.Listener
	tcpListener  *net.TCPListener // TCP listener which accepts connections
	Signals      chan int         // Channel used for signalling
	wgStatus     *sync.WaitGroup  // Optional Waitgroup to be notified after start
}

/*
SigShutdown is used to signal a request for shutdown
*/
const SigShutdown = 1

/*
ErrSigShutdown indicates that a signal was received
*/
var ErrSigShutdown = errors.New("Server was shut down")

/*
newSignalTCPListener wraps a given TCPListener.
*/
func newSignalTCPListener(l net.Listener, tl *net.TCPListener, wgStatus *sync.WaitGroup) *signalTCPListener {
	return &signalTCPListener{l, tl, make(chan int), wgStatus}
}

/*
Accept waits for a new connection. This accept call will check every
second if a signal or other shutdown event was received.
*/
func (sl *signalTCPListener) Accept() (net.Conn, error) {
	for {

		// Wait up to a second for a new connection

		sl.tcpListener.SetDeadline(time.Now().Add(time.Second))
		newConn, err := sl.Listener.Accept()

		// Notify wgStatus if it was specified

		if sl.wgStatus != nil {
			sl.wgStatus.Done()
			sl.wgStatus = nil
		}

		// Check for a received signal

		select {
		case sig := <-sl.Signals:

			// Check which signal was received

			if sig == SigShutdown {
				return nil, ErrSigShutdown
			}

			panic(fmt.Sprintf("Unknown signal received: %v", sig))

		default:

			netErr, ok := err.(net.Error)

			// If we got a connection or error at this point return it

			if (err != nil && (!ok || !(netErr.Timeout() && netErr.Temporary()))) || newConn != nil {
				return newConn, err
			}
		}
	}
}

/*
Shutdown sends a shutdown signal.
*/
func (sl *signalTCPListener) Shutdown() {
	sl.Signals <- SigShutdown
	close(sl.Signals)
}
