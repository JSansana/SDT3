---Laboratorio 3 Sistemas Distribuidos---

Integrantes: GRUPO 5
	Marcelo Cabezas Galvez : 201704546-9
	Alvaro Ortiz Hermosilla :201810523-6
	José Sansana Parra : 201773535-K

Orden de ejecución sugerido:

Se destaca que Ahsoka y Thrawn se ejecutan con comando make distinto

(Se debe abrir una terminal distinta para la ejecución de cada entidad,
independiente de que algunas se encuentre en la misma VM )

	1) Fulcrum 3
	   Entrar a VM dist23
	   $ cd SDT3/F3
	   $ make

	2) Fulcrum 2
	   Entrar a VM dist22
	   $ cd SDT3/F2
	   $ make

	3) Fulcrum 1
	   Entrar a VM dist21
	   $ cd SDT3/F1
	   $ make

	4) Broker Mos Eisley
	   Entrar a VM dist24
	   $ cd SDT3/BM
	   $ make

	5) Leia
	   Entrar a VM dist21
	   $ cd SDT3/L
	   $ make

	6) Ahsoka
	   Entrar a VM dist22
	   $ cd SDT3/INF
	   $ make A

	7) Thrawn
	   Entrar a VM dist23
	   $ cd SDT3/INF
	   $ make T	
	

Recomendación : 
	- No aplicar comandos al mismo tiempo que Fulcrum 1 señala que está realizando la consistencia Eventual


Consideraciones escenciales:
	
	- El nodo dominante escogido para realizar la consistencia eventual (merge) es el Fulcrum 1.
	- Los comandos se ingresan mediante una interfaz escogiendo el número indicado en pantalla e ingresando
	  a mano algunos datos cuando se soliciten (como nombre ciudad, planeta, soldados)
	- Si se quiere poner una ciudad que contenga un espacio, reemplazar este por guion bajo ( _ )
	- CTRL + C si se quiere terminar un proceso 
	
	Lógica y estructuras de datos usadas :

		Informantes: Slice [][4] Donde el primer valor es el nombre del planeta, el segundo el vector,
		el tercero la dirección y el cuarto el comando.
	
		Leia: Slice [][4] con la misma estructura que los informantes.
	
		Fulcrums: Slice[][2] donde el primer valor será el nombre del planeta y el segundo el vector.

	Comandos:

		AddCity:
			- Si el informante no tiene el planeta en su registro, el fulcrum lo redirige a un fulcrum aleatorio.
			- Si el informante tiene el planeta en su registro, el fulcrum comparará los relojes de informante y
	 		 fulcrums para verificar el que cumpla con read your writes(vector de informante menor o igual en todas
			 las posiciones respecto al vector de fulcrum).
			-Si el informante quiere agregar una ciudad que ya agregó, se comprueba el slice y si existe ya el registro de que
			 realizó Addcity el comando no se aplica.
			-Si está la misma ciudad en todos los Fulcrum el comando no aplica

		UpdateName:
    			-Caso 1: Informante encuentra un reloj dentro de su slice con los comandos anteriormente realizados para dicho planeta.
        		-El primer paso es que el informante acceda al ultimo reloj en el cual modifico dicha ciudad.
       			-El reloj accedido será enviado al Mos.
        		-Mos debe entonces comprobar los relojes de los fulcrum y entregar la dirección del fulcrum que cumpla read your writes.

    			-Caso 2: Informante no encuentra reloj (Accedió por primera vez):
        		- Como no existe reloj, el informante le mando un vector vacio(?
        		- Al ser una función update, Mos debe buscar el primer fulcrum que contenga ese planeta y ciudad.

		UpdateNumber:
    			Es lo mismo que el anterior.

		DeleteCity:
    			-Caso 1: Informante encuentra un reloj dentro de su slice con los comandos anteriormente realizados para dicho planeta.
        		-El primer paso es que el informante acceda al ultimo reloj en el cual modifico dicha ciudad.
        		-El reloj accedido será enviado al Mos.
        		-Mos debe entonces comprobar los relojes de los fulcrum y entregar la dirección del fulcrum que cumpla read your writes.

    			-Caso 2: Informante no encuentra reloj (Accedió por primera vez):
        		- Como no existe reloj, el informante le mando un vector vacio(?
        		- Al ser una función delete, Mos debe buscar el primer fulcrum que contenga ese planeta y ciudad.	

		GetNumberRebelds:

			-Siempre envía primero el reloj de vector desde leia al Broker, y luego este lo comparará con un Fulcrum
			Si el Fulcrum elegido no cumple MonotonicReads, pero tiene la ciudad buscada, no podrá acceder a los soldados
			Si el Fulcrum elegido cumple MonotonicReads y no tiene la ciudad buscada, accederá al fulcrum pero dirá que no se puede aplicar el comando
			Si el Fulcrum elegido cumple MonotonicReads y tiene la ciudad buscada, accederá y obtendrá la info necesaria.











