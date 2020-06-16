package main

import (
    "math"
    "math/rand"
	"fmt"
    "bufio"
    "os"
    "io/ioutil"
    "strconv"
    "strings"
	"time"
    "log"
	"github.com/shirou/gopsutil/mem"    
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

//Simulate external Device Value Generation
func GenerateVars() (temp, humidity, voltage, current int) {
    
        max := 50 
        min := 1
        rand.Seed(time.Now().UnixNano())
        temp = rand.Intn(rand.Intn(max - min +1) + min)
        humidity = rand.Intn(rand.Intn(max - min +1) + min)
        voltage = rand.Intn(rand.Intn(max - min +1) + min)
        current = rand.Intn(rand.Intn(max - min +1) + min)
        return
} 



//Store Values in the Database
func StoreValues(db *sql.DB, ival[4] int, fval[2] float64) () {
    
        statement, err := db.Prepare("INSERT INTO SystemData (CPU, RAM) VALUES (?, ?)")
        if err != nil {
            fmt.Println("Error:" , err)
            return           
            }
        statement.Exec(fval[0], fval[1])
        
        statement, err = db.Prepare("INSERT INTO DeviceData (temp , humidity, voltage, current) VALUES (?, ?, ?, ?)")
         if err != nil {
            fmt.Println("Error:" , err)
            return           
            }
        statement.Exec(ival[0], ival[1], ival[2], ival[3])
        
    return
}

//Retrieves Average vaue from each column in the Database

func getValues(db *sql.DB, variable string, NRows string, table int) () {
        
        var FloatValue float64
        var IntValue int
        text := " " 
        //Separate code by the different tables . 0= SystemData , 1=DeviceData
        if table == 0 {
            if variable == "1" { text = "CPU"
            }else if variable == "2" {  text = "RAM"
            }else{
                    fmt.Println("Variable doesn't exist.")                
                    return}
            rows, err :=db.Query("SELECT " + text +" FROM SystemData WHERE id > (SELECT MAX(id)  - " + NRows + "FROM SystemData)")
            if err != nil {
                    log.Fatalf("Error : %v", err)
                    return
                   }
            fmt.Println(text + " : ")            
            for rows.Next(){
                rows.Scan(&FloatValue)
                fmt.Println(strconv.FormatFloat(FloatValue, 'f', -1, 64))
                }
            rows.Close() 
            }
        
        if table == 1 {
            if variable == "3" { text = "temp" 
            }else if variable == "4" { text = "humidity"
            }else if variable == "5" { text = "voltage" 
            }else if variable == "6" { text = "current"
            }else{
                    fmt.Println("Variable doesn't exist.")                
                    return}            
            rows, err :=db.Query("SELECT " + text +" FROM DeviceData WHERE id > (SELECT MAX(id)  - " + NRows + "FROM DeviceData)")
            if err != nil {
                    log.Fatalf("Error : %v", err)
                    return
                   }
            fmt.Println(text + " : ")            
            for rows.Next()  {
                rows.Scan(&IntValue)
                fmt.Println(strconv.Itoa(IntValue))
                } 
            rows.Close()            
            }        
        return        
}

func averageValue (db *sql.DB, variable string, table int) () {
        var average float64
        text := " " 
        //Separate code by the different tables . 0= SystemData , 1=DeviceData
        if table == 0 {
            if variable == "1" { text = "CPU"
            }else if variable == "2" {  text = "RAM"
            }else{
                    fmt.Println("Variable doesn't exist.")                
                    return}
            rows, err :=db.Query("SELECT AVG(" + text + " ) FROM SystemData")
            if err != nil {
                    log.Fatalf("Error : %v", err)
                    return
                   }
            fmt.Println(text + " average : ")
            rows.Next()            
            rows.Scan(&average)
            fmt.Println(strconv.FormatFloat(average, 'f', -1, 64))
            rows.Close()    
            }

        if table == 1 {
            if variable == "3" { text = "temp" 
            }else if variable == "4" { text = "humidity"
            }else if variable == "5" { text = "voltage" 
            }else if variable == "6" { text = "current"
            }else{
                    fmt.Println("Variable doesn't exist.")                
                    return}            
            rows, err :=db.Query("SELECT AVG( " + text + " ) FROM DeviceData")
            if err != nil {
                    log.Fatalf("Error : %v", err)
                    return
                   }
            fmt.Println(text + " average : ")
            rows.Next()            
            rows.Scan(&average)
            fmt.Println(strconv.FormatFloat(average, 'f', -1, 64))
            rows.Close()                 
            }        
        return
    
}

func main() {
    
    reader := bufio.NewReader(os.Stdin)
    fmt.Println("Simple Shell")
    fmt.Println("-------------------")
    fmt.Println("Available Commands: ")
    fmt.Println("1. Get last n metrics for all variables;")
    fmt.Println("2. Get last n metrics for one or more variables;")
    fmt.Println("3. Get average value of one or more variables.")
    fmt.Println("4. Close application.")
    fmt.Println("Press 0 to see this menu again.")
    fmt.Printf("-> ")
    
    // Open DB and Create tables if they don't exist
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

    
    var CPU float64
    var RAM float64
    
    // Cycle to get reads every second
    idle0,total0 :=getCPUSample()
    
    for {
        text, _ := reader.ReadString('\n')
        text = strings.Replace(text, "\n", "", -1)
        time.Sleep(time.Second)  
        idle1, total1:= getCPUSample()

        idleTicks :=float64(idle1 -idle0)
        totalTicks := float64(total1 -total0)
        cpuUsage:=100* (totalTicks -idleTicks) /totalTicks

          
        //get Memory(RAM) percentage
	    v, _ := mem.VirtualMemory()
        
        //Get data from External Device
        temp, humidity, voltage, current := GenerateVars()
      
        StoreValues(database, [4]int{temp,humidity,voltage,current}, [2]float64{math.Round(cpuUsage*100)/100, math.Round(v.UsedPercent*100)/100})
        //Enter user interface state machine 
        if text=="1" {
                  fmt.Printf("Number of metrics to retrieve: ")
                  NRows, _ := reader.ReadString('\n')
                     
                  rowsSystem, _ :=database.Query("SELECT CPU, RAM FROM SystemData WHERE id > (SELECT MAX(id)  - " + NRows + "FROM SystemData)")
                  fmt.Println("CPU usage(%) | RAM Usage(%)")
                  for rowsSystem.Next()  {
                      rowsSystem.Scan(&CPU, &RAM)
                      fmt.Println(strconv.FormatFloat(CPU, 'f', -1, 64) +"          "+ strconv.FormatFloat(RAM, 'f', -1, 64))
                  }
                  
                  rowsDevice, _ :=database.Query("SELECT temp, humidity, voltage, current FROM DeviceData WHERE id > (SELECT MAX(id)  - " + NRows + "FROM DeviceData)")
                  fmt.Println("Temperature(ºC) | Humidity(%) | Voltage(V) | Current(mA)")
                  for rowsDevice.Next()  {
                     rowsDevice.Scan(&temp, &humidity, &voltage, &current)
                     fmt.Println(strconv.Itoa(temp)  +"                 "+ strconv.Itoa(humidity) + "             " + strconv.Itoa(voltage) + "            " + strconv.Itoa(current))
                     }
                  rowsSystem.Close()
                  rowsDevice.Close()
                  fmt.Printf("-> ")
             }     
         if text == "2"{ 
                     fmt.Printf("Number of metrics to retrieve: ")
                     NRows, _ := reader.ReadString('\n')
                     
                     fmt.Printf("Variables to display (Separate by commas): \n")
                     fmt.Println("1.CPU")
                     fmt.Println("2.RAM")
                     fmt.Println("3.Temperature")
                     fmt.Println("4.Humidity")
                     fmt.Println("5.Voltage")
                     fmt.Println("6.Current")
                     Variables, _ := reader.ReadString('\n')
                     Variables = strings.Replace(Variables, "\n", "", -1)
                     VarByte := strings.Split(Variables, ",")       
                     
                     for  index := 0; index< len(VarByte); index ++{
                         if VarByte[index] == "1" || VarByte[index] == "2" {
                            getValues(database, VarByte[index],NRows,0)
                            }
                         if VarByte[index] == "3" || VarByte[index] == "4" || VarByte[index] =="5" || VarByte[index] =="6" {
                            getValues(database, VarByte[index],NRows, 1)
                            }
                         }

                     fmt.Printf("-> ")                   
                    }  
          if text == "3" {
                     fmt.Printf("Variables to display (Separate by commas): \n")
                     fmt.Println("1-CPU")
                     fmt.Println("2-RAM")
                     fmt.Println("3-Temperature")
                     fmt.Println("4-Humidity")
                     fmt.Println("5-Voltage")
                     fmt.Println("6-Current")
                     Variables, _ := reader.ReadString('\n')
                     Variables = strings.Replace(Variables, "\n", "", -1)
                     VarByte := strings.Split(Variables, ",")       
                     
                     for  index := 0; index< len(VarByte); index ++{
                         if VarByte[index] == "1" || VarByte[index] == "2" {
                            averageValue(database, VarByte[index], 0)
                            }
                         if VarByte[index] == "3" || VarByte[index] == "4" || VarByte[index] =="5" || VarByte[index] =="6" {
                            averageValue(database, VarByte[index], 1)
                            }
                         }

                     fmt.Printf("-> ")
                     
                }
          if text == "0" {
                     fmt.Println("Simple Shell")
                     fmt.Println("-------------------")
                     fmt.Println("Available Commands: ")
                     fmt.Println("1. Get last n metrics for all variables;")
                     fmt.Println("2. Get last n metrics for one or more variables;")
                     fmt.Println("3. Get average value of one or more variables.")
                     fmt.Println("4. Close application.")
                     fmt.Println("Press 0 to see this menu again.")
                     fmt.Printf("-> ")
                     }
         if text == "4" {
                break            
                }
        idle0 = idle1
        total0 = total1
   }
    //Close database
   err = database.Close()
   if err != nil {
        fmt.Println("Error:" , err)
        return
    } 
}
