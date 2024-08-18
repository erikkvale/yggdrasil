FROM cgr.dev/chainguard/go as builder
WORKDIR /app
COPY . .
RUN go build -o batch-processing-api .

FROM cgr.dev/chainguard/glibc-dynamic
COPY --from=builder /app/batch-processing-api /usr/bin/
CMD ["/usr/bin/batch-processing-api"]
