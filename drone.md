Drone Config
============

To use with http://drone.io, configure the Drone commands as follows:

```
lsb_release -a
sudo start xvfb
sudo apt-get install -y telnet
go get
go build
cd testserver
./ibgwstart
sleep 15
cd ..
go test
./testserver/ibgwstop
```
