FROM golang:1.22 AS builder

WORKDIR /workspace
COPY . /workspace

# go-sqlite3 requires CGO
RUN CGO_ENABLED=1 go build -o /bin/quizmaker .

FROM busybox
COPY --from=builder /bin/quizmaker /bin/quizmaker

WORKDIR /app
COPY views/ /app/views

CMD ["quizmaker"]
