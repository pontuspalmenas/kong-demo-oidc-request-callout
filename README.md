# Kong Local Template

1. Create a Control Plane in Konnect
2. Create a `.env` file with your CP prefix:
```
KONG_CLUSTER_PREFIX=abc123def1
```
3. Create a new DP in Konnect to get the certificate.
4. Put certificate and key in `certs/tls.crt` and `tls.key` respectively.
5. Start Kong GW and all backends with `task compose`