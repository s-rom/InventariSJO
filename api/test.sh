#!/bin/bash

INTEGRATION=0
for arg in "$@"; do
    [[ "$arg" == "--integration" || "$arg" == "-i" ]] && INTEGRATION=1
done

if [[ $INTEGRATION -eq 1 ]]; then
    # Detectar socket de Podman (Fedora) o Docker
    systemctl --user start podman.socket 2>/dev/null
    if [[ -S "/run/user/$(id -u)/podman/podman.sock" ]]; then
        export DOCKER_HOST="unix:///run/user/$(id -u)/podman/podman.sock"
        export TESTCONTAINERS_RYUK_DISABLED=true
        echo "Usando Podman: $DOCKER_HOST"
    elif [[ -S "/var/run/docker.sock" ]]; then
        export DOCKER_HOST="unix:///var/run/docker.sock"
        echo "Usando Docker: $DOCKER_HOST"
    else
        echo "ERROR: No se encontró socket de Docker/Podman" >&2
        exit 1
    fi
    echo "=== Tests unitarios ==="
    go test ./internal/... -v
    echo ""
    echo "=== Tests de integración ==="
    go test ./tests/integration/... -tags integration -v
else
    go test ./... -v
fi


if [[ $? -ne 0 ]]; then
    echo "ERROR: Tests fallidos" >&2
    exit 1
else 
    echo "Todos los tests ok!"
fi
