FROM discoenv/golang-base:master

ENV CONF_TEMPLATE=/go/src/github.com/cyverse-de/permissions/permissions.yaml.tmpl
ENV CONF_FILENAME=permissions.yaml
ENV PROGRAM=permissions

COPY . /go/src/github.com/cyverse-de/permissions/

RUN go install github.com/cyverse-de/permissions/... \
    && cp /go/bin/permissions-server /bin/permissions

WORKDIR /
EXPOSE 60000

ARG git_commit=unknown
ARG version="2.9.0"
ARG descriptive_version=unknown

LABEL org.cyverse.git-ref="$git_commit"
LABEL org.cyverse.version="$version"
LABEL org.cyverse.descriptive-version="$descriptive_version"
LABEL org.label-schema.vcs-ref="$git_commit"
LABEL org.label-schema.vcs-url="https://github.com/cyverse-de/permissions"
LABEL org.label-schema.version="$descriptive_version"
f
