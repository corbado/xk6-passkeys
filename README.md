# xk6-passkeys

A k6 extension for testing passkey backends ("WebAuthn servers"). With this extension, you can load test your passkey registration and login flows.

## Build

For your convenience, we provide a Makefile to build the extension:

```bash
make build
```

This will create a `k6` binary in the current directory with the extension compiled in.

## Examples

We have implemented two example load tests in the `examples` directory: one for registration and one for login. Since a passkeys backend is required for testing, we provide a sample backend for load testing. To start the example, execute the following commands:

```bash
docker build -t passkeys-backend examples/backend
docker run -p 8080:8080 passkeys-backend
```

To load test the registration flow, run the following command:

```bash
k6 run examples/registration.js
```

To load test the login flow, run the following command:

```bash
k6 run examples/login.js
```