# Utiliser une image de base légère pour Go
FROM golang:1.22-alpine

# Définir le répertoire de travail
WORKDIR /app

# Copier les fichiers source
COPY . .

# Copier uniquement les fichiers go.mod et go.sum pour éviter de tout re-télécharger à chaque build
COPY go.mod go.sum ./
RUN go mod download

# Compiler l’application
RUN go build -o main .

# Exposer le port utilisé par l’API
EXPOSE 8080

# Lancer l’application
CMD ["./main"]
