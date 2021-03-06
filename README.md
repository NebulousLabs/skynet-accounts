# skynet-accounts

`skynet-accounts` is a service that stores [Skynet](https://siasky.net) user account data. It uses MongoDB for data
storage. It uses ORY Kratos for the actual account management.

## Setup steps

### The `.env` file

All local secrets are loaded from a `.env` file in the root directory of the project.

Those are (example values):

```.env
SKYNET_DB_HOST="localhost"
SKYNET_DB_PORT="27017"
SKYNET_DB_USER="username"
SKYNET_DB_PASS="password"
SKYNET_ACCOUNTS_PORT=3000
COOKIE_DOMAIN="siasky.net"
COOKIE_HASH_KEY="any thirty-two byte string is ok"
COOKIE_ENC_KEY="any thirty-two byte string is ok"
```

There are some optional ones, as well:
```.env
SKYNET_ACCOUNTS_LOG_LEVEL=trace
KRATOS_ADDR=localhost:4433
OATHKEEPER_ADDR=localhost:4456
```

## Recommended reading

- [JSON and BSON](https://www.mongodb.com/json-and-bson)
- [Using the official MongoDB Go driver](https://vkt.sh/go-mongodb-driver-cookbook/)
