/*
 * talkkonnect headless mumble client/gateway with lcd screen and channel control
 * Copyright (C) 2018-2019, Suvir Kumar <suvir@talkkonnect.com>
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 * Software distributed under the License is distributed on an "AS IS" basis,
 * WITHOUT WARRANTY OF ANY KIND, either express or implied. See the License
 * for the specific language governing rights and limitations under the
 * License.
 *
 * talkkonnect is the based on talkiepi and barnard by Daniel Chote and Tim Cooper
 *
 * The Initial Developer of the Original Code is
 * Suvir Kumar <suvir@talkkonnect.com>
 * Portions created by the Initial Developer are Copyright (C) Suvir Kumar. All Rights Reserved.
 *
 * talKKonnectContributor(s):
 *
 * Suvir Kumar <suvir@talkkonnect.com>
 *
 * My Blog is at www.talkkonnect.com
 * The source code is hosted at github.com/talkkonnect
 *
 * MQTT License Details Copyright (c) 2013 IBM Corp.
 *
 * This project is dual licensed under the Eclipse Public License 1.0 and the
 * Eclipse Distribution License 1.0 as described in the epl-v10 and edl-v10 files.
 * The EDL is copied below in order to pass the pkg.go.dev license check (https://pkg.go.dev/license-policy).
 * Eclipse Distribution License - v 1.0
 * Copyright (c) 2007, Eclipse Foundation, Inc. and its licensors.
 * All rights reserved.
 * Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:
 *
 * Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
 * Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
 * Neither the name of the Eclipse Foundation, Inc. nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 * All rights reserved. This program and the accompanying materials are made available under the terms of the Eclipse Public License v1.0
 * which accompanies this distribution, and is available at http://www.eclipse.org/legal/epl-v10.html
 *
 */

package talkkonnect

import (
	"crypto/tls"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func mqtttestpub() {

	if MQTTAction != "pub" {
		log.Println("error: Invalid setting for -action, must be pub for test pub")
		return
	}

	if MQTTTopic == "" {
		log.Println("error: Invalid setting for -topic, must not be empty")
		return
	}

	opts := MQTT.NewClientOptions()
	opts.AddBroker(MQTTBroker)
	opts.SetClientID(MQTTId)
	opts.SetUsername(MQTTUser)
	opts.SetPassword(MQTTPassword)
	opts.SetCleanSession(MQTTCleansess)

	if MQTTStore != ":memory:" {
		opts.SetStore(MQTT.NewFileStore(MQTTStore))
	}

	if MQTTAction == "pub" {

		log.Printf("info: action      : %s\n", MQTTAction)
		log.Printf("info: broker      : %s\n", MQTTBroker)
		log.Printf("info: clientid    : %s\n", MQTTId)
		log.Printf("info: user        : %s\n", MQTTUser)
		log.Printf("info: mqttpassword: %s\n", MQTTPassword)
		log.Printf("info: topic       : %s\n", MQTTTopic)
		log.Printf("info: message     : %s\n", MQTTPayload)
		log.Printf("info: qos         : %d\n", MQTTQos)
		log.Printf("info: cleansess   : %v\n", MQTTCleansess)
		log.Printf("info: num         : %d\n", MQTTNum)
		log.Printf("info: store       : %s\n", MQTTStore)

		client := MQTT.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}

		log.Println("info: Test MQTT Publisher Started")
		for i := 0; i < MQTTNum; i++ {
			log.Println("info: Publishing MQTT Message")
			token := client.Publish(MQTTTopic, byte(MQTTQos), false, MQTTPayload)
			token.Wait()
		}

		client.Disconnect(250)

		log.Println("info: Test MQTT Publisher Disconnected")
	}
}

func (b *Talkkonnect) mqttsubscribe() {

	log.Printf("info: MQTT Subscription Information")
	log.Printf("info: MQTT Broker      : %s\n", MQTTBroker)
	log.Printf("debug: MQTT clientid    : %s\n", MQTTId)
	log.Printf("debug: MQTT user        : %s\n", MQTTUser)
	log.Printf("debug: MQTT password    : %s\n", MQTTPassword)
	log.Printf("info: Subscribed topic : %s\n", MQTTTopic)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	connOpts := MQTT.NewClientOptions().AddBroker(MQTTBroker).SetClientID(MQTTId).SetCleanSession(true)
	if MQTTUser != "" {
		connOpts.SetUsername(MQTTUser)
		if MQTTPassword != "" {
			connOpts.SetPassword(MQTTPassword)
		}
	}
	tlsConfig := &tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert}
	connOpts.SetTLSConfig(tlsConfig)

	connOpts.OnConnect = func(c MQTT.Client) {
		if token := c.Subscribe(MQTTTopic, byte(MQTTQos), b.onMessageReceived); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
	}

	client := MQTT.NewClient(connOpts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	} else {
		log.Printf("info: Connected to     : %s\n", MQTTBroker)
	}

	<-c
}

func (b *Talkkonnect) onMessageReceived(client MQTT.Client, message MQTT.Message) {
	log.Printf("info: Received MQTT message on topic: %s Payload: %s\n", message.Topic(), message.Payload())

	switch string(message.Payload()) {
	case "DisplayMenu":
		log.Println("info: MQTT Display Menu Request Processed Successfully")
		b.cmdDisplayMenu()
	case "ChannelUp":
		log.Println("info: MQTT Channel Up Request Processed Successfully\n")
		b.cmdChannelUp()
	case "ChannelDown":
		log.Println("info: MQTT Channel Down Request Processed Successfully\n")
		b.cmdChannelDown()
	case "Mute-Toggle":
		log.Println("info: MQTT Mute/UnMute Speaker Request Processed Successfully\n")
		b.cmdMuteUnmute("toggle")
	case "Mute":
		log.Println("info: MQTT Mute/UnMute Speaker Request Processed Successfully\n")
		b.cmdMuteUnmute("mute")
	case "Unmute":
		log.Println("info: MQTT Mute/UnMute Speaker Request Processed Successfully\n")
		b.cmdMuteUnmute("unmute")
	case "ListChannels":
		log.Println("info: MQTT List Server Channels Request Processed Successfully\n")
		b.cmdListServerChannels()
	case "StartTransmitting":
		log.Println("info: MQTT Start Transmitting Request Processed Successfully\n")
		b.cmdStartTransmitting()
	case "StopTransmitting":
		log.Println("info: MQTT Stop Transmitting Request Processed Successfully\n")
		b.cmdStopTransmitting()
	case "ListOnlineUsers":
		log.Println("info: MQTT List Online Users Request Processed Successfully\n")
		b.cmdListOnlineUsers()
	case "Stream-Toggle":
		log.Println("info: MQTT Play/Stop Stream Request Processed Successfully\n")
		b.cmdPlayback()
	case "ClearScreen":
		log.Println("info: MQTT Clear Screen Processed Successfully\n")
		b.cmdClearScreen()
	case "PingServers":
		log.Println("info: MQTT Ping Servers Processed Successfully\n")
		b.cmdPingServers()
	case "ScanChannels":
		log.Println("info: MQTT Request Scan Processed Successfully\n")
		b.cmdScanChannels()
	case "ShowUptime":
		log.Println("info: MQTT Request Current Version Successfully\n")
		b.cmdShowUptime()

	// todo add other automation control for buttons, relays and leds here as needed in the future
	default:
		log.Printf("error: Undefined Command Received MQTT message on topic: %s Payload: %s\n", message.Topic(), message.Payload())
	}
	return
}
