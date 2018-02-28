
/*
 * Created by goland
 * User: Hasan UCAK <hasan.ucak@gmail.com>
 * Date: 2/21/2018
 * Time: 3:00 PM

functional test for proxySQL with percona mysql cluster 3.nodes

CREATE FUNCTION getServerName () returns varchar(255)
DETERMINISTIC
READS SQL DATA
begin
declare fieldresult varchar(255);
select substr(VARIABLE_VALUE,1,250) into fieldresult from performance_schema.session_variables where VARIABLE_NAME='wsrep_node_name';
return fieldresult;
end;


*/

package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"log"
)


//set database username
const username string = "kepcore"//"keprehber"
//set database password of username
const password string = "K3pC0r3$"//"k3ppr3hhb3Rr#"
//set proxySQL host ip and port number
const proxySQL_url string = "10.145.172.20:3306"
//set database name
const dbName string = "kepcore"
//total test conn number
const maxTry int = 100

var DBCon *sql.DB
var infolist map[string]int
var err error

func getNodeInfoSelect(){
	var out string
	//which connected node name
	err = DBCon.QueryRow("select VARIABLE_VALUE from performance_schema.session_variables where VARIABLE_NAME='wsrep_node_name'").Scan(&out)
	//Row := DBCon.QueryRow("show variables like 'wsrep_node_name'").Scan()
	if err != nil{
		log.Printf("Error : %s \n",err)
	}else {
		log.Printf("Node : %s \n",out)
		//inc node counter
		infolist[out]++
	}
}

func getNodeInfoInsert(){
	//which connected node name
	stmtIns, err := DBCon.Prepare("insert into test.ServerName(ServerName) values(test.getServerName())") // ? = placeholder
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates
	_, err = stmtIns.Exec() // Insert tuples (i, i^2)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
}

func main() {
	infolist = make(map[string]int)
	for index := 0; index<maxTry ; index++ {
		db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, proxySQL_url, dbName))
		DBCon = db
		if err != nil {
			log.Println(err)
		} else {
			log.Println("OK")
			getNodeInfoInsert()
		}
		defer db.Close()
	}
	// result
	fmt.Scanln()
	//getQuery(DBCon,"select serverName,count(*) as Total  from test.ServerName GROUP by 1")
	fmt.Println(infolist)
}
