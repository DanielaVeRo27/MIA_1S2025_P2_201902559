package comandos

import (
	"Backend/estructuras"
	"Backend/utils"
	"errors"
	"fmt"
	"regexp"
	"strconv" // convertir cadenas a otros tipos de datos
	"strings"
)

/*
	----->> FDISK --> comando que se encarga de la administración de las particiones en el archivo

*/

// Estructura para el comando FDISK -> parametros
type fdisk struct {
	size int    // tamaño de la particion
	unit string // unidad de medida
	path string // ruta del disco
	tipo string // tipo de particion -> P,E,L
	fit  string // ajuste de particiones
	name string // nombre de la particion
}

/*
 fdisk -size=300 -path=/home/Disco1.mia -name=Particion1
 fdisk -type=E -path=/home/Disco2.mia -Unit=K -name=Particion2
-size=300

*/

func ParserFdisk(tokens []string) (string, error) {

	comando := &fdisk{} // crea un nuevo comando fdisk -> instancio con estructur FDISK

	prm := strings.Join(tokens, " ") // pasa y une los parameytros a un string

	// --- Expresion regular para el comando fdisk
	expr := regexp.MustCompile(`-size=\d+|-unit=[kKmM]|-fit=[bBfFwW]{2}|-path="[^"]+"|-path=[^\s]+|-type=[pPeElL]|-name="[^"]+"|-name=[^\s]+`)

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
		fmt.Println("clave: ", clave)
		fmt.Println("valor: ", valor)
		switch clave {

		case "-size":
			// convierte el valor a entero
			size, err := strconv.Atoi(valor)
			if err != nil || size <= 0 {
				return "", errors.New("Error: valor incorrecto para size: %s" + valor)
			}
			comando.size = size
		case "-unit":
			// convierte el valor a minusculas
			unit := strings.ToLower(valor)
			if unit != "k" && unit != "m" {
				return "", errors.New("Error: valor incorrecto para unit: %s" + valor)
			}
			comando.unit = unit
		case "-fit":
			// convierte el valor a mayusculas
			fit := strings.ToUpper(valor)
			if fit != "BF" && fit != "FF" && fit != "WF" {
				return "", errors.New("Error: valor incorrecto para fit: el valor debe de ser BF, FF o WF no  %s" + valor)
			}
			comando.fit = fit
		case "-path":
			// que sea una ruta valida
			if valor == "" {
				return "", errors.New("Error: el path no puede estar vacio %s" + "")
			}
			comando.path = valor
		case "-type":
			// convierte el valor a mayusculas
			tipo := strings.ToUpper(valor)
			if tipo != "P" && tipo != "E" && tipo != "L" {
				return "", errors.New("el tipo debe ser P, E o L")
			}
			comando.tipo = tipo
		case "-name":
			// Verifica que el nombre no esté vacío
			if valor == "" {
				return "", errors.New("Error: el nombre no puede estar vacío")
			}
			comando.name = valor
		default:
			return "", fmt.Errorf("parámetro desconocido: %s", clave)
		}
	}
	// ----> Verificamos que los parametros obligatorios esten presentes
	if comando.size == 0 {
		return "", errors.New("Error: falta de parametros para: -size")
	}
	if comando.path == "" {
		return "", errors.New("Error: falta el parametro requerido para -path")
	}
	if comando.name == "" {
		return "", errors.New("Error: falta el parametro requerido para -name")
	}
	// ---> si los siguientes parametros no se especifican, se establecen por defecto
	if comando.tipo == "" {
		comando.tipo = "P"
	}
	if comando.unit == "" {
		comando.unit = "M"
	}
	if comando.fit == "" {
		comando.fit = "FF"
	}

	// ---->> Creacion de la particion con los parametros
	err := comando_fdisk(comando)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err

	}
	// retorna el comando fdisk
	return fmt.Sprintf("FDISK: Partición %s creada correctamente de %d%s de tipo %s con ajuste%s.", comando.name, comando.size, comando.unit, comando.tipo, comando.fit), nil
}

// ---->> Funcion para crear la particion
func comando_fdisk(comando *fdisk) error {
	fmt.Println("entrando a la funcion de crear particion ")

	// tamaño a byres
	sizeBytes, err := utils.Conversion(comando.size, comando.unit)
	if err != nil {
		fmt.Println("Error al convertir el tamaño ", err)
		return err
	}

	if comando.tipo == "P" {
		// ---> Particion primaria
		err = crear_particion_primaria(comando, sizeBytes)
		if err != nil {
			fmt.Println("Error al crear la particion primaria", err)
			return err
		}
	} else if comando.tipo == "E" {
		// ---> Particion extendida
		err = crear_particion_extendida(comando, sizeBytes)
		if err != nil {
			fmt.Println("Error al crear la particion extendida", err)
			return err
		}
	} else if comando.tipo == "L" {
		// ---> Particion logica
		err = crear_particion_logica(comando, sizeBytes)
		if err != nil {
			fmt.Println("Error al crear la particion logica", err)
			return err
		}
	}
	return nil
}

// --- CREAR PARTICION PRIMARIA
func crear_particion_primaria(comando *fdisk, sizeBytes int) error {

	var mbr estructuras.MBR

	// Desialiizar el MBR
	err := mbr.Deserializar(comando.path)
	if err != nil {
		fmt.Println("Error al deserializar el MBR", err)
		return err
	}

	fmt.Println("\nMBR original:")
	mbr.Imprimir_mbr()

	// obtener la 1era particion libre
	availablePartition, startPartition, indexPartition := mbr.GetFirstAvailablePartition()
	if availablePartition == nil {
		fmt.Println("No hay particiones disponibles.")
	}

	fmt.Println("\nPartición disponible:")
	availablePartition.PrintPartition()

	// --- Crear particion
	availablePartition.CrearParticion(startPartition, sizeBytes, comando.tipo, comando.fit, comando.name)

	fmt.Println("\nMBR modificado:")
	availablePartition.PrintPartition()

	// Colocar en el MBR
	if availablePartition != nil {
		mbr.Mbr_partitions[indexPartition] = *availablePartition
	}

	// Serializar el MBR
	err = mbr.Serializar(comando.path)
	if err != nil {
		fmt.Println("Error al serializar el MBR", err)
	}

	// --- Crear particion primaria
	return nil
}

// --- CREAR PARTICION EXTENDIDA
func crear_particion_extendida(comando *fdisk, sizeBytes int) error {
	// --- Crear particion extendida
	var mbr estructuras.MBR

	// Deserializar el MBR
	err := mbr.Deserializar(comando.path)
	if err != nil {
		fmt.Println("Error al deserializar el MBR", err)
	}

	fmt.Println("\nMBR original:")
	mbr.Imprimir_mbr()

	// verificamos que no existan otras particiones extendidas
	if mbr.ExisteParticionExtendida() {
		fmt.Println("Ya existe una partición extendida.")
		return errors.New("Error: ya existe una particion extendida")

	}

	// busaca la primera partición libre devuelve, puntero, ubicación y el índice
	availablePartition, startPartition, indexPartition := mbr.GetFirstAvailablePartition()
	if availablePartition == nil {
		return errors.New("No hay particiones disponibles para crear la partición extendida")
	}

	fmt.Println("\nPartición disponible:")
	availablePartition.PrintPartition()

	// asigna los valores a la partición
	availablePartition.CrearParticion(startPartition, sizeBytes, "E", comando.fit, comando.name)

	fmt.Println("\nPartición creada (modificada):")
	availablePartition.PrintPartition()

	// Colocar la partición en el MBR
	if availablePartition != nil {
		mbr.Mbr_partitions[indexPartition] = *availablePartition
	}

	fmt.Println(comando.name)
	var nameBytes [16]byte
	copy(nameBytes[:], comando.name)
	// Crear el primer EBR en la partición extendida
	ebr := &estructuras.EBR{
		Part_mount: [1]byte{0},              // 0-> no montada
		Part_fit:   [1]byte{comando.fit[0]}, // First Fit
		Part_start: int32(startPartition),
		Part_size:  int32(sizeBytes),
		Part_next:  int32(0), //  -1 en entero si lo cambias)
		Part_name:  nameBytes,
	}

	err = ebr.Serializar(comando.path, int64(startPartition))
	if err != nil {
		fmt.Println("Error al serializar el EBR", err)
		return err
	}

	// Guardar el MBR modificado
	err = mbr.Serializar(comando.path)
	if err != nil {
		fmt.Println("Error al serializar el MBR", err)
		return err
	}

	fmt.Println("Partición extendida creada correctamente")
	return nil
}

// -------->>  CREAR PARTICION LOGICA
func crear_particion_logica(comando *fdisk, sizeBytes int) error {

	var mbr estructuras.MBR
	// obtiene el MBR del disco -> accedemos al disco
	err := mbr.Deserializar(comando.path)
	if err != nil {
		fmt.Println("Error al deserializar el MBR", err)
		return err
	}

	// obtener particion extendida
	extendedPartition := mbr.GetExtendedPartition()
	if extendedPartition == nil {
		return errors.New("No existe una partición extendida en el disco %s ")
	}

	// Leer de EBRs
	ebrs, err := estructuras.LeerEBR(comando.path, extendedPartition.Part_inicio)
	if err != nil {
		return err
	}
	// Encontrar el primer espacio disponible
	startPartition := estructuras.EncontrarEspacioEBR(ebrs, extendedPartition.Part_inicio, extendedPartition.Part_size, int32(sizeBytes))
	fmt.Println("star espacio ", startPartition)

	if startPartition == -1 {
		return errors.New("No hay suficiente espacio en la partición extendida para la partición lógica")
	}
	// convertimos el nombre
	var nameBytes [16]byte
	copy(nameBytes[:], comando.name)
	// Crear nuevo EBR para la partición lógica
	nuevoEBR := estructuras.EBR{
		Part_mount: [1]byte{0},              // 0-> no montada
		Part_fit:   [1]byte{comando.fit[0]}, // First Fit
		Part_start: int32(startPartition),
		Part_size:  int32(sizeBytes),
		Part_next:  int32(0), //  -1 en entero si lo cambias)
		Part_name:  nameBytes,
	}
	// Escribir el nuevo EBR en el archivo
	err = nuevoEBR.Serializar(comando.path, int64(startPartition))
	if err != nil {
		fmt.Println("Error al serializar el nuevo EBR", err)
		return err
	}
	// Actualizar el EBR anterior si existe
	if len(ebrs) > 0 {
		ultimoEBR := &ebrs[len(ebrs)-1]
		ultimoEBR.Part_next = startPartition
		err = ultimoEBR.Serializar(comando.path, int64(ultimoEBR.Part_start))
		if err != nil {
			return err
		}
	}

	// --- Crear particion logica
	return nil
}
