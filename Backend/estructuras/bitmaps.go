package estructuras

import (
	"encoding/binary"
	"os"
)

// ---> Creacion de los Bitmaps de los inodos y bloques en el archivo
func (sb *SuperBlock) CrearBitMaps(path string) error {

	// Escribimos el Bitmap
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// bitmap de inodos -> movemos el puntero del archivo a la pos
	_, err = file.Seek(int64(sb.S_bm_inicio_inodo), 0)

	if err != nil {
		return err
	}
	// el buffer de n '0'
	buffer := make([]byte, sb.S_cont_inodos_libres)

	for i := range buffer {
		buffer[i] = '0'
	}

	// buffer escrito en el archivo
	err = binary.Write(file, binary.LittleEndian, buffer)

	if err != nil {
		return err
	}

	// Bitmap -> bloques | movemos el punteo del archivo a la pos
	_, err = file.Seek(int64(sb.S_bm_inicio_bloque), 0)
	if err != nil {
		return err
	}

	buffer = make([]byte, sb.S_cont_bloques_libres)
	for i := range buffer {
		buffer[i] = '0'
	}

	// --->> escribir en el archivo el buffer
	err = binary.Write(file, binary.LittleEndian, buffer)
	if err != nil {
		return err
	}

	return nil
}

// ---> Actualizar el Bitmap
func (sb *SuperBlock) UpdateBitMapInode(path string) error {

	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// mover el puntero del archivo a la pos del bit
	_, err = file.Seek(int64(sb.S_bm_inicio_inodo)+int64(sb.S_contador_inodos), 0)
	if err != nil {
		return err
	}

	// escribir en el archivo el bit
	_, err = file.Write([]byte{'1'})
	if err != nil {
		return err
	}

	return nil

}

func (sb *SuperBlock) UpdateBitMapBlock(path string) error {
	file, err := os.OpenFile(path, os.O_RDWR, 0644)

	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Seek(int64(sb.S_bm_inicio_bloque)+int64(sb.S_contador_bloques), 0)

	if err != nil {
		return err
	}

	_, err = file.Write([]byte{'X'})

	if err != nil {
		return err
	}
	return nil
}
