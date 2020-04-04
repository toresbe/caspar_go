package caspar_go

import (
    "testing"
    "net"
    "strings"
    "bufio"
//    "log"
)

func FakeServer(conn net.Conn) {
    scanner := bufio.NewScanner(conn)
    for scanner.Scan() {
        cmd := strings.Split(scanner.Text(), " ")
        switch {
        case cmd[0] == "PLAY":
            if cmd[2] == `"THE_ONLY_VALID_FILENAME"` {
                conn.Write([]byte("202 PLAY OK\r\n"))
            } else {
                conn.Write([]byte("404 PLAY FAILED\r\n"))
            }
        default:
            conn.Write([]byte("400 ERROR\r\n"))
        }
    }
}

func TestPlayWithExtantFileReturnsNilError(t *testing.T) {
    fake_client, fake_server := net.Pipe()
    go FakeServer(fake_server)
    Open(fake_client)
    err := Play("THE_ONLY_VALID_FILENAME")
    if err != nil {
        t.Error("Expected nil, got err:", err)
    }
}

func TestPlayWithNonextantFileReturnsError(t *testing.T) {
    fake_client, fake_server := net.Pipe()
    go FakeServer(fake_server)
    Open(fake_client)
    err := Play("not_valid_filename")
    if err == nil {
        t.Error("Expected error, got nil!")
    }
}
