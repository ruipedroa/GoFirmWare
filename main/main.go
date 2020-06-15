package main

import (
	"fmt"
    "io/ioutil"
    "strconv"
    "strings"
	"time"
	"github.com/shirou/gopsutil/mem"
)

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
    
    idle0,total0 :=getCPUSample()
    for {
        time.Sleep(time.Second)  
        idle1, total1:= getCPUSample()

        idleTicks :=float64(idle1 -idle0)
        totalTicks := float64(total1 -total0)
        cpuUsage:=100* (totalTicks -idleTicks) /totalTicks
	
        fmt.Printf("CPU Usage is %f%% [busy :%f, total: %f]\n", cpuUsage, totalTicks-idleTicks, totalTicks)
    
        //get Memory percentage
	    v, _ := mem.VirtualMemory()

	    fmt.Printf("Pecentage of RAM used:%f%%\n",v.UsedPercent)
        idle0 = idle1
        total0 = total1
        }
}
