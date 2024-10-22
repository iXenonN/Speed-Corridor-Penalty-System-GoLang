package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Vehicle struct { // Vehicle adında bir yapı tanımlanır
	Plate        string    // Aracın plaka bilgisi
	TimeIn       time.Time // Aracın yola giriş zamanı
	TimeOut      time.Time // Aracın yoldan çıkış zamanı
	AverageSpeed float64   // Aracın ortalama hızı
}

var enteredPlates []string // enteredPlates adında bir liste oluştur
var vehicles []Vehicle

const speedLimit = 70.0 // Hız sınırı

func main() {
	entrance()
	exit()
	calc()
}

func entrance() {
	filename := "entrance.txt"     // Dosya adını belirt
	file, err := os.Open(filename) // Dosyayı aç
	if err != nil {
		fmt.Println("Error While Trying To Open File:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file) // Dosyayı satır satır tara
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line) // Boşlukları böl ve ayır örn parts = ["06AU7456", "17", "40", "45"]

		// Satırları doğrula Plaka Saat Dakika Saniye
		if len(parts) < 4 {
			continue // Geçersiz satır
		}

		plate := parts[0]
		hour, _ := strconv.Atoi(parts[1])   // String değerini integer a çevir
		minute, _ := strconv.Atoi(parts[2]) // String değerini integer a çevir
		second, _ := strconv.Atoi(parts[3]) // String değerini integer a çevir

		timeIn := time.Date(0, 1, 1, hour, minute, second, 0, time.UTC) // Aracın yola girdiği zaman

		enteredPlates = append(enteredPlates, plate) // Giren arabaların plakalarını enteredPlates listesine ekle
		vehicles = append(vehicles, Vehicle{
			Plate:  plate,
			TimeIn: timeIn,
		})
	}
}

func exit() {
	filename := "exit.txt"         // Dosya adını belirt
	file, err := os.Open(filename) // Dosyayı Aç
	if err != nil {
		fmt.Println("Error While Trying To Open File:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file) // Dosyayı satır satır tara
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line) // Boşlukları böl ve ayır örn parts = ["06AU7456", "17", "40", "45"]

		// Satırları doğrula Plaka Saat Dakika Saniye
		if len(parts) < 4 {
			continue // Geçersiz satır
		}

		plate := parts[0]
		hour, _ := strconv.Atoi(parts[1])   // String değerini integer a çevir
		minute, _ := strconv.Atoi(parts[2]) // String değerini integer a çevir
		second, _ := strconv.Atoi(parts[3]) // String değerini integer a çevir

		timeOut := time.Date(0, 1, 1, hour, minute, second, 0, time.UTC) // Aracın yoldan çıktığı zaman

		for i := range vehicles {
			if vehicles[i].Plate == plate { // Giren ve çıkan araçların plakalarını karşılaştır
				vehicles[i].TimeOut = timeOut // Eşleşen plakanın çıkış saatini ekle
				break
			}
		}
	}
}

func calc() {
	totalDistance := 7.43 // Yol (km cinsinden)

	// Ceza dosyasını oluştur
	penaltyFile, err := os.Create("penalty.txt")
	if err != nil {
		fmt.Println("Error Creating Penalty File:", err)
		return
	}
	defer penaltyFile.Close()

	for _, vehicle := range vehicles {
		duration := vehicle.TimeOut.Sub(vehicle.TimeIn).Hours() // Aracın yoldan çıkış zamanını yola giriş zamanındna çıkar ve yolda geçirdiği zamanı bul
		if duration > 0 {
			averageSpeed := totalDistance / duration // Yol = Hız x Zaman Ortalama hızı hesapla (Hız = Yol / Zaman)
			vehicle.AverageSpeed = averageSpeed      // Ortalama hızı structure'a kaydet
		}

		fmt.Printf("%s numbered plate entered at %s and exited at %s with average speed: %.2f km/h\n",
			vehicle.Plate, vehicle.TimeIn.Format("15:04:05"), // Saat formatını belirle
			vehicle.TimeOut.Format("15:04:05"), vehicle.AverageSpeed)

		// Hız sınırını aşan araçlar için ceza yaz
		if vehicle.AverageSpeed > speedLimit {
			penaltyFile.WriteString(fmt.Sprintf("%s %.2f\n", vehicle.Plate, vehicle.AverageSpeed))
			fmt.Printf("Penalty issued for plate %s with speed %.2f km/h\n", vehicle.Plate, vehicle.AverageSpeed)
		}
	}
}
