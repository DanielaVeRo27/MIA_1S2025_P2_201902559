package comandos

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

/*
   ------> RMDISK comando que se encarga de eliminar un disco duro virtual del sistema

*/

// RMDISK Estructura del comando RMDISK
type rmdisk struct {
	path string
}

// rmdisk -path="/home/mis discos/Disco4.mia"
func ParserRmdisk(tokens []string) (string, error) {

	comando := &rmdisk{}             // instancia de comando rmdisk csu est
	prm := strings.Join(tokens, " ") // une los parametros en un string

	expr := regexp.MustCompile(`-path="[^"]+"|-path=[^\s]+`) // expresion del cmd

	// busca que hagan match
	matches := expr.FindAllString(prm, -1)

	if len(matches) != len(tokens) {
		for _, token := range tokens {
			if !expr.MatchString(token) {
				return "", fmt.Errorf("parámetro inválido: %s", token)
			}
		}
	}

	// recorre los matches
	for _, match := range matches {
		keyValue := strings.SplitN(match, "=", 2) // divide en clave y valor
		if len(keyValue) != 2 {
			return "", fmt.Errorf("Error: parametro incorrecto: %s", match)
		}
		clave, valor := strings.ToLower(keyValue[0]), keyValue[1] // extrae, guarda y convierte a minusculas

		// --- Eliminar las comillas
		if strings.HasPrefix(valor, "\"") && strings.HasSuffix(valor, "\"") {
			valor = strings.Trim(valor, "\"")
		}

		switch clave {
		case "-path":
			if valor == "" {
				return "", errors.New("el path no puede estar vacío")
			}
			comando.path = valor
			// elimina el archivo
			err := os.Remove(valor)
			if err != nil {
				fmt.Println(err)
			}
		default:
			return "", fmt.Errorf("Error: parametro incorrecto: %s", clave)
		}
	}

	if comando.path == "" {
		return "", errors.New("Error: faltan parametros obligatorios")
	}

	// Devuelve el comando RMDISK
	return "RMDISK: Disco eliminado correctamente.", nil
}
