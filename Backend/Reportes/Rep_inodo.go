package Reportes

import (
	structures "Backend/estructuras"
	utils "Backend/utils"
	"fmt"
	"os"
	"os/exec"
	"time"
)

// reporte
func GenerateInodeReport(sb *structures.SuperBlock, diskFile string, outputPath string) error {
	// Crear directorios padre si no existen
	err := utils.CreateParentDirs(outputPath)
	if err != nil {
		return err
	}

	// Obtener nombres de archivo base
	dotFile, imgFile := utils.GetFileNames(outputPath)

	// Iniciar contenido DOT
	dotData := `digraph G {
        node [shape=plaintext]
    `

	// Recorrer los inodos
	for idx := int32(0); idx < sb.S_contador_inodos; idx++ {
		inodeData := &structures.Inodo{}
		// Deserializar inodo
		err := inodeData.Deserializar(diskFile, int64(sb.S_inicio_inodo+(idx*sb.S_tamano_inodo)))
		if err != nil {
			return err
		}

		accessTime := time.Unix(int64(inodeData.I_atime), 0).Format(time.RFC3339)
		createTime := time.Unix(int64(inodeData.I_ctime), 0).Format(time.RFC3339)
		modTime := time.Unix(int64(inodeData.I_mtime), 0).Format(time.RFC3339)

		// Construir contenido DOT para el inodo actual
		dotData += fmt.Sprintf(`inode%d [label=<
            <table border="0" cellborder="1" cellspacing="0">
                <tr><td colspan="2"> REPORTE INODO %d </td></tr>
                <tr><td>i_uid</td><td>%d</td></tr>
                <tr><td>i_gid</td><td>%d</td></tr>
                <tr><td>i_size</td><td>%d</td></tr>
                <tr><td>i_atime</td><td>%s</td></tr>
                <tr><td>i_ctime</td><td>%s</td></tr>
                <tr><td>i_mtime</td><td>%s</td></tr>
                <tr><td>i_type</td><td>%c</td></tr>
                <tr><td>i_perm</td><td>%s</td></tr>
                <tr><td colspan="2">BLOQUES DIRECTOS</td></tr>
            `, idx, idx, inodeData.I_uid, inodeData.I_gid, inodeData.I_size, accessTime, createTime, modTime, rune(inodeData.I_type[0]), string(inodeData.I_perm[:]))

		// agreagar bloques directos
		for blkIdx, blkValue := range inodeData.I_block {
			if blkIdx > 11 {
				break
			}
			dotData += fmt.Sprintf("<tr><td>%d</td><td>%d</td></tr>", blkIdx+1, blkValue)
		}

		// Agregar bloques indirectos
		dotData += fmt.Sprintf(`
                <tr><td colspan="2">BLOQUE INDIRECTO</td></tr>
                <tr><td>%d</td><td>%d</td></tr>
                <tr><td colspan="2">BLOQUE INDIRECTO DOBLE</td></tr>
                <tr><td>%d</td><td>%d</td></tr>
                <tr><td colspan="2">BLOQUE INDIRECTO TRIPLE</td></tr>
                <tr><td>%d</td><td>%d</td></tr>
            </table>>];
        `, 13, inodeData.I_block[12], 14, inodeData.I_block[13], 15, inodeData.I_block[14])

		// Enlace al siguiente inodo si no es el último
		if idx < sb.S_contador_inodos-1 {
			dotData += fmt.Sprintf("inode%d -> inode%d;\n", idx, idx+1)
		}
	}

	// Cerrar DOT
	dotData += "}"

	// Crear archivo DOT
	dotFileHandle, err := os.Create(dotFile)
	if err != nil {
		return err
	}
	defer dotFileHandle.Close()

	// Escribir datos en el archivo DOT
	_, err = dotFileHandle.WriteString(dotData)
	if err != nil {
		return err
	}

	// Generar imagen con Graphviz
	cmd := exec.Command("dot", "-Tpng", dotFile, "-o", imgFile)
	err = cmd.Run()
	if err != nil {
		return err
	}

	fmt.Println("Imagen del reporte de inodos generada:", imgFile)
	return nil
}
