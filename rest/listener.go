// Stoppable tcp listener

package rest

import (
	nErrors "errors"
	"github.com/jonas747/fortia/errorcodes"
	"github.com/jonas747/fortia/errors"
	"net"
	"time"
)

type Listener struct {
	*net.TCPListener          //Wrapped listener
	stop             chan int //Channel used only to indicate listener should shutdown
}

func Listen(addr string) (*Listener, errors.FortiaError) {
	originalListener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, errors.Wrap(err, errorcodes.ErrorCode_ErrCreatingListener, "")
	}
	tcpL, ok := originalListener.(*net.TCPListener)
	if !ok {
		// This shouldnt be possible?
		return nil, errors.New(errorcodes.ErrorCode_CastErr, "Error casting listener to tcplistener")
	}
	retval := &Listener{}
	retval.TCPListener = tcpL
	retval.stop = make(chan int)

	return retval, nil
}

var StoppedError = nErrors.New("Listener stopped")

func (sl *Listener) Accept() (net.Conn, error) {

	for {
		//Wait up to one second for a new connection
		sl.SetDeadline(time.Now().Add(time.Second))

		newConn, err := sl.TCPListener.Accept()

		//Check for the channel being closed
		select {
		case <-sl.stop:
			return nil, StoppedError
		default:
			//If the channel is still open, continue as normal
		}

		if err != nil {
			netErr, ok := err.(net.Error)

			//If this is a timeout, then continue to wait for
			//new connections
			if ok && netErr.Timeout() && netErr.Temporary() {
				continue
			}
		}
		// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
		// connections. It's used by ListenAndServe and ListenAndServeTLS so
		// dead TCP connections (e.g. closing laptop mid-download) eventually
		// go away.
		tcpConn, ok := newConn.(*net.TCPConn)
		if !ok {
			return nil, nErrors.New("Error casting net.conn to tcpconn")
		}
		tcpConn.SetKeepAlive(true)
		tcpConn.SetKeepAlivePeriod(3 * time.Minute)

		return newConn, err
	}
}

func (sl *Listener) Stop() {
	close(sl.stop)
}
