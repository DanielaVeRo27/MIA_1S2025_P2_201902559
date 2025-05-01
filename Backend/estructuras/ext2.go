package estructuras

import (
	"time"
)

// Crear user en nuestro sistema de archivos
func (sb *SuperBlock) Crear_users_file(path string) error {

	// ---- Inodo raiz -----
	raiz_inodo := &Inodo{
		I_uid:   1,
		I_gid:   1,
		I_size:  0,
		I_atime: float32(time.Now().Unix()),
		I_ctime: float32(time.Now().Unix()),
		I_mtime: float32(time.Now().Unix()),
		I_block: [15]int32{sb.S_contador_bloques, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		I_type:  [1]byte{'0'},
		I_perm:  [3]byte{'7', '7', '7'},
	}

	// Serializar el inodo
	err := raiz_inodo.Serializar(path, int64(sb.S_primer_inodo))
	if err != nil {
		return err
	}

	// actualizar el bitmap
	err = sb.UpdateBitMapInode(path)
	if err != nil {
		return err
	}

	// Actuliza el I -> inodos
	sb.S_contador_inodos++
	sb.S_cont_inodos_libres--
	sb.S_primer_inodo += sb.S_tamano_inodo

	// ------- Bloque del Inodo raiz ------
	bloque_raiz := &FolderBlock{
		B_content: [4]FolderContent{
			{B_name: [12]byte{'.'}, B_inodo: 0},
			{B_name: [12]byte{'.', '.'}, B_inodo: 0},
			{B_name: [12]byte{'-'}, B_inodo: -1},
			{B_name: [12]byte{'-'}, B_inodo: -1},
		},
	}

	err = sb.UpdateBitMapBlock(path)
	if err != nil {
		return err
	}

	err = bloque_raiz.Serializar(path, int64(sb.S_primer_bloque))
	if err != nil {
		return err
	}

	// Actualizar el sb --> Super Bloque
	sb.S_contador_bloques++
	sb.S_cont_bloques_libres--
	sb.S_primer_bloque += sb.S_tamano_bloque

	// ----> creamos el txt de user

	text_Users := "1,G,root\n1,U,root,123\n"

	// Deserializar el inodo raíz
	err = raiz_inodo.Deserializar(path, int64(sb.S_bm_inicio_inodo+0))
	if err != nil {
		return err
	}

	raiz_inodo.I_atime = float32(time.Now().Unix())

	err = raiz_inodo.Serializar(path, int64(sb.S_bm_inicio_inodo+0))
	if err != nil {
		return err
	}

	bloque_raiz.B_content[2] = FolderContent{B_name: [12]byte{'u', 's', 'e', 'r', 's', '.', 't', 'x', 't'}, B_inodo: sb.S_contador_inodos}

	err = bloque_raiz.Serializar(path, int64(sb.S_bm_inicio_bloque+0))
	if err != nil {
		return err
	}

	usersInode := &Inodo{
		I_uid:   1,
		I_gid:   1,
		I_size:  int32(len(text_Users)),
		I_atime: float32(time.Now().Unix()),
		I_ctime: float32(time.Now().Unix()),
		I_mtime: float32(time.Now().Unix()),
		I_block: [15]int32{sb.S_contador_bloques, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		I_type:  [1]byte{'1'},
		I_perm:  [3]byte{'7', '7', '7'},
	}

	// Actualizar  inodos
	err = sb.UpdateBitMapInode(path)
	if err != nil {
		return err
	}

	err = usersInode.Serializar(path, int64(sb.S_primer_inodo))
	if err != nil {
		return err
	}

	// Actuliza el I -> inodos
	sb.S_contador_inodos++
	sb.S_cont_inodos_libres--
	sb.S_primer_inodo += sb.S_tamano_inodo

	// Creamos el bloque de users.txt
	block_user := &FileBlock{
		B_content: [64]byte{},
	}

	copy(block_user.B_content[:], text_Users)

	err = block_user.Serializar(path, int64(sb.S_primer_bloque))
	if err != nil {
		return err
	}

	err = sb.UpdateBitMapBlock(path)
	if err != nil {
		return err
	}

	// Actualizar el sb --> Super Bloque
	sb.S_contador_bloques++
	sb.S_cont_bloques_libres--
	sb.S_primer_bloque += sb.S_tamano_bloque

	return nil

}
