import (
        "fmt"
        "math/rand"
        "time"
    )



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


