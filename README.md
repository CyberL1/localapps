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

1. Create `access-url.txt` file using `echo "http://domain.tld" > ~/.config/localapps/access-url.txt`
2. Do `localapps up`
3. Go to the url you set to access localapps
