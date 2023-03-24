# checker-reporter
application for   links validation between distributed locations

on Linux, you must set:
setcap cap_net_raw=+ep /path/to/your/compiled/binary
for properly working.

Config file format:
update_interval=30    - Timeout between checks
check_hosts=["10.10.10.10","20.20.20.20"]  - array of ip addresses of checkers host, include all checkers, and one witch installed
metrics_port=2112  - port for scrape metrics
check_port=8080  - port for binding httpcheck handler, must be unfirewalled for other checkers