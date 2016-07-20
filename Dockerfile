FROM golang:1.6-alpine

ARG git_commit=unknown
LABEL org.cyverse.git-ref="$git_commit"

COPY . /go/src/github.com/cyverse-de/permissions
RUN apk update
RUN apk add git
RUN go get github.com/constabulary/gb/...
RUN cd /go/src/github.com/cyverse-de/permissions && \
	gb build && \
	cp bin/permissions-server /bin/permissions

EXPOSE 60000
ENTRYPOINT ["permissions"]
CMD ["--help"]
