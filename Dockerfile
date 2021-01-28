# Multi-stage builds
################################
# STEP 1 build executable binary
################################
FROM ubuntu:latest AS build

RUN apt-get update
RUN apt-get install -y wget git gcc g++ make bash

RUN wget -P /tmp https://dl.google.com/go/go1.15.1.linux-amd64.tar.gz

RUN tar -C /usr/local -xzf /tmp/go1.15.1.linux-amd64.tar.gz
RUN rm /tmp/go1.15.1.linux-amd64.tar.gz

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

WORKDIR $GOPATH/src/combotest

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go env && pwd && ls
RUN go build -o /bin/listener -v ./cmd/listener

COPY ./dist/closessl /lib
ENV LD_LIBRARY_PATH=/lib:$LD_LIBRARY_PATH

EXPOSE 8080 4000


ENTRYPOINT ["/bin/listener"]

# ERROR
# standard_init_linux.go:219: exec user process caused: no such file or directory

################################
# STEP 2 build a small image
############################
# FROM scratch AS final

# COPY --from=build /bin/listener /bin/listener
# COPY --from=build /lib /lib
# ENV LD_LIBRARY_PATH=/lib:$LD_LIBRARY_PATH

# EXPOSE 8080
# EXPOSE 4000

# Run the executable
# CMD ["/bin/listener"]