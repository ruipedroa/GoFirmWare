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

//Retrieves Values from the Database

func getValues(db *sql.DB, ival[4] int, fval[2] float64) () {
    
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

func main() {
    
    reader := bufio.NewReader(os.Stdin)
    fmt.Println("Simple Shell")
    fmt.Println("-------------------")
    fmt.Println("Available Commands: ")
    fmt.Println("1. Get last n metrics for all variables;")
    fmt.Println("2. Get last n metrics for one or more variables;")
    fmt.Println("3. Get average value of one or more variables.")
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
    
    idle0,total0 :=getCPUSample()
    // Cycle to get reads every second
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
                  fmt.Printf("-> ")
             }     
         if text == "2"{ 
                     fmt.Printf("Number of metrics to retrieve: ")
                     //NRows, _ := reader.ReadString('\n')
                     
                     fmt.Printf("Variables to display (Separate by commas): \n")
                     fmt.Println("1-CPU")
                     fmt.Println("2-RAM")
                     fmt.Println("3-Temperature")
                     fmt.Println("4-Humidity")
                     fmt.Println("5-Voltage")
                     fmt.Println("6-Current")
                    /* Variables, _ := reader.ReadString('\n')
                     Variables := strings.Split(SystemVariables, ",")
                     
                     rowsSystem, _ :=database.Query("SELECT"+ SystemVariables +" FROM SystemData WHERE id > (SELECT MAX(id)  - " + NRows + "FROM SystemData)")
                     fmt.Printf("Variables to display (use comma to separate) from System (temp, humidity, voltage, current): ")
                     DeviceVariables, _ := reader.ReadString('\n')
                     rowsDevice, _ :=database.Query("SELECT" + DeviceVariables +" FROM DeviceData WHERE id > (SELECT MAX(id)  - " + text + "FROM DeviceData)")
                     fmt.Println(SystemVariables + " : ")        
                     for rowsSystem.Next()  {
                         
                         rowsSystem.Scan(&CPU, &RAM)
                         fmt.Println(strconv.FormatFloat(CPU, 'f', -1, 64) +"          "+ strconv.FormatFloat(RAM, 'f', -1, 64))
                     }
                     fmt.Println(DeviceVariables + " : ") 
                     for rowsDevice.Next()  {
                        rowsDevice.Scan(&temp, &humidity, &voltage, &current)
                        fmt.Println(strconv.Itoa(temp)  +"                 "+ strconv.Itoa(humidity) + "             " + strconv.Itoa(voltage) + "            " + strconv.Itoa(current))
                     }  */
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
                     fmt.Printf("-> ")
                }
          if text == "0" {
                     fmt.Println("Simple Shell")
                     fmt.Println("-------------------")
                     fmt.Println("Available Commands: ")
                     fmt.Println("1. Get last n metrics for all variables;")
                     fmt.Println("2. Get last n metrics for one or more variables;")
                     fmt.Println("3. Get average value of one or more variables.")
                     fmt.Println("Press 0 to see this menu again.")
                     fmt.Printf("-> ")
                     }
        idle0 = idle1
        total0 = total1

   }
}
