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

### Locally

1. Do `localapps up`
2. Go to `http://localhost:8080` and you're good to go

### Remote (VPS)

1. Login to your vps
2. Create data directory for localapps using:
  ```bash
  mkdir -p .config/localapps
  ```

3. Create `access-url.txt` file inside it using:
  ```bash
  echo "http://example.com:8080" > .config/localapps/access-url.txt
  ```

2. Start localapps using:
  ```bash
  docker run -d --name localapps-server -v /var/run/docker.sock:/var/run/docker.sock -v .config/localapps:/root/.config/localapps -p 8080:8080 ghcr.io/cyberl1/localapps
  ```

3. Go to the url you set to access localapps
