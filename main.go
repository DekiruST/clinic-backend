package main

import (
	"clinic-backend/database"
	"clinic-backend/handlers"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error cargando .env")
	}

	// Conectar a Supabase
	if err := database.Connect(); err != nil {
		log.Fatal("Error conectando a Supabase:", err)
	}

	app := fiber.New()

	// Rutas para Pacientes
	pacienteGroup := app.Group("/pacientes")
	pacienteGroup.Post("/", handlers.CreatePaciente)
	pacienteGroup.Get("/", handlers.GetPacientes)
	pacienteGroup.Get("/:id", handlers.GetPaciente)
	pacienteGroup.Put("/:id", handlers.UpdatePaciente)
	pacienteGroup.Delete("/:id", handlers.DeletePaciente)

	// Rutas para Expedientes
	expedienteGroup := app.Group("/expedientes")
	expedienteGroup.Post("/", handlers.CreateExpediente)
	expedienteGroup.Get("/paciente/:id_paciente", handlers.GetExpedienteByPaciente)
	expedienteGroup.Put("/:id", handlers.UpdateExpediente)

	// Rutas para Consultorios
	consultorioGroup := app.Group("/consultorios")
	consultorioGroup.Post("/", handlers.CreateConsultorio)
	consultorioGroup.Get("/", handlers.GetConsultorios)
	consultorioGroup.Get("/:id", handlers.GetConsultorio)
	consultorioGroup.Put("/:id", handlers.UpdateConsultorio)
	consultorioGroup.Delete("/:id", handlers.DeleteConsultorio)

	// Rutas para Usuarios
	usuarioGroup := app.Group("/usuarios")
	usuarioGroup.Post("/", handlers.CreateUsuario)
	usuarioGroup.Get("/", handlers.GetUsuarios)
	usuarioGroup.Get("/:id", handlers.GetUsuario)
	usuarioGroup.Put("/:id", handlers.UpdateUsuario)
	usuarioGroup.Delete("/:id", handlers.DeleteUsuario)

	// Rutas para Consultas
	consultaGroup := app.Group("/consultas")
	consultaGroup.Post("/", handlers.CreateConsulta)
	consultaGroup.Get("/", handlers.GetConsultas)
	consultaGroup.Get("/paciente/:id_paciente", handlers.GetConsultasByPaciente)
	consultaGroup.Put("/:id", handlers.UpdateConsulta)
	consultaGroup.Delete("/:id", handlers.DeleteConsulta)

	// Rutas para Recetas
	recetaGroup := app.Group("/recetas")
	recetaGroup.Post("/", handlers.CreateReceta)
	recetaGroup.Get("/consulta/:id_consulta", handlers.GetRecetasByConsulta)

	// Rutas para Detalles de Receta
	detalleGroup := app.Group("/detalles-receta")
	detalleGroup.Post("/", handlers.CreateDetalleReceta)
	detalleGroup.Get("/receta/:id_receta", handlers.GetDetallesByReceta)

	// Rutas para Horarios
	horarioGroup := app.Group("/horarios")
	horarioGroup.Post("/", handlers.CreateHorario)
	horarioGroup.Get("/medico/:id_medico", handlers.GetHorariosByMedico)
	horarioGroup.Delete("/:id", handlers.DeleteHorario)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Fatal(app.Listen(":" + port))
}
