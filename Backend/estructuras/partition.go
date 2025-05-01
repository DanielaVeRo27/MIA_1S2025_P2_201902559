package estructuras

/*
	archivo en cual definimos la estructura de la partición
		--* la partición es un bloque de información que se almacena en el disco

----------->>> ESTADOS DE LA PARTICION
	-> N: Disponible
	-> 0: Creado
	-> 1: Montado
*/

import "fmt"

type Partition struct {
	Part_estado      [1]byte  // Estado de la partición
	Part_tipo        [1]byte  // Tipo de partición
	Part_fit         [1]byte  // Ajuste de la partición
	Part_inicio      int32    // Byte de inicio de la partición
	Part_size        int32    // Tamaño de la partición
	Part_nombre      [16]byte // Nombre de la partición
	Part_correlativo int32    // Correlativo de la partición
	Part_id          [4]byte  // ID de la partición
}

// Creacion de partición
func (p *Partition) CrearParticion(partStart, partSize int, partType, partFit, partName string) {

	// estado de la partición 0
	p.Part_estado[0] = '0' // partición creada

	// en que byte inicia la partición
	p.Part_inicio = int32(partStart)

	// tamaño de la partición
	p.Part_size = int32(partSize)

	// tipo de partición E o P
	if len(partType) > 0 {
		p.Part_tipo[0] = partType[0]
	}

	//  ajuste de la partición F, B o W
	if len(partFit) > 0 {
		p.Part_fit[0] = partFit[0]
	}

	// nombre de la partición
	copy(p.Part_nombre[:], partName)
}

// Montar la particion por id
func (p *Partition) MountParticion(correlative int, id string) error {
	// estado de la partición
	p.Part_estado[0] = '1' // 1 -> particion montada

	// Asignar correlativo a la partición
	p.Part_correlativo = int32(correlative)

	// id particion
	copy(p.Part_id[:], id)

	return nil
}

// Imprimir particion
func (p *Partition) PrintPartition() {
	fmt.Printf(("Estado de la partición: %c\n"), p.Part_estado[0])
	fmt.Printf(("Tipo de partición: %c\n"), p.Part_tipo[0])
	fmt.Printf(("Ajuste de la partición: %c\n"), p.Part_fit[0])
	fmt.Printf(("Inicio de la partición: %d\n"), p.Part_inicio)
	fmt.Printf(("Tamaño de la partición: %d\n"), p.Part_size)
	fmt.Printf(("Nombre de la partición: %s\n"), string(p.Part_nombre[:]))
	fmt.Printf(("Correlativo de la partición: %d\n"), p.Part_correlativo)
	fmt.Printf(("ID de la partición: %s\n"), string(p.Part_id[:]))

}
