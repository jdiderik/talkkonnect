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
 * Contributor(s):
 *
 * Suvir Kumar <suvir@talkkonnect.com>
 *
 * My Blog is at www.talkkonnect.com
 * The source code is hosted at github.com/talkkonnect
 *
 *
 */

package talkkonnect

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"runtime"
	"strconv"
	"time"
)

func (b *Talkkonnect) cmdDisplayMenu() {
	log.Println("debug: Delete Key Pressed Menu and Session Information Requested")
	b.talkkonnectMenu("\u001b[44;1m") // add blue background to banner reference https://www.lihaoyi.com/post/BuildyourownCommandLinewithANSIescapecodes.html#background-colors
}

func (b *Talkkonnect) cmdChannelUp() {
	log.Println("debug: F1 pressed Channel Up (+) Requested")
	b.ChannelUp()
}

func (b *Talkkonnect) cmdChannelDown() {
	log.Println("debug: F2 pressed Channel Down (-) Requested")
	b.ChannelDown()
}

func (b *Talkkonnect) cmdListServerChannels() {
	log.Println("debug: F7 pressed Channel List Requested")
	b.ListChannels(true)
}

func (b *Talkkonnect) cmdStartTransmitting() {
	log.Println("debug: F8 pressed TX Mode Requested (Start Transmitting)")
	log.Println("info: Start Transmitting")

	if IsPlayStream {
		IsPlayStream = false
		NowStreaming = false

		b.playIntoStream(StreamSoundFilenameAndPath, StreamSoundVolume)
	}

	if !b.IsTransmitting {
		time.Sleep(100 * time.Millisecond)
		b.TransmitStart()
	} else {
		log.Println("error: Already in Transmitting Mode")
	}
}

func (b *Talkkonnect) cmdStopTransmitting() {
	log.Println("debug: F9 pressed RX Mode Request (Stop Transmitting)")
	log.Println("info: Stop Transmitting")

	if IsPlayStream {
		IsPlayStream = false
		NowStreaming = false

		b.playIntoStream(StreamSoundFilenameAndPath, StreamSoundVolume)
	}

	if b.IsTransmitting {
		time.Sleep(100 * time.Millisecond)
		b.TransmitStop(true)
	} else {
		log.Println("info: Not Already Transmitting")
	}
}

func (b *Talkkonnect) cmdListOnlineUsers() {
	log.Println("debug: F10 pressed Online User(s) in Current Channel Requested")
	log.Println("info: F10 Online User(s) in Current Channel")

	log.Println(fmt.Sprintf("info: Channel %#v Has %d Online User(s)", b.Client.Self.Channel.Name, len(b.Client.Self.Channel.Users)))
	b.ListUsers()
}

func (b *Talkkonnect) cmdPlayback() {
	log.Println("debug: F11 pressed Start/Stop Stream Stream into Current Channel Requested")
	log.Println("info: Stream into Current Channel")

	if b.IsTransmitting {
		log.Println("alert: talkkonnect was already transmitting will now stop transmitting and start the stream")
		b.TransmitStop(false)
	}

	IsPlayStream = !IsPlayStream
	NowStreaming = IsPlayStream

	if IsPlayStream {
		b.SendMessage(fmt.Sprintf("%s Streaming", b.Username), false)
	}

	go b.playIntoStream(StreamSoundFilenameAndPath, StreamSoundVolume)

}

func (b *Talkkonnect) cmdQuitTalkkonnect() {
	log.Println("debug: Ctrl-C Terminate Program Requested")
	duration := time.Since(StartTime)
	log.Printf("info: Talkkonnect Now Running For %v \n", secondsToHuman(int(duration.Seconds())))

	ServerHop = true
	b.CleanUp()
}


func (b *Talkkonnect) cmdClearScreen() {
	reset()
	log.Println("debug: Ctrl-L Pressed Cleared Screen")
}

func (b *Talkkonnect) cmdPingServers() {
	log.Println("debug: Ctrl-O Pressed")
	log.Println("info: Ping Servers")
	b.pingServers()
}

// func (b *Talkkonnect) cmdPanicSimulation() {
// 	if !(IsConnected) {
// 		return
// 	}
// 	log.Println("debug: Ctrl-P Pressed")
// 	log.Println("info: Panic Button Start/Stop Simulation Requested")

// 	if PEnabled {

// 		if b.IsTransmitting {
// 			b.TransmitStop(false)
// 		} else {
// 			b.IsTransmitting = true
// 			b.SendMessage(PMessage, PRecursive)

// 		}

// 		if PSendIdent {
// 			b.SendMessage(fmt.Sprintf("My Username is %s and Ident is %s", b.Username, b.Ident), PRecursive)
// 		}

// 		IsPlayStream = false
// 		b.IsTransmitting = false
		
// 	}
// }
