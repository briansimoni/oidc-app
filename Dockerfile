FROM golang

ADD . /go/src/github.build.ge.com

WORKDIR /go/src/github.build.ge.com

RUN go build -o oidc-app

EXPOSE 4321

CMD ./oidc-app