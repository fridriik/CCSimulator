package main

import (
	"database/sql"
//	"fmt"
	_"github.com/lib/pq"
	"log"
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
	if err !=nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err =db.Exec(`create database tp`)
	if err != nil {
		log.Fatal(err)
	}
}

		
func deleteDatabase(){
	db, err := sql.Open("postgres", "user=postgres host=localhost dbname=postgres sslmode=disable")
	if err !=nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err =db.Exec(`drop database if exists tp`)
	if err != nil {
		log.Fatal(err)
	}
}


func definirPksYFks(db1 *sql.DB){

	db := db1
	var err error
	
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
    if err !=nil {
    	log.Fatal(err)
    }
}


func eliminarPksYFks(db1 *sql.DB){

	db := db1
	var err error
	
	 _, err = db.Exec(`alter table cliente drop constraint cliente_pk;
	 				   alter table tarjeta drop constraint tarjeta_pk;
	 				   alter table comercio drop constraint comercio_pk;
	 				   alter table compra drop constraint compra_pk;
	 				   alter table rechazo drop constraint rechazo_pk;
	 				   alter table cierre drop constraint cierre_pk;
	 				   alter table cabecera drop constraint cabecera_pk;
	 				   alter table detalle drop constraint detalle_pk;
	 				   alter table alerta drop constraint alerta_pk;

	 				   alter table tarjeta drop constraint tarjeta_nrocliente_fk;
	 				   alter table compra drop constraint compra_nrotarjeta_fk;
	 				   alter table compra drop constraint compra_nrocomerio_fk;
	 				   alter table rechazo drop constraint rechazo_nrotarjeta_fk;
	 				   alter table rechazo drop constraint rechazo_nrocomerio_fk;
	 				   alter table cabecera drop constraint cabecera_nrotarjeta_fk;
	 				   alter table alerta drop constraint alerta_nrotarjeta_fk;
	 				   alter table alerta drop constraint alerta_nrorechazoa_fk`)
    if err !=nil {
    	log.Fatal(err)
    }
}


func cargarClientes(db1 *sql.DB){
	db := db1
	var err error

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
    if err !=nil {
    	log.Fatal(err)
    }
}


func cargarTarjetas(db1 *sql.DB){
	db := db1
	var err error

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
    if err !=nil {
    	log.Fatal(err)
    }
}

func cargarCierres(db1 *sql.DB){

	db := db1
	var err error

	 _, err = db.Exec(`insert into cierre values(2022,01,0,'16-01-2022','15-01-2022','21-01-2022');
	 				   insert into cierre values(2022,02,0,'19-02-2022','18-02-2022','25-02-2022');
	 				   insert into cierre values(2022,03,0,'17-03-2022','16-03-2022','22-03-2022');
	 				   insert into cierre values(2022,04,0,'18-04-2022','17-04-2022','24-04-2022');
	 	    		   insert into cierre values(2022,05,0,'16-05-2022','15-05-2022','21-05-2022');
	 				   insert into cierre values(2022,06,0,'18-06-2022','17-06-2022','22-06-2022');
	 				   insert into cierre values(2022,07,0,'20-07-2022','19-07-2022','26-07-2022');
	 				   insert into cierre values(2022,08,0,'16-08-2022','15-08-2022','21-08-2022');
	 				   insert into cierre values(2022,09,0,'19-09-2022','18-09-2022','25-09-2022');
	 				   insert into cierre values(2022,10,0,'18-10-2022','17-10-2022','24-10-2022');
	 				   insert into cierre values(2022,11,0,'18-11-2022','17-11-2022','24-11-2022');
	 				   insert into cierre values(2022,12,0,'17-12-2022','16-12-2022','22-12-2022'); 

                       insert into cierre values(2022,01,1,'17-01-2022','16-01-2022','22-01-2022');
	 				   insert into cierre values(2022,02,1,'19-02-2022','18-02-2022','25-02-2022');
	 				   insert into cierre values(2022,03,1,'20-03-2022','19-03-2022','26-03-2022');
	 				   insert into cierre values(2022,04,1,'18-04-2022','17-04-2022','24-04-2022');
	 				   insert into cierre values(2022,05,1,'16-05-2022','15-05-2022','21-05-2022');
	 				   insert into cierre values(2022,06,1,'16-06-2022','15-06-2022','21-06-2022');
	 				   insert into cierre values(2022,07,1,'17-07-2022','16-07-2022','22-07-2022');
	 				   insert into cierre values(2022,08,1,'18-08-2022','17-08-2022','23-08-2022');
	 				   insert into cierre values(2022,09,1,'20-09-2022','19-09-2022','26-09-2022');
	 				   insert into cierre values(2022,10,1,'19-10-2022','18-10-2022','25-10-2022');
	 				   insert into cierre values(2022,11,1,'19-11-2022','18-11-2022','25-11-2022');
	 				   insert into cierre values(2022,12,1,'20-12-2022','19-11-2022','26-12-2022');

	 				   insert into cierre values(2022,01,2,'20-01-2022','19-01-2022','26-01-2022');
	 				   insert into cierre values(2022,02,2,'16-02-2022','15-02-2022','21-02-2022');
	 				   insert into cierre values(2022,03,2,'18-03-2022','17-03-2022','23-03-2022');
	 				   insert into cierre values(2022,04,2,'16-04-2022','15-04-2022','21-04-2022');
	 				   insert into cierre values(2022,05,2,'17-05-2022','16-05-2022','22-05-2022');
	 				   insert into cierre values(2022,06,2,'17-06-2022','16-06-2022','22-06-2022');
	 				   insert into cierre values(2022,07,2,'19-07-2022','18-07-2022','25-07-2022');
	 				   insert into cierre values(2022,08,2,'18-08-2022','17-08-2022','24-08-2022');
	 				   insert into cierre values(2022,09,2,'20-09-2022','19-09-2022','26-09-2022');
	 				   insert into cierre values(2022,10,2,'18-10-2022','17-10-2022','24-10-2022');
	 				   insert into cierre values(2022,11,2,'18-11-2022','17-11-2022','24-11-2022');
	 				   insert into cierre values(2022,12,2,'19-12-2022','18-12-2022','25-01-2022');

	 				   insert into cierre values(2022,01,3,'16-01-2022','15-01-2022','21-01-2022');
	 				   insert into cierre values(2022,02,3,'18-02-2022','17-02-2022','24-02-2022');
	 				   insert into cierre values(2022,03,3,'18-03-2022','17-03-2022','24-03-2022');
	 				   insert into cierre values(2022,04,3,'19-04-2022','18-04-2022','25-04-2022');
	 				   insert into cierre values(2022,05,3,'20-05-2022','19-05-2022','26-05-2022');
	 				   insert into cierre values(2022,06,3,'21-06-2022','20-06-2022','27-06-2022');
	 				   insert into cierre values(2022,07,3,'20-07-2022','19-07-2022','26-07-2022');
	 				   insert into cierre values(2022,08,3,'18-08-2022','17-08-2022','24-08-2022');
	 				   insert into cierre values(2022,09,3,'17-09-2022','16-09-2022','22-09-2022');
	 				   insert into cierre values(2022,10,3,'17-10-2022','16-10-2022','22-10-2022');
	 				   insert into cierre values(2022,11,3,'17-11-2022','15-11-2022','22-11-2022');
	 				   insert into cierre values(2022,12,3,'16-12-2022','15-12-2022','21-12-2022');

	 				   insert into cierre values(2022,01,4,'19-01-2022','18-01-2022','25-01-2022');
	 				   insert into cierre values(2022,02,4,'18-02-2022','17-02-2022','24-02-2022');
	 				   insert into cierre values(2022,03,4,'18-03-2022','17-03-2022','24-03-2022');
	 				   insert into cierre values(2022,04,4,'16-04-2022','15-04-2022','21-04-2022');
	 				   insert into cierre values(2022,05,4,'17-05-2022','16-05-2022','22-05-2022');
	 				   insert into cierre values(2022,06,4,'18-06-2022','17-06-2022','24-06-2022');
	 				   insert into cierre values(2022,07,4,'20-07-2022','19-07-2022','26-07-2022');
	 				   insert into cierre values(2022,08,4,'20-08-2022','19-08-2022','26-08-2022');
	 				   insert into cierre values(2022,09,4,'19-09-2022','18-09-2022','25-09-2022');
	 				   insert into cierre values(2022,10,4,'17-10-2022','16-10-2022','22-10-2022');
	 				   insert into cierre values(2022,11,4,'18-11-2022','17-11-2022','24-11-2022');
	 				   insert into cierre values(2022,12,4,'19-12-2022','18-12-2022','25-01-2022');

	 				   insert into cierre values(2022,01,5,'20-01-2022','19-01-2022','26-01-2022');
	 				   insert into cierre values(2022,02,5,'20-02-2022','19-02-2022','26-02-2022');
	 				   insert into cierre values(2022,03,5,'18-03-2022','17-03-2022','24-03-2022');
	 				   insert into cierre values(2022,04,5,'19-04-2022','18-04-2022','25-04-2022');
	 				   insert into cierre values(2022,05,5,'21-05-2022','20-05-2022','27-05-2022');
	 				   insert into cierre values(2022,06,5,'21-06-2022','20-06-2022','27-06-2022');
	 				   insert into cierre values(2022,07,5,'19-07-2022','18-07-2022','25-07-2022');
	 				   insert into cierre values(2022,08,5,'18-08-2022','17-08-2022','24-08-2022');
	 				   insert into cierre values(2022,09,5,'17-09-2022','16-09-2022','22-09-2022');
	 				   insert into cierre values(2022,10,5,'17-10-2022','16-10-2022','22-10-2022');
	 				   insert into cierre values(2022,11,5,'18-11-2022','17-11-2022','24-11-2022');
	 				   insert into cierre values(2022,12,5,'18-12-2022','18-12-2022','24-12-2022');

	 				   insert into cierre values(2022,01,6,'17-01-2022','16-01-2022','22-01-2022');
	 				   insert into cierre values(2022,02,6,'16-02-2022','15-02-2022','21-02-2022');
	 				   insert into cierre values(2022,03,6,'18-03-2022','17-03-2022','24-03-2022');
	 				   insert into cierre values(2022,04,6,'16-04-2022','15-04-2022','21-04-2022');
	 				   insert into cierre values(2022,05,6,'17-05-2022','16-05-2022','22-05-2022');
	 				   insert into cierre values(2022,06,6,'18-06-2022','17-06-2022','24-06-2022');
	 				   insert into cierre values(2022,07,6,'19-07-2022','18-07-2022','25-07-2022');
	 				   insert into cierre values(2022,08,6,'16-08-2022','15-08-2022','21-08-2022');
	 				   insert into cierre values(2022,09,6,'16-09-2022','15-09-2022','21-09-2022');
	 				   insert into cierre values(2022,10,6,'17-10-2022','16-10-2022','22-10-2022');
	 				   insert into cierre values(2022,11,6,'19-11-2022','18-11-2022','25-11-2022');
	 				   insert into cierre values(2022,12,6,'20-12-2022','19-12-2022','26-12-2022');

	 				   insert into cierre values(2022,01,7,'20-01-2022','19-01-2022','26-01-2022');
	 				   insert into cierre values(2022,02,7,'18-02-2022','17-02-2022','24-02-2022');
	 				   insert into cierre values(2022,03,7,'18-03-2022','17-03-2022','24-03-2022');
	 				   insert into cierre values(2022,04,7,'19-04-2022','18-04-2022','25-04-2022');
	 				   insert into cierre values(2022,05,7,'17-05-2022','16-05-2022','22-05-2022');
	 				   insert into cierre values(2022,06,7,'21-06-2022','20-06-2022','27-06-2022');
	 				   insert into cierre values(2022,07,7,'20-07-2022','19-07-2022','26-07-2022');
	 				   insert into cierre values(2022,08,7,'18-08-2022','17-08-2022','24-08-2022');
	 				   insert into cierre values(2022,09,7,'19-09-2022','18-09-2022','25-09-2022');
	 				   insert into cierre values(2022,10,7,'19-10-2022','18-10-2022','25-10-2022');
	 				   insert into cierre values(2022,11,7,'21-11-2022','20-11-2022','27-11-2022');
	 				   insert into cierre values(2022,12,7,'19-12-2022','18-12-2022','25-12-2022');
	 				   
	 				   insert into cierre values(2022,01,8,'16-01-2022','15-01-2022','21-01-2022');
	 				   insert into cierre values(2022,02,8,'17-02-2022','16-02-2022','22-02-2022');
	 				   insert into cierre values(2022,03,8,'17-03-2022','16-03-2022','22-03-2022');
	 				   insert into cierre values(2022,04,8,'19-04-2022','18-04-2022','25-04-2022');
	 				   insert into cierre values(2022,05,8,'17-05-2022','16-05-2022','22-05-2022');
	 				   insert into cierre values(2022,06,8,'18-06-2022','17-06-2022','24-06-2022');
	 				   insert into cierre values(2022,07,8,'19-07-2022','18-07-2022','25-07-2022');
	 				   insert into cierre values(2022,08,8,'18-08-2022','17-08-2022','23-08-2022');
	 				   insert into cierre values(2022,09,8,'17-09-2022','16-09-2022','22-09-2022');
	 				   insert into cierre values(2022,10,8,'18-10-2022','17-10-2022','24-10-2022');
	 				   insert into cierre values(2022,11,8,'18-11-2022','17-11-2022','24-11-2022');
	 				   insert into cierre values(2022,12,8,'20-12-2022','19-12-2022','26-12-2022');

	 				   insert into cierre values(2022,01,9,'16-01-2022','15-01-2022','21-01-2022');
	 				   insert into cierre values(2022,02,9,'16-02-2022','15-02-2022','21-02-2022');
	 				   insert into cierre values(2022,03,9,'18-03-2022','17-03-2022','24-03-2022');
	 				   insert into cierre values(2022,04,9,'19-04-2022','18-04-2022','25-04-2022');
	 				   insert into cierre values(2022,05,9,'17-05-2022','16-05-2022','22-05-2022');
	 				   insert into cierre values(2022,06,9,'21-06-2022','20-06-2022','27-06-2022');
	 				   insert into cierre values(2022,07,9,'20-07-2022','19-07-2022','26-07-2022');
	 				   insert into cierre values(2022,08,9,'18-08-2022','17-08-2022','23-08-2022');
	 				   insert into cierre values(2022,09,9,'16-09-2022','15-09-2022','21-09-2022');
	 				   insert into cierre values(2022,10,9,'18-10-2022','17-10-2022','24-10-2022');
	 				   insert into cierre values(2022,11,9,'19-11-2022','18-11-2022','25-11-2022');
	 				   insert into cierre values(2022,12,9,'19-12-2022','18-12-2022','25-12-2022');`)
    if err !=nil {
    	log.Fatal(err)
    }
}

func cargarComercio(db1 *sql.DB){
	db := db1
	var err error

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
    if err !=nil {
    	log.Fatal(err)
    }
}


func main(){

    deleteDatabase()
	createDatabase()

	db, err := sql.Open("postgres", "user=postgres host=localhost dbname=tp sslmode=disable")
	if err !=nil {
		log.Fatal(err)
	}
	defer db.Close()

	 _, err = db.Exec(`create table cliente(nrocliente int,
    										nombre text,
    										apellido text,
    										domicilio text,
    										telefono char(12))`)
    if err !=nil {
    	log.Fatal(err)
    }

    _, err = db.Exec(`create table tarjeta(nrotarjeta char(16),
    									   nrocliente int,
    									   validadesde char(6), 
    									   validahasta char(6),
    									   codseguridad char(4),
    									   limitecompra decimal(8,2),
    									   estado char(10))`)
    if err !=nil {
    	log.Fatal(err)
    }
    
    _, err = db.Exec(`create table comercio(nrocomercio int,
    										nombre text,
    										domicilio text,
    										codigopostal char(8),
    										telefono char(12))`)
     if err !=nil {
    	log.Fatal(err)
    }
    
    _, err = db.Exec(`create table compra(nrooperacion int,
    									  nrotarjeta char(16),
    									  nrocomercio int,
    									  fecha timestamp,
    									  monto decimal(7,2),
    									  pagado boolean)`)
    if err !=nil {
    	log.Fatal(err)
    }
    
    _, err = db.Exec(`create table rechazo(nrorechazo int,
    									   nrotarjeta char(16),
    									   nrocomercio int,
    									   fecha timestamp,
    									   monto decimal(7,2),
    									   motivo text)`)
    if err !=nil {
    	log.Fatal(err)
    }

    _, err = db.Exec(`create table cierre(año int,
    									  mes int,
    									  terminacion int,
    									  fechainicio date,
    									  fechacierre date,
    									  fechavto date)`)
    if err !=nil {
    	log.Fatal(err)
    }
    
    _, err = db.Exec(`create table cabecera(nroresumen int,
    										nombre text,
    										apellido text,
    										domicilio text,
    										nrotarjeta char(16),
    										desde date,
    										hasta date,
    										vence date,
    										total decimal(8,2))`)   
    if err !=nil {
    	log.Fatal(err)
    }

    _, err = db.Exec(`create table detalle(nroresumen int,
    									   nrolinea int,
    									   fecha date,
    									   nombrecomercio text,
    									   monto decimal(7,2))`)     
    if err !=nil {
    	log.Fatal(err)
    }    

    _, err = db.Exec(`create table alerta(nroalerta int,
    									  nrotarjeta char(16),
    									  fecha timestamp,
    									  nrorechazo int,
    									  codalerta int,
    									  descripcion text)`)
    if err !=nil {
    	log.Fatal(err)
    } 

    _, err = db.Exec(`create table consumo(nrotarjeta char(16),
    									   codseguridad char(4),
    									   nrocomercio int,
    									   monto decimal(7,2))`)
    if err !=nil {
    	log.Fatal(err)
    } 

    definirPksYFks(db)
    cargarClientes(db)
	cargarComercio(db)
	cargarTarjetas(db)
	cargarCierres(db)
}
