[Unit]
Description=fwdctl systemd service
After=network.target

[Service]
Type={{.ServiceType}}
{{if eq .ServiceType "oneshot"}}ExecStart={{.InstallationPath}}/fwdctl apply --file={{.RulesFile}}{{else if eq .ServiceType "fork"}}ExecStart={{.InstallationPath}}/fwdctl daemon start --file={{.RulesFile}}
ExecStop={{.InstallationPath}}/fwdctl daemon stop
Restart=always
RestartSec=5s{{end}}
StandardOutput=journal

[Install]
WantedBy=multi-user.target
