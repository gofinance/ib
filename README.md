go-trade
========

[![Build Status](https://drone.io/github.com/benalexau/go.trade/status.png)](https://drone.io/github.com/benalexau/go.trade/latest)

Go implementation of the Interactive Brokers TWS API  

Testing
=======

```go test``` requires IB Gateway be running at 127.0.0.1:4001. Always use a
demo or paper trade account, as the tests may modify your account state.

The easiest way to start IB Gateway with the demo account is to run
```testserver/ibgwstart``` (shutdown with ```ibgwstop```).

By default the tests are quiet. To view key engine communication logs during
test execution, set the ```IB_ENGINE_DUMP``` environment variable. For example,
```IB_ENGINE_DUMP=t go test```.
