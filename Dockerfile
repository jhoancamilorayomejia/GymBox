# ── Etapa 1: Build de Vue ─────────────────────────────────────────
FROM node:22-alpine AS vue-build

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm install

COPY index.html vite.config.js eslint.config.js ./
COPY src/ ./src/
COPY public/ ./public/

ARG VITE_MP_PUBLIC_KEY
ENV VITE_MP_PUBLIC_KEY=$VITE_MP_PUBLIC_KEY

RUN npm run build

# ── Etapa 2: Build de Go ──────────────────────────────────────────
FROM golang:alpine AS go-build

WORKDIR /app

ENV GOTOOLCHAIN=off

COPY go.mod go.sum ./
RUN go mod download

COPY main.go ./
COPY controllers/ ./controllers/
COPY db/ ./db/
COPY models/ ./models/

RUN go build -o server .

# ── Etapa 3: Imagen final ─────────────────────────────────────────
FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY --from=go-build /app/server .
COPY --from=vue-build /app/dist ./dist

ENV PORT=8080
EXPOSE 8080

CMD ["./server"]