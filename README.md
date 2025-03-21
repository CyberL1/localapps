# Localapps

Localapps is a project designed to simplify the management and deployment of local applications

## Features

- Easy setup and configuration
- Apps shut down when you don't use them

## Installation

### Requirements
 - Docker

Linux:
  ```
  curl -fsSL https://raw.githubusercontent.com/CyberL1/localapps/main/scripts/get.sh | sh
  ```

## Usage

1. Get [caddy](https://caddyserver.com)
2. Paste this config:
  ```
  apps.localhost, *.apps.localhost {
    reverse_proxy http://localhost:8080
  }
  ```

3. Do `localapps up`
4. Go to `apps.localhost` and you're good to go
