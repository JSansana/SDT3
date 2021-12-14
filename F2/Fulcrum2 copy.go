/*
package main

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	pb "github.com/JSansana/SD-T3/proto"
	"google.golang.org/grpc"
)

const (
	puertoMos        = ":50062"
	puertoInformante = ":50082"
	puertoFulcrum    = ":50072"
)

type InfoTo_FulcrumServer struct {
	pb.UnimplementedInfoTo_FulcrumServer
}
type Mos_FulcrumServer struct {
	pb.UnimplementedMos_FulcrumServer
}

type PlanetVector struct {
	Planeta     string
	VectorReloj []int
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func vaciar_archivo(planetin string) {
	arch := planetin + ".txt"
	/*
		var file, err = os.OpenFile(arch, os.O_RDWR, 0644)
		if isError(err) {
			return
		}
	if err := os.Truncate(arch, 0); err != nil {
		log.Printf("Failed to truncate: %v", err)
	}
	return
}

func ObtenerCiudades_Central(planetin string) ([]string, []int) {
	var Ciudades []string
	var Soldados []int
	arch := planetin + ".txt"
	input, err := ioutil.ReadFile(arch)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		linea_actual := strings.Split(line, " ")

		Ciudades = append(Ciudades, linea_actual[1])
		sold := strconv.Itoa(linea_actual[2])
		Soldados = append(Soldados, sold)

	}

	return Ciudades, Soldados
}

//Funcion para actualizar nombre o numero de soldados en determinado planeta
//Puede usarse en DELETE usando un espacio vacío en vez de comando
func actualizar_archivo_planeta(planeta string, ciudad string, comando string) {
	arch := planeta + ".txt"
	input, err := ioutil.ReadFile(arch)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, ciudad) {
			lines[i] = comando
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(arch, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

	return

}

//Función para comprobar que existe una ciudad. Retorna la cantidad de soldados si existe, de caso contrario retorna -1
func leer_archivo_planeta(planeta string, ciudad string) int {
	arch := planeta + ".txt"
	soldados := -1
	f, err := os.Open(arch)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), ciudad) {
			soldier := scanner.Scan()[len(scanner.Scan())-1]
			soldados, err = strconv.Atoi(soldier)
		}

	}
	defer f.close()
	return soldados
}

//Función que crea el archivo de planeta y agrega la primera ciudad con sus respectivos soldados. Usada en AddCity
func crear_archivo_planeta(planeta string, ciudad string, soldados int) {
	arch := planeta + ".txt"
	soldiers := strconv.Itoa(soldados)
	escritura := planeta + " " + ciudad + " " + soldiers + "\n"
	f, err := os.Create(arch)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Archivo: ", planeta, " creado.")
	_, err2 := f.WriteString(escritura)
	if err2 != nil {
		log.Fatal(err2)
	}
	defer f.Close()

	return
}

//Escribe en la ultima linea del archivo. Usada en AddCity
func abrir_escribir_archivo(planeta string, ciudad string, soldados int) {
	arch := planeta + ".txt"
	soldiers := strconv.Itoa(soldados)
	escritura := planeta + " " + ciudad + " " + soldiers + "\n"
	fmt.Println("Abriendo archivo: ", planeta)
	var file, err = os.OpenFile(arch, os.O_RDWR, 0644)
	if isError(err) {
		return
	}
	defer file.Close()
	_, err2 := f.WriteString(escritura)
	if err2 != nil {
		log.Fatal(err2)
	}
	defer f.Close()
	return
}

//Servicios Fulcrum Central - Fulcrum
//rpc GetPlanetas(SolicitudDominante) returns (RetornarPlanetas)
func (s *NodoDominante_NodoServer) GetPlanetas(ctx context.Context, in *pb.SolicitudDominante) (*pb.RetornarPlanetas, error) {

	return &pb.RetornarPlanetas{Vectores_Planetas: Planetas_Vectores}, nil
}

func (s *NodoDominante_NodoServer) GetCiudades(ctx context.Context, in *pb.SolicitudDominante) (*pb.RetornarCiudades, error) {
	planetin := in.GetPlaneta()
	ciudadines, soldadines := ObtenerCiudades_Central(planetin)
	vaciar_archivo(planetin)
	return &pb.RetornarCiudades{ciudades: ciudadines, soldados: soldadines}, nil
}

//rpc AgregarCiudad(MessageINF) returns (RespuestaMos)
func (s *NodoDominante_NodoServer) AgregarCiudad(ctx context.Context, in *pb.MessageINF) (*pb.RespuestaMos, error) {
	planetin := in.GetPlaneta()
	ciudadin := in.GetCiudad()
	soldadines := in.GetSoldados()
	booleano := true
	for i := 0; i < len(Planetas_Vectores); i++ {
		if Planetas_Vectores[i].Planeta == planetin {
			booleano := false
		}
	}

	if booleano {
		crear_archivo_planeta(planetin, ciudadin, soldadines)
	} else {
		abrir_escribir_archivo(planetin, ciudadin, soldadines)
	}

	return &pb.RespuestaMos{direccion: "Hola"}, nil
}

//
func (s *NodoDominante_NodoServer) ModificarVector(ctx context.Context, in *pb.VectorNuevo) (*pb.RespuestaMos, error) {
	planetin := in.GetPlaneta()
	vetorcin := in.GetNewVector()
	for i := 0; i < len(Planetas_Vectores); i++ {
		if Planetas_Vectores[i].Planeta == planetin {
			Planetas_Vectores[i].VectorReloj = vectorcin
		}
	}

	return &pb.RespuestaMos{direccion: "Hola"}, nil
}

// Servicios Mos - Fulcrum

//funcion para entregar la cantidad de soldados de determinada ciudad
func (s *Mos_FulcrumServer) GetSoldados(ctx context.Context, in *pb.MessageLeia) (*pb.LeiaResponse, error) {

	ciudadin := in.GetCiudad()
	planetin := in.GetPlaneta()
	VectoraGuardar := []int{-1, -1, -1}
	for i := 0; i < len(Planetas_Vectores); i++ {
		if Planetas_Vectores[i].Planeta == planetin {
			VectoraGuardar = Planetas_Vectores[i].VectorReloj
		}
	}

	soldiers := leer_archivo_planeta(planetin, ciudadin)
	return &pb.LeiaResponse{soldados: soldiers, vector: VectoraGuardar, direccion: "localhost:50061"}, nil

}

//Funcion para entregar el vector del planeta
func (s *Mos_FulcrumServer) GetVector(ctx context.Context, in *pb.Solicitud) (*pb.Reloj, error) {

	planetin := in.GetPlaneta()
	VectoraGuardar := []int{-1, -1, -1}
	for i := 0; i < len(Planetas_Vectores); i++ {
		if Planetas_Vectores[i].Planeta == planetin {
			VectoraGuardar = Planetas_Vectores[i].VectorReloj
		}
	}

	return &pb.Reloj{vector: VectoraGuardar, planeta: planetin}, nil

}

//Servicios Informantes - Fulcrum

func (s *InfoTo_FulcrumServer) AddCity(ctx context.Context, in *pb.MessageINF) (*pb.Reloj, error) {
	/*
			string planeta = 1;
		    string ciudad = 2;
		    string newcity = 3;
		    int32 soldados = 4;
	planetin := in.GetPlaneta()
	ciudadin := in.GetCiudad()
	soldados := in.GetSoldados()
	planetaEncontrado := false
	IndicePlaneta := 0
	for i := 0; i < len(Planetas_Vectores); i++ {
		if Planetas_Vectores[i].Planeta == planetin {
			planetaEncontrado = true
			IndicePlaneta = i
		}
	}

	if planetaEncontrado {
		//Si no encuentra la ciudad la agrega
		if leer_archivo_planeta(planetin, ciudadin) < 0 {
			abrir_escribir_archivo(planetin, ciudadin, soldados)
			Planetas_Vectores[IndicePlaneta].VectorReloj[0] += 1
			return &pb.Reloj{vector: Planetas_Vectores[IndicePlaneta].VectorReloj, planeta: planetin}, nil
		} else {
			VectoraGuardar := []int{-1, -1, -1}
			Println("La ciudad ya existe, comando AddCity no válido")
			return &pb.Reloj{vector: VectoraGuardar, planeta: planetin}, nil
		}
	}

	crear_archivo_planeta(planetin, ciudadin, soldados)
	VectoraGuardar := []int{0, 0, 0}
	Planetas_Vectores = append(Planetas_Vectores, PlanetVector{Planeta: planetin, VectorReloj: VectoraGuardar})
	return &pb.Reloj{vector: VectoraGuardar, planeta: planetin}, nil

}

func (s *InfoTo_FulcrumServer) UpdateName(ctx context.Context, in *pb.MessageINF) (*pb.Reloj, error) {
	planetin := in.GetPlaneta()
	ciudadin := in.GetCiudad()
	nuevaCiudad := in.GetNewcity()

	linea_texto := planetin + " " + nuevaCiudad + " "

	planetaEncontrado := false
	IndicePlaneta := 0
	for i := 0; i < len(Planetas_Vectores); i++ {
		if Planetas_Vectores[i].Planeta == planetin {
			planetaEncontrado = true
			IndicePlaneta = i
		}
	}

	if planetaEncontrado {
		//Si encuentra la ciudad la actualiza
		cantidad_soldados := leer_archivo_planeta(planetin, ciudadin)
		if cantidad_soldados >= 0 {
			linea_texto = linea_texto + strconv.Itoa(cantidad_soldados)
			actualizar_escribir_archivo(planetin, ciudadin, linea_texto)
			Planetas_Vectores[IndicePlaneta].VectorReloj[0] += 1
			Println("Nombre de la ciudad ", ciudadin, " actualizado a ", nuevaCiudad)
			return &pb.Reloj{vector: Planetas_Vectores[IndicePlaneta].VectorReloj, planeta: planetin}, nil
		} else {
			VectoraGuardar := []int{-1, -1, -1}
			Println("La ciudad ", ciudadin, " en ", planetin, " no existe, comando UpdateName no válido")
			return &pb.Reloj{vector: VectoraGuardar, planeta: planetin}, nil
		}
	}

	VectoraGuardar := []int{-1, -1, -1}
	Println("El planeta ", planetin, " no existe, comando UpdateName no válido")
	return &pb.Reloj{vector: VectoraGuardar, planeta: planetin}, nil

}

func (s *InfoTo_FulcrumServer) UpdateNumber(ctx context.Context, in *pb.MessageINF) (*pb.Reloj, error) {
	planetin := in.GetPlaneta()
	ciudadin := in.GetCiudad()
	nuevosSoldados := in.GetSoldados()

	linea_texto := planetin + " " + ciudadin + " " + strconv.Itoa(nuevosSoldados)

	planetaEncontrado := false
	IndicePlaneta := 0
	for i := 0; i < len(Planetas_Vectores); i++ {
		if Planetas_Vectores[i].Planeta == planetin {
			planetaEncontrado = true
			IndicePlaneta = i
		}
	}

	if planetaEncontrado {
		//Si encuentra la ciudad la actualiza
		if leer_archivo_planeta(planetin, ciudadin) >= 0 {
			actualizar_escribir_archivo(planetin, ciudadin, linea_texto)
			Planetas_Vectores[IndicePlaneta].VectorReloj[0] += 1
			Println("Numero de soldados de la ciudad ", ciudadin, " actualizado a ", nuevosSoldados)
			return &pb.Reloj{vector: Planetas_Vectores[IndicePlaneta].VectorReloj, planeta: planetin}, nil
		} else {
			VectoraGuardar := []int{-1, -1, -1}
			Println("La ciudad ", ciudadin, " en ", planetin, " no existe, comando UpdateNumber no válido")
			return &pb.Reloj{vector: VectoraGuardar, planeta: planetin}, nil
		}
	}

	VectoraGuardar := []int{-1, -1, -1}
	Println("El planeta ", planetin, " no existe, comando UpdateNumber no válido")
	return &pb.Reloj{vector: VectoraGuardar, planeta: planetin}, nil
}

func (s *InfoTo_FulcrumServer) DeleteCity(ctx context.Context, in *pb.MessageINF) (*pb.Reloj, error) {
	planetin := in.GetPlaneta()
	ciudadin := in.GetCiudad()

	planetaEncontrado := false
	IndicePlaneta := 0
	for i := 0; i < len(Planetas_Vectores); i++ {
		if Planetas_Vectores[i].Planeta == planetin {
			planetaEncontrado = true
			IndicePlaneta = i
		}
	}

	if planetaEncontrado {
		//Si encuentra la ciudad la actualiza
		if leer_archivo_planeta(planetin, ciudadin) >= 0 {
			actualizar_escribir_archivo(planetin, ciudadin, "")
			Planetas_Vectores[IndicePlaneta].VectorReloj[0] += 1
			Println("Ciudad ", ciudadin, " eliminada ")
			return &pb.Reloj{vector: Planetas_Vectores[IndicePlaneta].VectorReloj, planeta: planetin}, nil
		} else {
			VectoraGuardar := []int{-1, -1, -1}
			Println("La ciudad ", ciudadin, " en ", planetin, " no existe, comando DeleteCity no válido")
			return &pb.Reloj{vector: VectoraGuardar, planeta: planetin}, nil
		}
	}

	VectoraGuardar := []int{-1, -1, -1}
	Println("El planeta ", planetin, " no existe, comando DeleteCity no válido")
	return &pb.Reloj{vector: VectoraGuardar, planeta: planetin}, nil
}

func ServidorInformantes() {
	lis, err := net.Listen("tcp", puertoInformante)
	if err != nil {
		log.Fatalf("Fallo al escuchar: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterInfoTo_FulcrumServer(s, &InfoTo_FulcrumServer{})
	log.Printf("Servidor para informantes escuchando en %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Fallo en serve: %v", err)
	}
}

func ServidorMos() {
	lis, err := net.Listen("tcp", puertoMos)
	if err != nil {
		log.Fatalf("Fallo al escuchar: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterMos_FulcrumServer(s, &Mos_FulcrumServer{})
	log.Printf("Servidor para informantes escuchando en %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Fallo en serve: %v", err)
	}
}

func ServidorFulcrum_Central() {
	lis, err := net.Listen("tcp", puertoFulcrum)
	if err != nil {
		log.Fatalf("Fallo al escuchar: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterNodoDominante_NodoServer(s, &NodoDominante_NodoServer{})
	log.Printf("Servidor para informantes escuchando en %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Fallo en serve: %v", err)
	}
}

var Planetas_Vectores []PlanetVector

func main() {
	servidorMos()
	go servidorInformantes()
	go ServidorFulcrum_Central()
}*/
