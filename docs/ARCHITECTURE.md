# Arquitectura de Theming

## RootCSS y el Slot de un solo ganador

En el pipeline SSR de `tinywasm`, `assetmin` descubre los bloques `RootCSS()` del grafo de dependencias. El bloque `:root` (el vocabulario de tokens) es un **slot de un solo ganador por reemplazo**:

1. Si la aplicación (el proyecto raíz) declara `func RootCSS() *css.Stylesheet`, ese bloque **reemplaza por completo** el `RootCSS()` por defecto de la librería `tinywasm/css`.
2. `RenderCSS()` (la lógica de reglas y bindings) es **aditivo**: se concatenan las contribuciones de todos los módulos.

## Entrypoint: `Theme()`

Para facilitar el rebrand sin perder el catálogo completo de tokens (escalas de espacio, tipografía, etc.), la librería provee `css.Theme(...Override)`.

`Theme()` obtiene el catálogo por defecto mediante `RootCSS()` y añade un bloque `:root` final con los overrides proporcionados. Esto garantiza que:

- La aplicación no tenga que redeclarar tokens que no desea cambiar.
- Los cambios de la aplicación ganen en la cascada de CSS al aparecer al final del bloque.
- El catálogo se mantenga íntegro para que los componentes sigan funcionando.

## Type-Safety con `Override`

El tipo `Override` es opaco y solo puede construirse mediante `css.Set(Token, value)`. Esto impide estados ilegales como intentar inyectar propiedades CSS arbitrarias en el bloque `:root` a través del entrypoint de tema, forzando a que solo se sobreescriban tokens del catálogo tipado.

```go
// Uso correcto
css.Theme(css.Set(css.ColorPrimary, "#hex"))

// No compila (Set exige un Token)
// css.Theme(css.Set("padding", "20px"))
```
