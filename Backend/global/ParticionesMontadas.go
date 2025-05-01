package global

import (
	structures "Backend/estructuras"
	"errors"
	"fmt"
)

// carnet: 201902559
const Carnet string = "59"

var (
	ParticionesMontadas map[string]string = make(map[string]string)
)

func ObtenerParticionesMontadas(id string) (*structures.Partition, string, error) {
	path := ParticionesMontadas[id]

	if path == "" {
		return nil, "", errors.New("La partición no está montada")
	}
	var mbr structures.MBR

	err := mbr.Deserializar(path)
	if err != nil {
		return nil, "", err
	}
	partition, err := mbr.Get_Partition_ID(id)
	if partition == nil {
		return nil, "", err
	}
	return partition, path, nil
}

func ObtenerParticionesMontadas_SuperBlock(id string) (*structures.SuperBlock, *structures.Partition, string, error) {
	path := ParticionesMontadas[id]

	if path == "" {
		return nil, nil, "", errors.New("la partición no está montada")
	}

	var mbr structures.MBR
	err := mbr.Deserializar(path)

	if err != nil {
		return nil, nil, "", err
	}
	partition, err := mbr.Get_Partition_ID(id)

	if partition == nil {
		return nil, nil, "", err
	}
	var sb structures.SuperBlock

	if err != nil {
		return nil, nil, "", err
	}
	return &sb, partition, path, nil
}

func ObetenerParticionesMontadasRep(id string) (*structures.MBR, *structures.SuperBlock, string, error) {
	path := ParticionesMontadas[id]

	if path == "" {
		return nil, nil, "", errors.New("La partición no está montada")
	}
	var mbr structures.MBR
	err := mbr.Deserializar(path)

	if err != nil {
		return nil, nil, "", err
	}
	partition, err := mbr.Get_Partition_ID(id)

	if partition == nil {
		return nil, nil, "", err
	}
	var sb structures.SuperBlock

	err = sb.Deserializar(path, int64(partition.Part_inicio))

	if err != nil {
		return nil, nil, "", err
	}
	return &mbr, &sb, path, nil
}

// ObtenerListaParticionesMontadas devuelve una lista con los IDs de las particiones montadas
func ObtenerListaParticionesMontadas() []string {
	ids := make([]string, 0, len(ParticionesMontadas))

	for id := range ParticionesMontadas {
		ids = append(ids, id)
	}

	return ids
}

// MostrarParticionesMontadas imprime en consola las particiones montadas
func MostrarParticionesMontadas() {
	if len(ParticionesMontadas) == 0 {
		fmt.Println("No hay particiones montadas actualmente.")
		return
	}

	fmt.Println("Particiones montadas:")
	for id, path := range ParticionesMontadas {
		fmt.Printf("- ID: %s | Path: %s\n", id, path)
	}
}
