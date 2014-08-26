# DEPRECATED
This exporter was depricated in favor of [container-exporter](http://github.com/discordianfish/container-exporter).

## docker-exporter

This exposes metrics about docker and it's containers on /metrics.
This data is expected to be consumed by [prometheus](http://github.com/prometheus/prometheus)
but might be useful to you on it's own.


## Usage

    Usage of ./docker-exporter:
      -addr="unix:///var/run/docker.sock": Docker address to connect to
      -interval=15s: refresh interval
      -listen=":8080": Address to listen on
      -root="/sys/fs/cgroup": cgroup root
      -telemetry.abortonmisuse=false: abort if a semantic misuse is encountered (bool).
      -telemetry.debugregistration=false: display information about the metric registration process (bool).
      -telemetry.useaggressivesanitychecks=false: perform expensive validation of metrics (bool).
