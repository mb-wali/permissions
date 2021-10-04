FROM golang:1.16-alpine

RUN apk add --no-cache git
RUN go get -u github.com/jstemmer/go-junit-report

COPY . /permissions/
WORKDIR /permissions/
RUN go install ./... \
    && cp /go/bin/permissions-server /bin/permissions

WORKDIR /

# copy config file 
COPY permissions.yaml /etc/iplant/de/permissions.yaml

ENTRYPOINT ["permissions", "--host", "0.0.0.0", "--port", "60005"]
# CMD ["--help"]

EXPOSE 60005

ARG git_commit=unknown
ARG version="2.9.0"
ARG descriptive_version=unknown

LABEL org.cyverse.git-ref="$git_commit"
LABEL org.cyverse.version="$version"
LABEL org.cyverse.descriptive-version="$descriptive_version"
LABEL org.label-schema.vcs-ref="$git_commit"
LABEL org.label-schema.vcs-url="https://github.com/cyverse-de/permissions"
LABEL org.label-schema.version="$descriptive_version"

# build
# docker build -t mbwali/permissions:latest .

# run 
# docker rum -it -p 60005:60005 mbwali/permissions:latest

# config
# /etc/iplant/de/permissions.yaml
