# GO-Zabbix

Para integrar el modulo con el backend GO hay que añadir el archivo

`zabbix.go`

En la carpeta vizzeb del backend a integrar.

En el fichero cmd/main.go hay que añadir la linea 

`go vizzeb.StartZabbix()`

Dentro de la funcion main() idealmente despues de inicializar los modulos.

Para usar los servicios de LOG se llaman de la siguiente manera:

Para logs de utilidad:
`
go vizzeb.ZabbixErrorLog(msg)
`
Para logs criticos:
`
go vizzeb.ZabbixPanicLog(msg)
`

Ejemplos de uso: 

```
if res.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("Content Manager status: %v", res.StatusCode)
		go vizeb.ZabbixErrorLog(msg)
		return err
	}
```

```
func prepareDB() {
	str := cfg.GetString("db")

	d, err := dbx.MustOpen("oci8", str)
	if err != nil {
		msg := fmt.Sprintf("Failed to open db: %v", err)
		vizeb.ZabbixPanicLog(msg)
		log.Fatalf("Failed to open db: %v", err)
	}

	if env == EnvDev {
		d.LogFunc = lg.Sugar().Debugf //log.Printf
	}

	go pingForStayConnected(d)

	db = *d

	lg.Debug("db is ok")
}
```
