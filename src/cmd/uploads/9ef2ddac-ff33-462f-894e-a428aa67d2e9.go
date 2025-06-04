package main

import (
	ac "Bank/acoount"
	"bufio"
	"fmt"
	"os"
	"strings"
)

var acounnts []ac.Acounnt

func main() {
	fmt.Println("********** Welcome To Bank **********")
	for {
		fmt.Println(strings.Repeat("~", 40))
		fmt.Print("1) Make Acounnt\n2) Check Acounnt\n3) Send Money\n4) Recive Money\n\tchoise: ")
		var choise int
		fmt.Scanln(&choise)
		fmt.Println(strings.Repeat("~", 40))
		if choise == 1 {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("-enter your name: ")
			name, _ := reader.ReadString('\n')
			var (
				age   int
				walet float64
			)
			fmt.Print("-enter your age: ")
			fmt.Scanln(&age)
			fmt.Print("-enter your walet: ")
			fmt.Scanln(&walet)
			id := len(acounnts) + 1
			acounnts = append(acounnts, ac.Acounnt{
				ID:    id,
				Name:  name,
				Age:   age,
				Walet: walet,
			})
			fmt.Printf("your acounnt created with ID: %d\n", id)
		} else if choise == 2 {
			var ID int
			fmt.Print("enter your ID: ")
			fmt.Scanln(&ID)
			for _, acounnt := range acounnts {
				if acounnt.ID == ID {
					acounnt.ShowAcounnt()
					break
				}
				ac.CantShow()
			}
		} else if choise == 3 {
			var money_amount float64
			fmt.Print("enter money amount: ")
			fmt.Scanln(&money_amount)
			var (
				sender  int
				reciver int
			)
			fmt.Print("enter your ID: ")
			fmt.Scanln(&sender)
			fmt.Print("enter target ID: ")
			fmt.Scanln(&reciver)
			for _, sen := range acounnts {
				for _, rec := range acounnts {
					if sen.ID == sender && rec.ID == reciver {
						if sen.Walet >= money_amount {
							sen.Walet -= money_amount
							rec.Walet += money_amount
							fmt.Printf("%.2f$ sended from %d to %d\n", money_amount, sen.ID, rec.ID)
							break
						}
					}
				}
			}
		} else if choise == 4 {
			var money_amount float64
			fmt.Print("enter money amount: ")
			fmt.Scanln(money_amount)
			var ID int
			fmt.Print("enter your ID: ")
			fmt.Scanln(&ID)
			for _, acounnt := range acounnts {
				if acounnt.ID == ID {
					acounnt.Walet += money_amount
					fmt.Printf("%.2f$ you recived\n",money_amount)
					break
				}
			}
		}
	}
}
