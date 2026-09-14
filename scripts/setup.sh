#!/bin/sh
set -euo pipefail

echo "Creating /var/lib/gitria directory with proper permissions..."
sudo mkdir -pv /var/lib/gitria/repo
sudo chown -R "$(id -u):$(id -g)" /var/lib/gitria
