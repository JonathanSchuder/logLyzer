# logLyzer

Ein Go-Projekt zum Lernen von **Goroutines** und **Channels** anhand eines echten Anwendungsfalls: Das Analysieren einer großen Apache/Nginx Access-Log-Datei.

## Was macht das Programm?

Liest eine `access.log` Datei und zählt wie oft jeder HTTP-Statuscode vorkommt.

**Beispiel-Output:**
```
map[200:161984 301:40704 401:38208 404:37760 500:41344]
```

## Architektur - Producer/Consumer Pattern

```
┌─────────────────┐        ┌──────────────────┐
│ Scanner-Goroutine│        │  Worker 1        │
│ (Producer)       │──────▶│  Worker 2        │
│                  │channel │  Worker 3        │
│ liest Zeilen &   │       │  Worker 4        │
│ schickt sie rein │        │  Worker 5        │
└─────────────────┘        └──────────────────┘
                                    │
                              eigene Maps
                                    │
                             ┌──────▼──────┐
                             │   Merge     │
                             │  results    │
                             └─────────────┘
```

## Gelernte Konzepte

### Channel
Eine typisierte Pipe zwischen Goroutines. Wenn der Channel voll ist wartet der Sender, wenn er leer ist wartet der Empfänger.

```go
lines := make(chan string, 100) // buffered channel mit Kapazität 100

// Senden:
lines <- "eine zeile"

// Empfangen:
line := <-lines

// Über Channel iterieren (blockiert bis close()):
for line := range lines { ... }
```

### Goroutine
Eine leichtgewichtige Funktion die nebenläufig läuft. Einfach `go` davor schreiben.

```go
go func() {
    // läuft parallel zu main()
}()
```

### WaitGroup
Ein Zähler um auf das Ende mehrerer Goroutines zu warten.

```go
var wg sync.WaitGroup

wg.Add(1)           // Zähler +1 (vor dem Start der Goroutine!)
go func() {
    defer wg.Done() // Zähler -1 wenn Funktion endet
}()

wg.Wait()           // blockiert bis Zähler = 0
```

### Mutex (hier nicht mehr nötig, aber gelernt)
Schützt shared State vor gleichzeitigem Zugriff mehrerer Goroutines.

```go
var mu sync.Mutex

mu.Lock()
sharedMap[key]++  // nur eine Goroutine gleichzeitig
mu.Unlock()
```

**Besser:** Jede Goroutine arbeitet auf ihrer eigenen lokalen Map → kein Locking nötig → bessere Performance.

### defer-Falle in Loops
`defer` läuft erst wenn die **Funktion** endet, nicht wenn die Loop-Iteration endet!

```go
// FALSCH - Deadlock!
for line := range lines {
    mu.Lock()
    defer mu.Unlock() // läuft erst nach dem kompletten Loop
}

// RICHTIG
for line := range lines {
    mu.Lock()
    results[code]++
    mu.Unlock()
}
```

## Warum kein `os.ReadFile`?

`os.ReadFile` lädt die gesamte Datei in den RAM. Bei großen Dateien ist das ineffizient.
`bufio.Scanner` liest zeilenweise mit einem internen Buffer - konstanter Speicherverbrauch egal wie groß die Datei ist.

## Ausführen

```bash
go run main.go
```
