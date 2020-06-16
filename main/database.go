package database

import (
	"fmt"
    "io/ioutil"
    "strconv"
    "strings"
	"time"
    "database/sql"
_ "github.com/mattn/go-sqlite3"
)

var database* sql.DB

//Function to Create Database
func CreateDatabase() (database* sql.DB ) {
    
    database, err := sql.Open("sqlite3", "./local.db")
    if err != nil {
        fmt.Println("Error:" , err)
        return
    }
    //Create table for CPU and RAM values from OS
    statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS SystemData (id INTEGER PRIMARY KEY, CPU NUMERIC(5, 2), RAM NUMERIC(5, 2))")
    if err != nil {
        fmt.Println("Error:" , err)
        return
        }
    statement.Exec()
    statement1, err := database.Prepare("CREATE TABLE IF NOT EXISTS DeviceData (id INTEGER PRIMARY KEY, temp INTEGER, humidity INTEGER, voltage INTEGER, current INTEGER)")
    if err != nil {
        fmt.Println("Error:" , err)
        return
        }
    statement1.Exec()
    return

}

//Store Values in the Database
//Separate values by table
func StoreValues(database* sql.DB, table string, ival[4] int, fval[2] float64) () {
    
    if table == "SystemData" {
         statement, err := database.Prepare("INSERT INTO SystemData (CPU, RAM) VALUES (?, ?, ?, ?)")
         if err != nil {
            fmt.Println("Error:" , err)
            return           
            }
         statement.Exec(fval[0], fval[1])
        }
    if table == "DeviceData" {
        statement, err := database.Prepare("INSERT INTO DeviceData (temp , humidity, voltage, current) VALUES (?, ?, ?, ?)")
         if err != nil {
            fmt.Println("Error:" , err)
            return           
            }
        statement.Exec(ival[0], ival[1], ival[2], ival[3])
        }
    return
}
