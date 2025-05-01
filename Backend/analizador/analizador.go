package analizador

import (
	"Backend/comandos"
	global "Backend/global"
	"errors" // Importa el paquete "errors" para manejar errores
	"fmt"    // Importa el paquete "fmt" para formatear e imprimir texto
	"os"
	"os/exec"
	"strings" // Importa el paquete "strings" para manipulación de cadenas
)

/*
	Aqui se va analizar el comando enviaado y se va ejecutar
*/

func Analyzer2(input string) (string, error) {

	//Separa la entrada -> en tokens
	tokens := strings.Fields(input)

	// revisa que si hayan enviado tokens a evaluar y si no hay regresa erorf
	if len(tokens) == 0 {
		return "", errors.New("no se proporcionó ningún comando")
	}

	// Switch para manejar diferentes comandos
	switch tokens[0] {

	// --->> COMANDO 1  --* mkdisk -size=3000 -unit=K -path=/home/luisa/Música/ArchivoP/Disco1.mia​
	case "mkdisk": // -> crear discos
		return comandos.Parsermkdisk(tokens[1:])
		//return fmt.Println("función mkdisk", tokens[1:])

	// --->> COMANDO 2   --*  rmdisk -path="/home/mis discos/Disco4.mia"
	case "rmdisk": // ->  eliminar discos
		return comandos.ParserRmdisk(tokens[1:])

	// --->> COMANDO 3	--* fdisk -Size=300 -path=/home/Disco1.mia -name=Particion1 ​
	case "fdisk": // -> crea particiones  -= P,E,L
		return comandos.ParserFdisk(tokens[1:])

	// --->> COMANDO 4	--* mount -path=/home/Disco2.mia -name=Part2 #id=341A
	case "mount": // ->  montar particiones
		return comandos.Parse_mount(tokens[1:])

	// --->> COMANDO 5	--* mounted
	case "mounted": // ->  mostrar id  particiones montadas
		global.MostrarParticionesMontadas()
		return strings.Join(global.ObtenerListaParticionesMontadas(), ", "), nil

	case "mkfs":
		return comandos.ParserMkfs(tokens[1:])

	case "rep":

		return comandos.ParserRep(tokens[1:])

	case "clear":
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			return "", errors.New("no se pudo limpiar la terminal")
		}
		return "", nil

	default:
		// Si el comando no es reconocido, devuelve un error
		return "", fmt.Errorf("comando desconocido: %s", tokens[0])

	}
}
