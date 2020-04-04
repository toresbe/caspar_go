package caspar_go

import (
    "net"
    "bufio"
    "bytes"
    "strconv"
    "strings"
//    "log"
)

type AMCPResponse struct {
    Error bool
    Retcode uint
    Retdata string
}

func ParseResponse(data string) (response AMCPResponse) {
    retcode,err := strconv.ParseUint(data[0:3], 10, 16)
    if err == nil {
        response.Retcode = uint(retcode)
    }
    if (data[0:1] == "4" || data[0:1] == "5") {
        response.Error = true
    }
    response.Retdata = strings.Join(strings.Split(string(data), "\n")[1:], "\n")
    return
}

func AMCPReceiver(c net.Conn, return_ch chan AMCPResponse, error_ch chan error) {
    var buffer bytes.Buffer
    var in_line_continuation_mode bool
    scanner := bufio.NewScanner(c)

    for scanner.Scan() {
        data_line := scanner.Text()
        if in_line_continuation_mode {
            if data_line == "" {
                in_line_continuation_mode = false
                return_ch <- ParseResponse(buffer.String())
                buffer.Reset()
            } else {
                buffer.WriteString("\n" + data_line)
            }
        } else {
            response := ParseResponse(data_line)
            if (response.Retcode == 200 || response.Retcode == 201) {
                buffer.WriteString(data_line)
                in_line_continuation_mode = true
            } else {
                return_ch <- response
            }
        }
    }
    error_ch <- scanner.Err()
    close(return_ch)
    close(error_ch)
    return
}
