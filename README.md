# Certificate Status Page

This go WebApp is designed to be simple.

It will ONLY serve as a Status Page for the TLS Certificates of the tracked Endpoints.

It WILL:

- Show an Overview of All Configured Certificates with Expiry
- Show Details of Certificate Chain
- Show Domain Resolution and Address

It will NOT:

- Alert on Certificate Expiry
- Expose Metrics Externally - either by API or Prometheus as to Certifcate State
- Show History of Certificates

## Configuration

Endpoints can be configured with a YAML file (default `config/config.yml`, override with `CERTUI_CONFIG_PATH`):

```yaml
endpoints:
  - https://www.google.com
  - https://www.github.com
```

Alternatively, set the `ENDPOINTS` environment variable to a comma-separated list. When set, it takes precedence and no config file is needed:

```sh
ENDPOINTS=https://www.google.com,https://www.github.com
```
