# ── Etapa 1: Build de Vue ─────────────────────────────────────────
FROM node:22-alpine AS vue-build

WORKDIR /app

# Copiar archivos de dependencias
COPY package.json package-lock.json ./

# Instalar dependencias
RUN npm install

# Copiar el resto del código fuente de Vue
COPY index.html vite.config.js eslint.config.js ./
COPY src/ ./src/
COPY public/ ./public/

# ✅ Recibir la clave pública de Mercado Pago en tiempo de build
ARG VITE_MP_PUBLIC_KEY
ENV VITE_MP_PUBLIC_KEY=$VITE_MP_PUBLIC_KEY

# Compilar Vue para producción
RUN npm run build

# ── Etapa 2: Build de Go ──────────────────────────────────────────
FROM golang:1.25rc2-alpine AS go-build

WORKDIR /app

# Copiar módulos de Go
COPY go.mod go.sum ./
RUN go mod download

# Copiar código fuente de Go
COPY main.go ./
COPY controllers/ ./controllers/
COPY db/ ./db/
COPY models/ ./models/

# Compilar el binario
RUN go build -o server .

# ── Etapa 3: Imagen final ─────────────────────────────────────────
FROM alpine:latest

WORKDIR /app

# Certificados SSL para conexiones HTTPS (PostgreSQL)
RUN apk --no-cache add ca-certificates

# Copiar binario de Go
COPY --from=go-build /app/server .

# Copiar el dist de Vue generado en etapa 1
COPY --from=vue-build /app/dist ./dist

# Puerto dinámico de Railway
ENV PORT=8080
EXPOSE 8080

CMD ["./server"]