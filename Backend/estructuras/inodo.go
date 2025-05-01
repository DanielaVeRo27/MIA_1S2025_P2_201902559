package estructuras

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"time"
)

/*
	--->  i_type
		1 = Archivo
		0 = Carpeta
*/

type Inodo struct {
	I_uid   int32     // UID del usuario propietario del archivo o carpeta
	I_gid   int32     // GID del grupo al q pertenece el archivo o carpeta
	I_size  int32     // tamaño del archivo en bytes
	I_atime float32   // Ultima fecha que se leyo el inodo sin modificarlo
	I_ctime float32   // Fecha en la que se creo el inodo
	I_mtime float32   // Ultima fecha en la que se modifica
	I_block [15]int32 // array -> 12 registros son bloques directos
	I_type  [1]byte   // Indica si es carpeta o archivo
	I_perm  [3]byte   // Guardar los permisos del archivo o de la carpeta
}

// Serializar para escribir la struct del inodo en un archivo binario en la pos
func (inodo *Inodo) Serializar(path string, offset int64) error {

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)

	if err != nil {
		return err
	}

	defer file.Close()

	// Mover el puntero del archivo a la pos especificada
	_, err = file.Seek(offset, 0)
	if err != nil {
		return err
	}

	// Serializar la estructura Inodo directamente en el archivo
	err = binary.Write(file, binary.LittleEndian, inodo)
	if err != nil {
		return err
	}

	return nil
}

// Deserializar -> para leee la estructura Inode desde un archivo binario
func (inodo *Inodo) Deserializar(path string, offset int64) error {

	file, err := os.Open(path)
	if err != nil {
		return err
	}

	defer file.Close()

	// Mover el puntero del archivo a la posición
	_, err = file.Seek(offset, 0)
	if err != nil {
		return err
	}

	// Obtener el tamaño de la estructura Inodo
	inodoSize := binary.Size(inodo)
	if inodoSize <= 0 {
		return fmt.Errorf("invalido el tamaño del Inodo: %d", inodoSize)
	}

	// Leer solo la cantidad de bytes
	buffer := make([]byte, inodoSize)
	_, err = file.Read(buffer)
	if err != nil {
		return err
	}

	// Deserializar los bytes leídos en la estructura Inodo
	reader := bytes.NewReader(buffer)
	err = binary.Read(reader, binary.LittleEndian, inodo)
	if err != nil {
		return err
	}

	return nil
}

// ----> Imprimir los atributodos del inodo
func (inodo *Inodo) Imprimir() {
	atime := time.Unix(int64(inodo.I_atime), 0)
	ctime := time.Unix(int64(inodo.I_ctime), 0)
	mtime := time.Unix(int64(inodo.I_mtime), 0)

	fmt.Printf("I_uid: %d\n", inodo.I_uid)
	fmt.Printf("I_gid: %d\n", inodo.I_gid)
	fmt.Printf("I_size: %d\n", inodo.I_size)
	fmt.Printf("I_atime: %s\n", atime.Format(time.RFC3339))
	fmt.Printf("I_ctime: %s\n", ctime.Format(time.RFC3339))
	fmt.Printf("I_mtime: %s\n", mtime.Format(time.RFC3339))
	fmt.Printf("I_block: %v\n", inodo.I_block)
	fmt.Printf("I_type: %s\n", string(inodo.I_type[:]))
	fmt.Printf("I_perm: %s\n", string(inodo.I_perm[:]))
}
