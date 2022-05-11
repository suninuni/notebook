# Registry

## Deploy

### http

```shell
docker run --name registry -p 80:80 -d --restart=always -v /mnt/registry:/var/lib/registry \
    -e REGISTRY_HTTP_ADDR=0.0.0.0:80 \
    registry:2
```

```json daemon.json
{
  "insecure-registries" : ["local.com"]
}
```

### https

```shell
docker run --name registry -p 443:443 -d --restart=always -v /mnt/registry:/var/lib/registry \
    -v /home/user/soft/docker/certs:/certs \
    -e REGISTRY_HTTP_ADDR=0.0.0.0:443 \
    -e REGISTRY_HTTP_TLS_CERTIFICATE=/certs/domain.pem \
    -e REGISTRY_HTTP_TLS_KEY=/certs/domain.key \
    registry:2
```
