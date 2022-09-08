[Unit]
Description=fwdctl systemd service
After=network.target

[Service]
Type=oneshot
ExecStart={{.InstallationPath}}/fwdctl apply --rules-file={{.RulesFile}}
StandardOutput=journal

[Install]
WantedBy=multi-user.target