Datadog Coding Challenge - Tony Worm
====================================

[verdverm/dawg](https://github.com/verdverm/dawg)

Dependencies
------------

1. Docker
2. *nix

Running
------------

1. `git clone https://github.com/verdverm/dawg && cd dawg`
2. `sudo build.sh`
3. `sudo run.sh`

   `sudo` is required for the docker commands

Details
------------

### Sever

This docker contains a Python-Flask site.

### Client

This docker contains a http-client simulator.

### Monitor

This docker contains the CLI monitor program.


Enhancements
----------------------

- monitor multiple log files when there are multiple sites / server blocks
- there are race conditions on the statistics
  - there is a single reader and a single writer
  - shouldn't be too much of an issues since only one party reads and only one party writes

- weird '1' fill on some starts

- add history and alerts to errors

- ML triggers & alerts

Problem Statement
-----------------

### HTTP log monitoring console program

Create a simple console program that monitors HTTP traffic on your machine:

- Consume an actively written-to w3c-formatted HTTP access log
- Every 10s, display in the console the sections of the web site with the most hits (a section is defined as being what's before the second '/' in a URL. i.e. the section for "http://my.site.com/pages/create' is "http://my.site.com/pages"), as well as interesting summary statistics on the traffic as a whole.
- Make sure a user can keep the console app running and monitor traffic on their machine
- Whenever total traffic for the past 2 minutes exceeds a certain number on average, add a message saying that “High traffic generated an alert - hits = {value}, triggered at {time}”
- Whenever the total traffic drops again below that value on average for the past 2 minutes, add another message detailing when the alert recovered
- Make sure all messages showing when alerting thresholds are crossed remain visible on the page for historical reasons.
- Write a test for the alerting logic
- Explain how you’d improve on this application design


References
---------------

[Logging Control In W3C httpd - Logfile Format](http://www.w3.org/Daemon/User/Config/Logging.html#common-logfile-format)
