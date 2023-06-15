set -ex

/usr/sbin/ntpd -n &

sudo -u splunk-otel-collector -En /usr/lib/splunk-otel-collector/bin/otelcol --config /etc/config.yaml
