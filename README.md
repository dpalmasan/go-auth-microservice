# Authentication Service


## Generate Certificates (Public and Private RSA keys)

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

## To use JWK

Here is a [good starting point](https://www.baeldung.com/openssl-self-signed-cert) to understand what is next. On the other hand, it is a good idea checking RFC standards:

* [RFC 528: Internet X.509 Public Key Infrastructure Certificate and Certificate Revocation List (CRL) Profile](https://datatracker.ietf.org/doc/html/rfc528)
* [RFC 7519: JSON Web Token (JWT)](https://datatracker.ietf.org/doc/html/rfc7519)
* [RFC 7517: JSON Web Key (JWK)](https://datatracker.ietf.org/doc/html/rfc7517)
* [JWK parameters for RSA private keys](https://tools.ietf.org/id/draft-jones-jose-json-private-and-symmetric-key-00.html#rfc.section.3.2)

We can also run the following command to get info about key (prime factor 1, prime factor 2, modulos, quotient, etc; Which MAY be used for JWKs processing):

```
openssl rsa \
    -noout -text \
    -inform PEM \
    -in cert/private_key.pem
```

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
    -req -days 365 
    \-out cert/auth_service.crt
```

To view the certificate:

```
openssl x509 \
    -text -noout \
    -in domain.crt
```