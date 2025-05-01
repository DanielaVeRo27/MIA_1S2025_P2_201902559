package comandos

import (
	Reports "Backend/Reportes"
	"Backend/estructuras"
	global "Backend/global"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type rep struct {
	id           string
	path         string
	name         string
	path_file_ls string
}

func ParserRep(tokens []string) (string, error) {
	cmd := &rep{}
	args := strings.Join(tokens, " ")
	re := regexp.MustCompile(`-id=[^\s]+|-path="[^"]+"|-path=[^\s]+|-name=[^\s]+|-path_file_ls="[^"]+"|-path_file_ls=[^\s]+`)
	matches := re.FindAllString(args, -1)

	if len(matches) != len(tokens) {
		for _, token := range tokens {
			if !re.MatchString(token) {
				return "", fmt.Errorf("El parámetro es inválido: %s", token)
			}
		}
	}

	for _, match := range matches {
		kv := strings.SplitN(match, "=", 2)
		if len(kv) != 2 {
			return "", fmt.Errorf("El formato del parámetro es inválido: %s", match)
		}
		key, value := strings.ToLower(kv[0]), kv[1]

		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
			value = strings.Trim(value, "\"")
		}

		switch key {
		case "-id":
			if value == "" {
				return "", errors.New("El id no puede estar vacío")
			}
			cmd.id = value
		case "-path":
			if value == "" {
				return "", errors.New("El path no puede estar vacío")
			}
			cmd.path = value
		case "-name":
			validNames := []string{"mbr", "disk", "inode", "block", "bm_inode", "bm_block", "sb", "file", "ls"}

			if !contains(validNames, value) {
				return "", errors.New("Nombre inválido")
			}
			cmd.name = value
		case "-path_file_ls":
			cmd.path_file_ls = value
		default:
			return "", fmt.Errorf("parámetro desconocido: %s", key)
		}
	}

	if cmd.id == "" || cmd.path == "" || cmd.name == "" {
		return "", errors.New("Faltan parámetros requeridos: -id, -path o -name")
	}

	err := commandRep(cmd)

	if err != nil {
		return "", err
	}
	return "REP: Reporte generado correctamente", nil
}

func contains(list []string, value string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

func commandRep(rep *rep) error {
	mountedMBR, mountedSb, mountedDiskPath, err := global.ObetenerParticionesMontadasRep(rep.id)

	if err != nil {
		return err
	}

	switch rep.name {
	case "mbr":
		err = Reports.ReporteMBR(mountedMBR, rep.path)
		// Buscar inicio de la partición extendida en el MBR
		var inicioEBR int32 = -1

		for _, part := range mountedMBR.Mbr_partitions {
			if rune(part.Part_tipo[0]) == 'E' { // Si es partición extendida
				inicioEBR = part.Part_inicio
				break
			}
		}

		if inicioEBR == -1 {
			return fmt.Errorf("No se encontró una partición extendida")
		}

		// Leer los EBR desde el disco
		ebrs, err := estructuras.LeerEBR(mountedDiskPath, inicioEBR)
		if err != nil {
			return fmt.Errorf("Error al leer el EBR: %v", err)
		}
		fmt.Printf("ebrs: %v\n", ebrs)
		// Generar el reporte
		Reports.ReporteEBR(ebrs, rep.path)

		if err != nil {
			return fmt.Errorf("Error al generar el reporte EBR: %v", err)
		}

		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	case "inode":
		fmt.Println(" entranso al reporte inodo")
		err = Reports.GenerateInodeReport(mountedSb, mountedDiskPath, rep.path)

		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	case "bm_inode":
		err = Reports.ReporteMBR(mountedMBR, rep.path)

		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}
	fmt.Println(mountedSb, mountedDiskPath)
	return nil
}
