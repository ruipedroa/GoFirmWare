package main

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/cpu"
)

func main() {

	//Percent calculates the percentage of cpu used either per CPU or combined.
	percent, _ := cpu_windows.Percent(time.Second, true)
	fmt.Println("  User: %.2f\n", percent[cpu.CPUser])
	fmt.Println("  Nice: %.2f\n", percent[cpu.CPNice])
	fmt.Println("   Sys: %.2f\n", percent[cpu.CPSys])
	fmt.Println("  Intr: %.2f\n", percent[cpu.CPIntr])
	fmt.Println("  Idle: %.2f\n", percent[cpu.CPIdle])
	fmt.Println("States: %.2f\n", percent[cpu.CPUStates])
}
