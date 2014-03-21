GoIB
====

[![Build Status](https://drone.io/github.com/benalexau/go.trade/status.png)](https://drone.io/github.com/benalexau/go.trade/latest)
[![GoDoc](https://godoc.org/github.com/benalexau/go.trade?status.png)](https://godoc.org/github.com/benalexau/go.trade)

This is a pure Go interface to
[Interactive Brokers](https://www.interactivebrokers.com/)
[TWSAPI](http://interactivebrokers.github.io). Features include:

* 100% pure Go (no Java/C calls)
* Idiomatic design (Go naming conventions, channels, error handling etc)
* Choice of low-level types or our high-level [Manager](manager.go) types
* Extensively tested (lots of tests, CI via
  [Drone](https://drone.io/github.com/benalexau/go.trade/latest), local
  [test server](testserver/README.md))
* Documentation on [GoDoc](https://godoc.org/github.com/benalexau/go.trade)
* Reflects very up-to-date TWSAPI version

We welcome your involvement and contributions! If you like the project, please
click [star](../../star) to help others find it, or [fork](../../fork) it to
produce [pull requests](../../pulls) and join the [contributor
list](../../graphs/contributors).

Status
------

* The code presently supports TWSAPI 971.01 (latest as of March 2014)
* Most reply types (see [ereader.go](ereader.go)) are already handled
* Some request types (see [eclientsocket.go](eclientsocket.go)) require porting

Testing
-------

```go test``` requires IB Gateway be running at 127.0.0.1:4001. Always use a
demo or paper trade account, as the tests may modify your account.

The easiest way to start IB Gateway with a demo account is to use the test
server. Have a look at the [test server instructions](testserver/README.md) for
all the details.

By default the tests produce no output. If you'd like to view engine
communication logs during test execution, set the ```IB_ENGINE_DUMP```
environment variable to any value. For example, ```IB_ENGINE_DUMP=t go test```.

If you fork this project and would like to configure Drone for your fork, our
our [Drone instructions](drone.md) should be of help.

License
-------

This library is licensed under
[GNU Lesser General Public License](http://www.gnu.org/licenses/lgpl.html)
version 3.

**Static linking exception**: The copyright holders give you permission to link
this library with independent modules to produce an executable, regardless of
the license terms of these independent modules, and to copy and distribute
the resulting executable under terms of your choice, provided that you also
meet, for each linked independent module, the terms and conditions of the
license of that module. An independent module is a module which is not
derived from or based on this library. If you modify this library, you must
extend this exception to your version of the library.

* This library is safe for use in close-source applications. The LGPL
  share-alike terms do not apply to applications built on top of this library.
* You do not need a commercial license. The LGPL applies to the library's own
  source code, not your applications.
