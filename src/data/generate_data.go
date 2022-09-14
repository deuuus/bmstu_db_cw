package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/Pallinder/go-randomdata"
	"golang.org/x/crypto/bcrypt"
)

func encryptString(s string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

type Doctor struct {
	Name		      string `json:"name"`
	Surname           string `json:"surname"`
	Work_since        int    `json:"work_since"`
	Spec_id           int    `json:"spec_id"` 
	Email             string `json:"email"`
	Password          string `json:"password,omitempty"`
	EncryptedPassword string `json:"-"`
}


const (del = ",")

func main() {
	file_p, err := os.Create("patients.txt")
	if err != nil{
        fmt.Println("Unable to create file:", err) 
        os.Exit(1) 
    }
    defer file_p.Close() 

	password := "123456"
	var text_gender string

	doctors := []Doctor{}

	file_d, err := os.Create("doctors.txt")
	if err != nil{
        fmt.Println("Unable to create file:", err) 
        os.Exit(1) 
    }
    defer file_p.Close() 

	file_v, err := os.Create("visits.txt")
	if err != nil{
        fmt.Println("Unable to create file:", err) 
        os.Exit(1) 
    }
    defer file_v.Close() 

	file_r, err := os.Create("records.txt")
	if err != nil{
        fmt.Println("Unable to create file:", err) 
        os.Exit(1) 
    }
    defer file_r.Close()

	//doctors
	for i := 0; i < 100; i++ {
		gender := rand.Intn(2) + 1
		work_since := 1942 + rand.Intn(80)
		spec_id := rand.Intn(7) + 1
		enc, err := encryptString(password)
		if err != nil{
            fmt.Println("Unable to encrypt", err) 
            os.Exit(1) 
        }
		doctor := Doctor{Name: randomdata.FirstName(gender), Surname: randomdata.LastName(),
						Work_since: work_since, Spec_id: spec_id, Email: randomdata.Email(), EncryptedPassword: enc,}
		file_d.WriteString(doctor.Name + del + doctor.Surname +
						 del + strconv.Itoa(doctor.Work_since) + del + strconv.Itoa(doctor.Spec_id) + del + doctor.Email + del + enc + "\n")
		doctors = append(doctors, doctor)
	}

	vs := 0

	//patients
	for i := 0; i < 1000; i++ {
		gender := rand.Intn(2) + 1
		if gender == 1 {
			text_gender = "Женский"
		} else {
			text_gender = "Мужской"
		}
		birth_year := 1922 + rand.Intn(100)
		enc, err := encryptString(password)
		if err != nil{
            fmt.Println("Unable to encrypt", err) 
            os.Exit(1) 
        }
		phone := randomdata.PhoneNumber()
		phone = strings.ReplaceAll(phone, ",", "")
		phone = strings.ReplaceAll(phone, " ", "")
		file_p.WriteString(randomdata.FirstName(gender) + del + randomdata.LastName() +
						 del + string(text_gender) + del + strconv.Itoa(birth_year) + del + phone + del + randomdata.Email() + del + enc + "\n")

		for j := 0; j < 3; j++ {
			vs += 1
			k := rand.Intn(100) + 1
			file_v.WriteString("Active" + del + strconv.Itoa(i + 1) + del + strconv.Itoa(k) + "\n")

			vs += 1
			k = rand.Intn(100) + 1
			file_v.WriteString("Done" + del + strconv.Itoa(i + 1) + del + strconv.Itoa(k) + "\n")
			file_r.WriteString(strconv.Itoa(vs) + del + strconv.Itoa(rand.Intn(3) + 1 + 3 * (doctors[k-1].Spec_id - 1)) + del + strconv.Itoa(rand.Intn(9) + 1) + "\n")
		}
	}

	file_a, err := os.Create("admins.txt")
	if err != nil{
        fmt.Println("Unable to create file:", err) 
        os.Exit(1) 
    }
    defer file_p.Close() 

	//admins
	for i := 0; i < 10; i++ {
		gender := rand.Intn(2) + 1
		enc, err := encryptString(password)
		if err != nil{
            fmt.Println("Unable to encrypt", err) 
            os.Exit(1) 
        }
		file_a.WriteString(randomdata.FirstName(gender) + del + randomdata.LastName() +
							 del + randomdata.Email() + del + enc + "\n")
	}
}