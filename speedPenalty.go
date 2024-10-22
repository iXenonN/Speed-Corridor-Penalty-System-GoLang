package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Vehicle struct { // define a structure called Vehicle
	Plate        string    // license plate information of the vehicle
	TimeIn       time.Time // vehicle entry time
	TimeOut      time.Time // vehicle exit time
	AverageSpeed float64   // average speed of the vehicle
}

var enteredPlates []string // create a list named enteredPlates
var vehicles []Vehicle

const speedLimit = 70.0 // speed limit of the road

func main() {
	entrance()
	exit()
	calc()
}

func entrance() {
	filename := "entrance.txt"     // specify file name
	file, err := os.Open(filename) // open file
	if err != nil {
		fmt.Println("Error While Trying To Open File:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file) // scan the file line by line
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line) // divide and separate spaces e.g. parts = ["06AU7456", "17", "40", "45"]

		// validate rows LicensePlate Hour Minute Seconds
		if len(parts) < 4 {
			continue // invalid row
		}

		plate := parts[0]
		hour, _ := strconv.Atoi(parts[1])   // convert string value to integer
		minute, _ := strconv.Atoi(parts[2]) // convert string value to integer
		second, _ := strconv.Atoi(parts[3]) // convert string value to integer

		timeIn := time.Date(0, 1, 1, hour, minute, second, 0, time.UTC) // vehicle entry time

		enteredPlates = append(enteredPlates, plate) // add license plates of entering cars to the enteredPlates list
		vehicles = append(vehicles, Vehicle{
			Plate:  plate,
			TimeIn: timeIn,
		})
	}
}

func exit() {
	filename := "exit.txt"         // specify file name
	file, err := os.Open(filename) // open file
	if err != nil {
		fmt.Println("Error While Trying To Open File:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file) // scan the file line by line
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line) // divide and separate spaces e.g. parts = ["06AU7456", "17", "40", "45"]

		// validate rows LicensePlate Hour Minute Seconds
		if len(parts) < 4 {
			continue // invalid row
		}

		plate := parts[0]
		hour, _ := strconv.Atoi(parts[1])   // convert string value to integer
		minute, _ := strconv.Atoi(parts[2]) // convert string value to integer
		second, _ := strconv.Atoi(parts[3]) // convert string value to integer

		timeOut := time.Date(0, 1, 1, hour, minute, second, 0, time.UTC) // vehicle exit time

		for i := range vehicles {
			if vehicles[i].Plate == plate { // compare license plates of entering and exiting vehicles
				vehicles[i].TimeOut = timeOut // add exit time of matching license plate
				break
			}
		}
	}
}

func calc() {
	totalDistance := 7.43 // road (in km)

	// create penalty.txt file
	penaltyFile, err := os.Create("penalty.txt")
	if err != nil {
		fmt.Println("Error Creating Penalty File:", err)
		return
	}
	defer penaltyFile.Close()

	for _, vehicle := range vehicles {
		duration := vehicle.TimeOut.Sub(vehicle.TimeIn).Hours() // subtract the time the vehicle leaves the road from the time it enters the road and find the time it spends on the road.
		if duration > 0 {
			averageSpeed := totalDistance / duration // distance = Speed ​​x Time Calculate average speed (Speed ​​= Distance / Time)
			vehicle.AverageSpeed = averageSpeed      // save average speed to structure
		}

		fmt.Printf("%s numbered plate entered at %s and exited at %s with average speed: %.2f km/h\n",
			vehicle.Plate, vehicle.TimeIn.Format("15:04:05"), // set time format
			vehicle.TimeOut.Format("15:04:05"), vehicle.AverageSpeed)

		// Write penalty for vehicles exceeding the speed limit
		if vehicle.AverageSpeed > speedLimit {
			penaltyFile.WriteString(fmt.Sprintf("%s %.2f\n", vehicle.Plate, vehicle.AverageSpeed))
			fmt.Printf("Penalty issued for plate %s with speed %.2f km/h\n", vehicle.Plate, vehicle.AverageSpeed)
		}
	}
}
