package comandos

import (
	"Backend/estructuras"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"Backend/utils"
)

// Estructura para el comando mkdisk -> parametros
type mkdisk struct {
	size int    // tamaño del disco
	unit string // unidad de medida
	fit  string // ajuste de particiones
	path string // ruta del disco
}

// analizador de comando mkdisk
func Parsermkdisk(tokens []string) (string, error) {

	comando := &mkdisk{} // crea un nuevo comando mkdisk -> instancio con estructur MDISK

	prm := strings.Join(tokens, " ") // pasa y une los parameytros a un string

	// --- Expresion regular para el comando mkdisk
	expr := regexp.MustCompile(`-size=\d+|-unit=[kKmM]|-fit=[bBfFwW]{2}|-path="[^"]+"|-path=[^\s]+`)

	// --- Busca los parametros en la expresion regular y encuentra las coincidencias
	for _, match := range expr.FindAllString(prm, -1) {

		// dividir  en clave y valor
		keyValue := strings.SplitN(match, "=", 2)
		if len(keyValue) != 2 {
			return "", fmt.Errorf("Error: parametro incorrecto: %s", match)
		}
		clave, valor := strings.ToLower(keyValue[0]), keyValue[1]

		// --- Eliminar las comillas
		if strings.HasPrefix(valor, "\"") && strings.HasSuffix(valor, "\"") {
			valor = strings.Trim(valor, "\"")
		}

		// ----->> Siwtch para evaluar los parametros
		switch clave {
		case "-size":
			// convierte el valor a entero
			size, err := strconv.Atoi(valor)
			if err != nil || size <= 0 {
				return "", fmt.Errorf("Error: valor incorrecto para size: %s", valor)
			}
			comando.size = size
		case "-unit":
			// que sea K o M
			if valor != "K" && valor != "M" {
				return "", fmt.Errorf("Error: valor incorrecto para unit: %s", valor)
			}
			comando.unit = strings.ToUpper(valor)
		case "-fit":
			// que sea F, B o W
			if valor != "F" && valor != "B" && valor != "W" {
				return "", fmt.Errorf("Error: valor incorrecto para fit: %s", valor)
			}
			comando.fit = valor
		case "-path":
			// que sea una ruta valida
			if valor == "" {
				return "", errors.New("Error: path vacío")
			}
			comando.path = valor
		default:
			return "", fmt.Errorf("Error: parametro desconocido: %s", clave)
		}
	}

	// --- Verifica que los parametros obligatorios esten presentes
	if comando.size == 0 {
		return "", errors.New("Error: falta de parametros para size")
	}

	if comando.path == "" {
		return "", errors.New("Error: falta el parametro requerido para path")
	}

	// si el comando unit no se especifica, se establece en M
	if comando.unit == "" {
		comando.unit = "M"
	}

	// Si el comando fit no se especifica, se establece en "FF"
	if comando.fit == "" {
		comando.fit = "FF"
	}

	// ---->> Creacion del disco con los parametros
	err := comando_mkdisk(comando)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err

	}
	// Retorna el comando MKDISK creado
	return fmt.Sprintf("MKDISK: Disco creado exitosamente en la ruta %s de %d%s con ajuste %s.", comando.path, comando.size, comando.unit, comando.fit), nil
}

// --- Crear disco con los parametros especificados
func comando_mkdisk(mkdisk1 *mkdisk) error {
	sizeBytes, err := utils.Conversion(mkdisk1.size, mkdisk1.unit)
	if err != nil {
		fmt.Println("Error converting size:", err)
		return err
	}

	// --- Crear DISCO con el tamaño especificado
	err = crear_disco(mkdisk1, sizeBytes)
	if err != nil {
		fmt.Println("Error creating disk:", err)
		return err
	}

	// --- creacion del MBR
	err = crear_MBR(mkdisk1, sizeBytes)
	if err != nil {
		fmt.Println("Error al crear el MBR:", err)
		return err
	}
	return nil
}

func crear_disco(mkdisk1 *mkdisk, sizeBytes int) error {

	// extraer el directorio del path, crear las carpetas u subcarpetas si no existe
	err := os.MkdirAll(filepath.Dir(mkdisk1.path), os.ModePerm)
	if err != nil {
		fmt.Println("Error al crear  file: ", err)
		return err
	}

	// --- crear binario
	file, err := os.Create(mkdisk1.path)
	if err != nil {
		fmt.Println("Error al crear file: ", err)
		return err
	}
	defer file.Close()

	// --- escribir en el archivo
	buffer := make([]byte, 1024*1024) // 1MB espacio temp
	for sizeBytes > 0 {
		writeSize := len(buffer)
		if sizeBytes < len(buffer) {
			writeSize = sizeBytes // Configura el tamaño de escritura
		}
		if _, err := file.Write(buffer[:writeSize]); err != nil {
			fmt.Println("Error al escribir en el documento ", err)
			return err
		}
		sizeBytes -= writeSize // Resta el tamaño de escritura
	}
	return nil
}

func crear_MBR(mkdisk1 *mkdisk, sizeBytes int) error {

	// --- Tipo de ajuste
	var tipefite byte
	switch mkdisk1.fit {
	case "FF":
		tipefite = 'F'
	case "BF":
		tipefite = 'B'
	case "WF":
		tipefite = 'W'
	default:
		fmt.Println("Invalido el tipo de fit ")
		return nil
	}

	// Crear el MBR con los valores proporcionados
	mbr := &estructuras.MBR{
		Size_mbr:          int32(sizeBytes),           // tamaño del disco
		Creation_date_mbr: float32(time.Now().Unix()), // fecha de creacion
		Signature_mbr:     rand.Int31(),               // firma del disco
		Fit_mbr:           [1]byte{tipefite},          // tipo de ajuste F,W,B
		Mbr_partitions: [4]estructuras.Partition{

			// Particiones del MBR
			// -1 -> significa que no se ha asignado un valor
			// N -> significa que no se ha asignado un valor
			{Part_estado: [1]byte{'N'}, Part_tipo: [1]byte{'N'}, Part_fit: [1]byte{'N'}, Part_inicio: -1, Part_size: -1, Part_nombre: [16]byte{'N'}, Part_correlativo: -1, Part_id: [4]byte{'N'}},
			{Part_estado: [1]byte{'N'}, Part_tipo: [1]byte{'N'}, Part_fit: [1]byte{'N'}, Part_inicio: -1, Part_size: -1, Part_nombre: [16]byte{'N'}, Part_correlativo: -1, Part_id: [4]byte{'N'}},
			{Part_estado: [1]byte{'N'}, Part_tipo: [1]byte{'N'}, Part_fit: [1]byte{'N'}, Part_inicio: -1, Part_size: -1, Part_nombre: [16]byte{'N'}, Part_correlativo: -1, Part_id: [4]byte{'N'}},
			{Part_estado: [1]byte{'N'}, Part_tipo: [1]byte{'N'}, Part_fit: [1]byte{'N'}, Part_inicio: -1, Part_size: -1, Part_nombre: [16]byte{'N'}, Part_correlativo: -1, Part_id: [4]byte{'N'}},
		},
	}
	fmt.Println("\nMBR creado:")
	mbr.Imprimir_mbr()

	// Serializar --> MBR en el archivo
	err := mbr.Serializar(mkdisk1.path)
	if err != nil {
		fmt.Println("Error:", err)
	}

	return nil

}
