package comandos

import (
	"Backend/estructuras"
	global "Backend/global"
	utils "Backend/utils"

	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Estructura para el cmd mount -> parametros
type mount struct {
	path string
	name string
}

// mount -path=/home/Disco.mia -name=Part #id=341a

// analizador del los parametros del comando mount
func Parse_mount(token []string) (string, error) {

	comando := &mount{} // inst de mount

	// unir tokens
	prm := strings.Join(token, " ")

	expr := regexp.MustCompile(`-path="[^"]+"|-path=[^\s]+|-name="[^"]+"|-name=[^\s]+`)

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
		// Switch para manejar diferentes parámetros
		switch clave {
		case "-path":
			// Verifica que el path no esté vacío
			if valor == "" {
				return "", errors.New("el path no puede estar vacío")
			}
			comando.path = valor
		case "-name":
			// Verifica que el nombre no esté vacío
			if valor == "" {
				return "", errors.New("el nombre no puede estar vacío")
			}
			comando.name = valor
		default:
			return "", fmt.Errorf("parámetro desconocido: %s", clave)
		}
	}

	// ver que los parametros si hayan sido proporcionaddo
	if comando.path == "" {
		return "", errors.New(("faltan el párametro : -path"))
	}
	if comando.name == "" {
		return "", errors.New("faltaa el parametro: -name")
	}

	// montar particion
	id_partition, err := comandoMount(comando)

	if err != nil {
		return "", err
	}

	//return comando, nil

	return fmt.Sprintf("MOUNT: La partición %s fue montada correctamente con ID %d", comando.name, id_partition), nil

}
func comandoMount(mount *mount) (string, error) {
	var mbr estructuras.MBR
	err := mbr.Deserializar(mount.path)

	if err != nil {
		fmt.Println("Error al deserializar el MBR: ", err)
		return "", err
	}

	partition, index_partition := mbr.Get_partition_name(mount.name)

	if partition == nil {
		fmt.Println("Error: La partición no existe")
		return "", errors.New("La partición no existe")
	}

	fmt.Println("\nPartición disponible:")
	partition.PrintPartition()

	id_partition, partitionCorrelative, err := GenerateIdPartition(mount, index_partition)
	if err != nil {
		fmt.Println("Error al generar el ID de la partición: ", err)
		return "", err
	}
	// guardar particion montada
	global.ParticionesMontadas[id_partition] = mount.path

	// modificar la particion
	partition.MountParticion(partitionCorrelative, id_partition)

	fmt.Println("\nPartición montada (modificada):")
	partition.PrintPartition()

	fmt.Println("----------------------")
	fmt.Println("\n IDs de las particiones montadas")
	global.MostrarParticionesMontadas()

	mbr.Mbr_partitions[index_partition] = *partition

	err = mbr.Serializar(mount.path)
	if err != nil {
		fmt.Println("Error al serializar el MBR: ", err)
		return "", err
	}

	return id_partition, nil
}
func GenerateIdPartition(mount *mount, indexPartition int) (string, int, error) {
	letter, part_corr, err := utils.GetLetterAndPartitionCorrelative(mount.path)

	if err != nil {
		fmt.Println("Erro obteniendo la letra: ", err)
		return "", 0, err
	}

	id_partition := fmt.Sprintf("%s%d%s", global.Carnet, indexPartition+1, letter)
	// Crear id de partición
	//id_partition := fmt.Sprintf("%s%d%s",  global.Carnet, part_corr, letter)
	return id_partition, part_corr, nil
}
