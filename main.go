package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type Resultado struct {
	URL      string
	Estado   int
	Duracion time.Duration
	Err      error
}

func (r Resultado) OK() bool {
	return r.Err == nil && r.Estado < 400
}

func verificar(ctx context.Context, url string) Resultado {
	inicio := time.Now()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Resultado{URL: url, Err: err}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Resultado{URL: url, Duracion: time.Since(inicio), Err: err}
	}
	defer resp.Body.Close()

	return Resultado{URL: url, Estado: resp.StatusCode, Duracion: time.Since(inicio)}
}

func worker(urls <-chan string, resultados chan<- Resultado, timeout time.Duration) {
	for url := range urls {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		resultados <- verificar(ctx, url)
		cancel()
	}
}

func leerURLs(ruta string) ([]string, error) {
	f, err := os.Open(ruta)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var urls []string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		if linea := strings.TrimSpace(sc.Text()); linea != "" {
			urls = append(urls, linea)
		}
	}
	return urls, sc.Err()
}

func main() {
	archivo := flag.String("archivo", "sitios.txt", "archivo con una URL por línea")
	workers := flag.Int("workers", 5, "peticiones simultáneas")
	timeout := flag.Duration("timeout", 3*time.Second, "tiempo máximo por petición")
	flag.Parse()

	if *workers < 1 {
		fmt.Fprintln(os.Stderr, "error: -workers debe ser al menos 1")
		os.Exit(2)
	}

	urls, err := leerURLs(*archivo)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	inicio := time.Now()
	cola := make(chan string)
	resultados := make(chan Resultado)

	var wg sync.WaitGroup
	for i := 0; i < *workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker(cola, resultados, *timeout)
		}()
	}

	// Alimentamos la cola y cerramos los canales en su momento.
	go func() {
		for _, u := range urls {
			cola <- u
		}
		close(cola)
	}()
	go func() {
		wg.Wait()
		close(resultados)
	}()

	fallos := 0
	for r := range resultados {
		if r.OK() {
			fmt.Printf("OK     %d  %6v  %s\n", r.Estado, r.Duracion.Round(time.Millisecond), r.URL)
			continue
		}
		fallos++
		if r.Err != nil {
			fmt.Printf("FALLO  ---  %6v  %s  (%v)\n", r.Duracion.Round(time.Millisecond), r.URL, r.Err)
		} else {
			fmt.Printf("FALLO  %d  %6v  %s\n", r.Estado, r.Duracion.Round(time.Millisecond), r.URL)
		}
	}

	fmt.Printf("\n%d sitios revisados en %v: %d OK, %d con fallo\n",
		len(urls), time.Since(inicio).Round(time.Millisecond), len(urls)-fallos, fallos)

	if fallos > 0 {
		os.Exit(1)
	}
}
