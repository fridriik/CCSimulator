package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	_"github.com/lib/pq"
	"github.com/dixonwille/wmenu/v5"
)


type cliente struct {
	nrocliente int
    nombre string
    apellido string
	domicilio string
	telefono [12] rune
}

type tarjeta struct {
    nrotarjeta [16] rune
    nrocliente int
    validadesde [6] rune
    validadhasta [6] rune
    codseguridad [4] rune
    limitecompra float32    //decimal(8,2)
    estado [10] rune
}

type comercio struct {
    nrocomercio int
    nombre string
    domicilio string
    codigopostal [8] rune
    telefono [12] rune
}

type compra struct {
    nrooperacion int
    nrotarjeta [16] rune
    nrocomercio int
    fecha string    //time.Time
    monto float32   //decimal(7,2)
    pagado bool
}

type rechazo struct {
    nrorechazo int
    nrotarjeta [16] rune
    nrocomercio int
    fecha string    //time.Time
    monto float32   //decimal(7,2)
    motivo string
}

type cierre struct {
    anio int
    mes int
    terminacion int
    fechainicio string  //time.Time
    fechacierre string  //time.Time
    fechavto string   //time.Time
}

type cabecera struct {
    nroresumen int
    nombre string
    apellido string
    domicilio string
    nrotarjeta [16] rune
    desde string    //time.Time
    hasta string    //time.Time
    vence string    //time.Time
    total float32   //decimal(8,2)
}

type detalle struct {
    nroresumen int
    nrolinea int
    fecha string    //time.Time
    nombrecomercio string
    monto float32   //decimal(7,2)
}

type alerta struct {
    nroalerta int
    nrotarjeta [16] rune
    fecha string    //time.Time
    nrorechazo int
    codalerta int   //0: rechazo, 1: compra 1 min, 5: compra 5 min, 32: limite
    descripcion string
}

type consumo struct {
    nrotarjeta [16] rune
    codseguridad [4] rune
    nrocomercio int
    monto float32   //decimal(7,2)
}


func createDatabase() {

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

		
func deleteDatabase() {
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



func crearTablas() {

	db, err := sql.Open("postgres", "user=postgres host=localhost dbname=tp sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	
	_, err = db.Exec(
 `create sequence aumentoCompra;
  create sequence aumentoRechazo;
  create table cliente(
							nrocliente int,
							nombre text,
    						apellido text,
    						domicilio text,
    						telefono char(12));
  create table tarjeta(
   							nrotarjeta char(16),
							nrocliente int,
							validadesde char(6), 
							validahasta char(6),
							codseguridad char(4),
							limitecompra decimal(8,2),
							estado char(10));
  create table comercio(
   							nrocomercio int,
   							nombre text,
    						domicilio text,
    						codigopostal char(8),
    						telefono char(12));
  create table compra(
  							nrooperacion int not null default nextval('aumentoCompra'),
    						nrotarjeta char(16),
    						nrocomercio int,
    						fecha timestamp,
    						monto decimal(7,2),
    						pagado boolean);
  create table rechazo(
  							nrorechazo int not null default nextval('aumentoRechazo'),
    						nrotarjeta char(16),
    						nrocomercio int,
    						fecha timestamp,
    						monto decimal(7,2),
    						motivo text);
  create table cierre(
  							año int,
    						mes int,
    						terminacion int,
    						fechainicio date,
    						fechacierre date,
							fechavto date);
  create table cabecera(
  							nroresumen int,
    						nombre text,
    						apellido text,
    						domicilio text,
    						nrotarjeta char(16),
    						desde date,
    						hasta date,
    						vence date,
    						total decimal(8,2));   
  create table detalle(
 							nroresumen int,
    						nrolinea int,
    						fecha date,
    						nombrecomercio text,
    						monto decimal(7,2));     
  create table alerta(
  							nroalerta int,
    						nrotarjeta char(16),
    						fecha timestamp,
    						nrorechazo int,
    						codalerta int,
    						descripcion text);
  create table consumo(
   							nrotarjeta char(16),
    						codseguridad char(4),
    						nrocomercio int,
    						monto decimal(7,2));
  alter sequence aumentoCompra increment 1 start 1 cache 1 owned by compra.nrooperacion;
  alter sequence aumentoRechazo increment 1 start 1 cache 1 owned by rechazo.nrorechazo;`)
    if err != nil {
    	log.Fatal(err)
    } 
}
	


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
    		 				   insert into tarjeta values('1111323453543433', 8520147, '201606', '202907', '1456', 131500.00, 'vigente')`)	//tarjeta vencida: 8986664678589100 de cliente 8845660 Barney
    if err != nil {
    	log.Fatal(err)
    }
}


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


func cargarComercio() {

	db, err := sql.Open("postgres", "user=postgres host=localhost dbname=tp sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	 _, err = db.Exec(`insert into comercio values(1935485,'McDonalds','Concejal Tribulato 636','1744','6541-6542');
	 				   insert into comercio values(3455465,'Peluqueria Adomo','Av. Victorica 421','1614','6785-1354');
	 				   insert into comercio values(9869845,'Supermercado Li','Crisóstomo Álvarez 2825','1406','2251-5268');
	 				   insert into comercio values(2314575,'Santeria Espacio Afrodita','Av.Juan Domingo Peron 1522','1663','4846-7639');
	 				   insert into comercio values(2609429,'Burger King','Santa Rosa 1680','1714','6245-9214');
	 				   insert into comercio values(6532468,'Garbarino','Av. Gral. Juan Manuel de Rosas 658','1712','5634-1480');
	 				   insert into comercio values(8712384,'Fravega','Av. Rivadavia 11626','1408','46864-5674');
	 				   insert into comercio values(3466632,'YPF','Av. Ricardo Balbin 1897','1650','4753-1745');
	 				   insert into comercio values(5632132,'Maxiconsumo','Gaona Acceso Oeste 8676','1744','3214-8413');
	 				   insert into comercio values(5570712,'Farmacia Fernandez','Bartolome Mitre 800','1742','9871-3587');
	 				   insert into comercio values(7860698,'Mostaza','Arturo Jauretche 978','1969','6984-45289');
	 				   insert into comercio values(2545384,'Coppel','Belgrano 3231','1650','3216-6512');
	 				   insert into comercio values(4163701,'Libreria Rodriguez','Independencia 4647','1653','8431-3218');
	 				   insert into comercio values(0923934,'Supermercado Dia','De la tradicion 185','1713','2451-9871');
	 				   insert into comercio values(2105614,'Peluqueria paty','Gral. Lavalle 848','1714','8721-9852');
	 				   insert into comercio values(4334530,'Santeria la paz','Vidal 1769','1426','8922-3265');
	 				   insert into comercio values(2054706,'Musimundo','Av. Lope de Vega 1520','1407','2154-3285');
	 				   insert into comercio values(1287436,'Shell','Sta Rosa 2489','1712','8970-8132');
	 				   insert into comercio values(9836840,'Casa del Audio','Rivadavia 2198','1714','4015-9872');
	 				   insert into comercio values(5419987,'KFC','Av. Bartolome Mitre','1744','0454-1134')`)
    if err != nil {
    	log.Fatal(err)
    }
}


func cargarConsumos() {

	db, err := sql.Open("postgres", "user=postgres host=localhost dbname=tp sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`insert into consumo values('4455674512546534','2020',1935485,11320.00);
					  insert into consumo values('4455674512546534','2020',9836840,15840.50);	 
					  insert into consumo values('1435471512346032','1212',1935485,8900.5); 
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
				v decimal(7,2);
			begin
				perform * from tarjeta where nrotarjeta = nrotarjetaAux;
				if not found then
					insert into rechazo values(nextval('aumentoRechazo'),null,nrocomercioAux,current_timestamp,montoAux,'tarjeta inexistente');
					return false;
				else
					select * into resultado from tarjeta where nrotarjeta = nrotarjetaAux and codseguridad = codseguridadAux;
					if not found then
						insert into rechazo values(nextval('aumentoRechazo'),nrotarjetaAux,nrocomercioAux,current_timestamp,montoAux,'codigo de seguridad incorrecto');
						return false;
					else
						total := 0;
						for v in select monto from compra where compra.nrotarjeta = nrotarjetaAux and compra.pagado = true loop
							total := total + v;
						end loop;
						total := total + montoAux;
						select * into resultado from tarjeta t where t.nrotarjeta = nrotarjetaAux and t.limitecompra > total;
						if not found then
							insert into rechazo values(nextval('aumentoRechazo'),nrotarjetaAux,nrocomercioAux,current_timestamp,montoAux,'supera limite tarjeta');
							return false;
						else
							select * into resultado from tarjeta where nrotarjeta = nrotarjetaAux and to_date(validahasta, 'YYYYMM') >= to_date('202201', 'YYYYMM');
							if not found then
								insert into rechazo values(nextval('aumentoRechazo'),nrotarjetaAux,nrocomercioAux,current_timestamp,montoAux,'plazo de vigencia expirado');
								return false;
							else
								select * into resultado from tarjeta where nrotarjeta = nrotarjetaAux and estado = 'vigente';
								if not found then
									insert into rechazo values(nextval('aumentoRechazo'),nrotarjetaAux,nrocomercioAux,current_timestamp,montoAux,'la tarjeta se encuentra suspendida');
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


func probarConsumos() {

	db, err := sql.Open("postgres", "user=postgres host=localhost dbname=tp sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	
	_, err = db.Query(`create or replace function probarconsumo() returns void as $$
			declare
					v record;
					a record;
			begin
					for v in select * from consumo loop
						select autorizarCompra(v.nrotarjeta,v.codseguridad,v.nrocomercio,v.monto) into a;
					end loop;
			end;
			$$ language plpgsql;`)

	if err != nil {
		log.Fatal(err)
	}
}


func llamarConsumos() {
	db, err := sql.Open("postgres", "user=postgres host=localhost dbname=tp sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Query(`select probarconsumo();`)
	
	if err != nil {
		log.Fatal(err)
	}
}


func menuPrincipal() *wmenu.Menu {

	menu := wmenu.NewMenu("Seleccione una opcion")
	menu.Action(func (opts []wmenu.Opt) error {fmt.Print(opts[0].Text + "ha sido seleccionada"); return nil})
	menu.Option("Crear base de datos", nil, false, func(opt wmenu.Opt) error {
		fmt.Println("")
		createDatabase()
		mv := menuVolver()
		return mv.Run()
	})
	menu.Option("Eliminar base de datos", nil, false, func(opt wmenu.Opt) error {
		fmt.Println("")
		deleteDatabase()
		mv := menuVolver()
		return mv.Run()
	})		
	menu.Option("Crear tablas", nil, false, func(opt wmenu.Opt) error {
		fmt.Println("")
		crearTablas()
		mv := menuVolver()
		return mv.Run()
	})
	menu.Option("Agregar PKs y FKs", nil, false, func(opt wmenu.Opt) error {
		fmt.Println("")
		definirPksYFks()
		mv := menuVolver()
		return mv.Run()
	})
	menu.Option("Eliminar PKs y FKs", nil, false, func(opt wmenu.Opt) error {
		fmt.Println("")
		eliminarPksYFks()
		mv := menuVolver()
		return mv.Run()
	})
	menu.Option("Rellenar con valores tablas esenciales", nil, false, func(opt wmenu.Opt) error {
		fmt.Println("")
		cargarClientes()
		cargarComercio()
		cargarTarjetas()
		cargarCierres()
		cargarConsumos()
		mv := menuVolver()
		return mv.Run() 
	})
	menu.Option("Cargar funciones", nil, false, func(opt wmenu.Opt) error {
		fmt.Println("")	
		autorizarCompra()
		probarConsumos()
		llamarConsumos()
		mv := menuVolver()
		return mv.Run()
	})
	menu.Option("Salir", nil, false, func(opt wmenu.Opt) error {
		fmt.Println("")	
		os.Exit(0)
		return nil
	})
	return menu
}


func menuVolver() *wmenu.Menu {
	menu := wmenu.NewMenu("")
	menu.Action(func (opts []wmenu.Opt) error {fmt.Print(opts[0].Text + "ha sido seleccionada"); return nil})
	menu.Option("Volver al menu principal", nil, false, func(opt wmenu.Opt) error {
		fmt.Println("")	
		mp := menuPrincipal()
		return mp.Run()
	})
	menu.Option("Salir para verificar en PostgreSQL", nil, false, func(opt wmenu.Opt) error {
		fmt.Println("")	
		os.Exit(0)
		return nil
	})
	return menu
}


func main() {
	m := menuPrincipal()
	err := m.Run()
	if err != nil {
		log.Fatal(err)
		fmt.Println(err)
		os.Exit(1)
	}
}
