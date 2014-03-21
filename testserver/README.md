TestServer
==========

The files in this directory are not part of GoIB, but they will help you test it.

* IBController provides automated TWS/IB Gateway start/stop/management. Its
  home page is http://sourceforge.net/projects/ibcontroller/ and is GPLv3
  licensed.
* Trader Workstation (in ```unixmacosx-*.jar```) is the official, unmodified
  distribution JAR from Interactive Brokers. Its home page is
  http://www.interactivebrokers.com/en/pagemap/pagemap_APISolutions.php and its
  license is shown by loading TWS and clicking  Help > About Trader Workstation.

The ```ibcontroller-*.ini``` has been configured to automatically login with
IB's demo account. This is adequate for most IB tests. Ports configured:

* 4001 is the TWSAPI port. It is what the tests will use.
* 7462 is the IBController telnet control port (bound to 127.0.0.1).

Run ``ibgwstart`` from the ``testserver`` directory to run the gateway. The
script will create a clean (ie default) setup for IB Gateway and run it.

To cleanly terminate the server via IBController, use ``ibgwstop``.

Note that IB Gateway will connect to actual IB backends, and these backends are
regularly reset at fixed times each day, cycle through demo accounts, apply rate
limits and regularly timeout. Test failures are frequently related to these
conditions as opposed to errors in the GoIB trading library. If you receive a
test failure, try re-running the test suite later.
