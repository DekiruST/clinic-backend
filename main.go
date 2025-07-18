package main

import (
	"log"
	"os"
	"time"

	"clinic-backend/auth"
	"clinic-backend/database"
	"clinic-backend/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error cargando .env")
	}
	if err := database.Connect(); err != nil {
		log.Fatal("Error conectando a Supabase:", err)
	}

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:4200",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// 🔐 Auth routes con Rate Limit
	authGroup := app.Group("/auth", limiter.New(limiter.Config{
		Max:        5,               // Máximo 10 peticiones
		Expiration: 1 * time.Minute, // Cada 10 minutos
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(429).JSON(fiber.Map{
				"success": false,
				"message": "Demasiados intentos de autenticación. Espera un momento.",
			})
		},
	}))

	authGroup.Post("/register", auth.Register)
	authGroup.Post("/login", auth.Login)
	authGroup.Post("/refresh", auth.Refresh)
	authGroup.Post("/totp/generate", auth.ValidateJWT, auth.GenerateTOTPSecret)
	authGroup.Get("/totp/qrcode", auth.ValidateJWT, auth.GetTOTPQRCode)
	authGroup.Get("/totp/qrcode/view/:id_usuario", auth.ServeTOTPQRViewByID)
	authGroup.Post("/totp/verify", auth.ValidateJWT, auth.VerifyTOTP)

	// Pacientes
	pacienteGroup := app.Group("/pacientes", auth.ValidateJWT, auth.OnlyPaciente)
	pacienteGroup.Post("/", auth.RequirePermission("pacientes:create"), handlers.CreatePaciente)
	pacienteGroup.Get("/", auth.RequirePermission("pacientes:read"), handlers.GetPacientes)
	pacienteGroup.Get("/me", auth.RequirePermission("pacientes:read"), handlers.GetPacienteByToken)
	pacienteGroup.Get("/:id", auth.RequirePermission("pacientes:read"), handlers.GetPaciente)
	pacienteGroup.Put("/:id", auth.RequirePermission("pacientes:update"), handlers.UpdatePaciente)
	pacienteGroup.Delete("/:id", auth.RequirePermission("pacientes:delete"), handlers.DeletePaciente)
	pacienteGroup.Get("/consultas", auth.RequirePermission("consultas:read"), handlers.GetConsultasByPaciente)

	// Expedientes
	app.Get("/expedientes", handlers.GetAllExpedientes)

	// Consultorios
	consultorioGroup := app.Group("/consultorios", auth.ValidateJWT)
	consultorioGroup.Post("/", auth.RequirePermission("consultorios:create"), handlers.CreateConsultorio)
	consultorioGroup.Get("/", auth.RequirePermission("consultorios:read"), handlers.GetConsultorios)
	consultorioGroup.Get("/:id", auth.RequirePermission("consultorios:read"), handlers.GetConsultorio)
	consultorioGroup.Put("/:id", auth.RequirePermission("consultorios:update"), handlers.UpdateConsultorio)
	consultorioGroup.Delete("/:id", auth.RequirePermission("consultorios:delete"), handlers.DeleteConsultorio)

	// Usuarios
	usuarioGroup := app.Group("/usuarios", auth.ValidateJWT)
	usuarioGroup.Post("/", auth.RequirePermission("usuarios:create"), handlers.CreateUsuario)
	usuarioGroup.Get("/", auth.RequirePermission("usuarios:read"), handlers.GetUsuarios)
	usuarioGroup.Get("/:id", auth.RequirePermission("usuarios:read"), handlers.GetUsuario)
	usuarioGroup.Put("/:id", auth.RequirePermission("usuarios:update"), handlers.UpdateUsuario)
	usuarioGroup.Delete("/:id", auth.RequirePermission("usuarios:delete"), handlers.DeleteUsuario)

	// Consultas
	consultaGroup := app.Group("/consultas", auth.ValidateJWT)
	consultaGroup.Post("/", auth.RequirePermission("consultas:create"), handlers.CreateConsulta)
	consultaGroup.Get("/", auth.RequirePermission("consultas:read"), handlers.GetConsultas)
	consultaGroup.Get("/paciente/:id_paciente", auth.RequirePermission("consultas:read"), handlers.GetConsultasByPaciente)
	consultaGroup.Put("/:id", auth.RequirePermission("consultas:update"), handlers.UpdateConsulta)
	consultaGroup.Delete("/:id", auth.RequirePermission("consultas:delete"), handlers.DeleteConsulta)
	consultaGroup.Get("/paciente-id/:id_paciente", auth.RequirePermission("consultas:read"), handlers.GetConsultasByPacienteIDParam)
	consultaGroup.Get("/paciente", auth.RequirePermission("consultas:read"), handlers.GetConsultasByPaciente)

	// Recetas
	detalleGroup := app.Group("/detalles-receta", auth.ValidateJWT)
	detalleGroup.Post("/", auth.RequirePermission("detalles_receta:create"), handlers.CreateDetalleReceta)
	detalleGroup.Get("/receta/:id_receta", auth.RequirePermission("detalles_receta:read"), handlers.GetDetallesByReceta)
	recetaGroup := app.Group("/recetas", auth.ValidateJWT)
	recetaGroup.Get("/paciente/:id_paciente", handlers.GetRecetasByPaciente)
	consultas := app.Group("/consultas")
	consultas.Get("/medico/:id_medico", auth.RequirePermission("consultas:read"), handlers.GetConsultasMedico)
	// Detalles de Receta
	detalleGroup = app.Group("/detalles-receta", auth.ValidateJWT)
	detalleGroup.Post("/", handlers.CreateDetalleReceta)
	detalleGroup.Get("/receta/:id_receta", handlers.GetDetallesByReceta)

	// Horarios
	horarioGroup := app.Group("/horarios", auth.ValidateJWT)
	horarioGroup.Post("/", auth.RequirePermission("horarios:create"), handlers.CreateHorario)
	horarioGroup.Get("/medico/:id_medico", auth.RequirePermission("horarios:read"), handlers.GetHorariosByMedico)
	horarioGroup.Delete("/:id", auth.RequirePermission("horarios:delete"), handlers.DeleteHorario)

	// Enfermera
	enfermeraGroup := app.Group("/enfermera", auth.ValidateJWT)
	enfermeraGroup.Get("/pacientes", auth.RequirePermission("pacientes:read"), handlers.GetPacientesParaEnfermera)

	// Médico
	medicoGroup := app.Group("/medico", auth.ValidateJWT)
	medicoGroup.Get("/consultas", auth.RequirePermission("consultas:read"), handlers.GetConsultasMedico)
	medicoGroup.Get("/recetas", auth.RequirePermission("recetas:read"), handlers.GetRecetasMedico)
	medicoGroup.Get("/expedientes", auth.RequirePermission("expedientes:read"), handlers.GetExpedientesMedico)

	authGroup.Post("/totp/reset", auth.ValidateJWT, auth.ResetTOTPSecret)
	dashboardPacienteGroup := app.Group("/dashboard-paciente", auth.ValidateJWT)
	authGroup.Post("/login/init", auth.InitLogin)
	dashboardPacienteGroup.Get("/citas", handlers.GetMisCitas)
	dashboardPacienteGroup.Get("/citas/historico", handlers.GetHistoricoCitas)
	dashboardPacienteGroup.Get("/recetas", handlers.GetMisRecetas)
	dashboardPacienteGroup.Post("/citas/solicitar", handlers.SolicitarCita)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Fatal(app.Listen(":" + port))
}
