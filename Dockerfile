FROM golang:1.16-alpine as builder

ENV GO111MODULE on
ENV GOPROXY https://proxy.golang.org/

# Installing nodejs
RUN apk add --update nodejs curl bash build-base python2

# Installing Yarn
RUN curl -o- -L https://yarnpkg.com/install.sh | bash
ENV PATH="$PATH:/root/.yarn/bin:/root/.config/yarn/global/node_modules"

# Installing ox
RUN go install github.com/wawandco/ox/cmd/ox@v0.11.2

WORKDIR /kathenovivno
ADD . .

# Building the application binary in bin/app 
RUN ox build --static -o bin/app --tags timetzdata

# Building bin/cli with the tooling
RUN go build -o ./bin/cli -ldflags '-linkmode external -extldflags "-static"' ./cmd/ox 

FROM alpine

# Binaries
COPY --from=builder /kathenovivno/bin/app /bin/
COPY --from=builder /kathenovivno/bin/cli /bin/

# For migrations use 
# CMD cli db migrate up; app 
CMD cli db migrate up; app

