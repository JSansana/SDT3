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
	address = "localhost:50052"
)

func arrayToString(a []int32, delim string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(a), " ", delim, -1), "[]")
}

func main() {
	fmt.Println("¿Cómo está Doña Leia?, bienvenida!")

	var monotonic_reads [][3]string
	var planetin string
	var ciudadin string
	Contador := 0
	Bandera := 0
	Indice := 0
	choice := 0

	for {
		log.Printf("---------------------------------")
		fmt.Println("1) Pedir Planeta, ciudad")
		fmt.Println("2) Salir")
		fmt.Scan(&choice)

		switch choice {
		case 1:
			log.Printf("---------------------------------")
			fmt.Println("Ingrese el planeta:")
			fmt.Scan(&planetin)
			fmt.Println("Ingrese la ciudad:")
			fmt.Scan(&ciudadin)

			vector_reloj := []int32{-1, -1, -1}
			if len(monotonic_reads) != 0 {
				for c := 0; c < len(monotonic_reads); c++ {
					if monotonic_reads[c][0] == planetin {
						el_vector := strings.Split(monotonic_reads[c][1], ",")
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
			c := pb.NewLeiaToMosClient(conn)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			//Envía un 1 de solicitud a Broker para que este se
			r, err := c.GetNumberRebelds(ctx, &pb.MessageLeia{Planeta: planetin, Ciudad: ciudadin, Vector: vector_reloj})
			if err != nil {
				log.Fatalf("No se pudo enviar solicitud: %v", err)
			}

			soldados := r.GetSoldados()
			vector := r.GetVector()
			direccionNueva := r.GetDireccion()

			if soldados > -1 {
				fmt.Println("Doña Leia, hay %d soldados.", soldados)
				if Bandera == 0 {
					monotonic_reads[Contador][0] = planetin
					monotonic_reads[Contador][1] = arrayToString(vector, ",")
					monotonic_reads[Contador][2] = direccionNueva
					Contador = Contador + 1
				} else {

					monotonic_reads[Indice][1] = arrayToString(vector, ",")
					monotonic_reads[Indice][2] = direccionNueva
					Bandera = 0
				}
			} else {

				fmt.Println("No existe la ciudad donde se buscan rebeldes.")
			}
		case 2:
			//conn.Close()
			break
		default:
			fmt.Println("Opción inválida")
			continue
		}

	}

}
