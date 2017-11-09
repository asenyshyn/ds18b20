# ds18b20
Read sensor data from ds18b20 for Raspberry PI

## Example

```go
package main

import (
	"log"
	"sync"
	"time"

	"github.com/asenyshyn/ds18b20"
)

func main() {
	sensors, err := ds18b20.Sensors()
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup

	for {
		wg.Add(len(sensors))
		for _, s := range sensors {
			go func(s ds18b20.Sensor) {
				t, err := s.Reading()
				if err != nil {
					log.Fatal(err)
				}
				log.Println(s.ID, t.Value)
				wg.Done()
			}(s)
		}
		wg.Wait()
		time.Sleep(time.Second * 5)
	}
}
```
