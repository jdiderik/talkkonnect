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
 * talkkonnect.go -> function in talkkonnect for printing banners to screen
 */

package talkkonnect

import (
	"log"
	"os"
	"strconv"
)

func talkkonnectBanner(backgroundcolor string) {
	var backgroundreset string = "\u001b[0m"
	log.Println("info: " + backgroundcolor + "┌────────────────────────────────────────────────────────────────┐" + backgroundreset)
	log.Println("info: " + backgroundcolor + "│  _        _ _    _                               _             │" + backgroundreset)
	log.Println("info: " + backgroundcolor + "│ | |_ __ _| | | _| | _____  _ __  _ __   ___  ___| |_           │" + backgroundreset)
	log.Println("info: " + backgroundcolor + "│ | __/ _` | | |/ / |/ / _ \\| '_ \\| '_ \\ / _ \\/ __|  __|         │" + backgroundreset)
	log.Println("info: " + backgroundcolor + "│ | || (_| | |   <|   < (_) | | | | | | |  __/ (__| |_           │" + backgroundreset)
	log.Println("info: " + backgroundcolor + "│  \\__\\__,_|_|_|\\_\\_|\\_\\___/|_| |_|_| |_|\\___|\\_ _|\\__|          │" + backgroundreset)
	log.Println("info: " + backgroundcolor + "├────────────────────────────────────────────────────────────────┤" + backgroundreset)
	log.Println("info: " + backgroundcolor + "│A Flexible Headless Mumble Transceiver/Gateway for RPi/PC/VM    │" + backgroundreset)
	log.Println("info: " + backgroundcolor + "├────────────────────────────────────────────────────────────────┤" + backgroundreset)
	log.Println("info: " + backgroundcolor + "│Created By : Suvir Kumar  <suvir@talkkonnect.com>               │" + backgroundreset)
	log.Println("info: " + backgroundcolor + "├────────────────────────────────────────────────────────────────┤" + backgroundreset)
	log.Println("info: " + backgroundcolor + "│Press the <Del> key for Menu or <Ctrl-c> to Quit talkkonnect    │" + backgroundreset)
	log.Println("info: " + backgroundcolor + "│Additional Modifications Released under MPL 2.0 License         │" + backgroundreset)
	log.Println("info: " + backgroundcolor + "│Blog at www.talkkonnect.com, source at github.com/talkkonnect   │" + backgroundreset)
	log.Println("info: " + backgroundcolor + "└────────────────────────────────────────────────────────────────┘" + backgroundreset)
	log.Printf("info: Talkkonnect Version %v Released %v", talkkonnectVersion, talkkonnectReleased)
	log.Printf("info: ")
}

func talkkonnectAcknowledgements(backgroundcolor string) {
	var backgroundreset string = "\u001b[0m"
	log.Println("info: " + backgroundcolor + "┌──────────────────────────────────────────────────────────────────────────────────────────────┐" + backgroundreset)
	log.Println("info: " + backgroundcolor + "│Acknowledgements & Inspriation from the talkkonnect team of developers, maintainers & testers │" + backgroundreset)
	log.Println("info: " + backgroundcolor + "│talkkonnect is based on the works of many people and many open source projects                │" + backgroundreset)
	log.Println("info: " + backgroundcolor + "├──────────────────────────────────────────────────────────────────────────────────────────────┤" + backgroundreset)
	log.Println("info: " + backgroundcolor + "│Thanks to Organizations :-                                                                    │" + backgroundreset)
	log.Println("info: " + backgroundcolor + "│The Mumble Development team, Raspberry Pi Foundation, Developers and Maintainers of Debian    │" + backgroundreset)
	log.Println("info: " + backgroundcolor + "│The Creators and Maintainers of Golang and all the libraries available on github.com          │" + backgroundreset)
	log.Println("info: " + backgroundcolor + "│Global Coders Co., Ltd. For Sponsoring this project                                           │" + backgroundreset)
	log.Println("info: " + backgroundcolor + "│                                                                                              │" + backgroundreset)
	log.Println("info: " + backgroundcolor + "│Thanks to Individuals :-                                                                      │" + backgroundreset)
	log.Println("info: " + backgroundcolor + "│Daniel Chote Creator of talkiepi and Tim Cooper Creator of Barnard and gumble library         │" + backgroundreset)
	log.Println("info: " + backgroundcolor + "│Zoran Dimitrijevic for his commitment, building, testing, docummentation and kind feedback    │" + backgroundreset)
	log.Println("info: " + backgroundcolor + "│enabling us to take talkkonnect to use cases never originally imagined                        │" + backgroundreset)
	log.Println("info: " + backgroundcolor + "├──────────────────────────────────────────────────────────────────────────────────────────────┤" + backgroundreset)
	log.Println("info: " + backgroundcolor + "│visit us at www.talkkonnect.com and github.com/talkkonnect                                    │" + backgroundreset)
	log.Println("info: " + backgroundcolor + "│talkkonnect was created by Suvir Kumar <suvir@talkkonnect.com> & Released under MPLV2 License │" + backgroundreset)
	log.Println("info: " + backgroundcolor + "└──────────────────────────────────────────────────────────────────────────────────────────────┘" + backgroundreset)
}

func (b *Talkkonnect) talkkonnectMenu(backgroundcolor string) {
	var backgroundreset string = "\u001b[0m"
	log.Println("info: " + backgroundcolor + "┌──────────────────────────────────────────────────────────────┐" + backgroundreset)
	log.Println("info: " + backgroundcolor + "│     _ __ ___   __ _(_)_ __    _ __ ___   ___ _ __  _   _     │" + backgroundreset)
	log.Println("info: " + backgroundcolor + "│    | '_ ` _ \\ / _` | | '_ \\  | '_ ` _ \\ / _ \\ '_ \\| | | |    │" + backgroundreset)
	log.Println("info: " + backgroundcolor + "│    | | | | | | (_| | | | | | | | | | | |  __/ | | | |_| |    │" + backgroundreset)
	log.Println("info: " + backgroundcolor + "│    |_| |_| |_|\\__,_|_|_| |_| |_| |_| |_|\\___|_| |_|\\__,_|    │" + backgroundreset)
	log.Println("info: " + backgroundcolor + "├─────────────────────────────┬────────────────────────────────┤" + backgroundreset)
	log.Println("info: " + backgroundcolor + "│ <Del> to Display this Menu  | <Ctrl-C> to Quit talkkonnect   │" + backgroundreset)
	log.Println("info: " + backgroundcolor + "├─────────────────────────────┼────────────────────────────────┤" + backgroundreset)
	log.Println("info: " + backgroundcolor + "│ <F1>  Channel Up (+)        │ <F2>  Channel Down (-)         │" + backgroundreset)
	log.Println("info: " + backgroundcolor + "│ <F3>  Mute/Unmute Speaker   │ <F4>  Current Volume Level     │" + backgroundreset)
	log.Println("info: " + backgroundcolor + "│ <F5>  Digital Volume Up (+) │ <F6>  Digital Volume Down (-)  │" + backgroundreset)
	log.Println("info: " + backgroundcolor + "│ <F7>  List Server Channels  │ <F8>  Start Transmitting       │" + backgroundreset)
	log.Println("info: " + backgroundcolor + "│ <F9>  Stop Transmitting     │ <F10> List Online Users        │" + backgroundreset)
	log.Println("info: " + backgroundcolor + "│ <F11> Playback/Stop Stream  │ <F12> For GPS Position         │" + backgroundreset)
	log.Println("info: " + backgroundcolor + "├─────────────────────────────┼────────────────────────────────┤" + backgroundreset)
	log.Println("info: " + backgroundcolor + "│<Ctrl-D> Debug Stacktrace    │                                │" + backgroundreset)
	log.Println("info: " + backgroundcolor + "├─────────────────────────────┼────────────────────────────────┤" + backgroundreset)
	log.Println("info: " + backgroundcolor + "│<Ctrl-E> Send Email          │<Ctrl-N> Conn Next Server       │" + backgroundreset)
	log.Println("info: " + backgroundcolor + "│<Ctrl-F> Conn Previous Server│<Ctrl-P> Panic Simulation       │" + backgroundreset)
	log.Println("info: " + backgroundcolor + "│<Ctrl-G> Send Repeater Tone  │<Ctrl-S> Scan Channels          │" + backgroundreset)
	log.Println("info: " + backgroundcolor + "│<Ctrl-V> Display Version     │<Ctrl-T> Thanks/Acknowledgements│" + backgroundreset)
	log.Println("info: " + backgroundcolor + "├─────────────────────────────┼────────────────────────────────┤" + backgroundreset)
	log.Println("info: " + backgroundcolor + "│<Ctrl-L> Clear Screen        │<Ctrl-O> Ping Servers           │" + backgroundreset)
	log.Println("info: " + backgroundcolor + "│<Ctrl-R> Repeat TX Loop Test │<Ctrl-X> Dump XML Config        │" + backgroundreset)
	log.Println("info: " + backgroundcolor + "├─────────────────────────────┼────────────────────────────────┤" + backgroundreset)
	log.Println("info: " + backgroundcolor + "│<Ctrl-I> Traffic Record      │<Ctrl-J> Mic Record             │" + backgroundreset)
	log.Println("info: " + backgroundcolor + "│<Ctrl-K> Traffic & Mic Record│<Ctrl-U> Show Uptime            │" + backgroundreset)
	log.Println("info: " + backgroundcolor + "├─────────────────────────────┼────────────────────────────────┤" + backgroundreset)
	log.Println("info: " + backgroundcolor + "│  Visit us at www.talkkonnect.com and github.com/talkkonnect  │" + backgroundreset)
	log.Println("info: " + backgroundcolor + "│  Thanks to Global Coders Co., Ltd. for their sponsorship     │" + backgroundreset)
	log.Println("info: " + backgroundcolor + "└──────────────────────────────────────────────────────────────┘" + backgroundreset)

	log.Println("info: IP Address & Session Information")
	b.pingconnectedserver()
	localAddresses()
	log.Println("info: Internet WAN IP Is", getOutboundIP())

	macaddress, err := getMacAddr()
	if err != nil {
		log.Println("error: Could Not Get Network Interface MAC Address")
	} else {
		for i, a := range macaddress {
			log.Println("info: Network Interface MAC Address (" + strconv.Itoa(i) + "): " + a)
		}
	}

	hostname, err1 := os.Hostname()
	if err1 != nil {
		log.Printf("warn: Cannot Get Hostname\n")
	} else {
		log.Printf("info: Hostname is %s\n", hostname)
	}

	log.Printf("info: Talkkonnect Version %v Released %v\n", talkkonnectVersion, talkkonnectReleased)
}
