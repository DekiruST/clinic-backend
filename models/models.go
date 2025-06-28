package models

type RolUsuario string
type TipoConsulta string
type Turno string
type CategoriaConsultorio string

const (
	RolMedico    RolUsuario = "médico"
	RolEnfermera RolUsuario = "enfermera"
	RolAdmin     RolUsuario = "admin"
	RolPaciente  RolUsuario = "paciente"
)

type TipoConsultaEnum string

const (
	ConsultaGeneral      TipoConsultaEnum = "general"
	ConsultaEspecialidad TipoConsultaEnum = "especialidad"
	ConsultaUrgencia     TipoConsultaEnum = "urgencia"
	ConsultaControl      TipoConsultaEnum = "control"
)

type TurnoEnum string

const (
	TurnoManana TurnoEnum = "mañana"
	TurnoTarde  TurnoEnum = "tarde"
	TurnoNoche  TurnoEnum = "noche"
)

type CategoriaConsultorioEnum string

const (
	CategoriaPrimaria      CategoriaConsultorioEnum = "primaria"
	CategoriaEspecializada CategoriaConsultorioEnum = "especializada"
	CategoriaDiagnostico   CategoriaConsultorioEnum = "diagnóstico"
)

type Paciente struct {
	ID       int     `json:"id_paciente"`
	Nombre   string  `json:"nombre"`
	Seguro   *string `json:"seguro,omitempty"`
	Contacto string  `json:"contacto"`
}

type Expediente struct {
	ID               int    `json:"id_expediente"`
	Antecedentes     string `json:"antecedentes,omitempty"`
	HistorialClinico string `json:"historial_clinico"`
	IDPaciente       int    `json:"id_paciente"`
}

type Consultorio struct {
	ID           int                      `json:"id_consultorio"`
	Nombre       string                   `json:"nombre"`
	Ubicacion    string                   `json:"ubicacion"`
	Especialidad string                   `json:"especialidad"`
	Categoria    CategoriaConsultorioEnum `json:"categoria"`
}

type Usuario struct {
	ID           int        `json:"id_usuario"`
	Rol          RolUsuario `json:"rol"`
	Especialidad *string    `json:"especialidad,omitempty"`
	IDPaciente   *int       `json:"id_paciente,omitempty"`
}

type Consulta struct {
	ID            int              `json:"id_consulta"`
	Tipo          TipoConsultaEnum `json:"tipo"`
	Horario       string           `json:"horario"`
	Diagnostico   *string          `json:"diagnostico,omitempty"`
	Costo         float64          `json:"costo"`
	IDConsultorio *int             `json:"id_consultorio,omitempty"`
	IDPaciente    int              `json:"id_paciente"`
	IDMedico      *int             `json:"id_medico,omitempty"`
}

type Receta struct {
	ID         int    `json:"id_receta"`
	Fecha      string `json:"fecha"`
	IDConsulta int    `json:"id_consulta"`
}

type DetalleReceta struct {
	ID          int    `json:"id_detalle"`
	Medicamento string `json:"medicamento"`
	Dosis       string `json:"dosis"`
	IDReceta    int    `json:"id_receta"`
}

type Horario struct {
	ID            int       `json:"id_horario"`
	Turno         TurnoEnum `json:"turno"`
	IDConsultorio int       `json:"id_consultorio"`
	IDMedico      int       `json:"id_medico"`
}
