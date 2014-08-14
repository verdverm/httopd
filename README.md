httopd - top for httpd logs
====================================

pronounced "hopped" like a gopher

![boing!](https://raw.github.com/verdverm/httopd/master/glenda_space_medium.jpg)


Dependencies
------------

1. Docker
2. *nix

Installation
-------------
`go get github.com/verdverm/httopd/httopd`

Running
-----------

watching a single log file
`httopd -fn="/abs/path/to/log/file"`

watching a list of log files
`httopd -fnList="path/to/list/file"`

list file format
```
/abs/path/to/log/access.log
/abs/path/to/site/log/access.log
/abs/path/to/site/log/access.log
```


Simulator
------------

1. `git clone https://github.com/verdverm/httopd && cd httopd`
2. `sudo build.sh`
3. `sudo run.sh`


`sudo` is required for the docker commands (unless you run a non-sudo docker setup)

Enhancements / Issues / Todo's
------------------------------

Want to help out? Add or remove items from the following list.

- log file format
  -- multiple files / domains... how to handle?
  -- only handling access.log default from nginx
  -- what about error.log ?
  -- other log providers ?
- monitor multiple log files when there are multiple sites / server blocks
- there are race conditions on the statistics
  -- there is a single reader and a single writer
  -- shouldn't be too much of an issues since only one party reads and only one party writes
- page stats only update when new line_data shows up for that page
- add history and alerts to errors
- backfill with 10 minutes of history on startup
- sort by columns
- config file or directory inspection for log files / log directories
- ML triggers & alerts
- more configurable triggers / set from CLI / save to file?
- better log line parser
  -- simpler and more flexible
  -- only tested with nginx, need to check apache and others
- stats should be kept in something like an R data frame
  -- so that aggregates can be calculated more easily
  -- can rely on an external library for stats calculations
- there may be some errors now resulting from multiple log files
  -- multiple watchers are firing linedata to the stats over channels
  -- if they arrive in a 'bad' order (reverse alternations over a minute boundary)
  -- we will switch bins, but the next arrival (last minute) will end up in the current bin, not the last one
  -- should really determine the bin by the time stamp, and then index to it
  -- right now we are just adding to the current bin
  -- this will be an issue with backfilling, which once resolved, we can backfile file by file, and not worry about time

Subdir Details
------------

#### server

This docker contains a Python-Flask site.

#### client

This docker contains a http-client simulator.

#### httpod

This folder contains the httpod program code.

#### logs

a temporary directory created and destroyed by the simulator

References
---------------

1. [Logging Control In W3C httpd - Logfile Format](http://www.w3.org/Daemon/User/Config/Logging.html#common-logfile-format)
2. [termbox-go](https://github.com/nsf/termbox-go)
