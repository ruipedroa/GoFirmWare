package main

import (
	"fmt"
    "io/ioutil"
    "strconv"
    "strings"
	"time"
    //"log"
	"github.com/shirou/gopsutil/mem"    
	//"github.com/boltdb/bolt"
    "database/sql"
_ "github.com/mattn/go-sqlite3"
)

// Function to extract info from CPU and get total ticks and Idle time ticks
func getCPUSample() (idle, total uint64) {
    contents, err:= ioutil.ReadFile("/proc/stat")
    if err != nil {
        return
    }
    lines := strings.Split(string(contents), "\n")
    for _,line := range(lines) {
        fields:= strings.Fields(line)
        if fields[0] == "cpu"  {
            numFields := len(fields)
            for i:= 1; i< numFields; i++ {
                val, err :=strconv.ParseUint(fields[i], 10, 64)
                if err != nil {
                    fmt.Println("Error: ", i, fields[i], err)                
                }
                total +=val //sum all numbers to get total ticks
                if i == 4 {
                    idle = val //idle is the 5th field in cpu           
                }
            }
            return
        }    
    }
    return    
}

func main() {
    
    database, _ := sql.Open("sqlite3", "./local.db")
    //Create table for CPU and RAM values from OS
    statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS SystemData (id INTEGER PRIMARY KEY, CPU DECIMAL(5, 2), RAM DECIMAL(5, 2))")
    statement.Exec()
    //Create table for Data obtained from external device
    statement1, _ := database.Prepare("CREATE TABLE IF NOT EXISTS DeviceData (id INTEGER PRIMARY KEY, temp INTEGER, humidity INTEGER, voltage INTEGER, current INTEGER)")
    statement1.Exec()
    statement, _ = database.Prepare("INSERT INTO SystemData (CPU , RAM) VALUES (?, ?)")
    statement.Exec(14.5,25.65)

    //Insert Random Values in Device Data
    statement, _ = database.Prepare("INSERT INTO DeviceData (temp , humidity, voltage, current) VALUES (?, ?, ?, ?)")
    statement.Exec(21, 9 , 5, 1)
    rows, _ :=database.Query("SELECT id, CPU, RAM FROM SystemData")
    rows2, _ :=database.Query("SELECT id, temp, humidity, voltage, current FROM DeviceData")

    var id, temp, humidity, voltage, current int
    var CPU float64
    var RAM float64
    
    for rows.Next()  {
        rows.Scan(&id, &CPU, &RAM)
        fmt.Println(strconv.Itoa(id) + ":" + strconv.FormatFloat(CPU, 'f', -1, 64) +" "+ strconv.FormatFloat(RAM, 'f', -1, 64))
    }

    for rows2.Next()  {
        rows2.Scan(&id, &temp, &humidity, &voltage, &current)
        fmt.Println(strconv.Itoa(id) + ":" + strconv.Itoa(temp)  +" "+ strconv.Itoa(humidity) + " " + strconv.Itoa(voltage) + " " + strconv.Itoa(current))
    }
    
    
    idle0,total0 :=getCPUSample()
    // Cycle to get reads every second
    for {
        time.Sleep(time.Second)  
        idle1, total1:= getCPUSample()

        idleTicks :=float64(idle1 -idle0)
        totalTicks := float64(total1 -total0)
        cpuUsage:=100* (totalTicks -idleTicks) /totalTicks
	
        fmt.Printf("CPU Usage is %f%% [busy :%f, total: %f]\n", cpuUsage, totalTicks-idleTicks, totalTicks)
    
        //get Memory(RAM) percentage
	    v, _ := mem.VirtualMemory()

	    fmt.Printf("Pecentage of RAM used:%f%%\n",v.UsedPercent)
        idle0 = idle1
        total0 = total1
        }
}
