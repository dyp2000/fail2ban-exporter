## Fail2ban exporter for Prometheus

Yet another Fail2ban exporter

### Installation

#### From source

You can clone this repository from git https://github.com/dyp2000/fail2ban-exporter.git.
Then all you have to do is run `make all` and you're done!

Required packages:

 - github.com/prometheus/client_golang/prometheus
 - github.com/prometheus/client_golang/prometheus/promauto
 - github.com/prometheus/client_golang/prometheus/promhttp

### Configuration

| param | Default | Description |
|-------|---------|-------------|
| -port | 9635 | Listen port |
| -engine | ip2c | Set engine for GeoIP |
| -version |  | Show version |
| -h |  | Show help |

`-engine` key supported values:

- `ip2c` - get info from ip2c.org. Unlimited, just be reasonable. Currently ip2c.org can sustain a maximum of about 30 million per day.
- `geoiplookup` - use geoiplookup console util (fastest method. Linux only).
- `freegeoip` - get info from https://freegeoip.app. Not released.

### Configure Systemd daemon

Copy file `fail2ban-exporter.service` from project to `/etc/systemd/system` directory.

Find `ExecStart=/home/prometheus/fail2ban-exporter/fail2ban-exporter"` string and replace path to real path where located your file `fail2ban-exporter`.

Execute next commands in terminal:

```bash
systemctl enable fail2ban-exporter
systemctl start fail2ban-exporter
systemctl status fail2ban-exporter
```

### Metrics for Prometheus

| Metric | Description |
|--------|-------------|
| fail2ban_banned_current | Number of currently banned IP addresses. |
| fail2ban_banned_total | Number of total banned IP addresses. |
| fail2ban_failed_current | Number of current failed connections. |
| fail2ban_failed_total | Number of total failed total. |
| fail2ban_hackers_locations | Location of hackers on world map. |


For Grafana Worldmap Panel use query: `sum(fail2ban_hackers_locations{jail="sshd"}) by (location)`

Install Worldmap Panel see [there](https://grafana.com/grafana/plugins/grafana-worldmap-panel/installation)
