# Golang Open ID Connect Demo Application

You can configure it using the following environment variables:

```
TOKEN_URL
AUTH_URL
CLIENT_ID
CLIENT_SECRET
REDIRECT_URI
CODE_EXCHANGE_URL
USER_INFO_URL
PORT
```

Example running with Docker:

```
docker run --rm \
-p 8888:8888 \
-e TOKEN_URL=https://pingfederate:9031/token \
-e AUTH_URL=https://localhost:9031/as/authorization.oauth2 \
-e CLIENT_ID=test-client \
-e CLIENT_SECRET=test-secret \
-e REDIRECT_URI=http://localhost:8888/oidc-app/redirect \
-e CODE_EXCHANGE_URL=https://localhost:9031/as/token.oauth2 \
-e USER_INFO_URL=https://localhost:9031/idp/userinfo.openid \
-e PORT=8888 \
briansimoni/oidc-app
```
