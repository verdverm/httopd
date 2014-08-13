httopd - top for httpd logs
====================================

[verdverm/httpod](https://github.com/verdverm/httpod)
pronounced "hopped" like a gopher

Dependencies
------------

1. Docker
2. *nix

Installation
-------------
`go get github.com/verdverm/httopd/httopd`

Running Simulator
------------

1. `git clone https://github.com/verdverm/httopd && cd httopd`
2. `sudo build.sh`
3. `sudo run.sh`
4. `cd httopd`
5. `go build`
6. `httopd -fn="/abs/path/to/log/file"`

   `sudo` is required for the docker commands (unless you run a non-sudo docker setup)

Details
------------

### server

This docker contains a Python-Flask site.

### client

This docker contains a http-client simulator.

### httpod

This docker contains the httpod program.


Enhancements
----------------------

- monitor multiple log files when there are multiple sites / server blocks

- there are race conditions on the statistics
  - there is a single reader and a single writer
  - shouldn't be too much of an issues since only one party reads and only one party writes

- page stats only update when new line_data shows up for that page

- add history and alerts to errors

- ML triggers & alerts


References
---------------

[Logging Control In W3C httpd - Logfile Format](http://www.w3.org/Daemon/User/Config/Logging.html#common-logfile-format)
[termbox-go](https://github.com/nsf/termbox-go)
