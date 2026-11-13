# Verificador de URLs en Go

Código de ejemplo del artículo [«Construye un verificador de URLs en Go: concurrencia en la práctica»](https://juredev.com/blog/2026/11/verificador-urls-go-concurrencia-practica/).

Construye paso a paso un verificador de URLs concurrente en unas 130 líneas usando únicamente la librería estándar de Go (`goroutines`, `channels`, `sync.WaitGroup` y `context`).

## Requisitos

- Go 1.22 o superior.
- Conexión a internet para `sitios.txt` y `testdata/sitios-reales.txt` (la prueba con el servidor local no la necesita).

## Estructura del repositorio

```text
verificador-go/
├── .gitignore
├── README.md
├── go.mod                  # Módulo único (go 1.22)
├── main.go                 # Versión final (idéntica a pasos/paso3)
├── sitios.txt              # Lista de URLs de ejemplo
├── probar.sh               # Script para comparar los 3 pasos
├── pasos/
│   ├── paso1/main.go       # Paso 1: versión secuencial
│   ├── paso2/main.go       # Paso 2: una goroutine por sitio
│   └── paso3/main.go       # Paso 3: worker pool + context con timeout
├── herramientas/
│   └── servidor/main.go    # Servidor HTTP local para pruebas repetibles
└── testdata/
    ├── sitios-locales.txt  # 10 URLs apuntando al servidor local
    └── sitios-reales.txt   # 9 URLs públicas para pruebas en vivo
```

## Uso rápido (versión final)

Ejecutar directamente desde la raíz con `sitios.txt`:

```bash
go run . -archivo sitios.txt -workers 5 -timeout 3s
```

Compilar un binario:

```bash
go build -o verificador                   # binario para tu sistema
CGO_ENABLED=0 go build -o verificador     # binario estático, sin dependencias
./verificador -archivo sitios.txt -workers 5 -timeout 3s
```

Para un servidor Linux, desde cualquier sistema operativo:

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o verificador-linux
```

El programa termina con código de salida 1 si algún sitio falla, así que puedes usarlo en scripts o tareas programadas.

## Ejecutar cada paso por separado

```bash
# Paso 1: Secuencial
go run ./pasos/paso1 sitios.txt

# Paso 2: Una goroutine por sitio
go run ./pasos/paso2 sitios.txt

# Paso 3: Worker pool y context
go run ./pasos/paso3 -archivo sitios.txt -workers 5 -timeout 3s
```

Con `sitios.txt`, el paso 1 tarda unos 7–8 s y el paso 2 unos 6 s, porque `httpbin.org/delay/6` no tiene límite de tiempo hasta el paso 3. Es lo esperado y es justo lo que el artículo quiere mostrar.

## Reproducir las mediciones del artículo

El script `./probar.sh` ejecuta los tres pasos seguidos para comparar sus tiempos. Sin argumentos, levanta antes el servidor local de pruebas en `127.0.0.1:8099`:

```bash
# Contra el servidor local simulado (~16s -> ~10s -> ~4s)
./probar.sh

# Contra sitios reales de prueba (testdata/sitios-reales.txt)
./probar.sh reales
```

Los tiempos con sitios reales dependen de tu red. httpbin.org, además, a veces responde más lento de lo normal, así que tus números pueden variar.
