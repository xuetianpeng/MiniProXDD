package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/alexbrainman/odbc"
	"github.com/axgle/mahonia"
	_ "github.com/denisenkom/go-mssqldb"
	"io"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
)

func main() {
	p, _ := filepath.Abs("C://mini//images//")
	http.Handle("/", http.FileServer(http.Dir(p)))
	http.HandleFunc("/image/", GetHeadImage)
	http.HandleFunc("/pt/", GetPosttype)
	err := http.ListenAndServe(":8088", nil)
	if err != nil {
		fmt.Println(err)
	}
	//err := http.ListenAndServeTLS(":443", "crt/server.crt", "crt/server.key", nil)
	//if err != nil {
	//	fmt.Println(err)
	//}
}
func GetHeadImage(w http.ResponseWriter, req *http.Request) {
	parm := req.URL.Query()
	a := parm.Get("key")
	db, err := sql.Open("odbc", "driver={sql server};server=47.96.12.245;port=1433;uid=sa;pwd=rt-MART16;database=XDDHM")
	if err != nil {
		fmt.Printf(err.Error())
	}
	querystr := `SELECT ImageID, ImageKey, ImageText, ImagePath, ImageType FROM T_Images WHERE ImageType = ` + a
	rows, err := db.Query(querystr)
	if err != nil {
		log.Fatal("Query failed:", err.Error())
	}
	defer rows.Close()
	var rowsData []*Images
	var headImage []*HeadImagePostType
	for rows.Next() {
		var row = new(Images)
		var hi = new(HeadImagePostType)
		rows.Scan(&row.ImageID, &row.ImageKey, &row.ImageText, &row.ImagePath, &row.ImageType)
		enc := mahonia.NewDecoder("gbk")
		row.ImageText = enc.ConvertString(string(row.ImageText))
		row.ImagePath = "posttype/" + row.ImagePath
		hi.PostType = row.ImageText
		hi.ImageUrl = row.ImagePath
		rowsData = append(rowsData, row)
		headImage = append(headImage, hi)
	}
	data, err := json.Marshal(headImage)

	if err != nil {
		fmt.Println("json.marshal failed, err:", err)
		return
	}

	//fmt.Printf("%s\n", string(data))
	io.WriteString(w, string(data))
}

type HeadImagePostType struct {
	PostType string
	ImageUrl string
}

func GetPosttype(w http.ResponseWriter, req *http.Request) {
	enc := mahonia.NewDecoder("gbk")
	u, err := url.Parse(req.URL.String())
	if err != nil {
		log.Fatal(err)
	}
	parm := u.Query()
	println(parm.Encode())
	a := parm.Get("ptid")
	db, err := sql.Open("odbc", "driver={sql server};server=47.96.12.245;port=1433;uid=sa;pwd=rt-MART16;database=XDDHM")
	if err != nil {
		fmt.Printf(err.Error())
	}
	if len(parm) < 1 {
		querystr := `SELECT PT_ID ,PT_Name,MPT_ID FROM T_PostType where MPT_ID = 1`
		rows, err := db.Query(querystr)
		if err != nil {
			log.Fatal("Query failed:", err.Error())
		}
		defer rows.Close()
		var ps []*Posts
		var lee int
		var lee2 int
		for rows.Next() {
			var row = new(Posts)
			lee = lee + 1
			rows.Scan(&row.PT_ID, &row.PT_Name, &row.MPT_ID)
			row.PT_Name = enc.ConvertString(string(row.PT_Name))
			querystr2 := `SELECT PT_ID ,PT_Name,MPT_ID FROM T_PostType where MPT_ID =  ` + strconv.Itoa(row.PT_ID)
			rows2, err2 := db.Query(querystr2)
			if err2 != nil {
				log.Fatal("Query failed:", err2.Error())
			}
			defer rows2.Close()
			for rows2.Next() {
				var row2 PostType
				lee2 = lee2 + 1
				rows2.Scan(&row2.PT_ID, &row2.PT_Name, &row2.MPT_ID)
				row2.PT_Name = enc.ConvertString(string(row2.PT_Name))
				row.PostType = append(row.PostType, row2)
			}
			//println(strconv.Itoa(lee2))
			ps = append(ps, row)
		}
		//println(strconv.Itoa(lee))

		data, err := json.Marshal(ps)
		if err != nil {
			fmt.Println("json.marshal failed, err:", err)
			return
		}

		fmt.Printf("%s\n", string(data))
		io.WriteString(w, string(data))

	} else {
		if IsNum(a) {
			b, _ := strconv.Atoi(a)
			if b > 44 {
				querystr := `SELECT PT_ID ,PT_Name,MPT_ID FROM T_PostType where PT_ID  = ` + strconv.Itoa(b)
				rows, err := db.Query(querystr)
				if err != nil {
					log.Fatal("Query failed:", err.Error())
				}
				defer rows.Close()
				var row = new(PostType)
				for rows.Next() {

					rows.Scan(&row.PT_ID, &row.PT_Name, &row.MPT_ID)
					row.PT_Name = enc.ConvertString(string(row.PT_Name))
				}
				data, err := json.Marshal(row)

				if err != nil {
					fmt.Println("json.marshal failed, err:", err)
					return
				}

				fmt.Printf("%s\n", string(data))
				io.WriteString(w, string(data))
			} else {
				querystr := `SELECT PT_ID ,PT_Name,MPT_ID FROM T_PostType where PT_ID  = ` + strconv.Itoa(b)
				rows, err := db.Query(querystr)
				if err != nil {
					log.Fatal("Query failed:", err.Error())
				}
				defer rows.Close()
				var ps []*Posts
				var lee int
				var lee2 int
				for rows.Next() {
					var row = new(Posts)
					lee = lee + 1
					rows.Scan(&row.PT_ID, &row.PT_Name, &row.MPT_ID)
					row.PT_Name = enc.ConvertString(string(row.PT_Name))
					querystr2 := `SELECT PT_ID ,PT_Name,MPT_ID FROM T_PostType where MPT_ID =  ` + strconv.Itoa(row.PT_ID)
					rows2, err2 := db.Query(querystr2)
					if err2 != nil {
						log.Fatal("Query failed:", err2.Error())
					}
					defer rows2.Close()
					for rows2.Next() {
						var row2 PostType
						lee2 = lee2 + 1
						rows2.Scan(&row2.PT_ID, &row2.PT_Name, &row2.MPT_ID)
						row2.PT_Name = enc.ConvertString(string(row2.PT_Name))
						row.PostType = append(row.PostType, row2)
					}
					//println(strconv.Itoa(lee2))
					ps = append(ps, row)
				}
				//println(strconv.Itoa(lee))

				data, err := json.Marshal(ps)
				if err != nil {
					fmt.Println("json.marshal failed, err:", err)
					return
				}

				fmt.Printf("%s\n", string(data))
				io.WriteString(w, string(data))
			}
		} else {
			fmt.Println("Key值错误或参数错误")
			io.WriteString(w, "null")
		}

	}

}

type Images struct {
	ImageID   int
	ImageKey  string
	ImageText string
	ImagePath string
	ImageType int
}

type PostType struct {
	PT_ID   int
	PT_Name string
	MPT_ID  int
}
type Posts struct {
	PT_ID    int
	PT_Name  string
	MPT_ID   int
	PostType []PostType
}

func IsNum(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}
