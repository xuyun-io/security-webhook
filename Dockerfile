FROM public.ecr.aws/docker/library/golang:1.20-bullseye AS build
COPY . /src
RUN cd /src & & go mod tidy
RUN go env -w CGO_ENABLED=0
RUN cd /src && go build -o security-webhook

FROM public.ecr.aws/docker/library/debian:bullseye as publish
WORKDIR /app
COPY --from=build src/security-webhook /app/security-webhook
ENV TZ=Asia/Tokyo
RUN echo $TZ > /etc/timezon

RUN chmod +x /app/security-webhook
ENTRYPOINT ["/app/security-webhook"]
