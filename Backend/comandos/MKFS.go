package comandos

import (
	estructura "Backend/estructuras"
	global "Backend/global"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"regexp"
	"strings"
	"time"
)

type MKFS struct {
	id   string
	tipo string
}

/*
	mkfs -type=full -id=341A
	mkfs -id=342A
*/

// Parser del comando MKFS
func ParserMkfs(tokens []string) (string, error) {

	comando := &MKFS{}

	prm := strings.Join(tokens, " ")

	expr := regexp.MustCompile((`-id=[^\s]+|-type=[^\s]+`))

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

		switch clave {
		case "-id":
			if valor == "" {
				return "", errors.New("Error -> el id no puede estar vacío")
			}
			comando.id = valor

		case "-type":
			if valor != "full" {
				return "", errors.New("el tipo debe ser full")
			}
			comando.tipo = valor

		default:
			return "", fmt.Errorf("parámetro desconocido: %s", clave)
		}
	}

	if comando.id == "" {
		return "", errors.New("faltan parámetros requeridos: -id")
	}

	if comando.tipo == "" {
		comando.tipo = "full"
	}

	err := comando_Mkfs(comando)
	if err != nil {
		fmt.Println("Error:", err)
	}
	return fmt.Sprintf("MKFS: Partición %s formateada correctamente con sistema de archivos %s.", comando.id, comando.tipo), nil

}

func comando_Mkfs(mkfs *MKFS) error {

	// obtenemos la par-- montada
	particion_Mount, path_Partition, err := global.ObtenerParticionesMontadas(mkfs.id)

	if err != nil {
		return err
	}

	fmt.Println("\nPatición montada:")
	particion_Mount.PrintPartition()

	// Valor de n
	n := calcularN(particion_Mount)

	fmt.Println("\nValor de n:", n)

	super_block := crear_superBlock(particion_Mount, n)

	fmt.Println("\nSuperBlock:")
	super_block.Print()

	err = super_block.CrearBitMaps(path_Partition)
	if err != nil {
		return err
	}

	// crea el archivo txt de users
	err = super_block.Crear_users_file(path_Partition)
	if err != nil {
		return err
	}

	fmt.Println("\nSuperBlock actualizado:")
	super_block.Print()

	// Serializar el superbloque
	err = super_block.Serialize(path_Partition, int64(particion_Mount.Part_inicio))
	if err != nil {
		return err
	}

	return nil

}

func calcularN(particion *estructura.Partition) int32 {

	numerador := int(particion.Part_size) - binary.Size(estructura.SuperBlock{})
	denominador := 4 + binary.Size(estructura.Inodo{}) + 3*binary.Size(estructura.SuperBlock{})

	n := math.Floor(float64(numerador) / float64(denominador))

	return int32(n)

}

func crear_superBlock(partition *estructura.Partition, n int32) *estructura.SuperBlock {

	bm_inodo_inicio := partition.Part_inicio + int32(binary.Size(estructura.SuperBlock{}))
	bm_block_inicio := bm_inodo_inicio + n

	inodo_inicio := bm_block_inicio + (3 * n)

	bloque_inicio := inodo_inicio + (int32(binary.Size(estructura.Inodo{})) * n)

	// Crea un nuevo super bloque
	super_Block := &estructura.SuperBlock{
		S_tipo_archivosistema: 2,
		S_contador_inodos:     0,
		S_contador_bloques:    0,
		S_cont_inodos_libres:  int32(n),
		S_cont_bloques_libres: int32(n * 3),
		S_mtiempo:             float32(time.Now().Unix()),
		S_umtiempo:            float32(time.Now().Unix()),
		S_contador_mont:       1,
		S_magic:               0xEF53,
		S_tamano_inodo:        int32(binary.Size(estructura.Inodo{})),
		S_tamano_bloque:       int32(binary.Size(estructura.FileBlock{})),
		S_primer_inodo:        inodo_inicio,
		S_primer_bloque:       bloque_inicio,
		S_bm_inicio_inodo:     bm_inodo_inicio,
		S_bm_inicio_bloque:    bm_block_inicio,
		S_inicio_inodo:        inodo_inicio,
		S_inicio_bloque:       bloque_inicio,
	}

	return super_Block

}
