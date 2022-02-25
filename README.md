# Koala

koala is a simple browser frontend for editing zonefiles and applying the changes made. It supports basic authentication via JWT tokens built-in, usage in a production environment is not yet recommended.

## Dependencies
koala requires Go 1.16 and NodeJS 12.X or newer.

## Configuration
```toml
addr = ":8081"          # address the server will be listening on
debug = true            # debug mode
cors = false            # enable cors
apply_cmd = "sleep 1"   # command to execute after zonefile change
api_root = "/api"       # path to use for api

[[zones]]
path = "zone"           # zonefile to be edited
origin = "home.arpa."   # zone to be edited
ttl = 3600              # default TTL

[security]
mode = "none"           # guard to use, either 'none', 'oidc' or 'jwt'

[security.oidc]
client_id = ""          # OpenID connect client ID
identity_server = ""    # URL of identity provider

[security.jwt]
secret = ""             # Auth secret for JWT tokens
```

## Installation
**Step 1:** Download one of the binary packages from the release site
```bash
# Linux amd64
curl -O -L https://github.com/lnsp/koala/releases/latest/download/koala_linux_amd64
# Linux arm
curl -O -L https://github.com/lnsp/koala/releases/latest/download/koala_linux_arm64
# macOS amd64
curl -O -L https://github.com/lnsp/koala/releases/latest/download/koala_darwin_amd64
```

**Step 2:** Copy binary to bin dir
```bash
cp koala_* /usr/local/bin/
```

**Step 3:** Create configuration folder and customize configuration
```bash
mkdir -p /etc/koala
cp config.toml /etc/koala/config.toml
nano config.toml
```

**Step 4:** Install a startup script, you should customize it though. Please remember to
protect yourself from unauthorized access.
```bash
cat > /etc/systemd/system/koala.service << EOF
[Unit]
Description=Koala DNS editing frontend
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/etc/koala/
ExecStart=/usr/local/bin/koala -config config.toml
Restart=on-abort

[Install]
WantedBy=multi-user.target
EOF
systemctl daemon-reload
systemctl enable koala
systemctl start koala
```
