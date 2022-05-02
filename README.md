# DEV-SERVER

A simple dev server.

## Install

+ download [releases](https://github.com/TMaize/dev-server/releases)

+ install

    ```shell
    # ${GOPATH}/bin
    go install github.com/TMaize/dev-server@latest
    ```

+ build

    ```shell
    git clone https://github.com/TMaize/dev-server.git
    cd dev-server
    make
    ```

## Feature

- [x] static server

- [x] https support

- [ ] reverse proxy, cors

- [ ] [trust Root CA](https://github.com/FiloSottile/mkcert/blob/master/truststore_windows.go)

- [ ] mock data

## Usage

```
dev-server start .
dev-server start -p 8443 --https --domain test.com /site/demo 
```

## Acknowledgments

[spf13/cobra](https://github.com/spf13/cobra)

[https://github.com/kubernetes/client-go/blob/master/util/cert/cert.go](https://github.com/kubernetes/client-go/blob/master/util/cert/cert.go#L84)

[https://github.com/vitejs/vite/blob/main/packages/vite/src/node/certificate.ts](https://github.com/vitejs/vite/blob/main/packages/vite/src/node/certificate.ts)
