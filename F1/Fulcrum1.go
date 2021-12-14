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
	"time"

	pb "github.com/JSansana/SDT3/proto"
	"google.golang.org/grpc"
)

const (
	puertoMos        = ":50061"
	puertoInformante = ":50081"
	addressFulcrum2  = "dist21:50072"
	addressFulcrum3  = "dist21:50073"
)

type InfoTo_FulcrumServer struct {
	pb.UnimplementedInfoTo_FulcrumServer
}
type Mos_FulcrumServer struct {
	pb.UnimplementedMos_FulcrumServer
}

type PlanetVector struct {
	Planeta     string
	VectorReloj []int32
}

func Existe_Planeta(planetin string) int {
	if len(Planetas_Vectores) == 0 {
		return -1
	}

	for i := 0; i < len(Planetas_Vectores); i++ {
		if Planetas_Vectores[i].Planeta == planetin {
			return i
		}
	}
	return -1
}
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func ModificarVector_Fulcrum(direccion string, planetin string, reloj []int32) {
	conn, err := grpc.Dial(direccion, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("No conectó: %v", err)
	}
	defer conn.Close()
	c := pb.NewNodoDominante_NodoClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	c.ModificarVector(ctx, &pb.VectorNuevo{NewVector: reloj, Planeta: planetin})
	return

}

func AgregarCiudades(direccion string, planetin string, ciudadin string, soldadines int) {
	conn, err := grpc.Dial(direccion, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("No conectó: %v", err)
	}
	defer conn.Close()
	c := pb.NewNodoDominante_NodoClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	c.AgregarCiudad(ctx, &pb.MessageINF{Planeta: planetin, Ciudad: ciudadin, Newcity: "Hola", Soldados: int32(soldadines)})

	return

}

func ObtenerCiudades_Central(planetin string) ([]string, []int) {
	var Ciudades []string
	var Soldados []int
	arch := planetin + ".txt"
	file, err := os.Open(arch)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		linea_actual := strings.Split(scanner.Text(), " ")
		if len(linea_actual) == 3 {
			Ciudades = append(Ciudades, linea_actual[1])
			sold, _ := strconv.Atoi(linea_actual[2])
			Soldados = append(Soldados, sold)
		} else {
			Ciudades = append(Ciudades, "Hola")
			Soldados = append(Soldados, -1)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return Ciudades, Soldados
}

func ObtenerCiudades_Fulcrum(direccion string, planetin string) ([]string, []int32) {
	conn, err := grpc.Dial(direccion, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("No conectó: %v", err)
	}
	defer conn.Close()
	c := pb.NewNodoDominante_NodoClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.GetCiudades(ctx, &pb.SolicitudDominante{Solicitud: 1, Planeta: planetin})
	Ciudades := r.GetCiudades()
	Soldados := r.GetSoldados()

	return Ciudades, Soldados
}

func ObtenerPlanetas_Fulcrum(direccion string) []PlanetVector {

	conn, err := grpc.Dial(direccion, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("No conectó: %v", err)
	}
	defer conn.Close()
	c := pb.NewNodoDominante_NodoClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.GetPlanetas(ctx, &pb.SolicitudDominante{Solicitud: 1, Planeta: "Hola"})
	Planetines := r.GetNombresplanetas()
	Relojines := r.GetRelojes()
	var Retorno []PlanetVector
	for i := 0; i < len(Planetines); i++ {
		var aux []int32
		txt := strings.Split(Relojines[i], ",")
		for j := 0; j < len(txt); j++ {
			Numero, _ := strconv.Atoi(txt[j])
			aux = append(aux, int32(Numero))
		}
		Retorno = append(Retorno, PlanetVector{Planeta: Planetines[i], VectorReloj: aux})
	}

	return Retorno

}

func Consistencia_Eventual() {
	//agregar conexion a fulcrum
	//Planetas_Vectores
	//		Nombre del planeta
	//		Vector
	//Ciudades_fulcrum
	//		Arreglo de ciudades
	//Soldados fulcrum
	//		Arreglo de numero
	fmt.Println("Empieza la consistencia eventual")
	for {
		time.Sleep(120 * time.Second)
		fmt.Println("Pasaron 2 minutos, se realiza merge entre replicas")
		Planetas_Vectores_Fulcrum2 := ObtenerPlanetas_Fulcrum("dist21:50072")
		Planetas_Vectores_Fulcrum3 := ObtenerPlanetas_Fulcrum("dist21:50073")
		//Ciudades_Fulcrum2, Soldados_Fulcrum2 := ObtenerCiudades_Fulcrum("dist21:50072")
		//Ciudades_Fulcrum3, Soldados_Fulcrum3 := ObtenerCiudades_Fulcrum("dist21:50073")
		if len(Planetas_Vectores) == 0 {
			for j := 0; j < len(Planetas_Vectores_Fulcrum2); j++ {
				Planetas_Vectores = append(Planetas_Vectores, Planetas_Vectores_Fulcrum2[j])
				Ciudades_Fulcrum2, Soldados_Fulcrum2 := ObtenerCiudades_Fulcrum("dist21:50072", Planetas_Vectores_Fulcrum2[j].Planeta)
				fmt.Println("Datos: ", Ciudades_Fulcrum2, Soldados_Fulcrum2)
				if len(Ciudades_Fulcrum2) > 0 {
					crear_archivo_planeta(Planetas_Vectores_Fulcrum2[j].Planeta, Ciudades_Fulcrum2[0], Soldados_Fulcrum2[0])
				}
				if len(Ciudades_Fulcrum2) > 1 {
					for c := 1; c < len(Ciudades_Fulcrum2); c++ {
						abrir_escribir_archivo(Planetas_Vectores_Fulcrum2[j].Planeta, Ciudades_Fulcrum2[c], Soldados_Fulcrum2[c])
					}
				}
			}
		} else {
			for j := 0; j < len(Planetas_Vectores_Fulcrum2); j++ {
				fmt.Println("Planeta en Fulcrum2: ", Planetas_Vectores_Fulcrum2[j].Planeta)
				i := Existe_Planeta(Planetas_Vectores_Fulcrum2[j].Planeta)
				if i > -1 {
					for c := 0; c < 3; c++ {
						if Planetas_Vectores[i].VectorReloj[c] < Planetas_Vectores_Fulcrum2[j].VectorReloj[c] {
							Planetas_Vectores[i].VectorReloj[c] = Planetas_Vectores_Fulcrum2[j].VectorReloj[c]
						}
					}
					//Se solicitan las ciudades del fulcrum central y del fulcrum 2 para hacer un merge entre ellas
					Ciudades_Fulcrum1, _ := ObtenerCiudades_Central(Planetas_Vectores[i].Planeta)
					Ciudades_Fulcrum2, Soldados_Fulcrum2 := ObtenerCiudades_Fulcrum("dist21:50072", Planetas_Vectores[i].Planeta)
					if len(Ciudades_Fulcrum2) > 0 {
						fmt.Println("Ciudades en Fulcrum2: ", Ciudades_Fulcrum2)
						var Ciudades_NoEstan []string
						var Soldados_NoEstan []int32

						//se identifican las ciudades que no están en fulcrum1 y se almacenan
						for k := 0; k < len(Ciudades_Fulcrum2); k++ {
							if !stringInSlice(Ciudades_Fulcrum2[k], Ciudades_Fulcrum1) {
								//Se agrega a fulcrum 1 la ciudad que no esté
								Ciudades_NoEstan = append(Ciudades_NoEstan, Ciudades_Fulcrum2[k])
								Soldados_NoEstan = append(Soldados_NoEstan, Soldados_Fulcrum2[k])
							}
						}
						//se guardan las ciudades que no están al archivo del fulcrum1
						for l := 0; l < len(Ciudades_NoEstan); l++ {
							abrir_escribir_archivo(Planetas_Vectores[i].Planeta, Ciudades_NoEstan[l], Soldados_NoEstan[l])
						}
					}
				} else { // Si el planeta no existe en fulcrum 1, pero si en fulcrum 2, entonces se crea el archivo y se agregan las ciudades

					Planetas_Vectores = append(Planetas_Vectores, Planetas_Vectores_Fulcrum2[j])
					Ciudades_Fulcrum2, Soldados_Fulcrum2 := ObtenerCiudades_Fulcrum("dist21:50072", Planetas_Vectores_Fulcrum2[j].Planeta)
					fmt.Println("Datos: ", Ciudades_Fulcrum2, Soldados_Fulcrum2)
					if len(Ciudades_Fulcrum2) > 0 {
						crear_archivo_planeta(Planetas_Vectores_Fulcrum2[j].Planeta, Ciudades_Fulcrum2[0], Soldados_Fulcrum2[0])
					}
					if len(Ciudades_Fulcrum2) > 1 {
						for c := 1; c < len(Ciudades_Fulcrum2); c++ {
							abrir_escribir_archivo(Planetas_Vectores_Fulcrum2[j].Planeta, Ciudades_Fulcrum2[c], Soldados_Fulcrum2[c])
						}
					}
				}
			}
		}

		if len(Planetas_Vectores) == 0 {
			for j := 0; j < len(Planetas_Vectores_Fulcrum3); j++ {
				Planetas_Vectores = append(Planetas_Vectores, Planetas_Vectores_Fulcrum3[j])
				Ciudades_Fulcrum3, Soldados_Fulcrum3 := ObtenerCiudades_Fulcrum("dist21:50073", Planetas_Vectores_Fulcrum3[j].Planeta)
				fmt.Println("Datos: ", Ciudades_Fulcrum3, Soldados_Fulcrum3)
				if len(Ciudades_Fulcrum3) > 0 {
					crear_archivo_planeta(Planetas_Vectores_Fulcrum3[j].Planeta, Ciudades_Fulcrum3[0], Soldados_Fulcrum3[0])
				}
				if len(Ciudades_Fulcrum3) > 1 {
					for c := 1; c < len(Ciudades_Fulcrum3); c++ {
						abrir_escribir_archivo(Planetas_Vectores_Fulcrum3[j].Planeta, Ciudades_Fulcrum3[c], Soldados_Fulcrum3[c])
					}
				}
			}
		} else {
			for j := 0; j < len(Planetas_Vectores_Fulcrum3); j++ {
				fmt.Println("Planeta en Fulcrum3: ", Planetas_Vectores_Fulcrum3[j].Planeta)
				i := Existe_Planeta(Planetas_Vectores_Fulcrum3[j].Planeta)
				if i != -1 {
					for c := 0; c < 3; c++ {
						if Planetas_Vectores[i].VectorReloj[c] < Planetas_Vectores_Fulcrum3[j].VectorReloj[c] {
							Planetas_Vectores[i].VectorReloj[c] = Planetas_Vectores_Fulcrum3[j].VectorReloj[c]
						}
					}
					Ciudades_Fulcrum1, _ := ObtenerCiudades_Central(Planetas_Vectores[i].Planeta)
					Ciudades_Fulcrum3, Soldados_Fulcrum3 := ObtenerCiudades_Fulcrum("dist21:50073", Planetas_Vectores[i].Planeta)
					if len(Ciudades_Fulcrum3) > 0 {
						fmt.Println("Ciudades en Fulcrum3: ", Ciudades_Fulcrum3)
						var Ciudades_NoEstan []string
						var Soldados_NoEstan []int32

						//se identifican las ciudades que no están en fulcrum1 y se almacenan
						for k := 0; k < len(Ciudades_Fulcrum3); k++ {
							if !stringInSlice(Ciudades_Fulcrum3[k], Ciudades_Fulcrum1) {
								//Se agrega a fulcrum 1 la ciudad que no esté
								Ciudades_NoEstan = append(Ciudades_NoEstan, Ciudades_Fulcrum3[k])
								Soldados_NoEstan = append(Soldados_NoEstan, Soldados_Fulcrum3[k])
							}
						}
						//se guardan las ciudades que no están al archivo del fulcrum1
						for l := 0; l < len(Ciudades_NoEstan); l++ {
							abrir_escribir_archivo(Planetas_Vectores[i].Planeta, Ciudades_NoEstan[l], Soldados_NoEstan[l])
						}
					}
				} else {
					Planetas_Vectores = append(Planetas_Vectores, Planetas_Vectores_Fulcrum3[j])
					Ciudades_Fulcrum3, Soldados_Fulcrum3 := ObtenerCiudades_Fulcrum("dist21:50073", Planetas_Vectores_Fulcrum3[j].Planeta)
					fmt.Println("Datos: ", Ciudades_Fulcrum3, Soldados_Fulcrum3)
					if len(Ciudades_Fulcrum3) > 0 {
						crear_archivo_planeta(Planetas_Vectores_Fulcrum3[j].Planeta, Ciudades_Fulcrum3[0], Soldados_Fulcrum3[0])
					}
					if len(Ciudades_Fulcrum3) > 1 {
						for c := 1; c < len(Ciudades_Fulcrum3); c++ {
							abrir_escribir_archivo(Planetas_Vectores_Fulcrum3[j].Planeta, Ciudades_Fulcrum3[c], Soldados_Fulcrum3[c])
						}
					}
				}
			}
		}
		for i := 0; i < len(Planetas_Vectores); i++ {

			Ciudades_Fulcrum, Soldados_Fulcrum := ObtenerCiudades_Central(Planetas_Vectores[i].Planeta)
			fmt.Println(Planetas_Vectores[i].Planeta, Ciudades_Fulcrum)
			for j := 0; j < len(Ciudades_Fulcrum); j++ {

				AgregarCiudades("dist21:50072", Planetas_Vectores[i].Planeta, Ciudades_Fulcrum[j], Soldados_Fulcrum[j])
				AgregarCiudades("dist21:50073", Planetas_Vectores[i].Planeta, Ciudades_Fulcrum[j], Soldados_Fulcrum[j])

			}

			ModificarVector_Fulcrum("dist21:50072", Planetas_Vectores[i].Planeta, Planetas_Vectores[i].VectorReloj)
			ModificarVector_Fulcrum("dist21:50073", Planetas_Vectores[i].Planeta, Planetas_Vectores[i].VectorReloj)
		}
	}
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
	soldados := -1
	arch := planeta + ".txt"
	Existe := false
	for i := 0; i < len(Planetas_Vectores); i++ {
		if Planetas_Vectores[i].Planeta == planeta {
			Existe = true
		}
	}
	if Existe {
		file, err := os.Open(arch)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		// optionally, resize scanner's capacity for lines over 64K, see next example
		for scanner.Scan() {
			linea_actual := strings.Split(scanner.Text(), " ")
			if len(linea_actual) == 3 {
				if linea_actual[1] == ciudad {
					sold, _ := strconv.Atoi(linea_actual[2])
					soldados = sold
				}
			}
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}

	return soldados
}

//Función que crea el archivo de planeta y agrega la primera ciudad con sus respectivos soldados. Usada en AddCity
func crear_archivo_planeta(planeta string, ciudad string, soldados int32) {
	arch := planeta + ".txt"
	soldiers := strconv.Itoa(int(soldados))
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
func abrir_escribir_archivo(planeta string, ciudad string, soldados int32) {
	arch := planeta + ".txt"
	soldiers := strconv.Itoa(int(soldados))
	escritura := planeta + " " + ciudad + " " + soldiers + "\n"
	fmt.Println("Abriendo archivo: ", planeta)
	file, err := os.OpenFile(arch, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	_, err2 := file.WriteString(escritura)
	if err2 != nil {
		log.Fatal(err2)
	}
	defer file.Close()
	return
}

// Servicios Mos - Fulcrum

//funcion para entregar la cantidad de soldados de determinada ciudad
func (s *Mos_FulcrumServer) GetSoldados(ctx context.Context, in *pb.Solicitud) (*pb.LeiaResponse, error) {

	ciudadin := in.GetCiudad()
	planetin := in.GetPlaneta()
	VectoraGuardar := []int32{-1, -1, -1}
	for i := 0; i < len(Planetas_Vectores); i++ {
		if Planetas_Vectores[i].Planeta == planetin {
			VectoraGuardar = Planetas_Vectores[i].VectorReloj
		}
	}

	soldiers := int32(leer_archivo_planeta(planetin, ciudadin))
	return &pb.LeiaResponse{Soldados: soldiers, Vector: VectoraGuardar, Direccion: "dist21:50061"}, nil

}

//Funcion para entregar el vector del planeta
func (s *Mos_FulcrumServer) GetVector(ctx context.Context, in *pb.Solicitud) (*pb.Reloj, error) {

	planetin := in.GetPlaneta()
	VectoraGuardar := []int32{-1, -1, -1}
	for i := 0; i < len(Planetas_Vectores); i++ {
		if Planetas_Vectores[i].Planeta == planetin {
			VectoraGuardar = Planetas_Vectores[i].VectorReloj
		}
	}

	return &pb.Reloj{Vector: VectoraGuardar, Planeta: planetin}, nil

}

//Servicios Informantes - Fulcrum

func (s *InfoTo_FulcrumServer) AddCity(ctx context.Context, in *pb.MessageINF) (*pb.Reloj, error) {
	/*
			string planeta = 1;
		    string ciudad = 2;
		    string newcity = 3;
		    int32 soldados = 4;*/
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
		fmt.Println("leer_archivo_planeta: ", leer_archivo_planeta(planetin, ciudadin))
		if leer_archivo_planeta(planetin, ciudadin) < 0 {
			abrir_escribir_archivo(planetin, ciudadin, soldados)
			Planetas_Vectores[IndicePlaneta].VectorReloj[0] += 1
			return &pb.Reloj{Vector: Planetas_Vectores[IndicePlaneta].VectorReloj, Planeta: planetin}, nil
		} else {
			VectoraGuardar := []int32{-1, -1, -1}
			fmt.Println("La ciudad ya existe, comando AddCity no válido")
			return &pb.Reloj{Vector: VectoraGuardar, Planeta: planetin}, nil
		}
	}

	crear_archivo_planeta(planetin, ciudadin, soldados)
	VectoraGuardar := []int32{0, 0, 0}
	Planetas_Vectores = append(Planetas_Vectores, PlanetVector{Planeta: planetin, VectorReloj: VectoraGuardar})
	return &pb.Reloj{Vector: VectoraGuardar, Planeta: planetin}, nil

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
			actualizar_archivo_planeta(planetin, ciudadin, linea_texto)
			Planetas_Vectores[IndicePlaneta].VectorReloj[0] += 1
			fmt.Println("Nombre de la ciudad ", ciudadin, " actualizado a ", nuevaCiudad)
			return &pb.Reloj{Vector: Planetas_Vectores[IndicePlaneta].VectorReloj, Planeta: planetin}, nil
		} else {
			VectoraGuardar := []int32{-1, -1, -1}
			fmt.Println("La ciudad ", ciudadin, " en ", planetin, " no existe, comando UpdateName no válido")
			return &pb.Reloj{Vector: VectoraGuardar, Planeta: planetin}, nil
		}
	}

	VectoraGuardar := []int32{-1, -1, -1}
	fmt.Println("El planeta ", planetin, " no existe, comando UpdateName no válido")
	return &pb.Reloj{Vector: VectoraGuardar, Planeta: planetin}, nil

}

func (s *InfoTo_FulcrumServer) UpdateNumber(ctx context.Context, in *pb.MessageINF) (*pb.Reloj, error) {
	planetin := in.GetPlaneta()
	ciudadin := in.GetCiudad()
	nuevosSoldados := in.GetSoldados()

	linea_texto := planetin + " " + ciudadin + " " + strconv.Itoa(int(nuevosSoldados))

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
			actualizar_archivo_planeta(planetin, ciudadin, linea_texto)
			Planetas_Vectores[IndicePlaneta].VectorReloj[0] += 1
			fmt.Println("Numero de soldados de la ciudad ", ciudadin, " actualizado a ", nuevosSoldados)
			return &pb.Reloj{Vector: Planetas_Vectores[IndicePlaneta].VectorReloj, Planeta: planetin}, nil
		} else {
			VectoraGuardar := []int32{-1, -1, -1}
			fmt.Println("La ciudad ", ciudadin, " en ", planetin, " no existe, comando UpdateNumber no válido")
			return &pb.Reloj{Vector: VectoraGuardar, Planeta: planetin}, nil
		}
	}

	VectoraGuardar := []int32{-1, -1, -1}
	fmt.Println("El planeta ", planetin, " no existe, comando UpdateNumber no válido")
	return &pb.Reloj{Vector: VectoraGuardar, Planeta: planetin}, nil
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
			actualizar_archivo_planeta(planetin, ciudadin, "")
			Planetas_Vectores[IndicePlaneta].VectorReloj[0] += 1
			fmt.Println("Ciudad ", ciudadin, " eliminada ")
			return &pb.Reloj{Vector: Planetas_Vectores[IndicePlaneta].VectorReloj, Planeta: planetin}, nil
		} else {
			VectoraGuardar := []int32{-1, -1, -1}
			fmt.Println("La ciudad ", ciudadin, " en ", planetin, " no existe, comando DeleteCity no válido")
			return &pb.Reloj{Vector: VectoraGuardar, Planeta: planetin}, nil
		}
	}

	VectoraGuardar := []int32{-1, -1, -1}
	fmt.Println("El planeta ", planetin, " no existe, comando DeleteCity no válido")
	return &pb.Reloj{Vector: VectoraGuardar, Planeta: planetin}, nil
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
	log.Printf("Servidor para Broker Mos escuchando en %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Fallo en serve: %v", err)
	}
}

var Planetas_Vectores []PlanetVector

func main() {
	go ServidorMos()
	go ServidorInformantes()
	Consistencia_Eventual()
}
