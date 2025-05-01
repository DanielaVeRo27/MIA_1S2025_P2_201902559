package estructuras

import (
	"encoding/binary"
	"errors"
	"os"
)

/*
	archivo en cual definimos la estructura del EBR -> para manejar las particiones logicas
	0->  no montada
	1 -> montada
*/

// Estructura del EBR
type EBR struct {
	Part_mount [1]byte  // Estado de la partición
	Part_fit   [1]byte  // Ajuste de la partición
	Part_start int32    // Byte de inicio de la partición
	Part_size  int32    // Tamaño de la partición
	Part_next  int32    // Byte de inicio de la partición siguiente
	Part_name  [16]byte // Nombre de la partición

}

// Serializar escribe el EBR en el archivo binario en la posición especificada
func (ebr *EBR) Serializar(path string, startPartition int64) error {
	// Abrir el archivo en modo lectura/escritura
	file, err := os.OpenFile(path, os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Mover el puntero del archivo a la posición de inicio de la partición lógica
	_, err = file.Seek(startPartition, 0)
	if err != nil {
		return err
	}

	// Escribir la estructura EBR en el archivo binario
	err = binary.Write(file, binary.LittleEndian, ebr)
	if err != nil {
		return err
	}

	return nil
}

func LeerEBR(path string, start int32) ([]EBR, error) {
	var ebrs []EBR

	file, err := os.Open(path)
	if err != nil {
		return nil, errors.New("Error al abrir el archivo del disco")
	}
	defer file.Close()

	var ebr EBR
	currentStart := start

	for {
		file.Seek(int64(currentStart), 0)
		err := binary.Read(file, binary.LittleEndian, &ebr)
		if err != nil {
			return nil, errors.New("Error al leer el EBR")
		}

		if ebr.Part_start == 0 || ebr.Part_size == 0 {
			break // Si el EBR no tiene datos, terminamos la lectura
		}

		ebrs = append(ebrs, ebr)

		if ebr.Part_next == -1 {
			break // Último EBR en la lista
		}

		currentStart = ebr.Part_size
	}

	return ebrs, nil
}

func EncontrarEspacioEBR(ebrs []EBR, startExtended int32, sizeExtended int32, sizeNewPartition int32) int32 {
	var currentStart int32 = startExtended
	var availableStart int32 = -1

	// Recorrer la lista de EBRs para encontrar un espacio disponible
	for _, ebr := range ebrs {
		espacioDisponible := ebr.Part_start - currentStart
		if espacioDisponible >= sizeNewPartition {
			availableStart = currentStart
			break
		}
		currentStart = ebr.Part_start + ebr.Part_size
	}

	// Verificar si hay espacio al final de la partición extendida
	if availableStart == -1 && (startExtended+sizeExtended)-(currentStart) >= sizeNewPartition {
		availableStart = currentStart
	}

	return availableStart
}
