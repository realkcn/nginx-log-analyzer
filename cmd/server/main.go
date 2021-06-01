package main

import (
	"fmt"
	"gopkg.in/mcuadros/go-syslog.v2"
)

func main() {
	channel := make(syslog.LogPartsChannel)
	handler := syslog.NewChannelHandler(channel)

	server := syslog.NewServer()
	//server.SetFormat(syslog.RFC5424)
	server.SetHandler(handler)
	server.ListenUDP("127.0.0.1:5501")
	server.SetFormat(NewNginxLogFormatter("$remote_addr - $remote_user [$time_local] \"$request\" $status $body_bytes_sent \"$http_referer\" \"$http_user_agent\" \"$http_x_forwarded_for\""))
	//server.SetFormat(NewNginxLogFormatter("$remote_addr "))
	server.Boot()

	go func(channel syslog.LogPartsChannel) {
		for logParts := range channel {
			fmt.Println(logParts)
		}
	}(channel)

	server.Wait()
}
