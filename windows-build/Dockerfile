ARG GO_VERSION=latest
FROM x1unix/go-mingw:${GO_VERSION}

COPY ./entrypoint.sh /entrypoint.sh

ENTRYPOINT [ "/entrypoint.sh" ]
