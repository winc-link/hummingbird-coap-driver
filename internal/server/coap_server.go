/*******************************************************************************
 * Copyright 2017.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *******************************************************************************/

package server

import (
	"github.com/winc-link/hummingbird-sdk-go/service"
	"log"
	"net"

	"github.com/dustin/go-coap"
)

type CoapServer struct {
	sd *service.DriverService
}

func NewCoapService(sd *service.DriverService) *CoapServer {
	sd.GetCustomParam()
	return &CoapServer{
		sd: sd,
	}
}

func (c *CoapServer) handleA(l *net.UDPConn, a *net.UDPAddr, m *coap.Message) *coap.Message {
	log.Printf("Got message in handleB: path=%q: %#v from %v", m.Path(), m, a)
	if m.IsConfirmable() {
		res := &coap.Message{
			Type:      coap.Acknowledgement,
			Code:      coap.Content,
			MessageID: m.MessageID,
			Token:     m.Token,
			Payload:   m.Payload,
		}
		res.SetOption(coap.ContentFormat, coap.TextPlain)

		log.Printf("Transmitting from B %#v", res)
		return res
	}
	return nil
}

func (c *CoapServer) Start() {
	mux := coap.NewServeMux()
	mux.Handle("/a", coap.FuncHandler(c.handleA))
	c.sd.GetLogger().Error(coap.ListenAndServe("udp", ":5683", mux))
}
