package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"reflect"
	"strconv"
	"time"

	pb "github.com/JSansana/SDT3/proto"
	"google.golang.org/grpc"
)

const (
	PortInformante = ":50051"
	PortLeia       = ":50052"

	AdressFulcrum1 = "localhost:50061"
	AdressFulcrum2 = "localhost:50062"
	AdressFulcrum3 = "localhost:50063"
)

type LeiaToMosServer struct {
	pb.UnimplementedLeiaToMosServer
}
type InfoToMosServer struct {
	pb.UnimplementedInfoToMosServer
}

func Slices_MayorOIgual(vector1 []int32, vector2 []int32) bool {

	cont := 0

	for i := 0; i < 3; i++ {
		if vector1[i] <= vector2[i] {
			cont++
		}
	}
	if cont == 3 {
		return true
	}

	return false
}

func Conexion_Mos_Fulcrum(direccion string, planetin string) []int32 {
	//REALIZA LA CONEXIÓN CON EL FULCRUM
	conn, err := grpc.Dial(direccion, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("No conectó: %v", err)
	}
	defer conn.Close()
	c := pb.NewMos_FulcrumClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.GetVector(ctx, &pb.Solicitud{Planeta: planetin, Ciudad: "Murphy"})
	if err != nil {
		log.Fatalf("No se pudo enviar solicitud: %v", err)
	}
	el_vector := r.GetVector()

	return el_vector
}

func Conexion_Mos_Fulcrum_Leia(direccion string, planetin string, ciudadin string) (int32, []int32, string) {

	conn, err := grpc.Dial(direccion, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("No conectó: %v", err)
	}
	defer conn.Close()
	c := pb.NewMos_FulcrumClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.GetSoldados(ctx, &pb.Solicitud{Planeta: planetin, Ciudad: ciudadin})
	if err != nil {
		log.Fatalf("No se pudo enviar solicitud: %v", err)
	}
	soldados := r.GetSoldados()
	vector := r.GetVector()
	direccionNueva := r.GetDireccion()

	return soldados, vector, direccionNueva
}

//Implementación del método definido en proto
func (s *InfoToMosServer) AskAdress(ctx context.Context, in *pb.Reloj) (*pb.RespuestaMos, error) {
	fmt.Printf("RECIBIDOOOOOOOOOOOOOOOO")
	vectorInformante := in.GetVector()
	planetin := in.GetPlaneta()
	Merge := false

	//Vectores obtenidos de cada Fulcrum para comparar
	vector1 := Conexion_Mos_Fulcrum(AdressFulcrum1, planetin)
	vector2 := Conexion_Mos_Fulcrum(AdressFulcrum2, planetin)
	vector3 := Conexion_Mos_Fulcrum(AdressFulcrum3, planetin)

	//Si los tres vectores de reloj son iguales, entonces hubo un merge, entonces la asignación de fulcrum será aleatoria
	if reflect.DeepEqual(vector1, vector2) && reflect.DeepEqual(vector2, vector3) {
		Merge = true
	}

	// Si hubo un merge anteriormente o el vector de informantes tiene -1
	// significa que el broker escoge una dirección al azar para informantes
	if vectorInformante[0] != -1 && !Merge {
		if Slices_MayorOIgual(vectorInformante, vector1) == true {

			return &pb.RespuestaMos{Direccion: "localhost:50081"}, nil

		} else if Slices_MayorOIgual(vectorInformante, vector2) == true {

			return &pb.RespuestaMos{Direccion: "localhost:50082"}, nil

		} else if Slices_MayorOIgual(vectorInformante, vector3) == true {

			return &pb.RespuestaMos{Direccion: "localhost:50083"}, nil
		}

	}
	rand.Seed(time.Now().UnixNano())
	min := 50081
	max := 50083
	var numero int = int(rand.Intn(max-min+1) + min)
	dir := "localhost:" + strconv.Itoa(numero)
	return &pb.RespuestaMos{Direccion: dir}, nil
}

func (s *LeiaToMosServer) GetNumberRebelds(ctx context.Context, in *pb.MessageLeia) (*pb.LeiaResponse, error) {
	vectorLeia := in.GetVector()
	planetin := in.GetPlaneta()
	ciudadin := in.GetCiudad()
	Merge := false

	//Vectores obtenidos de cada Fulcrum para comparar
	soldados1, vector1, direccion1 := Conexion_Mos_Fulcrum_Leia(AdressFulcrum1, planetin, ciudadin)
	soldados2, vector2, direccion2 := Conexion_Mos_Fulcrum_Leia(AdressFulcrum2, planetin, ciudadin)
	soldados3, vector3, direccion3 := Conexion_Mos_Fulcrum_Leia(AdressFulcrum3, planetin, ciudadin)

	if reflect.DeepEqual(vector1, vector2) && reflect.DeepEqual(vector2, vector3) {
		Merge = true
	}

	if vectorLeia[0] == -1 || Merge {
		rand.Seed(time.Now().UnixNano())
		min := 50061
		max := 50063
		var numero int32 = int32(rand.Intn(max-min+1) + min)
		dir := "localhost:" + strconv.Itoa(int(numero))
		if dir == direccion1 {
			return &pb.LeiaResponse{Soldados: soldados1, Vector: vector1, Direccion: dir}, nil

		} else if dir == direccion2 {
			return &pb.LeiaResponse{Soldados: soldados2, Vector: vector2, Direccion: dir}, nil

		} else if dir == direccion3 {
			return &pb.LeiaResponse{Soldados: soldados3, Vector: vector3, Direccion: dir}, nil
		}

	} else {
		if Slices_MayorOIgual(vectorLeia, vector1) {

			return &pb.LeiaResponse{Soldados: soldados1, Vector: vector1, Direccion: AdressFulcrum1}, nil

		} else if Slices_MayorOIgual(vectorLeia, vector2) {

			return &pb.LeiaResponse{Soldados: soldados2, Vector: vector2, Direccion: AdressFulcrum2}, nil

		} else if Slices_MayorOIgual(vectorLeia, vector3) {

			return &pb.LeiaResponse{Soldados: soldados3, Vector: vector3, Direccion: AdressFulcrum1}, nil

		}

	}

	return &pb.LeiaResponse{Soldados: 1, Vector: vector1, Direccion: "hola"}, nil

}

func ServidorInformantes() {
	lis, err := net.Listen("tcp", PortInformante)
	if err != nil {
		log.Fatalf("Fallo al escuchar: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterInfoToMosServer(s, &InfoToMosServer{})
	log.Printf("Servidor para informantes escuchando en %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Fallo en serve: %v", err)
	}

}

func ServidorLeia() {
	lis, err := net.Listen("tcp", PortLeia)
	if err != nil {
		log.Fatalf("Fallo al escuchar: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterLeiaToMosServer(s, &LeiaToMosServer{})
	log.Printf("Servidor para Leia escuchando en %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Fallo en serve: %v", err)
	}

}

func main() {
	fmt.Println("Bienvenido a Mos Eisley, vaya con cuidado.")
	go ServidorInformantes()
	ServidorLeia()

}
