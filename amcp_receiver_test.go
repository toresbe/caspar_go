package caspar_go

import (
    "testing"
    "net"
//    "log"
)

func TestAMCPHandles400(t *testing.T) {
    fake_client, fake_server := net.Pipe()
    data_channel := make(chan AMCPResponse)
    error_channel := make(chan error)

    go AMCPReceiver(fake_client, data_channel, error_channel)

    fake_server.Write([]byte("400 PLAY FAILED\r\n"))

    response := <-data_channel
    if response.Retcode != 400 {
        t.Error("Wrong response code returned!")
    }
    if response.Error == false {
        t.Error("Falsely identified as success!")
    }
}

func TestAMCPHandles200(t *testing.T) {
    fake_client, fake_server := net.Pipe()
    data_channel := make(chan AMCPResponse)
    error_channel := make(chan error)

    go AMCPReceiver(fake_client, data_channel, error_channel)

    fake_server.Write([]byte("200 INFO OK\r\n1 720p5000 PLAYING\r\n\r\n"))

    response := <-data_channel
    if response.Retcode != 200 {
        t.Errorf("Wrong response code returned!")
    }
    if response.Retdata != "1 720p5000 PLAYING" {
        t.Errorf("Wrong return data returned!")
    }
    if response.Error != false {
        t.Errorf("Falsely identified as error!")
    }
}

func TestAMCPHandles201(t *testing.T) {
    fake_client, fake_server := net.Pipe()
    data_channel := make(chan AMCPResponse)
    error_channel := make(chan error)

    go AMCPReceiver(fake_client, data_channel, error_channel)

    fake_server.Write([]byte("201 INFO OK\r\nDummy line 1\r\nDummy line 2\r\n\r\n"))

    response := <-data_channel
    if response.Retcode != 201 {
        t.Errorf("Wrong response code returned!")
    }
    if response.Retdata != "Dummy line 1\nDummy line 2" {
        t.Errorf("Wrong return data returned!")
    }
    if response.Error != false {
        t.Errorf("Falsely identified as error!")
    }
}

func TestAMCPGivesNilErrorOnDisconnect(t *testing.T) {
    fake_client, fake_server := net.Pipe()
    data_channel := make(chan AMCPResponse)
    error_channel := make(chan error)

    go AMCPReceiver(fake_client, data_channel, error_channel)
    fake_server.Close()
    if <-error_channel != nil {
        t.Errorf("Expected nil error on socket close")
    }
}
