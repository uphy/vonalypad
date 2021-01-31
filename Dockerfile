FROM golang:1.13 AS backend-builder

WORKDIR /work
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /vonalypad

FROM node:15-alpine AS frontend-builder

WORKDIR /work
COPY frontend/package.json .
COPY frontend/package-lock.json .
RUN npm install
COPY frontend/ ./
RUN npm run build

FROM alpine:3.12.3

# RUN apk --no-cache add tzdata

WORKDIR /vonalypad/
COPY --from=backend-builder /vonalypad .
RUN chmod +x vonalypad
COPY --from=frontend-builder /work/dist static
RUN ls -al .
RUN mkdir /recipe

ENTRYPOINT [ "/vonalypad/vonalypad", "--recipe-dir", "/recipe", "start", "--static-dir", "/vonalypad/static" ]