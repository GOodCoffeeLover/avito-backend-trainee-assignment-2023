FROM golang:1.20 as base

WORKDIR /app

COPY go.mod go.sum ./

RUN ["go", "mod", "download"]

COPY . ./

ENV CGO_ENABLED=0

RUN  ["go", "build", "-o", "api", "./cmd/api/main.go"]

# ---- 
FROM scratch

COPY --from=base /app/api .

EXPOSE 7001

CMD ["./api"]