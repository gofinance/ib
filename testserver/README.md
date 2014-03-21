TestServer
==========

The files in this directory are not part of GoIB, but they will help you test it.

* [Trader Workstation](http://www.interactivebrokers.com/en/pagemap/pagemap_APISolutions.php)
  (contained in file ``unixmacosx-*.jar``) is the official, unmodified
  distribution JAR from Interactive Brokers. IB API is accessed by locally
  running either Trader Workstation (TWS) or IB Gateway, both of which are
  shipped in the Trader Workstation distribution. The license can be viewed by
  loading TWS and clicking Help > About Trader Workstation.
* [IBController](http://sourceforge.net/projects/ibcontroller/) provides
  automated management of IB Gateway. It is GPLv3 licensed.

The ``ibcontroller-*.ini`` has been configured to automatically load IB
Gateway and login with Interactive Brokers' demo account. This is adequate for
most GoIB tests. Two ports are bound:

* 4001 is the IB API port. This is what the tests will use.
* 7462 is the IBController telnet control port (bound to 127.0.0.1). Tests do
  not use this port, but the ``ibgwstop`` command will.

Run ``ibgwstart`` from the ``testserver`` directory to start the gateway. The
script will create a clean (ie default) setup for IB Gateway and run it. You
can use the same IB Gateway instance for repeated tests.

To cleanly terminate the server via IBController, use ``ibgwstop``.

Note that IB Gateway will connect to actual IB backends, and these backends are
regularly reset at fixed times each day, cycle through demo accounts, apply rate
limits and occasionally timeout. Test failures are frequently related to these
conditions as opposed to errors in the GoIB trading library. If you receive a
test failure, try re-running the test suite later.
