# Cyberrange

Automated virtual machine deployment solution.

## Deployment

There are two services, `webserver` and `manager`. Webserver deploys a website with a MySQL/MariaDB backend. It makes calls to the VM Manager. The VM Manager makes calls to the oVirt API to manage (create, delete, snapshot, start, etc.) the VMs.

### Webserver

1. `git clone https://github.com/CUCyber/cyberrange.git`
2. Update the configuration information in `services/webserver/db/config.yaml`.
3. `make webserver`

### Manager

1. `git clone https://github.com/CUCyber/cyberrange.git`
2. Update the configuration information in `services/manager/git/config.yaml` and `services/manager/ovirt/config.yaml`.
3. `make manager`

## Development

1. `git clone https://github.com/CUCyber/cyberrange.git`
2. Make your changes
3. Submit a PR
