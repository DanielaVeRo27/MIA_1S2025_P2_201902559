// Paquete para manejo de estructuras de datos
package estructuras

/*
    archivo en cual definimos la estructura del MBR -> para manejar los discos
	y las funciones para serializar y deserializar
	buscamos particiones oir nomvre o id
	imprimir los valores del MBR
	imprimir las particiones del MBR
*/

import (
	"bytes"           // -> manipulación de buffers
	"encoding/binary" // -> codificación y decodificación de datos binarios
	"errors"
	"fmt"
	"os" // -> para funciones del sistema operativo
	"strings"
	"time"
)

// Estructura del MBR
type MBR struct {
	Size_mbr          int32        // tamaño
	Creation_date_mbr float32      // fecha de creación
	Signature_mbr     int32        // firma del disco
	Fit_mbr           [1]byte      // tipo de ajuste F,W,B
	Mbr_partitions    [4]Partition // Particiones del MBR
}

// Serialize MBR	-> guarda el MBR en el archivo binario
func (mbr *MBR) Serializar(path string) error {
	// Abrir el archivo en modo escritura y creación
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	err = binary.Write(file, binary.LittleEndian, mbr)
	if err != nil {
		return err
	}

	return nil
}

// Deserializar MBR -> lee la estructura MBR desde el archivo binario
func (mbr *MBR) Deserializar(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// tamaño de la estructura MBR
	mbrSize := binary.Size(mbr)
	if mbrSize <= 0 {
		return fmt.Errorf("invalido MBR size: %d", mbrSize)
	}

	// Leer los bytes del archivo
	buffer := make([]byte, mbrSize)
	_, err = file.Read(buffer)
	if err != nil {
		return err
	}

	// Deserializar los bytes leídos en la estructura MBR
	reader := bytes.NewReader(buffer)
	err = binary.Read(reader, binary.LittleEndian, mbr)
	if err != nil {
		return err
	}

	return nil
}

// Obtener partición disponible
func (mbr *MBR) GetFirstAvailablePartition() (*Partition, int, int) {

	offset := binary.Size(mbr) // Inicializar el offset con el tamaño del MBR

	for i := 0; i < len(mbr.Mbr_partitions); i++ {

		if mbr.Mbr_partitions[i].Part_inicio == -1 {
			return &mbr.Mbr_partitions[i], offset, i
		} else {
			offset += int(mbr.Mbr_partitions[i].Part_size)
		}
	}
	return nil, -1, -1
}

// Busca una partición por nombre
func (mbr *MBR) Get_partition_name(nombre string) (*Partition, int) {

	for i, partition := range mbr.Mbr_partitions {
		// extraer el nombre de la partición y eliminar los caracteres nulos
		partition_nom := strings.Trim(string(partition.Part_nombre[:]), "\x00 ")
		// elimina los caracteres nulos
		name_i := strings.Trim(nombre, "\x00 ")
		if strings.EqualFold(partition_nom, name_i) {
			return &partition, i
		}
	}
	return nil, -1
}

// busca y extrae partición por ID
func (mbr *MBR) Get_Partition_ID(id string) (*Partition, error) {
	for i := 0; i < len(mbr.Mbr_partitions); i++ {
		//eliminar los caracteres nulos deL id de la partición
		IDpartition := strings.Trim(string(mbr.Mbr_partitions[i].Part_id[:]), "\x00 ")
		// eliminar los caracteres nulos del id
		ID_i := strings.Trim(id, "\x00 ")

		// comparar los id de las particiones
		if strings.EqualFold(IDpartition, ID_i) {
			return &mbr.Mbr_partitions[i], nil
		}
	}
	return nil, errors.New("partición no encontrada")
}

// Imprimir los valores del MBR
func (mbr *MBR) Imprimir_mbr() {

	creationTime := time.Unix(int64(mbr.Creation_date_mbr), 0)

	// Convertir Mbr_disk_fit a char
	discoFit := rune(mbr.Fit_mbr[0])

	fmt.Printf("MBR Size: %d\n", mbr.Size_mbr)
	fmt.Printf("Creation Date: %s\n", creationTime.Format(time.RFC3339))
	fmt.Printf("Disk Signature: %d\n", mbr.Signature_mbr)
	fmt.Printf("Disk Fit: %c\n", discoFit)

}

// ExisteParticionExtendida verifica si ya hay una partición extendida en el MBR
func (mbr *MBR) ExisteParticionExtendida() bool {
	for _, partition := range mbr.Mbr_partitions {
		if partition.Part_tipo[0] == 'E' { // Verifica si el tipo de partición es 'E' (Extendida)
			return true
		}
	}
	return false
}

// Obtener la partición extendida
func (mbr *MBR) GetExtendedPartition() *Partition {
	for i := 0; i < len(mbr.Mbr_partitions); i++ {
		if mbr.Mbr_partitions[i].Part_tipo[0] == 'E' {
			return &mbr.Mbr_partitions[i]
		}
	}
	return nil
}

// Imprimir las particiones del MBR
func (mbr *MBR) PrintPartitions() {

	for i, partition := range mbr.Mbr_partitions {
		// concertir los valores a char
		estadoP := rune(partition.Part_estado[0])
		tipoP := rune(partition.Part_tipo[0])
		fitP := rune(partition.Part_fit[0])
		nombreP := string(partition.Part_nombre[:])
		idP := string(partition.Part_id[:])

		fmt.Printf("Partition %d:\n", i+1)
		fmt.Printf("  Status: %c\n", estadoP)
		fmt.Printf("  Type: %c\n", tipoP)
		fmt.Printf("  Fit: %c\n", fitP)
		fmt.Printf("  Start: %d\n", partition.Part_inicio)
		fmt.Printf("  Size: %d\n", partition.Part_size)
		fmt.Printf("  Name: %s\n", nombreP)
		fmt.Printf("  Correlative: %d\n", partition.Part_correlativo)
		fmt.Printf("  ID: %s\n", idP)

	}
}
