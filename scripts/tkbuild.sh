#!/bin/bash

## talkkonnect headless mumble client/gateway with lcd screen and channel control
## Copyright (C) 2018-2019, Suvir Kumar <suvir@talkkonnect.com>
##
## This Source Code Form is subject to the terms of the Mozilla Public
## License, v. 2.0. If a copy of the MPL was not distributed with this
## file, You can obtain one at http://mozilla.org/MPL/2.0/.
##
## Software distributed under the License is distributed on an "AS IS" basis,
## WITHOUT WARRANTY OF ANY KIND, either express or implied. See the License
## for the specific language governing rights and limitations under the
## License.
##
## The Initial Developer of the Original Code is
## Suvir Kumar <suvir@talkkonnect.com>
## Portions created by the Initial Developer are Copyright (C) Suvir Kumar. All Rights Reserved.
##
## Contributor(s):
##
## Suvir Kumar <suvir@talkkonnect.com>
##
## My Blog is at www.talkkonnect.com
## The source code is hosted at github.com/talkkonnect


## Installation BASH Script for talkkonnect on fresh install of raspbian
## Please RUN this Script as root user

## If this script is run after a fresh install of raspbian you man want to update the 2 lines below
apt-get update
apt-get -y dist upgrade

## Add talkkonnect user to the system
adduser --disabled-password --disabled-login --gecos "" talkkonnect
usermod -a -G cdrom,audio,video,plugdev,users,dialout,dip,input,gpio talkkonnect

## Install the dependencies required for talkkonnect
apt-get -y install libopenal-dev libopus-dev libasound2-dev git ffmpeg omxplayer screen

## Create the necessary directory structure under /home/talkkonnect/
cd /home/talkkonnect/
mkdir -p /home/talkkonnect/gocode
mkdir -p /home/talkkonnect/bin

## Create the log file
touch /var/log/talkkonnect.log

cd /usr/local
wget https://golang.org/dl/go1.15.linux-armv6l.tar.gz
tar -zxvf go1.15.linux-armv6l.tar.gz

echo export PATH=$PATH:/usr/local/go/bin >>  ~/.bashrc
echo export GOPATH=/home/talkkonnect/gocode >>  ~/.bashrc
echo export GOBIN=/home/talkkonnect/bin >>  ~/.bashrc
echo export GO111MODULE="auto" >>  ~/.bashrc
echo "alias tk='cd ~/go/src/github.com/jdiderik/talkkonnect/'" >>  ~/.bashrc


## Set up GOENVIRONMENT
export PATH=$PATH:/usr/local/go/bin
export GOPATH=/home/talkkonnect/gocode
export GOBIN=/home/talkkonnect/bin
export GO111MODULE="auto"

## Get the latest source code of talkkonnect from githu.com
go get -v github.com/jdiderik/talkkonnect

## Build talkkonnect as binary
cd $GOPATH/src/github.com/jdiderik/talkkonnect
/usr/local/go/bin/go build -o /home/talkkonnect/bin/talkkonnect cmd/talkkonnect/main.go

## Notify User
echo "=> Finished building TalKKonnect"
echo "=> talkkonnect binary is in /home/talkkonect/bin"
echo "=> Now enter Mumble server connectivity details"
echo "talkkonnect.xml from $GOPATH/src/github.com/jdiderik/talkkonnect"
echo "and configure talkkonnect features. Happy talkkonnecting!!"

exit


