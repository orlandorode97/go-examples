package main

import (
	"fmt"
	"log"
	"syscall"
)

func main() {
	// 1. Crear socket (AF_INET = IPv4, SOCK_STREAM = TCP)
	// Todo en UNIX es un archivo
	// syscall.Socket retorna un file descriptor el cual es un identificador numerico que el sistema
	// operativo asigna a un archivo, socket y otro recurso de entra/salida cuando se abre. En este caso, un socket.
	// AF_INET representa una direccion que es usada para designar el tipo de direccion que el socket puede comunicarse, IPv3
	// SOCK_STREAM representa que el socket proporcione un canal de comunicacion orientado a una conexion y confiable donde se manden
	// un flujo continuo de bytes, sin perdida, ni duplicacion y en order.
	// IPPROTO_TCP indica al kernel que protocol de transporte debe usar el socket, TCP en este caso.
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
	if err != nil {
		log.Fatalf("Unable to create socket: %v", err)
		return
	}
	// syscall.Close mata el file descriptor o el identificador numerico.
	defer syscall.Close(fd)

	fmt.Printf("Socker created, fd = %v\n", fd)

	// Reutiliza el puerto para evitar address already in use
	err = syscall.SetsockoptInet4Addr(fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, [4]byte{1})
	if err != nil {
		log.Fatalf("Error setsockopt: %v", err)
	}

	// Asociar a la direccion 0.0.0.0:8080
	addr := syscall.SockaddrInet4{Port: 8080}
	copy(addr.Addr[:], []byte{0, 0, 0, 0})

	if err := syscall.Bind(fd, &addr); err != nil {
		log.Fatalf("Error in bind: %v\n", err)
	}

	// A la escucha de la direccion y puerto 0.0.0.0:8080
	// Cambia el socke a LISTEN
	if err := syscall.Listen(fd, 5); err != nil {
		log.Fatalf("Error en listen: %v", err)
	}

	fmt.Println("Escuchando en puerto 8080...")

	for {
		// Acepta el socket de conexiocion, no crea la conexion.
		// El cliente manda un SYN (syncronization)
		// El kernel responde con un SYN+ACK (syncronization0-aknowledge)
		// Y luego un nuevo socket (nuevo FD) entra al proceso
		connFd, sa, err := syscall.Accept(fd)
		if err != nil {
			log.Printf("Error en accept: %v", err)
			continue
		}
		go handleConnectionFd(connFd, sa)
	}
}

func handleConnectionFd(fd int, sa syscall.Sockaddr) {
	defer syscall.Close(fd) // El kernel envia FIN

	buf := make([]byte, 1024)

	for {
		n, err := syscall.Read(fd, buf)
		if err != nil {
			log.Fatalf("Unable to read: %v", err)
			return
		}

		if n == 0 {
			fmt.Println("Closing connection")
			return
		}

		data := buf[:n]

		fmt.Printf("Getting :%s", string(data))

		syscall.Write(fd, []byte("Echo: "+string(data)))
	}
}
