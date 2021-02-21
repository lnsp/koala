# Koala

koala is a simple browser frontend for editing zonefiles and applying the changes made. It supports basic authentication via JWT tokens built-in, usage in a production environment is not yet recommended.

## Dependencies
koala requires Go 1.16 and NodeJS 12.X or newer.

## Configuration
KEY                         | TYPE             | DEFAULT                  | REQUIRED    | DESCRIPTION
----------------------------|------------------|--------------------------|-------------|----------------------------------------------------
KOALA_ADDR                  | String           | :8080                    |             | Address the server will be listening on
KOALA_ZONEFILE              | String           |                          | true        | Zonefile to be edited
KOALA_ORIGIN                | String           | .                        |             | Zone to be edited
KOALA_TTL                   | Integer          | 3600                     |             | Default TTL for records
KOALA_APPLYCMD              | String           | sleep 1                  |             | Command executed after applying zonefile changes
KOALA_DEBUG                 | True or False    | false                    |             | Enable debug logging
KOALA_CORS                  | True or False    | false                    |             | Enable support for CORS
KOALA_SECURITY              | String           | none                     |             | Security guard to use [none|oidc|jwt]
KOALA_OIDCCLIENTID          | String           |                          |             | OpenID Connect Client ID
KOALA_OIDCIDENTITYSERVER    | String           |                          |             | URL of identity provider
KOALA_JWTSECRET             | String           |                          |             | Auth secret for JWT tokens

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

**Step 3:** Install a startup script, you should customize it though. Please remember to
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
