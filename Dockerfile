FROM public.ecr.aws/docker/library/golang:1.21-bullseye AS build
COPY . /src
RUN cd /src && go mod tidy
RUN go env -w CGO_ENABLED=0
RUN cd /src && go build -a -buildvcs=false -o security-webhook

FROM public.ecr.aws/docker/library/debian:bullseye as publish

RUN groupadd -r security-webhook -g 1000 && useradd -r -u 1000 -g security-webhook security-webhook

WORKDIR /app

COPY --from=build src/security-webhook /app/security-webhook

RUN chown -R security-webhook:security-webhook /app

ENV TZ=Asia/Tokyo
RUN echo $TZ > /etc/timezon
RUN chmod +x /app/security-webhook

USER security-webhook
ENTRYPOINT ["/app/security-webhook"]
