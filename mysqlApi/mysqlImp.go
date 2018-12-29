package mysqlApi

import (
	"database/sql"
	"fmt"
	"log"
)
import _ "github.com/go-sql-driver/mysql"


type MsgSource struct {
	Src_ID     int
	Src_Firm   string
	Src_Proto  string
	Src_Access string
	Src_Prov   string
	Src_City   string
	Src_Area   string
	Src_Info   string
}


var db = &sql.DB{}



func QueryMsgSource() ([]*MsgSource,int){
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/netdb?charset=utf8")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT SRC_NO,FIRM_NAME,PROTOCOL_TYPE,ACCESS_CODE,SRC_PROV,SRC_CITY,SRC_AREA,SRC_INFO FROM MSG_SOURCE")
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	var msg []*MsgSource
	i := 0

	for rows.Next() {
		fmt.Println(i)
		var (
			srcid     int
			srcfrim   string
			srcproto  string
			srcassess string
			srcprov   string
			srccity   string
			srcarea   string
			srcinfo   string
		)

		if err := rows.Scan(&srcid,&srcfrim,&srcproto,&srcassess,&srcprov,&srccity,&srcarea,&srcinfo); err != nil {
			log.Fatal(err)
		}

		msg = append(msg,&MsgSource{Src_ID:srcid,Src_Firm:srcfrim,Src_Proto:srcproto,Src_Access:srcassess,Src_Prov:srcprov,Src_City:srccity,Src_Area:srcarea,Src_Info:srcinfo})
		//fmt.Printf("%d ,%s\n", msg.Src_ID,msg.Src_Firm)
		i = i+ 1
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return  msg,i
}


func QueryMsgType() ([]*MsgSource,int){
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/netdb?charset=utf8")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT SRC_NO,FIRM_NAME,PROTOCOL_TYPE,ACCESS_CODE,SRC_PROV,SRC_CITY,SRC_AREA,SRC_INFO FROM MSG_SOURCE")
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	var msg []*MsgSource
	i := 0

	for rows.Next() {
		fmt.Println(i)
		var (
			srcid     int
			srcfrim   string
			srcproto  string
			srcassess string
			srcprov   string
			srccity   string
			srcarea   string
			srcinfo   string
		)

		if err := rows.Scan(&srcid,&srcfrim,&srcproto,&srcassess,&srcprov,&srccity,&srcarea,&srcinfo); err != nil {
			log.Fatal(err)
		}

		msg = append(msg,&MsgSource{Src_ID:srcid,Src_Firm:srcfrim,Src_Proto:srcproto,Src_Access:srcassess,Src_Prov:srcprov,Src_City:srccity,Src_Area:srcarea,Src_Info:srcinfo})
		//fmt.Printf("%d ,%s\n", msg.Src_ID,msg.Src_Firm)
		i = i+ 1
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return  msg,i
}