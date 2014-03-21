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
