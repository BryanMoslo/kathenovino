FROM golang:1.19.3 as builder

# Installing nodejs and build-base
RUN apt-get update
RUN apt-get install -y nodejs npm curl bash build-essential git

# Installing Yarn
RUN curl -o- -L https://yarnpkg.com/install.sh | bash -s -- --version 1.22.10
ENV PATH="$PATH:/root/.yarn/bin:/root/.config/yarn/global/node_modules"

# Installing ox from pre-built binary
RUN wget https://github.com/wawandco/ox/releases/download/v0.13.1/ox_0.13.1_Linux_x86_64.tar.gz
RUN tar -xzf ox_0.13.1_Linux_x86_64.tar.gz
RUN mv ox /bin/ox

WORKDIR /kathenovino
ADD . .

# Building the application binary in bin/app 
RUN ox build --static -o bin/app

# Building bin/cli with the tooling
RUN go build -o ./bin/cli ./cmd/ox 

FROM debian:latest

# Binaries
COPY --from=builder /kathenovino/bin/app /bin/app
COPY --from=builder /kathenovino/bin/cli /bin/cli

# Binaries
COPY --from=builder /kathenovino/bin/* /bin/

ENV ADDR=0.0.0.0
EXPOSE 3000

# For migrations use
# CMD /bin/cli db migrate; /bin/app