#!/bin/sh

echo "=== System Information ==="
uname -a
cat /etc/os-release

echo "\n=== Dovecot Version ==="
dovecot --version

echo "\n=== Dovecot Configuration ==="
doveconf -n

echo "\n=== SSL Certificate and Key ==="
ls -l /etc/dovecot/dovecot.pem /etc/dovecot/dovecot.key

echo "\n=== Dovecot Configuration File ==="
cat /etc/dovecot/dovecot.conf

echo "\n=== System SSL Certificates ==="
ls -l /etc/ssl/certs

echo "\n=== Dovecot Processes ==="
ps aux | grep dovecot

echo "\n=== System Logs ==="
tail -n 50 /var/log/messages /var/log/syslog 2>/dev/null

echo "\n=== Dovecot Logs ==="
tail -n 50 /var/log/dovecot.log 2>/dev/null

echo "\n=== File Permissions ==="
ls -l /etc/dovecot
ls -l /var/run/dovecot

echo "\n=== Network Status ==="
netstat -tuln

echo "\n=== Dovecot Start Attempt ==="
dovecot -F