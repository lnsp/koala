# ![koala](https://raw.githubusercontent.com/lnsp/koala/master/docs/logo.png)

koala is a simple web-service for editing local zonefiles. It supports basic authentication via JWT tokens built-in, usage in a production environment is not yet recommended.

## Dependencies
koala requires a recent Go version (tested on `>=1.10`) and npm.

## Development setup

To enable easy and configuration-less local testing, this repository includes a Vagrantfile to help you setup a working
testing environment. Vagrant is an open-source tool by HashiCorp that enables fast VM setup.
The development environment includes a recent version of Vim, cURL, Golang, NodeJS and Bind9.
After startup the VM is reachable under the IP `192.168.100.100`.

```bash
git clone git@github.com:lnsp/koala
cd koala
vagrant up
vagrant ssh
```

We recommend using the installation guidelines below to achieve a similar-to-prod environment configuration.

## Configuration
```
KOALA_ADDR=":8080"
KOALA_ZONEFILE=
KOALA_APPLYCMD="sleep 1"
KOALA_DEBUG=false
KOALA_JWTSECRET=
KOALA_CERTIFICATE=
KOALA_PRIVATEKEY=
KOALA_CORS=
```

## Installation
**Step 1:** Download one of the binary packages from the release site
```bash
# Linux amd64
curl -O -L https://github.com/lnsp/koala/releases/download/v0.4.2/koala-v0.4.2-darwin-amd64.tar.gz
# Linux arm
curl -O -L https://github.com/lnsp/koala/releases/download/v0.4.2/koala-v0.4.2-linux-arm.tar.gz
# macOS amd64
curl -O -L https://github.com/lnsp/koala/releases/download/v0.4.2/koala-v0.4.2-darwin-amd64.tar.gz
```

**Step 2:** Extract the contents to a target location
```bash
tar -C /opt -xzvf koala-v0.4.2-*.tar.gz
```

**Step 3:** *(Optional)* Create link to binary
```bash
ln -sf /opt/koala/koala /usr/local/bin/koala
```

**Step 4:** *(Optional, Linux only)* Install a startup script, you should customize it though. Please remember to
protect yourself from unauthorized access.
```bash
cat > /etc/systemd/system/koala.service << EOF
[Unit]
Description=Koala DNS editing frontend
After=network.target

[Service]
Type=simple
User=root
Environment=KOALA_ADDR=localhost:8000
Environment=KOALA_ZONEFILE=/etc/coredns/home.db
Environment=KOALA_APPLYCMD=echo
WorkingDirectory=/root/
ExecStart=/usr/local/bin/koala
Restart=on-abort

[Install]
WantedBy=multi-user.target
EOF
systemctl daemon-reload
systemctl enable koala
systemctl start koala
```

**Step 5**: *(Optional, Linux only)* Route requests using Nginx reverse proxy
```bash
apt-get update && apt-get install -y nginx
cat > /etc/nginx/sites-available/default <<\EOF
server {
  listen 80 default_server;
  listen [::]:80 default_server;

  root /opt/koala/webui;
  index index.html;
  server_name _;

  location / {
      try_files $uri $uri/ =404;
  }

  location /api {
      return 302 /api/;
  }

  location /api/ {
      proxy_set_header Host $host;
      proxy_set_header X-Real-IP $remote_addr;
      proxy_pass http://localhost:8000/;
  }
}
EOF
systemctl reload nginx
```
