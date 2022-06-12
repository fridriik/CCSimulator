package main


import (
	"database/sql"
	"fmt"
	"log"
	"os"
	_"github.com/lib/pq"
	"github.com/dixonwille/wmenu/v5"
    "encoding/json"
	"strconv"
	bolt "github.com/boltdb/bolt"
)


/*
Se crean los structs y variable de boltDB para utilizar en la base de datos NoSQL
*/
type Cliente struct {
	Nrocliente int
    Nombre string
    Apellido string
	Domicilio string
	Telefono string
}


type Tarjeta struct {
    Nrotarjeta string
    Nrocliente int
    Validadesde string
    Validadhasta string
    Codseguridad string
    Limitecompra float32
    Estado string
}


type Comercio struct {
    Nrocomercio int
    Nombre string
    Domicilio string
    Codigopostal string
    Telefono string
}


type Compra struct {
    Nrooperacion int
    Nrotarjeta string
    Nrocomercio int
    Fecha string
    Monto float32
    Pagado bool
}


var dbbolt *bolt.DB


/*
Crea la base de datos
*/
func crearBD() {

	db, err := sql.Open("postgres", "user=postgres host=localhost dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`create database tp`)
	if err != nil {
		log.Fatal(err)
	}
}


/*
Elimina la base de datos
*/
func eliminarBD() {
	db, err := sql.Open("postgres", "user=postgres host=localhost dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`drop database if exists tp`)
	if err != nil {
		log.Fatal(err)
	}
}


/*
Crea las tablas y secuencias en la base de datos
*/
func crearTablasYSecuencias() {

	db, err := sql.Open("postgres", "user=postgres host=localhost dbname=tp sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`create sequence aumentoCompra;

  					  create sequence aumentoRechazo;

  					  create sequence aumentoCabecera;			  

  					  create sequence aumentoAlerta;

  					  create table cliente(nrocliente int,
										   nombre text,
										   apellido text,
										   domicilio text,
										   telefono char(12));

					  create table tarjeta(nrotarjeta char(16),
					  					   nrocliente int,
					  					   validadesde char(6),
					  					   validahasta char(6),
					  					   codseguridad char(4),
					  					   limitecompra decimal(8,2),
					  					   estado char(10));

					  create table comercio(nrocomercio int,
					  						nombre text,
					  						domicilio text,
					  						codigopostal char(8),
					  						telefono char(12));

					  create table compra(nrooperacion int not null default nextval('aumentoCompra'),
					  					  nrotarjeta char(16),
					  					  nrocomercio int,
					  					  fecha timestamp,
					  					  monto decimal(7,2),
					  					  pagado boolean);

					  create table rechazo(nrorechazo int not null default nextval('aumentoRechazo'),
					  					   nrotarjeta char(16),
					  					   nrocomercio int,
					  					   fecha timestamp,
					  					   monto decimal(7,2),
					  					   motivo text);

					  create table cierre(año int,
					  					  mes int,
					  					  terminacion int,
					  					  fechainicio date,
					  					  fechacierre date,
					  					  fechavto date);

					  create table cabecera(nroresumen int not null default nextval('aumentoCabecera'),
					  						nombre text,
					  						apellido text,
					  						domicilio text,
					  						nrotarjeta char(16),
					  						desde date,
					  						hasta date,
					  						vence date,
					  						total decimal(8,2));   

					  create table detalle(nroresumen int,
					  					   nrolinea int,
					  					   fecha date,
					  					   nombrecomercio text,
					  					   monto decimal(7,2));     

					  create table alerta(nroalerta int not null default nextval('aumentoAlerta'),
					  					  nrotarjeta char(16),
					  					  fecha timestamp,
					  					  nrorechazo int,
					  					  codalerta int,
					  					  descripcion text);

					  create table consumo(nrotarjeta char(16),
					  					   codseguridad char(4),
					  					   nrocomercio int,
					  					   monto decimal(7,2));

					  alter sequence aumentoCompra 
					  		increment 1 
					  		start 1 
					  		cache 1 
					  		owned by compra.nrooperacion;
					  		
					  alter sequence aumentoRechazo 
					  		increment 1 
					  		start 1 
					  		cache 1 
					  		owned by rechazo.nrorechazo;

					  alter sequence aumentoCabecera 
					  		increment 1 
					  		start 1 
					  		cache 1 
					  		owned by cabecera.nroresumen;					  		

					  alter sequence aumentoAlerta 
					  		increment 1 
					  		start 1 
					  		cache 1 
					  		owned by alerta.nroalerta;`)
    if err != nil {
    	log.Fatal(err)
    } 
}


/*
Define las Primary Keys y las Foreign Keys en la base de datos
*/
func definirPksYFks() {

	db, err := sql.Open("postgres", "user=postgres host=localhost dbname=tp sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	 _, err = db.Exec(`alter table cliente add constraint cliente_pk primary key (nrocliente);
	 				   alter table tarjeta add constraint tarjeta_pk primary key (nrotarjeta);
	 				   alter table comercio add constraint comercio_pk primary key (nrocomercio);
	 				   alter table compra add constraint compra_pk primary key (nrooperacion);
	 				   alter table rechazo add constraint rechazo_pk primary key (nrorechazo);
	 				   alter table cierre add constraint cierre_pk primary key (año,mes,terminacion);
	 				   alter table cabecera add constraint cabecera_pk primary key (nroresumen);
	 				   alter table detalle add constraint detalle_pk primary key (nroresumen, nrolinea);
	 				   alter table alerta add constraint alerta_pk primary key (nroalerta);

	 				   alter table tarjeta add constraint tarjeta_nrocliente_fk foreign key (nrocliente) references cliente(nrocliente);
	 				   alter table compra add constraint compra_nrotarjeta_fk foreign key (nrotarjeta) references tarjeta(nrotarjeta);
	 				   alter table compra add constraint compra_nrocomerio_fk foreign key (nrocomercio) references comercio(nrocomercio);
	 				   alter table rechazo add constraint rechazo_nrotarjeta_fk foreign key (nrotarjeta) references tarjeta(nrotarjeta);
	 				   alter table rechazo add constraint rechazo_nrocomerio_fk foreign key (nrocomercio) references comercio(nrocomercio);
	 				   alter table cabecera add constraint cabecera_nrotarjeta_fk foreign key (nrotarjeta) references tarjeta(nrotarjeta);
	 				   alter table alerta add constraint alerta_nrotarjeta_fk foreign key (nrotarjeta) references tarjeta(nrotarjeta);
	 				   alter table alerta add constraint alerta_nrorechazoa_fk foreign key (nrorechazo) references rechazo(nrorechazo);`)
    if err != nil {
    	log.Fatal(err)
    }
}


/*
Elimina las Primary Keys y las Foreign Keys de la base de datos
*/
func eliminarPksYFks() {

	db, err := sql.Open("postgres", "user=postgres host=localhost dbname=tp sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	 _, err = db.Exec(`alter table tarjeta drop constraint tarjeta_nrocliente_fk;
	 	 			   alter table compra drop constraint compra_nrotarjeta_fk;
	 	 			   alter table compra drop constraint compra_nrocomerio_fk;
	 	 			   alter table rechazo drop constraint rechazo_nrotarjeta_fk;
	 	 			   alter table rechazo drop constraint rechazo_nrocomerio_fk;
	 	 			   alter table cabecera drop constraint cabecera_nrotarjeta_fk;
	 	 			   alter table alerta drop constraint alerta_nrotarjeta_fk;
	 	 			   alter table alerta drop constraint alerta_nrorechazoa_fk;

	 	 			   alter table cliente drop constraint cliente_pk;
	 				   alter table tarjeta drop constraint tarjeta_pk;
	 				   alter table comercio drop constraint comercio_pk;
	 				   alter table compra drop constraint compra_pk;
	 				   alter table rechazo drop constraint rechazo_pk;
	 				   alter table cierre drop constraint cierre_pk;
	 				   alter table cabecera drop constraint cabecera_pk;
	 				   alter table detalle drop constraint detalle_pk;
	 				   alter table alerta drop constraint alerta_pk`)
    if err != nil {
    	log.Fatal(err)
    }
}


/*
Inserta los datos de los clientes en la base de datos
*/
func cargarClientes() {

	db, err := sql.Open("postgres", "user=postgres host=localhost dbname=tp sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	 _, err = db.Exec(`insert into cliente values(2981775,'Juan','Peron','Sarmiento 362','4331-1775');
	 				   insert into cliente values(2965465,'Nestor','Kirchner','Suipacha 1422','4327-0228');
	 				   insert into cliente values(8979845,'Cristina','Kirchner','Av. San Juan 328','5299-2010');
	 				   insert into cliente values(2313575,'Maria','Martinez','Av. Figueroa Alcorta 3415','4808-6500');
	 				   insert into cliente values(6267429,'Antonio','Cafiero','Av. del Libertador 8151','5280-0750');
	 				   insert into cliente values(5618468,'Carlos','Menem','Av. del Libertador 999','4800-1888');
	 				   insert into cliente values(1568484,'Alberto','Fernandez','Chacabuco 955','4362-5963');
	 				   insert into cliente values(1568432,'Jose','Gioja','Guevara 492','4553-9440');
	 				   insert into cliente values(6549832,'Luis','Barrionuevo','Av. Angel Gallardo 470','4822-8340');
	 				   insert into cliente values(8785512,'Eva','Peron','Tomás de Anchorena 1660','4982-6595');
	 				   insert into cliente values(7878798,'Homero','Simpson','Jean Jaures 735','4784-4040');
	 				   insert into cliente values(2135484,'Marge','Simpson','Av. Santa Fe 702','4342-3001');
	 				   insert into cliente values(0211541,'Bartolomeo','Simpson','Juramento 2291','4774-9452');
	 				   insert into cliente values(5421054,'Lisa','Simpson','Av. San Juan 350','4301-1080');
	 				   insert into cliente values(2161054,'Apu','Nahasapeesasolomoesolomoe','Av. del Libertador 2373','4361-4419');
	 				   insert into cliente values(3487910,'Martin','Price','Av. Infanta Isabel 555','4433-3396');
	 				   insert into cliente values(0216546,'Selma','Bouvie','Av. Pedro de Mendoza 1843','4343-2123');
	 				   insert into cliente values(2105646,'Patty','Bouvie','Av. España 1701','4370-6105');
	 				   insert into cliente values(8845660,'Barney','Gumble','Pujol 644','4893-0322');
	 				   insert into cliente values(8520147,'Waylon','Smithers','Defensa 219','4362-1100')`)
    if err != nil {
    	log.Fatal(err)
    }
}


/*
Inserta los datos de los tarjetas en la base de datos
*/
func cargarTarjetas() {

	db, err := sql.Open("postgres", "user=postgres host=localhost dbname=tp sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	 _, err = db.Exec(`insert into tarjeta values('4455674512546534', 2981775, '201106', '202306', '2020', 100000.00, 'vigente');
    		 		   insert into tarjeta values('1435471512346032', 2965465, '201510', '202812', '1212', 150000.00, 'vigente');
    		 		   insert into tarjeta values('9438541511146093', 8979845, '201207', '202411', '8807', 110000.00, 'vigente');
    		 		   insert into tarjeta values('2988781555176533', 2313575, '201609', '202910', '1331', 140000.00, 'vigente');
    		 		   insert into tarjeta values('5610556817653930', 6267429, '201301', '202703', '6854', 120000.00, 'vigente');
    		 		   insert into tarjeta values('2783277803985152', 5618468, '201402', '202504', '8902', 130000.00, 'vigente');
    		 		   insert into tarjeta values('4905123783542322', 1568484, '201112', '202401', '6723', 105000.00, 'vigente');
    		 		   insert into tarjeta values('7125367183482819', 1568432, '201202', '202903', '0913', 155000.00, 'vigente');
    		 		   insert into tarjeta values('8172631238129381', 6549832, '201305', '202408', '6512', 115000.00, 'vigente');
    		 		   insert into tarjeta values('9182743719349715', 8785512, '201410', '202606', '0782', 165000.00, 'vigente');
    		 		   insert into tarjeta values('4812346279123678', 7878798, '201303', '202411', '9991', 125000.00, 'vigente');
    		 		   insert into tarjeta values('9823478185734782', 2135484, '201703', '203010', '8172', 175000.00, 'vigente');
    		 		   insert into tarjeta values('6767676712371263', 0211541, '202009', '203512', '1119', 160000.00, 'vigente');
    		 		   insert into tarjeta values('9009723487324881', 5421054, '201801', '203308', '7667', 101000.00, 'vigente');
    		 		   insert into tarjeta values('1554388976675265', 2161054, '201004', '202301', '1865', 191000.00, 'vigente');
    		 		   insert into tarjeta values('8123612763817657', 3487910, '201912', '203312', '6661', 121000.00, 'vigente');
    		 		   insert into tarjeta values('7612376287677872', 0216546, '201111', '202301', '3942', 171000.00, 'vigente');
    		 		   insert into tarjeta values('7777612376765651', 2105646, '202102', '204309', '1942', 111000.00, 'vigente');
    		 		   insert into tarjeta values('8986664678589100', 8845660, '201108', '202112', '8888', 111100.00, 'vigente');
    		 		   insert into tarjeta values('8008555687165299', 8845660, '201708', '202608', '6666', 181500.00, 'vigente');
    		 		   insert into tarjeta values('3200111161616232', 8520147, '201403', '202804', '4423', 121500.00, 'vigente');
    		 		   insert into tarjeta values('1111323453543433', 8520147, '201606', '202907', '1456', 131500.00, 'vigente')`)	
    		 		   //tarjeta vencida: 8986664678589100 de cliente 8845660 Barney
    if err != nil {
    	log.Fatal(err)
    }
}


/*
Inserta los datos de los comercios en la base de datos
*/
func cargarComercio() {

	db, err := sql.Open("postgres", "user=postgres host=localhost dbname=tp sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	 _, err = db.Exec(`insert into comercio values(1935485,'McDonalds','Concejal Tribulato 636','1714','4541-6542');
	 				   insert into comercio values(3455465,'Peluqueria Adomo','Av. Victorica 421','1514','4785-1354');
	 				   insert into comercio values(9869845,'Supermercado Li','Crisóstomo Álvarez 2825','1406','4251-5268');
	 				   insert into comercio values(2314575,'Santeria Espacio Afrodita','Av.Juan Domingo Peron 1522','1663','4846-7639');
	 				   insert into comercio values(2609429,'Burger King','Santa Rosa 1680','1714','4245-9214');
	 				   insert into comercio values(6532468,'Garbarino','Av. Gral. Juan Manuel de Rosas 658','1712','4634-1480');
	 				   insert into comercio values(8712384,'Fravega','Av. Rivadavia 11626','1408','4664-5674');
	 				   insert into comercio values(3466632,'YPF','Av. Ricardo Balbin 1897','1650','4753-1745');
	 				   insert into comercio values(5632132,'Maxiconsumo','Gaona Acceso Oeste 8676','1744','4214-8413');
	 				   insert into comercio values(5570712,'Farmacia Fernandez','Bartolome Mitre 800','1742','4871-3587');
	 				   insert into comercio values(7860698,'Mostaza','Arturo Jauretche 978','1969','4984-45289');
	 				   insert into comercio values(2545384,'Coppel','Belgrano 3231','1650','4216-6512');
	 				   insert into comercio values(4163701,'Libreria Rodriguez','Independencia 4647','1653','4431-3218');
	 				   insert into comercio values(0923934,'Supermercado Dia','De la tradicion 185','1713','4451-9871');
	 				   insert into comercio values(2105614,'Peluqueria paty','Gral. Lavalle 848','1714','4721-9852');
	 				   insert into comercio values(4334530,'Santeria la paz','Vidal 1769','1426','4922-3265');
	 				   insert into comercio values(2054706,'Musimundo','Av. Lope de Vega 1520','1407','4154-3285');
	 				   insert into comercio values(1287436,'Shell','Sta Rosa 2489','1712','4970-8132');
	 				   insert into comercio values(9836840,'Casa del Audio','Rivadavia 2198','1714','4015-9872');
	 				   insert into comercio values(5419987,'KFC','Av. Bartolome Mitre','1744','4454-1134')`)
    if err != nil {
    	log.Fatal(err)
    }
}


/*
Inserta los datos de los cierres
*/
func cargarCierres() {

	db, err := sql.Open("postgres", "user=postgres host=localhost dbname=tp sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	 _, err = db.Exec(`insert into cierre values(2022,01,0,'2022-01-16','2022-01-15','2022-01-21');
	 				   insert into cierre values(2022,02,0,'2022-02-19','2022-02-18','2022-02-25');
	 				   insert into cierre values(2022,03,0,'2022-03-17','2022-03-16','2022-03-22');
	 				   insert into cierre values(2022,04,0,'2022-04-18','2022-04-17','2022-04-24');
	 	    		   insert into cierre values(2022,05,0,'2022-05-16','2022-05-15','2022-05-25');
	 				   insert into cierre values(2022,06,0,'2022-06-18','2022-06-17','2022-06-22');
	 				   insert into cierre values(2022,07,0,'2022-07-20','2022-07-19','2022-07-26');
	 				   insert into cierre values(2022,08,0,'2022-08-16','2022-08-15','2022-08-21');
	 				   insert into cierre values(2022,09,0,'2022-09-19','2022-09-18','2022-09-25');
	 				   insert into cierre values(2022,10,0,'2022-10-18','2022-10-17','2022-10-24');
	 				   insert into cierre values(2022,11,0,'2022-11-18','2022-11-17','2022-11-24');
	 				   insert into cierre values(2022,12,0,'2022-12-17','2022-12-16','2022-12-22'); 


                       insert into cierre values(2022,01,1,'2022-01-17','2022-01-16','2022-01-22');
	 				   insert into cierre values(2022,02,1,'2022-02-19','2022-02-18','2022-02-25');
	 				   insert into cierre values(2022,03,1,'2022-03-20','2022-03-19','2022-03-26');
	 				   insert into cierre values(2022,04,1,'2022-04-18','2022-04-17','2022-04-24');
	 				   insert into cierre values(2022,05,1,'2022-05-16','2022-05-15','2022-05-21');
	 				   insert into cierre values(2022,06,1,'2022-06-16','2022-06-15','2022-06-21');
	 				   insert into cierre values(2022,07,1,'2022-07-17','2022-07-16','2022-07-22');
	 				   insert into cierre values(2022,08,1,'2022-08-18','2022-08-17','2022-08-23');
	 				   insert into cierre values(2022,09,1,'2022-09-20','2022-09-19','2022-09-26');
	 				   insert into cierre values(2022,10,1,'2022-10-19','2022-10-18','2022-10-25');
	 				   insert into cierre values(2022,11,1,'2022-11-19','2022-11-18','2022-11-25');
	 				   insert into cierre values(2022,12,1,'2022-12-20','2022-12-19','2022-12-26');

	 				   insert into cierre values(2022,01,2,'2022-01-20','2022-01-19','2022-01-26');
	 				   insert into cierre values(2022,02,2,'2022-02-16','2022-02-15','2022-02-21');
	 				   insert into cierre values(2022,03,2,'2022-03-18','2022-03-17','2022-03-23');
	 				   insert into cierre values(2022,04,2,'2022-04-16','2022-04-15','2022-04-21');
	 				   insert into cierre values(2022,05,2,'2022-05-17','2022-05-16','2022-05-22');
	 				   insert into cierre values(2022,06,2,'2022-06-17','2022-06-16','2022-06-22');
	 				   insert into cierre values(2022,07,2,'2022-07-19','2022-07-18','2022-07-25');
	 				   insert into cierre values(2022,08,2,'2022-08-19','2022-08-17','2022-08-24');
	 				   insert into cierre values(2022,09,2,'2022-09-20','2022-09-19','2022-09-26');
	 				   insert into cierre values(2022,10,2,'2022-10-18','2022-10-17','2022-10-24');
	 				   insert into cierre values(2022,11,2,'2022-11-18','2022-11-17','2022-11-24');
	 				   insert into cierre values(2022,12,2,'2022-12-19','2022-12-18','2022-01-25');

	 				   insert into cierre values(2022,01,3,'2022-01-16','2022-01-15','2022-01-21');
	 				   insert into cierre values(2022,02,3,'2022-02-18','2022-02-17','2022-02-24');
	 				   insert into cierre values(2022,03,3,'2022-03-18','2022-03-17','2022-03-24');
	 				   insert into cierre values(2022,04,3,'2022-04-19','2022-04-18','2022-04-25');
	 				   insert into cierre values(2022,05,3,'2022-05-20','2022-05-19','2022-05-26');
	 				   insert into cierre values(2022,06,3,'2022-06-21','2022-06-20','2022-06-27');
	 				   insert into cierre values(2022,07,3,'2022-07-20','2022-07-19','2022-07-26');
	 				   insert into cierre values(2022,08,3,'2022-08-18','2022-08-17','2022-08-24');
	 				   insert into cierre values(2022,09,3,'2022-09-17','2022-09-16','2022-09-22');
	 				   insert into cierre values(2022,10,3,'2022-10-17','2022-10-16','2022-10-22');
	 				   insert into cierre values(2022,11,3,'2022-11-17','2022-11-15','2022-11-22');
	 				   insert into cierre values(2022,12,3,'2022-12-16','2022-12-15','2022-12-21');

	 				   insert into cierre values(2022,01,4,'2022-01-19','2022-01-18','2022-01-25');
	 				   insert into cierre values(2022,02,4,'2022-02-18','2022-02-17','2022-02-24');
	 				   insert into cierre values(2022,03,4,'2022-03-18','2022-03-17','2022-03-24');
	 				   insert into cierre values(2022,04,4,'2022-04-16','2022-04-15','2022-04-21');
	 				   insert into cierre values(2022,05,4,'2022-05-17','2022-05-16','2022-05-22');
	 				   insert into cierre values(2022,06,4,'2022-06-18','2022-06-17','2022-06-24');
	 				   insert into cierre values(2022,07,4,'2022-07-20','2022-07-19','2022-07-26');
	 				   insert into cierre values(2022,08,4,'2022-08-20','2022-08-19','2022-08-26');
	 				   insert into cierre values(2022,09,4,'2022-09-19','2022-09-18','2022-09-25');
	 				   insert into cierre values(2022,10,4,'2022-10-17','2022-10-16','2022-10-22');
	 				   insert into cierre values(2022,11,4,'2022-11-18','2022-11-17','2022-11-24');
	 				   insert into cierre values(2022,12,4,'2022-12-19','2022-12-18','2022-12-25');

	 				   insert into cierre values(2022,01,5,'2022-01-20','2022-01-19','2022-01-26');
	 				   insert into cierre values(2022,02,5,'2022-02-20','2022-02-19','2022-02-26');
	 				   insert into cierre values(2022,03,5,'2022-03-18','2022-03-17','2022-03-24');
	 				   insert into cierre values(2022,04,5,'2022-04-19','2022-04-18','2022-04-25');
	 				   insert into cierre values(2022,05,5,'2022-05-21','2022-05-20','2022-05-27');
	 				   insert into cierre values(2022,06,5,'2022-06-21','2022-06-20','2022-06-27');
	 				   insert into cierre values(2022,07,5,'2022-07-19','2022-07-18','2022-07-25');
	 				   insert into cierre values(2022,08,5,'2022-08-18','2022-08-17','2022-08-24');
	 				   insert into cierre values(2022,09,5,'2022-09-17','2022-09-16','2022-09-22');
	 				   insert into cierre values(2022,10,5,'2022-10-17','2022-10-16','2022-10-22');
	 				   insert into cierre values(2022,11,5,'2022-11-18','2022-11-17','2022-11-24');
	 				   insert into cierre values(2022,12,5,'2022-12-18','2022-12-18','2022-12-24');

	 				   insert into cierre values(2022,01,6,'2022-01-17','2022-01-16','2022-01-22');
	 				   insert into cierre values(2022,02,6,'2022-02-16','2022-02-15','2022-02-21');
	 				   insert into cierre values(2022,03,6,'2022-03-18','2022-03-17','2022-03-24');
	 				   insert into cierre values(2022,04,6,'2022-04-16','2022-04-15','2022-04-21');
	 				   insert into cierre values(2022,05,6,'2022-05-17','2022-05-16','2022-05-22');
	 				   insert into cierre values(2022,06,6,'2022-06-18','2022-06-17','2022-06-24');
	 				   insert into cierre values(2022,07,6,'2022-07-19','2022-07-18','2022-07-25');
	 				   insert into cierre values(2022,08,6,'2022-08-16','2022-08-15','2022-08-21');
	 				   insert into cierre values(2022,09,6,'2022-09-16','2022-09-15','2022-09-21');
	 				   insert into cierre values(2022,10,6,'2022-10-17','2022-10-16','2022-10-22');
	 				   insert into cierre values(2022,11,6,'2022-11-19','2022-11-18','2022-11-25');
	 				   insert into cierre values(2022,12,6,'2022-12-20','2022-12-19','2022-12-26');

	 				   insert into cierre values(2022,01,7,'2022-01-20','2022-01-19','2022-01-26');
	 				   insert into cierre values(2022,02,7,'2022-02-18','2022-02-17','2022-02-24');
	 				   insert into cierre values(2022,03,7,'2022-03-18','2022-03-17','2022-03-24');
	 				   insert into cierre values(2022,04,7,'2022-04-19','2022-04-18','2022-04-25');
	 				   insert into cierre values(2022,05,7,'2022-05-17','2022-05-16','2022-05-22');
	 				   insert into cierre values(2022,06,7,'2022-06-21','2022-06-20','2022-06-27');
	 				   insert into cierre values(2022,07,7,'2022-07-20','2022-07-19','2022-07-26');
	 				   insert into cierre values(2022,08,7,'2022-08-18','2022-08-17','2022-08-24');
	 				   insert into cierre values(2022,09,7,'2022-09-19','2022-09-18','2022-09-25');
	 				   insert into cierre values(2022,10,7,'2022-10-19','2022-10-18','2022-10-25');
	 				   insert into cierre values(2022,11,7,'2022-11-21','2022-11-20','2022-11-27');
	 				   insert into cierre values(2022,12,7,'2022-12-19','2022-12-18','2022-12-25');
	 				   
	 				   insert into cierre values(2022,01,8,'2022-01-16','2022-01-15','2022-01-21');
	 				   insert into cierre values(2022,02,8,'2022-02-17','2022-02-16','2022-02-22');
	 				   insert into cierre values(2022,03,8,'2022-03-17','2022-03-16','2022-03-22');
	 				   insert into cierre values(2022,04,8,'2022-04-19','2022-04-18','2022-04-25');
	 				   insert into cierre values(2022,05,8,'2022-05-17','2022-05-16','2022-05-22');
	 				   insert into cierre values(2022,06,8,'2022-06-18','2022-06-17','2022-06-24');
	 				   insert into cierre values(2022,07,8,'2022-07-19','2022-07-18','2022-07-25');
	 				   insert into cierre values(2022,08,8,'2022-08-18','2022-08-17','2022-08-23');
	 				   insert into cierre values(2022,09,8,'2022-09-17','2022-09-16','2022-09-22');
	 				   insert into cierre values(2022,10,8,'2022-10-18','2022-10-17','2022-10-24');
	 				   insert into cierre values(2022,11,8,'2022-11-18','2022-11-17','2022-11-24');
	 				   insert into cierre values(2022,12,8,'2022-12-20','2022-12-19','2022-12-26');

	 				   insert into cierre values(2022,01,9,'2022-01-16','2022-01-15','2022-01-21');
	 				   insert into cierre values(2022,02,9,'2022-02-16','2022-02-15','2022-02-21');
	 				   insert into cierre values(2022,03,9,'2022-03-18','2022-03-17','2022-03-24');
	 				   insert into cierre values(2022,04,9,'2022-04-19','2022-04-18','2022-04-25');
	 				   insert into cierre values(2022,05,9,'2022-05-17','2022-05-16','2022-05-22');
	 				   insert into cierre values(2022,06,9,'2022-06-21','2022-06-20','2022-06-27');
	 				   insert into cierre values(2022,07,9,'2022-07-20','2022-07-19','2022-07-26');
	 				   insert into cierre values(2022,08,9,'2022-08-18','2022-08-17','2022-08-23');
	 				   insert into cierre values(2022,09,9,'2022-09-16','2022-09-15','2022-09-21');
	 				   insert into cierre values(2022,10,9,'2022-10-18','2022-10-17','2022-10-24');
	 				   insert into cierre values(2022,11,9,'2022-11-19','2022-11-18','2022-11-25');
	 				   insert into cierre values(2022,12,9,'2022-12-19','2022-12-18','2022-12-25');`)
    if err != nil {
    	log.Fatal(err)
    }
}


/*
Inserta los datos de los consumos para probar
*/
func cargarConsumos() {

	db, err := sql.Open("postgres", "user=postgres host=localhost dbname=tp sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`insert into consumo values('4455674512546534','2020',1935485,11320.00);
					  insert into consumo values('4455674512546534','2020',9836840,15840.50);	 
					  insert into consumo values('1435471512346032','1212',1935485,8900.50); 
					  insert into consumo values('1435471512346032','1212',3455465,5750.00);
					  insert into consumo values('1435471512346032','1213',3455465,100.00);
					  insert into consumo values('7777612376765651','1942',3455465,2005.00); 
					  insert into consumo values('8986664678589100','8888',3455465,24015.63);
					  insert into consumo values('4455674512546534','2020',1935485,99999.98);
					  insert into consumo values('4455674512546534','2020',3455465,99999.99);
					  insert into consumo values('1234567898765432','7069',3455465,1111.10);`)
	if err != nil {
		log.Fatal(err)
	}
}


/*
Se encarga de las validaciones para autorizar las compras
1) Que la tarjeta exista, si no cumple se carga un rechazo con el mensaje: "Tarjeta inexistente"
2) Que el codigo de seguridad sea correcto, si no cumple se carga un rechazo con el mensaje: "Codigo de seguridad invalido"
3) Que el monto total de compras pendientes de pago no supere el limite de compra de la tarjeta, si no cumple se carga un rechazo con el mensaje: "Supera limite de tarjeta"
4) Que la tarjeta no se encuentre vencida, si no cumple se carga un rechazo con el mensaje: "Plazo de vigencia expirado"
5) Que la tarjeta no se encuentre suspendida, si no cumple se carga un rechazo con el mensaje: "La tarjeta se encuentra suspendida"
*/
func autorizarCompra() {

	db, err := sql.Open("postgres", "user=postgres host=localhost dbname=tp sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Query(`create or replace function autorizarCompra(nrotarjetaAux char(16), codseguridadAux char(4), nrocomercioAux int, montoAux decimal(7, 2)) returns boolean as $$

					   declare
					   		resultado record;
					   		parcial decimal(7, 2);
					   		total decimal(8, 2);

					   begin
					   		perform * from tarjeta where nrotarjeta = nrotarjetaAux;
					   		if not found then
					   			insert into rechazo values(nextval('aumentoRechazo'),null,nrocomercioAux,current_timestamp,montoAux,'Tarjeta inexistente');
					   			return false;
					   		else
								select * into resultado from tarjeta where nrotarjeta = nrotarjetaAux and codseguridad = codseguridadAux;
								if not found then
									insert into rechazo values(nextval('aumentoRechazo'),nrotarjetaAux,nrocomercioAux,current_timestamp,montoAux,'Codigo de seguridad invalido');
									return false;
								else
									total := 0;
									for parcial in select monto from compra where compra.nrotarjeta = nrotarjetaAux and compra.pagado = true loop
										total := total + parcial;
									end loop;
									total := total + montoAux;
									select * into resultado from tarjeta t where t.nrotarjeta = nrotarjetaAux and t.limitecompra > total;
									if not found then
										insert into rechazo values(nextval('aumentoRechazo'),nrotarjetaAux,nrocomercioAux,current_timestamp,montoAux,'Supera limite de tarjeta');
										return false;
									else
										select * into resultado from tarjeta where nrotarjeta = nrotarjetaAux and to_date(validahasta, 'YYYYMM') >= to_date('202201', 'YYYYMM');
										if not found then
											update tarjeta set estado = 'anulada' where nrotarjeta = nrotarjetaAux;
											insert into rechazo values(nextval('aumentoRechazo'),nrotarjetaAux,nrocomercioAux,current_timestamp,montoAux,'Plazo de vigencia expirado');
											return false;
										else
											select * into resultado from tarjeta where nrotarjeta = nrotarjetaAux and estado = 'vigente';
											if not found then
												insert into rechazo values(nextval('aumentoRechazo'),nrotarjetaAux,nrocomercioAux,current_timestamp,montoAux,'La tarjeta se encuentra suspendida');
												return false;
											else
												insert into compra values(nextval('aumentoCompra'),nrotarjetaAux,nrocomercioAux,current_timestamp,montoAux,true);
												return true;
											end if;		
										end if;	
									end if;
								end if;
							end if;	
						end;
						$$ language plpgsql;;`)
	if err != nil {
		log.Fatal(err)
	}
}


/*
Se encarga de generar el resumen con las compras que realizo cada uno de los clientes
*/
func generarResumen() {

	db, err := sql.Open("postgres", "user=postgres host=localhost dbname=tp sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Query(`create or replace function generarResumen(numcliente int, periodo_mes int, periodo_anio int) returns void as $$

					   declare
					   		total decimal(7,2);
					   		monto_parcial decimal(7,2);
					   		nro_tarjeta char(16);
					   		nombre_cliente text;
					   		apellido_cliente text;
					   		domicilio_cliente text;
					   		auxiliar record;
					   		ultimoNro int;
					   		fecha_inicio date;
					   		fecha_cierre date;
					   		fecha_vto date;

					   begin
					   		perform * from cliente where nrocliente = numcliente;
					   		if (not found) then
					   			raise 'El nro de cliente % es invalido',numcliente;
					   		end if;
					   		
					   		if (periodo_anio!=2022 or periodo_mes<0 or periodo_mes>12) then 
						   		raise 'El mes % es invalido',periodo_mes;
						   		raise 'El año % es invalido',periodo_anio;
					   		end if;

					   		select * into auxiliar from cliente where nrocliente = numcliente;
					   		nombre_cliente := auxiliar.nombre;
					   		apellido_cliente := auxiliar.apellido;
					   		domicilio_cliente := auxiliar.domicilio;
			   		
					   		select * into auxiliar  from tarjeta t where t.nrocliente = numcliente and t.estado != 'anulada';
					   		nro_tarjeta := auxiliar.nrotarjeta;
					   		ultimoNro := right (nro_tarjeta,1);
				   		
					   		select * into auxiliar from cierre where año = periodo_anio and mes = periodo_mes and terminacion = ultimoNro;
					   		fecha_inicio:= auxiliar.fechainicio;
					   		fecha_cierre:= auxiliar.fechacierre;
					   		fecha_vto:= auxiliar.fechavto;
					   		
					   		total := 0.00;
					   		for monto_parcial in select c.monto from compra c
					   		where (select extract (month from c.fecha) = periodo_mes and c.nrotarjeta=nro_tarjeta) loop
					   			total := total+ monto_parcial;
					   		end loop;
				   		
					   		insert into cabecera values (nextval('aumentoCabecera'),nombre_cliente, apellido_cliente, domicilio_cliente,
					   									nro_tarjeta,fecha_inicio,fecha_cierre, fecha_vto, total);
					   		end;
					   		$$ language plpgsql;`)
	if err != nil {
		log.Fatal(err)
	}
}


/*
Utiliza la funcion de SQL autorizarCompra(char, char, int, decimal(7,2)) para probar los consumos
*/
func probarConsumos() {

	db, err := sql.Open("postgres", "user=postgres host=localhost dbname=tp sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Query(`create or replace function probarConsumos() returns void as $$

					   declare
					   		v record;
					   		resultado record;

					   begin
					   		for v in select * from consumo loop
								select autorizarCompra(v.nrotarjeta,v.codseguridad,v.nrocomercio,v.monto) into resultado;
							end loop;
					   end;
					   $$ language plpgsql;`)
	if err != nil {
		log.Fatal(err)
	}
}


/*
Test para probar la funcion de probarConsumos()
*/
func llamarConsumos() {

	db, err := sql.Open("postgres", "user=postgres host=localhost dbname=tp sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Query(`select probarConsumos();`)
	if err != nil {
		log.Fatal(err)
	}
}


/*
Utiliza la funcion de SQL generarResumen(int, date, date) para probar la generacion de resumenes
*/
func probarResumen() {

	db, err := sql.Open("postgres", "user=postgres host=localhost dbname=tp sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Query(`create or replace function probarResumen() returns void as $$ 

					   declare
					   		v record;
					   		a record;
					   		resultado record;
					   		ultimoNro int;

					   begin
					   		for v in select * from tarjeta t where t.estado != 'anulada' loop
					   			ultimoNro := right (v.nrotarjeta,1);
					   			for a in select * from cierre where terminacion = ultimoNro loop
					   				select generarResumen(v.nrocliente, a.mes, a.año) into resultado;
					   			end loop;
					   		end loop;
					   end;
					   $$ language plpgsql;`)	
	if err != nil {
		log.Fatal(err)
	}
}


/*
Test para probar la funcion probarResumen()
*/
func llamarResumen() {

	db, err := sql.Open("postgres", "user=postgres host=localhost dbname=tp sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Query(`select probarResumen();`)
	if err != nil {
		log.Fatal(err)
	}
}


/*
Trigger para ingresar a la tabla de alertas todos los rechazos automaticamente 
*/
func ingresoRechazoAlerta() {

	db, err := sql.Open("postgres", "user=postgres host=localhost dbname=tp sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Query(`create or replace function ingresoRechazoAlerta() returns trigger as $$

					   begin
					   		insert into alerta(nroalerta,nrotarjeta,fecha,nrorechazo,codalerta,descripcion) 
					   					values(nextval('aumentoAlerta'),new.nrotarjeta, new.fecha,new.nrorechazo,0,'Se produjo un nuevo rechazo');
					   		return new;
					   end;
					   $$ language plpgsql;

					   create trigger ingresoRechazoAlerta_trg 
					   after insert on rechazo for each row execute procedure ingresoRechazoAlerta();`)
	if err != nil {
		log.Fatal(err)
	}
}


/*
Trigger para ingresar a la tabla de alertas si una tarjeta registra 2 compras en menos de 1 minuto
en comercios ubicados en el mismo codigo postal
*/
func dosComprasMenosUnMinutoMismoCPAlerta() {

	db, err := sql.Open("postgres", "user=postgres host=localhost dbname=tp sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Query(`create or replace function dosComprasMenosUnMinutoMismoCPAlerta() returns trigger as $$
				   
					   declare
					   		ultima record;
					   		tiempo decimal;
					   		cp1 record;
					   		cp2 record;

					   begin
					   		select * into ultima from compra where nrotarjeta = new.nrotarjeta order by nrotarjeta desc limit 1;
					   		if not found then
					   			return new;
					   		end if;
			   		
							select into tiempo extract(epoch from new.fecha - ultima.fecha);
							select codigopostal into cp1 from comercio where nrocomercio = ultima.nrocomercio;
							select codigopostal into cp2 from comercio where nrocomercio = new.nrocomercio;

							if tiempo < 60 and ultima.nrocomercio != new.nrocomercio and cp1 = cp2 then
								insert into alerta (nroalerta,nrotarjeta,fecha,nrorechazo,codalerta,descripcion)
											values(nextval('aumentoAlerta'), new.nrotarjeta, new.fecha, null, 1, 'Compra en menos de 1 minuto en mismo CP');
							end if;
							return new;
					   	end;
					   	$$ language plpgsql;
				   
					   	create trigger dosComprasMenosUnMinutoMismoCPAlerta_trg 
					   	before insert on compra for each row execute procedure dosComprasMenosUnMinutoMismoCPAlerta();`)
	if err != nil {
		log.Fatal(err)
	}
}


/*
Trigger para ingresar a la tabla de alertas si una tarjeta registra 2 compras en menos de 5 minutos
en comercios ubicados en diferentes codigos postales
*/
func dosComprasMenosCincoMinutosDifCPAlerta() {

	db, err := sql.Open("postgres", "user=postgres host=localhost dbname=tp sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Query(`create or replace function dosComprasMenosCincoMinutosDifCPAlerta() returns trigger as $$
					   
					   declare
					   		ultima record;
					   		tiempo decimal;
					   		cp1 record;
					   		cp2 record;

					   begin
					   		select * into ultima from compra where nrotarjeta = new.nrotarjeta order by nrotarjeta desc limit 1;
					   		if not found then
					   			return new;
					   		end if;
				   		
							select into tiempo extract(epoch from new.fecha - ultima.fecha);
							select codigopostal into cp1 from comercio where nrocomercio = ultima.nrocomercio;
							select codigopostal into cp2 from comercio where nrocomercio = new.nrocomercio;

							if tiempo < 300 and ultima.nrocomercio != new.nrocomercio and cp1 != cp2 then
								insert into alerta (nroalerta,nrotarjeta,fecha,nrorechazo,codalerta,descripcion)
											values(nextval('aumentoAlerta'), new.nrotarjeta, new.fecha, null, 5, 'Compra en menos de 5 minutos en diferente CP');
							end if;
							return new;
					   	end;
					   	$$ language plpgsql;
				   
					   	create trigger dosComprasMenosCincoMinutosDifCPAlerta_trg 
					   	before insert on compra for each row execute procedure dosComprasMenosCincoMinutosDifCPAlerta();`)
	if err != nil {
		log.Fatal(err)
	}
}


/*
Trigger para ingresar a la tabla de alertas si una tarjeta registra 2 rechazos en el mismo dia y suspenderla preventivamente
*/
func excesoLimiteAlerta() {

	db, err := sql.Open("postgres", "user=postgres host=localhost dbname=tp sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Query(`create or replace function excesoLimiteAlerta() returns trigger as $$

					   declare
					   		ultimo record;
					   		anio decimal;
					   		mes decimal;
					   		dia decimal;
				   		
					   begin
					    	select * into ultimo from rechazo where nrotarjeta = new.nrotarjeta and motivo = 'Supera limite de tarjeta' order by nrotarjeta desc limit 1;
					    	if not found then
					    		return new;
					    	end if;
				    	
					    	select into dia extract(day from new.fecha - ultimo.fecha);
					    	select into mes extract(month from new.fecha - ultimo.fecha);
					    	select into anio extract(year from new.fecha - ultimo.fecha);
					    	
					    	if dia < 1 and mes < 1 and anio < 1 then
					    		update tarjeta set estado = 'suspendida' where nrotarjeta = new.nrotarjeta;
					    		insert into alerta (nroalerta,nrotarjeta,fecha,nrorechazo,codalerta,descripcion)
					    					values(nextval('aumentoAlerta'), new.nrotarjeta, new.fecha, null, 32, 'Tarjeta suspendida por 2 excesos de limite en el mismo dia');
					    	end if;
					    	return new;
					    end;
					    $$ language plpgsql;

					    create trigger excesoLimiteAlerta_trg 
					    before insert on rechazo for each row execute procedure excesoLimiteAlerta();`)
	if err != nil {
		log.Fatal(err)
	}
}


/*
Marshalea los datos de Cliente, los escribe en la base de datos NoSQL, los lee y luego los muestra por pantalla
*/
func clienteBB(nrocliente int, nombre string, apellido string, domicilio string, telefono string) {

	cliente := Cliente{nrocliente, nombre, apellido, domicilio, telefono}
	data, err := json.Marshal(cliente)
	if err != nil {
		log.Fatal(err)
	}

	CreateUpdate(dbbolt, "Cliente", []byte(strconv.Itoa(cliente.Nrocliente)), data)

	read, _ := ReadUnique(dbbolt, "Cliente",[]byte(strconv.Itoa(cliente.Nrocliente)))
	var lec string = fmt.Sprintf("%s\n", read)
	fmt.Println(lec)
}


/*
Marshalea los datos de Tarjeta, los escribe en la base de datos NoSQL, los lee y luego los muestra por pantalla
*/
func tarjetaBB(nrotarjeta string, nrocliente int, validadesde string, validahasta string, codseguridad string, limitecompra float32, estado string) {

	tarjeta := Tarjeta{nrotarjeta, nrocliente, validadesde, validahasta, codseguridad, limitecompra, estado}
	data, err := json.Marshal(tarjeta)
	if err != nil {
		log.Fatal(err)
	}

	CreateUpdate(dbbolt, "Tarjeta", []byte(tarjeta.Nrotarjeta), data)

	read, _ := ReadUnique(dbbolt, "Tarjeta",[]byte(tarjeta.Nrotarjeta))
	var lec string = fmt.Sprintf("%s\n", read)
	fmt.Println(lec)
}


/*
Marshalea los datos de Comercio, los escribe en la base de datos NoSQL, los lee y luego los muestra por pantalla
*/
func comercioBB(nrocomercio int, nombre string, domicilio string, codigopostal string, telefono string) {

	comercio := Comercio{nrocomercio, nombre, domicilio, codigopostal, telefono}
	data, err := json.Marshal(comercio)
	if err != nil {
		log.Fatal(err)
	}

	CreateUpdate(dbbolt, "Comercio", []byte(strconv.Itoa(comercio.Nrocomercio)), data)

	read, _ := ReadUnique(dbbolt, "Comercio",[]byte(strconv.Itoa(comercio.Nrocomercio)))
	var lec string = fmt.Sprintf("%s\n", read)
	fmt.Println(lec)
}


/*
Marshalea los datos de Compra, los escribe en la base de datos NoSQL, los lee y luego los muestra por pantalla
*/
func compraBB(nrooperacion int, nrotarjeta string, nrocomercio int, fecha string, monto float32, pagado bool) {

	compra := Compra{nrooperacion, nrotarjeta, nrocomercio, fecha, monto, pagado}
	data, err := json.Marshal(compra)
	if err != nil {
		log.Fatal(err)
	}
	
	CreateUpdate(dbbolt, "Compra", []byte(strconv.Itoa(compra.Nrooperacion)), data)

	read, _ := ReadUnique(dbbolt, "Compra",[]byte(strconv.Itoa(compra.Nrooperacion)))
	var lec string = fmt.Sprintf("%s\n", read)
	fmt.Println(lec)
}


/*
Transaccion de escritura
*/
func CreateUpdate(dbbolt *bolt.DB, bucketName string, key []byte, val []byte) error {

	tx, err := dbbolt.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	b, _ := tx.CreateBucketIfNotExists([]byte(bucketName))
	err = b.Put(key, val)
	if err != nil {
		return err
	}
	
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}


/*
Transaccion de lectura
*/
func ReadUnique(dbbolt *bolt.DB, bucketName string, key []byte) ([]byte, error) {

	var buf []byte

	err := dbbolt.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		buf = b.Get(key)
		return nil
	})
	return buf, err
}


/*
Creamos los diccionarios
*/
func datosBB() {

	var err error
	dbbolt, err = bolt.Open("tp.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer dbbolt.Close()

	clienteBB(2981775, "Juan", "Peron", "Sarmiento 362", "4331-1775")
	clienteBB(2965465, "Nestor", "Kirchner", "Suipacha 1422", "4327-0228")
	clienteBB(8979845, "Cristina", "Kirchner", "Av. San Juan", "5299-2010")

	tarjetaBB("4455674512546534", 2981775, "201106", "202306", "2020", 100000.00, "vigente")
	tarjetaBB("1435471512346032", 2965465, "201510", "202812", "1212", 150000.00, "vigente")
	tarjetaBB("9438541511146093", 8979845, "201207", "202411", "8807", 110000.00, "vigente")

	comercioBB(1935485, "McDonalds", "Concejal Tribulato 636", "1714", "4541-6542")
	comercioBB(3455465, "Peluqueria Adomo", "Av. Victorica 421", "1514", "4785-1354")
	comercioBB(9836840, "Casa del Audio", "Rivadavia 2198", "1714", "4015-9872")

	compraBB(1, "4455674512546534", 1935485, "2022-06-12 01:19:24.435615", 11320.00, true)
	compraBB(2, "4455674512546534", 9836840, "2022-06-12 01:19:24.435615", 15840.50, true)
	compraBB(3, "1435471512346032", 1935485, "2022-06-12 01:19:24.435615", 8900.50, true)
}


/*
Menu principal creado con "github.com/dixonwille/wmenu/v5"
Dentro de cada opcion se encuentran las funciones utilizadas en el proyecto
*/
func menuPrincipal() *wmenu.Menu {

	menu := wmenu.NewMenu("Seleccione una opcion")
	menu.Option("Crear base de datos", nil, false, func(opt wmenu.Opt) error {
		fmt.Println("")
		crearBD()
		mv := menuVolver()
		return mv.Run()
	})
	menu.Option("Eliminar base de datos", nil, false, func(opt wmenu.Opt) error {
		fmt.Println("")
		eliminarBD()
		mv := menuVolver()
		return mv.Run()
	})		
	menu.Option("Crear tablas y secuencias", nil, false, func(opt wmenu.Opt) error {
		fmt.Println("")
		crearTablasYSecuencias()
		mv := menuVolver()
		return mv.Run()
	})
	menu.Option("Agregar Primary Keys y Foreign Keys", nil, false, func(opt wmenu.Opt) error {
		fmt.Println("")
		definirPksYFks()
		mv := menuVolver()
		return mv.Run()
	})
	menu.Option("Eliminar Primary Keys y Foreign Keys", nil, false, func(opt wmenu.Opt) error {
		fmt.Println("")
		eliminarPksYFks()
		mv := menuVolver()
		return mv.Run()
	})
	menu.Option("Cargar datos en las tablas", nil, false, func(opt wmenu.Opt) error {
		fmt.Println("")
		cargarClientes()
		cargarComercio()
		cargarTarjetas()
		cargarCierres()
		cargarConsumos()
		mv := menuVolver()
		return mv.Run() 
	})
	menu.Option("Cargar funciones y triggers", nil, false, func(opt wmenu.Opt) error {
		fmt.Println("")
		ingresoRechazoAlerta()
		autorizarCompra()
		dosComprasMenosUnMinutoMismoCPAlerta()
		dosComprasMenosCincoMinutosDifCPAlerta()
		excesoLimiteAlerta()
		probarConsumos()
		llamarConsumos()
		generarResumen()
		probarResumen()
		llamarResumen()
		mv := menuVolver()
		return mv.Run()
	})
	menu.Option("Agregar datos en base de datos NoSQL y mostrarlos por pantalla", nil, false, func(opt wmenu.Opt) error {
		fmt.Println("")
		datosBB()
		mv := menuVolver()
		return mv.Run()
		})
	menu.Option("Salir", nil, false, func(opt wmenu.Opt) error {
		fmt.Println("")
		fmt.Println("Gracias por usar nuestro programa")	
		os.Exit(0)
		return nil
	})
	return menu
}


/*
Menu secundario para poder volver al menu principal creado con "github.com/dixonwille/wmenu/v5"
*/
func menuVolver() *wmenu.Menu {

	menu := wmenu.NewMenu("")
	menu.Option("Volver al menu principal", nil, false, func(opt wmenu.Opt) error {
		fmt.Println("")
		mp := menuPrincipal()
		return mp.Run()
	})
	menu.Option("Salir", nil, false, func(opt wmenu.Opt) error {
		fmt.Println("")
		fmt.Println("Gracias por usar nuestro programa")	
		os.Exit(0)
		return nil
	})
	return menu
}


/*
Main donde se inicia el menu
*/
func main() {

	m := menuPrincipal()
	err := m.Run()
	if err != nil {
		log.Fatal(err)
		fmt.Println(err)
		os.Exit(1)
	}
}
