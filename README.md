# Authentication Service

## Generate Certificates (Public and Private RSA keys)

Certificates for access tokens (JWT tokens used for API calls):

```
openssl genrsa \
    -passout pass:12345678 \
    -out cert/private_key.pem \
    2048
    
openssl rsa \
    -passin pass:12345678 \
    -in cert/private_key.pem \
    -pubout > cert/public_key.pub
```

For the refresh tokens (JWT tokens used to issue new access tokens):

```
openssl genrsa \
    -passout pass:12345678 \
    -out cert/refresh_private_key.pem \
    2048
    
openssl rsa \
    -passin pass:12345678 \
    -in cert/refresh_private_key.pem \
    -pubout > cert/refresh_public_key.pub
```

## Using JWK

Here is a [good starting point](https://www.baeldung.com/openssl-self-signed-cert) to understand what is next. On the other hand, it is a good idea checking RFC standards:

* [RFC 528: Internet X.509 Public Key Infrastructure Certificate and Certificate Revocation List (CRL) Profile](https://datatracker.ietf.org/doc/html/rfc528)
* [RFC 7519: JSON Web Token (JWT)](https://datatracker.ietf.org/doc/html/rfc7519)
* [RFC 7517: JSON Web Key (JWK)](https://datatracker.ietf.org/doc/html/rfc7517)
* [JWK parameters for RSA private keys](https://tools.ietf.org/id/draft-jones-jose-json-private-and-symmetric-key-00.html#rfc.section.3.2)

We can also run the following command to get info about key (prime factor 1, prime factor 2, modulos, quotient, etc; Which MAY be used for JWKs processing):

```
openssl rsa -pubin \
    -in cert/public_key.pub \
    -text -noout 
```

I have also implemented a script utility that will help you get the modulo (Assuming the exponent is `65537`):

```
scripts/jwk_modulus_base64url --infile ./cert/public_key.pub
```

Then you can use the output in your `JWK` public key modulo property (`n`). Note that padding was removed in the output string.

### Certificate signing request

```
openssl req \
    -key cert/private_key.pem \
    -new -out cert/auth_service.csr
```

### Self Signed Certificate (x509)

This certificate is the one to be used in JWKs certificates chain (e.g. `x5c` attribute):

```
openssl x509 \
    -signkey cert/private_key.pem \
    -in cert/auth_service.csr \
    -req -days 365 \
    -out cert/auth_service.crt
```

To view the certificate:

```
openssl x509 \
    -text -noout \
    -in cert/auth_service.crt
```

## Running with Docker

You can build a docker image:

```
docker build -t auth-service .
```

Assuming you have a mongo server and a redis server in your local, you can do:

```
docker run -e MONGO_URI="mongodb://host.docker.internal/auth" -e REDIS_URL="redis://host.docker.internal:6379/1" -p 4000:4000 auth-service
```

The service will be listening on port 4000 in this case.

## Running with docker-compose

You can also run the service without having to install and set mongo and redis. For development purposes, to execute the service you can run:

```
docker-compose up --build
```

## Running on minikube

```
# Install dependencies (charts)
cd helm && helm dependency build && cd -

# Run application
helm \
    --kube-context=minikube \
    --create-namespace \
    --namespace=auth-svc \
    upgrade --install --force --recreate-pods \
    auth-svc \
    helm
```

I also added `dnsutils.yaml` for debugging purposes, as I was stuck with some developments (`redis` endpoint as an example and other DNS mappings). As an example, this command was useful:

```
kubectl exec -i -t dnsutils -n auth-svc -- nslookup 10.108.136.235
```