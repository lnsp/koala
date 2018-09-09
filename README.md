<p align="center">
<img src="https://github.com/lnsp/koala/raw/master/koala.png" alt="koala">
</p>
<hr>
koala is a simple web-service for editing local zonefiles. It HAS NO security built-in and should not be used outside any test environment.

## Dependencies
koala requires a recent Go version (tested on `>=1.10`) and npm.

## Configuration
```
KOALA_ADDR=":8080"
KOALA_ZONEFILE=
KOALA_STATICDIR="web/dist"
KOALA_APPLYCMD="sleep 1"
```
## Installation
**Step 1:** Download one of the binary packages from the release site
```bash
# for linux amd64
curl -O https://github.com/lnsp/koala/releases/download/v0.2.1/koala-v0.2.1-linux-amd64.tar.gz
# for linux arm
curl -O https://github.com/lnsp/koala/releases/download/v0.2.1/koala-v0.2.1-linux-arm.tar.gz
# for macOS amd64
curl -O https://github.com/lnsp/koala/releases/download/v0.2.1/koala-v0.2.1-darwin-amd64.tar.gz
```
**Step 2:** Extract the contents to a target location
```bash
tar -C /usr/local -xzvf koala-v0.2.1-$GOOS-$GOARCH.tar.gz
```
**Step 3:** *(Optional)* Create link to binary
```bash
ln -sf /usr/local/koala/koala /usr/local/bin/koala
```
**Step 4:** *(Optional, Linux only)* Install a startup script, you should customize it though.
```bash
cat > /etc/systemd/system/koala.service << EOF
[Unit]
Description=Koala DNS editing frontend
After=network.target

[Service]
Type=simple
User=root
Environment=KOALA_ADDR=:80
Environment=KOALA_ZONEFILE=/etc/bind/zones/home.arpa.zone
Environment=KOALA_STATICDIR=/usr/local/koala/web
Environment=KOALA_APPLYCMD="systemctl reload bind9"
WorkingDirectory=/root/
ExecStart=/usr/local/bin/koala
Restart=on-abort

[Install]
WantedBy=multi-user.target
EOF
```
