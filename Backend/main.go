package main

import (
	analizador "Backend/analizador"
	"time"

	//analyzer "backend/analyzer" // Importa el paquete "analyzer" desde el directorio "backend/analyzer"
	"fmt"
	"log"     // Importa el paquete "log" para registrar mensajes de error
	"strings" // Importa el paquete "strings" para manipulación de cadenas

	"github.com/gofiber/fiber/v2"                 // Importa el paquete Fiber para crear la API
	"github.com/gofiber/fiber/v2/middleware/cors" // Importa el middleware CORS para manejar CORS
)

func main() {
	// Crear una nueva instancia de Fiber
	app := fiber.New()

	// Configurar el middleware CORS
	app.Use(cors.New())
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"time":   time.Now(),
		})
	})

	// Definir la ruta POST para recibir el comando del usuario
	app.Post("/analyze", func(c *fiber.Ctx) error {
		// Estructura para recibir el JSON
		type Request struct {
			Command string `json:"command"`
		}

		// Crear una instancia de Request
		var req Request

		// Parsear el cuerpo de la solicitud como JSON
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid JSON",
			})
		}

		// Obtener el comando del cuerpo de la solicitud
		input := req.Command
		fmt.Println("input: ", input)

		// Separar el comando en líneas
		lines := strings.Split(input, "\n")

		// Lista para acumular los resultados de salida
		var results []string

		// Analizar cada línea
		for _, line := range lines {
			// Ignorar líneas vacías
			if strings.TrimSpace(line) == "" {
				continue
			}

			// Llamar a la función Analyzer del paquete analyzer para analizar la línea
			result, err := analizador.Analyzer2(line)
			if err != nil {
				// Si hay un error, almacenar el mensaje de error en lugar del resultado
				result = fmt.Sprintf("Error: %s", err.Error())
			}

			// Acumular los resultados
			results = append(results, result)
		}

		// Devolver una respuesta JSON con la lista de resultados
		return c.JSON(fiber.Map{
			"results": results,
		})
	})

	// Iniciar el servidor en el puerto 3000
	log.Fatal(app.Listen(":3000"))
}

/*
import (
	"Backend/analizador"

	"bufio" // para operaciones de buffer de entrada/salida
	"fmt"   // "fmt" para formatear e imprimir texto
	"os"    // para interactuar con el sistema operativo
)

func main() {

	scanner := bufio.NewScanner(os.Stdin)

	// for  para leer los comandos ingresados
	for {

		fmt.Print(">>> ") // prompt para el usuario

		// Lee la siguiente línea de entrada

		if !scanner.Scan() { // Si no hay más líneas para leer, rompe el bucle
			break
		}

		input := scanner.Text() // Obtiene el texto ingresado por el usuario

		// Llama a la función Analyzer del paquete analyzer para analizar el comando ingresado
		_, err := analizador.Analyzer2(input)
		if err != nil {
			// Si hay un error al analizar el comando, imprime el error y continúa con el siguiente comando
			fmt.Println("Error:", err)
			continue
		}

		// Comentado: Aquí podrías imprimir el comando analizado
		// fmt.Printf("Parsed Command: %+v\n", cmd)
	}

	// Verifica si hubo algún error al leer la entrada
	if err := scanner.Err(); err != nil {
		// Si hubo un error al leer la entrada, lo imprime
		fmt.Println("Error al leer:", err)
	}
}
*/
