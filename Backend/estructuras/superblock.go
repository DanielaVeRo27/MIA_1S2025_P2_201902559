package estructuras

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"time"
)

type SuperBlock struct {
	S_tipo_archivosistema int32
	S_contador_inodos     int32
	S_contador_bloques    int32
	S_cont_inodos_libres  int32
	S_cont_bloques_libres int32
	S_mtiempo             float32
	S_umtiempo            float32
	S_contador_mont       int32
	S_magic               int32
	S_tamano_inodo        int32
	S_tamano_bloque       int32
	S_primer_inodo        int32
	S_primer_bloque       int32
	S_bm_inicio_inodo     int32
	S_bm_inicio_bloque    int32
	S_inicio_inodo        int32
	S_inicio_bloque       int32
}

func (sb *SuperBlock) Serialize(path string, offset int64) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Seek(offset, 0)

	if err != nil {
		return err
	}

	err = binary.Write(file, binary.LittleEndian, sb)
	if err != nil {
		return err
	}
	return nil
}

func (sb *SuperBlock) Deserializar(path string, offset int64) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Seek(offset, 0)

	if err != nil {
		return err
	}

	sbSize := binary.Size(sb)
	if sbSize <= 0 {
		return fmt.Errorf("Tamaño de SuperBloque inválido: %d", sbSize)
	}

	buffer := make([]byte, sbSize)
	_, err = file.Read(buffer)

	if err != nil {
		return err
	}

	reader := bytes.NewReader(buffer)
	err = binary.Read(reader, binary.LittleEndian, sb)

	if err != nil {
		return err
	}
	return nil
}

func (sb *SuperBlock) Print() {
	mountTime := time.Unix(int64(sb.S_mtiempo), 0)
	unmountTime := time.Unix(int64(sb.S_umtiempo), 0)

	fmt.Printf("Filesystem Type: %d\n", sb.S_tipo_archivosistema)
	fmt.Printf("Inodes Count: %d\n", sb.S_contador_inodos)
	fmt.Printf("Blocks Count: %d\n", sb.S_contador_bloques)
	fmt.Printf("Free Inodes Count: %d\n", sb.S_cont_inodos_libres)
	fmt.Printf("Free Blocks Count: %d\n", sb.S_cont_bloques_libres)
	fmt.Printf("Mount Time: %s\n", mountTime.Format(time.RFC3339))
	fmt.Printf("Unmount Time: %s\n", unmountTime.Format(time.RFC3339))
	fmt.Printf("Mount Count: %d\n", sb.S_contador_mont)
	fmt.Printf("Magic: %d\n", sb.S_magic)
	fmt.Printf("Inode Size: %d\n", sb.S_tamano_inodo)
	fmt.Printf("Block Size: %d\n", sb.S_tamano_bloque)
	fmt.Printf("First Inode: %d\n", sb.S_primer_inodo)
	fmt.Printf("First Block: %d\n", sb.S_primer_bloque)
	fmt.Printf("Bitmap Inode Start: %d\n", sb.S_bm_inicio_inodo)
	fmt.Printf("Bitmap Block Start: %d\n", sb.S_bm_inicio_bloque)
	fmt.Printf("Inode Start: %d\n", sb.S_inicio_inodo)
	fmt.Printf("Block Start: %d\n", sb.S_inicio_bloque)
}
