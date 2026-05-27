// Copyright (c) quickfixengine.org  All rights reserved.
//
// This file may be distributed under the terms of the quickfixengine.org
// license as defined by quickfixengine.org and appearing in the file
// LICENSE included in the packaging of this file.
//
// This file is provided AS IS with NO WARRANTY OF ANY KIND, INCLUDING
// THE WARRANTY OF DESIGN, MERCHANTABILITY AND FITNESS FOR A
// PARTICULAR PURPOSE.
//
// See http://www.quickfixengine.org/LICENSE for licensing information.
//
// Contact ask@quickfixengine.org if any conditions of this licensing
// are not clear to you.

package quickfix

import (
	"net"
	"time"
)

// writeLoop closes done when it returns, signaling that the connection is
// closed so a blocked sender can escape.
func writeLoop(connection net.Conn, messageOut chan []byte, timeout time.Duration, done chan<- struct{}, log Log) {
	defer close(done)

	for msg := range messageOut {
		if timeout > 0 {
			if err := connection.SetWriteDeadline(time.Now().Add(timeout)); err != nil {
				log.OnEvent(err.Error())
				return
			}
		}
		if _, err := connection.Write(msg); err != nil {
			log.OnEvent(err.Error())
			return
		}
	}
}

func readLoop(parser *parser, msgIn chan fixIn, log Log) {
	defer close(msgIn)

	for {
		msg, err := parser.ReadMessage()
		if err != nil {
			log.OnEvent(err.Error())
			return
		}
		msgIn <- fixIn{msg, parser.lastRead}
	}
}
