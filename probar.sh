#!/usr/bin/env bash
# Ejecuta los tres pasos del artículo contra el servidor local de pruebas.
# Uso: ./probar.sh          (servidor local, testdata/sitios-locales.txt)
#      ./probar.sh reales   (los 3 pasos contra testdata/sitios-reales.txt)
set -u
cd "$(dirname "$0")"

if ! diff -q main.go pasos/paso3/main.go >/dev/null 2>&1; then
	echo "AVISO: main.go y pasos/paso3/main.go no son idénticos." >&2
fi

if [ "${1:-}" = "reales" ]; then
	echo "=== Paso 1: secuencial contra sitios reales"
	go run ./pasos/paso1 testdata/sitios-reales.txt

	echo; echo "=== Paso 2: una goroutine por sitio contra sitios reales"
	go run ./pasos/paso2 testdata/sitios-reales.txt

	echo; echo "=== Paso 3: worker pool + context (5 workers, timeout 5s)"
	go run ./pasos/paso3 -archivo testdata/sitios-reales.txt -workers 5 -timeout 5s
	echo "código de salida: $?"
	exit 0
fi

echo "Compilando servidor de pruebas..."
go build -o servidor-bin ./herramientas/servidor || exit 1
./servidor-bin &
SRV=$!
trap 'kill $SRV 2>/dev/null; rm -f servidor-bin' EXIT
sleep 1

echo; echo "=== Paso 1: secuencial (esperado ~16 s)"
go run ./pasos/paso1 testdata/sitios-locales.txt

echo; echo "=== Paso 2: una goroutine por sitio (esperado ~10 s)"
go run ./pasos/paso2 testdata/sitios-locales.txt

echo; echo "=== Paso 3: worker pool + context, 3 workers, timeout 2s (esperado ~4 s)"
go run ./pasos/paso3 -archivo testdata/sitios-locales.txt -workers 3 -timeout 2s
echo "código de salida: $? (1 = hubo fallos)"
