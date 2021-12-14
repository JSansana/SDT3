package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	pb "github.com/JSansana/SDT3/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	address = "dist23:50051"
)

func arrayToString(a []int32, delim string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(a), " ", delim, -1), "[]")
}

//ESTA FUNCION DEBE RETORNAR LAS VARIABLES PARA REALIZAR EL READ YOUR WRITES
func Conexion_Fulcrum(direccion string, numero_funcion int, planetin string) ([]int32, string) {

	//REALIZA LA CONEXIÓN CON EL FULCRUM

	var ciudadin string

	var vector []int32
	var comando string

	switch numero_funcion {
	case 1:
		comando = "AddCity "
		nuevo_valor := 0
		log.Printf("---------------------------------")
		comando = comando + planetin + " "
		fmt.Println("Ingrese la ciudad a añadir:")
		fmt.Scan(&ciudadin)
		comando = comando + ciudadin + " "
		fmt.Println("¿Nuevo valor de rebeldes? Si: 1, No: 0")
		fmt.Scan(&nuevo_valor)
		if nuevo_valor == 1 {
			fmt.Println("Ingrese nuevo valor de rebeldes")
			fmt.Scan(&nuevo_valor)
			conn, err := grpc.Dial(direccion, grpc.WithInsecure(), grpc.WithBlock())
			if err != nil {
				log.Fatalf("No conectó: %v", err)
			}
			defer conn.Close()
			c := pb.NewInfoTo_FulcrumClient(conn)
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()
			comando = comando + strconv.Itoa(nuevo_valor)
			r, err := c.AddCity(ctx, &pb.MessageINF{Planeta: planetin, Ciudad: ciudadin, Newcity: "Hola", Soldados: int32(nuevo_valor)})
			if err != nil {
				log.Fatalf("No se pudo enviar solicitud: %v", err)
			}
			vector = r.GetVector()

		} else {
			comando = comando + strconv.Itoa(nuevo_valor)
			conn, err := grpc.Dial(direccion, grpc.WithInsecure(), grpc.WithBlock())
			if err != nil {
				log.Fatalf("No conectó: %v", err)
			}
			defer conn.Close()
			c := pb.NewInfoTo_FulcrumClient(conn)
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()
			r, err := c.AddCity(ctx, &pb.MessageINF{Planeta: planetin, Ciudad: ciudadin, Newcity: "Hola", Soldados: int32(nuevo_valor)})
			if err != nil {
				log.Fatalf("No se pudo enviar solicitud: %v", err)
			}
			vector = r.GetVector()

		}
		//Recibe el vector de un fulcrum
		//vector = r.AddCity()

	case 2:
		comando = "UpdateName "
		var nuevo_nombre string
		log.Printf("---------------------------------")
		comando = comando + planetin + " "
		fmt.Println("Ingrese la ciudad:")
		fmt.Scan(&ciudadin)
		comando = comando + ciudadin + " "
		fmt.Println("Ingrese el nuevo nombre de la ciudad:")
		fmt.Scan(&nuevo_nombre)
		comando = comando + nuevo_nombre
		conn, err := grpc.Dial(direccion, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("No conectó: %v", err)
		}
		defer conn.Close()
		c := pb.NewInfoTo_FulcrumClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		r, err := c.UpdateName(ctx, &pb.MessageINF{Planeta: planetin, Ciudad: ciudadin, Newcity: nuevo_nombre, Soldados: 0})
		if err != nil {
			log.Fatalf("No se pudo enviar solicitud: %v", err)
		}
		vector = r.GetVector()

	case 3:
		var nuevo_valor int
		comando = "UpdateNumber "
		log.Printf("---------------------------------")
		comando = comando + planetin + " "
		fmt.Println("Ingrese la ciudad:")
		fmt.Scan(&ciudadin)
		comando = comando + ciudadin + " "
		fmt.Println("Ingrese el nuevo valor de rebeldes:")
		fmt.Scan(&nuevo_valor)
		comando = comando + strconv.Itoa(nuevo_valor)
		conn, err := grpc.Dial(direccion, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("No conectó: %v", err)
		}
		defer conn.Close()
		c := pb.NewInfoTo_FulcrumClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		r, err := c.UpdateNumber(ctx, &pb.MessageINF{Planeta: planetin, Ciudad: ciudadin, Newcity: "Hola", Soldados: int32(nuevo_valor)})
		if err != nil {
			log.Fatalf("No se pudo enviar solicitud: %v", err)
		}
		vector = r.GetVector()
	case 4:
		comando = "DeleteCity "
		log.Printf("---------------------------------")
		comando = comando + planetin + " "
		fmt.Println("Ingrese la ciudad :")
		fmt.Scan(&ciudadin)
		comando = comando + ciudadin
		conn, err := grpc.Dial(direccion, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("No conectó: %v", err)
		}
		defer conn.Close()
		c := pb.NewInfoTo_FulcrumClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		r, err := c.DeleteCity(ctx, &pb.MessageINF{Planeta: planetin, Ciudad: ciudadin, Newcity: "Hola", Soldados: 0})
		if err != nil {
			log.Fatalf("No se pudo enviar solicitud: %v", err)
		}
		vector = r.GetVector()

	}

	return vector, comando
}

func main() {
	fmt.Println("Hola Thrawn, ¿Qué tal el imperio? ¡Espero que muy bien!")

	//Planeta, Reloj, Comando, Direccion
	var Read_Writes [][4]string
	var planetin string
	Bandera := 0
	Indice := 0
	//nuevo cliente

	var choice int
	for {
		log.Printf("---------------------------------")
		fmt.Println("1) Añadir Ciudad")
		fmt.Println("2) Actualizar nombre de una ciudad")
		fmt.Println("3) Actualizar cantidad de soldados en una ciudad")
		fmt.Println("4) Borrar ciudad")
		fmt.Println("5) Salir")
		fmt.Scan(&choice)

		switch choice {
		case 1:
			fmt.Println("Ingrese el planeta:")
			fmt.Scan(&planetin)

			vector_reloj := []int32{-1, -1, -1}
			if len(Read_Writes) != 0 {
				for c := 0; c < len(Read_Writes); c++ {
					if Read_Writes[c][0] == planetin {
						el_vector := strings.Split(Read_Writes[c][1], ",")
						Bandera = 1
						Indice = c
						for j := 0; j < 3; j++ {

							auxiliar, _ := strconv.Atoi(el_vector[j])
							vector_reloj[j] = int32(auxiliar)
						}
					}
				}
			}
			conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
			if err != nil {
				log.Fatalf("No conectó: %v", err)
			}
			defer conn.Close()
			c := pb.NewInfoToMosClient(conn)
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()
			//Envía un 1 de solicitud a Broker para recibir la dirección determinada
			r, err := c.AskAdress(ctx, &pb.Reloj{Vector: vector_reloj, Planeta: planetin})
			if err != nil {
				log.Fatalf("No se pudo enviar solicitud: %v ", err)
			}
			defer conn.Close()

			Direccion_Fulcrum := r.GetDireccion()
			fmt.Println("Direccion: %s ", Direccion_Fulcrum)
			vector, comando := Conexion_Fulcrum(Direccion_Fulcrum, 1, planetin)
			NombrePlaneta := strings.Split(comando, " ")[1]

			if vector[0] != -1 {
				fmt.Println("Ciudad añadida! en Fulcrum con dirección: %s ", Direccion_Fulcrum)
				if Bandera == 0 {
					datos := [4]string{NombrePlaneta, arrayToString(vector, ","), comando, Direccion_Fulcrum}
					Read_Writes = append(Read_Writes, datos)

				} else {
					Read_Writes[Indice][1] = arrayToString(vector, ",")
					Read_Writes[Indice][2] = comando
					Read_Writes[Indice][3] = Direccion_Fulcrum
					Bandera = 0

				}
			} else {
				fmt.Println("La ciudad ya existe!")
			}
		case 2:
			fmt.Println("Ingrese el planeta:")
			fmt.Scan(&planetin)

			vector_reloj := []int32{-1, -1, -1}
			if len(Read_Writes) != 0 {
				for c := 0; c < len(Read_Writes); c++ {
					if Read_Writes[c][0] == planetin {
						el_vector := strings.Split(Read_Writes[c][1], ",")
						Bandera = 1
						Indice = c
						for j := 0; j < 3; j++ {
							auxiliar, _ := strconv.Atoi(el_vector[j])
							vector_reloj[j] = int32(auxiliar)
						}
					}
				}
			}
			conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
			if err != nil {
				log.Fatalf("No conectó: %v", err)
			}
			defer conn.Close()
			c := pb.NewInfoToMosClient(conn)
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()
			//Envía un 2 de solicitud a Broker para recibir la dirección determinada
			r, err := c.AskAdress(ctx, &pb.Reloj{Vector: vector_reloj, Planeta: planetin})
			if err != nil {
				log.Fatalf("No se pudo enviar solicitud: %v ", err)
			}
			Direccion_Fulcrum := r.GetDireccion()
			vector, comando := Conexion_Fulcrum(Direccion_Fulcrum, 2, planetin)
			NombrePlaneta := strings.Split(comando, " ")[1]

			if vector[0] != -1 {
				fmt.Println("Nombre actualizado! en Fulcrum con dirección: %s ", Direccion_Fulcrum)
				if Bandera == 0 {
					datos := [4]string{NombrePlaneta, arrayToString(vector, ","), comando, Direccion_Fulcrum}
					Read_Writes = append(Read_Writes, datos)

				} else {
					Read_Writes[Indice][1] = arrayToString(vector, ",")
					Read_Writes[Indice][2] = comando
					Read_Writes[Indice][3] = Direccion_Fulcrum
					Bandera = 0

				}
			} else {
				fmt.Println("La ciudad no existe!")
			}

		case 3:
			fmt.Println("Ingrese el planeta:")
			fmt.Scan(&planetin)

			vector_reloj := []int32{-1, -1, -1}
			if len(Read_Writes) != 0 {
				for c := 0; c < len(Read_Writes); c++ {
					if Read_Writes[c][0] == planetin {
						el_vector := strings.Split(Read_Writes[c][1], ",")
						Bandera = 1
						Indice = c
						for j := 0; j < 3; j++ {
							auxiliar, _ := strconv.Atoi(el_vector[j])
							vector_reloj[j] = int32(auxiliar)
						}
					}
				}
			}
			conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
			if err != nil {
				log.Fatalf("No conectó: %v", err)
			}
			defer conn.Close()
			c := pb.NewInfoToMosClient(conn)
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()
			//Envía un 3 de solicitud a Broker para recibir la dirección determinada
			r, err := c.AskAdress(ctx, &pb.Reloj{Vector: vector_reloj, Planeta: planetin})
			if err != nil {
				log.Fatalf("No se pudo enviar solicitud: %v ", err)
			}
			Direccion_Fulcrum := r.GetDireccion()
			vector, comando := Conexion_Fulcrum(Direccion_Fulcrum, 3, planetin)
			NombrePlaneta := strings.Split(comando, " ")[1]

			if vector[0] != -1 {
				fmt.Println("Numero de soldados actualizado! en Fulcrum con dirección: %s \n", Direccion_Fulcrum)
				if Bandera == 0 {
					datos := [4]string{NombrePlaneta, arrayToString(vector, ","), comando, Direccion_Fulcrum}
					Read_Writes = append(Read_Writes, datos)

				} else {
					Read_Writes[Indice][1] = arrayToString(vector, ",")
					Read_Writes[Indice][2] = comando
					Read_Writes[Indice][3] = Direccion_Fulcrum
					Bandera = 0

				}
			} else {
				fmt.Println("La ciudad no existe!")
			}

		case 4:
			fmt.Println("Ingrese el planeta:")
			fmt.Scan(&planetin)

			vector_reloj := []int32{-1, -1, -1}
			if len(Read_Writes) != 0 {
				for c := 0; c < len(Read_Writes); c++ {
					if Read_Writes[c][0] == planetin {
						el_vector := strings.Split(Read_Writes[c][1], ",")
						Bandera = 1
						Indice = c
						for j := 0; j < 3; j++ {
							auxiliar, _ := strconv.Atoi(el_vector[j])
							vector_reloj[j] = int32(auxiliar)
						}
					}
				}
			}
			//Envía un 4 de solicitud a Broker para recibir la dirección determinada
			conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
			if err != nil {
				log.Fatalf("No conectó: %v", err)
			}
			defer conn.Close()
			c := pb.NewInfoToMosClient(conn)
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()
			r, err := c.AskAdress(ctx, &pb.Reloj{Vector: vector_reloj, Planeta: planetin})
			if err != nil {
				log.Fatalf("No se pudo enviar solicitud: %v ", err)
			}
			defer conn.Close()
			Direccion_Fulcrum := r.GetDireccion()
			vector, comando := Conexion_Fulcrum(Direccion_Fulcrum, 4, planetin)
			NombrePlaneta := strings.Split(comando, " ")[1]
			if vector[0] != -1 {
				fmt.Println("Ciudad eliminada! en Fulcrum con dirección: %s \n", Direccion_Fulcrum)
				if Bandera == 0 {
					datos := [4]string{NombrePlaneta, arrayToString(vector, ","), comando, Direccion_Fulcrum}
					Read_Writes = append(Read_Writes, datos)

				} else {
					Read_Writes[Indice][1] = arrayToString(vector, ",")
					Read_Writes[Indice][2] = comando
					Read_Writes[Indice][3] = Direccion_Fulcrum
					Bandera = 0

				}
			} else {
				fmt.Println("La ciudad no existe!")
			}
		case 5:
			//conn.Close()
			break

		default:
			fmt.Println("Opción inválida")
			continue
		}
	}
}
