# Systemd Email Alerter

Trigger email Alerter for failed systemd services

## Build
```sh
goreleaser --rm-dist --snapshot && mv dist/**/systemd-alerter /usr/bin/
```

## Copy Unit file
Copy systemd unit file called "alert-email@service" on /etc/systemd/system/ or /lib/systemd/system

## Register on your unit service
```
[Unit]
OnFailure=alert-email@%i.service
...
```

Reload systemd unit files
```sh 
systemctl daemon-reload
```