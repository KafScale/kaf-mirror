#!/bin/bash

# Check for the existence of the binaries
if [ ! -f "kaf-mirror" ] || [ ! -f "admin-cli" ] || [ ! -f "mirror-cli" ]; then
  echo "Binaries not found in the current directory."
  read -p "Please enter the path to the binaries: " BINARY_PATH
  if [ ! -f "$BINARY_PATH/kaf-mirror" ] || [ ! -f "$BINARY_PATH/admin-cli" ] || [ ! -f "$BINARY_PATH/mirror-cli" ]; then
    echo "Binaries not found in the specified path. Exiting."
    exit 1
  fi
  cd $BINARY_PATH
fi

# Check for sudo privileges
if [ "$EUID" -ne 0 ]; then
  SUDO=sudo
else
  SUDO=
fi

# Create a dedicated user and group for the service
$SUDO groupadd kaf-mirror
$SUDO useradd -r -g kaf-mirror -s /bin/false kaf-mirror

# Create the necessary directories
$SUDO mkdir -p /opt/kaf-mirror
$SUDO mkdir -p /var/log/kaf-mirror
$SUDO chown -R kaf-mirror:kaf-mirror /opt/kaf-mirror
$SUDO chown -R kaf-mirror:kaf-mirror /var/log/kaf-mirror

# Copy the binaries and configuration to the correct locations
$SUDO cp kaf-mirror /opt/kaf-mirror/
$SUDO cp admin-cli /usr/local/bin/
$SUDO cp mirror-cli /usr/local/bin/
$SUDO cp -r configs /opt/kaf-mirror/
$SUDO cp -r data /opt/kaf-mirror/

# Copy the service file to the systemd directory
$SUDO cp scripts/kaf-mirror.service /etc/systemd/system/

# Reload the systemd daemon, enable the service, and start it
$SUDO systemctl daemon-reload
$SUDO systemctl enable kaf-mirror.service
$SUDO systemctl start kaf-mirror.service

echo "kaf-mirror service has been installed and started."
