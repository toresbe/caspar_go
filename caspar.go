package caspar_go

import (
    "fmt"
    "errors"
    "net"
)

var receive_ch chan AMCPResponse
var receive_error_ch chan error
var conn net.Conn

func Play(media_uri string) (err error) {
    conn.Write([]byte(fmt.Sprintf("PLAY 1-50 \"%v\"\r\n", media_uri)))
    result := <-receive_ch
    if result.Error {
        return errors.New("Received error while trying to play")
    }
    return nil
}

func Stop() (err error) {
    conn.Write([]byte("STOP 1-50\r\n"))
    result := <-receive_ch
    if result.Error {
        return errors.New("Received error while trying to stop")
    }
    return nil
}

func Open(open_conn net.Conn) {
    receive_ch = make(chan AMCPResponse)
    receive_error_ch = make(chan error)
    go AMCPReceiver(open_conn, receive_ch, receive_error_ch)
}
