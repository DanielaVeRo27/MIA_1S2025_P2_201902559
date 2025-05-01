package Reportes

import (
	estructura "Backend/estructuras"
	util "Backend/utils"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func ReporteMBR(mbr *estructura.MBR, path string) error {

	fmt.Println("ENTRANDO AL REPORTE DEL MBR")

	err := util.CreateParentDirs(path)

	if err != nil {
		return err
	}

	dotFileName, outputImage := util.GetFileNames(path)
	dotFileName += "_EBR.dot"
	outputImage += "_EBR.png"

	dotContent := fmt.Sprintf(`digraph G {
        node [shape=plaintext]
        tabla [label=<
            <table border="0" cellborder="1" cellspacing="0">
                <tr><td colspan="2" bgcolor="pink" > REPORTE MBR </td></tr>
                <tr><td bgcolor="lightblue">mbr_tamano</td><td bgcolor="lemonchiffon">%d</td></tr>
                <tr><td  bgcolor="lightblue" >mrb_fecha_creacion</td><td bgcolor="lemonchiffon">%s</td></tr>
                <tr><td  bgcolor="lightblue">mbr_disk_signature</td><td bgcolor="lemonchiffon">%d</td></tr>
            `, mbr.Size_mbr, time.Unix(int64(mbr.Creation_date_mbr), 0), mbr.Signature_mbr)

	for i, part := range mbr.Mbr_partitions {
		partName := strings.TrimRight(string(part.Part_nombre[:]), "\x00")
		partStatus := rune(part.Part_estado[0])
		partType := rune(part.Part_tipo[0])
		partFit := rune(part.Part_fit[0])

		dotContent += fmt.Sprintf(`
				<tr><td colspan="2" bgcolor="deepskyblue" > PARTICIÓN %d </td></tr>
				<tr><td bgcolor="lavender" >part_status</td><td bgcolor="paleturquoise">%c</td></tr>
				<tr><td bgcolor="paleturquoise" >part_type</td><td bgcolor="orchid">%c</td></tr>
				<tr><td bgcolor="lavender">part_fit</td><td bgcolor="paleturquoise">%c</td></tr>
				<tr><td bgcolor="paleturquoise">part_start</td><td bgcolor="orchid">%d</td></tr>
				<tr><td bgcolor="lavender">part_size</td><td bgcolor="paleturquoise">%d</td></tr>
				<tr><td bgcolor="paleturquoise">part_name</td><td bgcolor="orchid">%s</td></tr>
			`, i+1, partStatus, partType, partFit, part.Part_inicio, part.Part_size, partName)
	}

	dotContent += "</table>>] }"

	file, err := os.Create(dotFileName)

	if err != nil {
		return fmt.Errorf("Ocurrió un error al crear el archivo: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(dotContent)

	if err != nil {
		return fmt.Errorf("Error al escribir en el archivo %b", err)
	}
	cmd := exec.Command("dot", "-Tpng", dotFileName, "-o", outputImage)
	err = cmd.Run()

	if err != nil {
		return fmt.Errorf("Error al ejecutar el comando Graphviz: %v", err)
	}
	return nil
}
func ReporteEBR(ebrs []estructura.EBR, path string) error {
	fmt.Println("ENTRANDO AL REPORTE DEL EBR")

	err := util.CreateParentDirs(path)
	if err != nil {
		return err
	}

	dotFileName, outputImage := util.GetFileNames(path)

	dotContent := `digraph G {
        node [shape=plaintext]
        tabla [label=<
            <table border="0" cellborder="1" cellspacing="0">
                <tr><td colspan="2" bgcolor="pink"> REPORTE EBR </td></tr>`

	for i, ebr := range ebrs {
		partName := strings.TrimRight(string(ebr.Part_name[:]), "\x00")
		partMount := "-"
		partFit := "-"
		if len(ebr.Part_mount) > 0 && ebr.Part_mount[0] != 0 {
			partMount = string(ebr.Part_mount[0])
		}

		//partFit := "-"
		if len(ebr.Part_fit) > 0 && ebr.Part_fit[0] != 0 {
			partFit = string(ebr.Part_fit[0])
		}

		dotContent += fmt.Sprintf(`
                <tr><td colspan="2" bgcolor="deepskyblue"> PARTICIÓN LÓGICA %d </td></tr>
                <tr><td bgcolor="lavender">part_mount</td><td bgcolor="paleturquoise">%c</td></tr>
                <tr><td bgcolor="paleturquoise">part_fit</td><td bgcolor="orchid">%c</td></tr>
                <tr><td bgcolor="lavender">part_start</td><td bgcolor="paleturquoise">%d</td></tr>
                <tr><td bgcolor="paleturquoise">part_size</td><td bgcolor="orchid">%d</td></tr>
                <tr><td bgcolor="lavender">part_next</td><td bgcolor="paleturquoise">%d</td></tr>
                <tr><td bgcolor="paleturquoise">part_name</td><td bgcolor="orchid">%s</td></tr>
            `, i+1, partMount, partFit, ebr.Part_start, ebr.Part_size, ebr.Part_next, partName)
	}

	dotContent += "</table>>] }"

	file, err := os.Create(dotFileName)
	if err != nil {
		return fmt.Errorf("Ocurrió un error al crear el archivo: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(dotContent)
	if err != nil {
		return fmt.Errorf("Error al escribir en el archivo: %v", err)
	}

	cmd := exec.Command("dot", "-Tpng", dotFileName, "-o", outputImage)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("Error al ejecutar el comando Graphviz: %v", err)
	}

	return nil
}
