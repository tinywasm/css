# PLAN: Token único con light-dark — API mínima

## Problema

`Token.Fallback` guarda `"light-dark(#FFFFFF, #0D1117)"` como string opaco. No se puede acceder a `.Light` / `.Dark` desde Go.

## Solución: 3 campos, sin duplicación

### Tipo

```go
type Token struct {
    Name       string
    Light, Dark string   // Dark solo = static fallback; ambos = light-dark pair
}
```

- **Static**: `Light == ""` → fallback es `Dark`
- **Theme**: `Light != ""` → fallback es `light-dark(Light, Dark)`

### Métodos

```go
func (t Token) Var() string {
    if t.Light == "" {
        return "var(" + t.Name + "," + t.Dark + ")"
    }
    return "var(" + t.Name + ",light-dark(" + t.Light + ", " + t.Dark + "))"
}

func (t Token) GetFallback() string {
    if t.Light == "" {
        return t.Dark
    }
    return "light-dark(" + t.Light + ", " + t.Dark + ")"
}
```

### SetTheme

```go
func SetTheme(t Token, light, dark string) Override {
    return Override{t, "light-dark(" + light + ", " + dark + ")"}
}
```

Uso:

```go
css.Theme(
    css.SetTheme(css.ColorBackground, "#FAFAFA", "#121212"),
)
```

## Uso

### Declaración — campos nombrados siempre

En Go no se puede inicializar un struct de 3 campos con solo 2 argumentos posicionales. Todos los tokens usan campos nombrados:

```go
// Static — solo Dark
ColorPrimary = Token{Name: "--color-primary", Dark: "#1b5d8c"}

// Theme — Light y Dark
ColorBackground = Token{Name: "--color-background", Light: "#FFFFFF", Dark: "#0D1117"}
```

### Introspección desde Go

```go
ColorBackground.Light  // "#FFFFFF"
ColorBackground.Dark   // "#0D1117"
```

### Pair — sin cambios

```go
type Pair struct{ Bg, Fg Token }
```

## Separación de archivos

| Archivo | Responsabilidad |
|---|---|
| `tokens.go` | Tipos (`Token`, `Pair`, `NamedPair`, `ValueGetter`), métodos (`Var`, `GetFallback`), `AllPairs()` |
| `catalog.go` | Solo el catálogo `var(...)` — añadir/quitar tokens no toca la lógica |
| `css.go` | `SetTheme()` |

## Cambios totales

| Archivo | Cambio |
|---|---|
| `tokens.go` | Token: `Fallback` → `Light`/`Dark`; se extrae el catálogo `var(...)` |
| `catalog.go` | **Nuevo** — catálogo `var(...)` con tokens static y theme |
| `css.go` | Se añade `SetTheme()` |
| `docs/MIGRATION.md` | Recipe A: `Set(ColorBackground, "light-dark(...)")` → `SetTheme(ColorBackground, "#FAFAFA", "#121212")`; tabla Theme/Canvas tiene columnas Light/Dark separadas |
