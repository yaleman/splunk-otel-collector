ARG SPLUNK_OTEL_COLLECTOR_TARBALL=splunk-otel-collector_0.78.1_amd64.tar.gz
FROM python:3.9.16-bullseye
ARG SPLUNK_OTEL_COLLECTOR_TARBALL

RUN groupadd -g 123 splunk-otel-collector
RUN useradd splunk-otel-collector -u 123 -g 123 -s /usr/sbin/nologin -d /nonexistent -M

ENV DEBIAN_FRONTEND noninteractive
RUN apt-get update -y && apt-get install ntp sudo -y

WORKDIR /tmp
COPY ${SPLUNK_OTEL_COLLECTOR_TARBALL} ./
WORKDIR /usr/lib
RUN tar -xvf /tmp/"${SPLUNK_OTEL_COLLECTOR_TARBALL}"
RUN chown -R splunk-otel-collector:splunk-otel-collector /usr/lib/splunk-otel-collector
RUN chmod a+rx /usr/lib

USER splunk-otel-collector
RUN /usr/lib/splunk-otel-collector/agent-bundle/bin/patch-interpreter /usr/lib/splunk-otel-collector/agent-bundle

COPY --chown=123 upstat /var/collectd-python/upstat
COPY --chown=123 config.yaml /etc/config.yaml

ENV PLUGIN_FOLDER /var/collectd-python/upstat

USER root
# Make sure that splunk-otel-collector cannot access any other python libs
RUN chmod -R o-rwx /usr/local/lib
COPY entrypoint.sh /opt/entrypoint.sh
CMD ["bash", "/opt/entrypoint.sh"]

