APP_ROOT=/opt/barracuda
SERVICE_FILE=/etc/systemd/system/barracuda.service
SERVICE_EXISTS=0

read -r -d '' SERVICE_TEMPLATE << EOF
[Unit]
Description=Barracuda Server
After=network.target

[Service]
ExecStart=%s/bin/barracuda
User=%s
Type=simple
Restart=always

[Install]
WantedBy=multi-user.target\n
EOF

if [[ $(id -u) > 0 ]]; then
    echo "script must be run as root" >&2
    exit 1
fi

if [[ -e "$SERVICE_FILE" ]]; then
    echo "stopping service"
    SERVICE_EXISTS=1
    # systemctl stop barracuda
else
    echo "creating service"
    printf "$SERVICE_TEMPLATE" $APP_ROOT $USER > $SERVICE_FILE
fi

echo "building application"
go build -o ${APP_ROOT}/bin/barracuda ./cmd
echo "build success"

if [[ $SERVICE_EXISTS == 1 ]]; then
    echo "restarting service"
    systemctl restart barracuda
    exit 0
fi

echo "starting service"
systemctl start barracuda

