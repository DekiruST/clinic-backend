# Sistema de Gestión Hospitalaria - Backend

[![Go](https://img.shields.io/badge/Go-1.20+-00ADD8?logo=go)](https://golang.org/)
[![Fiber](https://img.shields.io/badge/Fiber-2.50.0-00ADD8)](https://gofiber.io/)
[![Supabase](https://img.shields.io/badge/Supabase-3.0.0-3ECF8E?logo=supabase)](https://supabase.io/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Backend para el sistema de gestión de citas y reportes de un hospital, desarrollado en Go con Fiber y Supabase como base de datos PostgreSQL en la nube. Implementa prácticas de seguridad avanzadas para el manejo de datos médicos sensibles.

## Características Principales

- 🏥 Gestión completa de pacientes, expedientes médicos y citas
- 🔐 Autenticación JWT 
- 📚 25 endpoints RESTful documentados

## Estructura del Proyecto
CLINIC-BACKEND/
├── database/
│ └── database.go # Conexión a Supabase
├── handlers/
│ ├── consulta.go # Controlador de consultas
│ ├── consultorio.go # Controlador de consultorios
│ ├── detalle_receta.go # Controlador de detalles de recetas
│ ├── expediente.go # Controlador de expedientes
│ ├── horario.go # Controlador de horarios
│ ├── paciente.go # Controlador de pacientes
│ ├── receta.go # Controlador de recetas
│ └── usuario.go # Controlador de usuarios
├── models/
│ └── models.go # Modelos de datos
├── utils/
│ └── response.go # Funciones de respuesta API
├── .env # Variables de entorno (no incluido en repo)
├── go.mod # Dependencias de Go
├── go.sum # Checksums de dependencias
├── main.go # Punto de entrada principal
└── README.md # Este archivo


## Requisitos Previos

- Go 1.20 o superior
- Cuenta de [Supabase](https://supabase.io/)
- PostgreSQL 15+
- Variables de entorno configuradas (.env)

## Configuración Inicial

1. **Clonar repositorio**:
   ```bash
   git clone https://github.com/DekiruST/clinic-backend.git
   cd clinic-backend

## Cambios Recientes
Para ver el historial completo de cambios, consulta el [CHANGELOG.md](CHANGELOG.md)
