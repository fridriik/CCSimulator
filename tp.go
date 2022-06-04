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

	 				   alter table tarjeta add constraint tarjeta_nrocliente_fk foreign key (nrocliente) references cliente (nrocliente);
	 				   alter table compra add constraint compra_nrotarjeta_fk foreign key (nrotarjeta) references tarjeta(nrotarjeta);
	 				   alter table compra add constraint compra_nrocomerio_fk foreign key (nrocomercio) references comercio(nrocomercio);
	 				   alter table rechazo add constraint rechazo_nrotarjeta_fk foreign key (nrotarjeta) references tarjeta(nrotarjeta);
	 				   alter table rechazo add constraint rechazo_nrocomerio_fk foreign key (nrocomercio) references comercio(nrocomercio);
	 				   alter table cabecera add constraint cabecera_nrotarjeta_fk foreign key (nrotarjeta) references tarjeta (nrotarjeta);
	 				   alter table alerta add constraint alerta_nrotarjeta_fk foreign key (nrotarjeta) references tarjeta (nrotarjeta);
	 				   alter table alerta add constraint alerta_nrorechazoa_fk foreign key (nrorechazo) references rechazo (nrorechazo)`)
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
	 				   insert into cliente values(2161054,'Apu','Nahasapeemape','Av. del Libertador 2373','4361-4419');
	 				   insert into cliente values(3487910,'Martin','Price','Av. Infanta Isabel 555','4433-3396');
	 				   insert into cliente values(0216546,'Selma','Bouvie','Av. Pedro de Mendoza 1843','4343-2123');
	 				   insert into cliente values(2105646,'Patty','Bouvie','Av. España 1701','4370-6105');
	 				   insert into cliente values(8845660,'Barney','Gumble','Pujol 644','4893-0322');
	 				   insert into cliente values(8520147,'Waylon','Smithers','Defensa 219','4362-1100');`)
    if err !=nil {
    	log.Fatal(err)
    }
}

func cargarComercio(db1 *sql.DB){
	db := db1
	var err error

	 _, err = db.Exec(`
	 insert into comercio values(1935485,'McDonalds','Concejal Tribulato 636','1744','6541-6542');
	 insert into comercio values(3455465,'Peluqueria Adomo','Av. Victorica 421','1614','6785-1354');
	 insert into comercio values(9869845,'Supermercado Li','Crisóstomo Álvarez 2825','1406,'2251-5268');
	 insert into comercio values(2314575,'Santeria Espacio Afrodita','Av.Juan Domingo Peron 1522','1663','4846-7639');
	 insert into comercio values(2609429,'Burger King','Santa Rosa 1680','1714','6245-9214');
	 insert into comercio values(6532468,'Garbarino','Av. Gral. Juan Manuel de Rosas 658','1712','5634-1480');
	 insert into comercio values(8712384,'Fravega','Av. Rivadavia 11626','1408','46864-5674');
	 insert into comercio values(3466632,'YPF','Av. Ricardo Balbin 1897','1650','4753-1745');
	 insert into comercio values(5632132,'Maxiconsumo','Gaona Acceso Oeste 8676','1744','3214-8413');
	 insert into comercio values(5570712,'Farmacia Fernandez','Bartolome Mitre 800','1742','9871-3587');
	 
	 insert into comercio values(7860698,'McDonalds','Arturo Jauretche 978','1969','6984-45289');
	 insert into comercio values(2545384,'Coppel','Belgrano 3231','1650','3216-6512');
	 insert into comercio values(4163701,'Libreria Rodriguez','Independencia 4647','1653','8431-3218');
	 insert into comercio values(0923934,'Supermercado Dia','De la tradicion 185','1713','2451-9871');
	 insert into comercio values(2105614,'Peluqueria paty','Gral. Lavalle 848','1714','8721-9852');
	 insert into comercio values(4334530,'Santeria la paz','Vidal 1769','1426','8922-3265');
	 insert into comercio values(2054706,'Fravega','Av. Lope de Vega 1520','1407','2154-3285');
	 insert into comercio values(1287436,'YPF','Sta Rosa 2489','1712','8970-8132');
	 insert into comercio values(9836840,'Cada del Audio','Rivadavia 2198','1714','4015-9872');
	 insert into comercio values(5419987,'Burger King','Av. Bartolome Mitre','1744','0454-1134');
	 `)
	  
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
}
