package main 

import (
  "html/template"
  "net/http"
  "strconv"
  "fmt"
  "reflect"
  "regexp"
  "math"
)

type Coin struct {
	Label string
	Loose string
	Rolled string
}
type Bill struct {
	Label string
	Name string
}
type FCoin struct {
	Label string
	Loose float64
	Rolled float64
}
type FBill struct {
	Label string
	Value float64
}
type CoinTable struct {
	PenniesL float64
	PenniesR float64
	NickelsL float64
	NickelsR float64
	DimesL float64
	DimesR float64
	QuartersL float64
	QuartersR float64
	Other float64
}
type BillTable struct {
	Ones float64
	Fives float64
	Tens float64
	Twenties float64
	Fifties float64
	Hundreds float64
}
type Till struct {
	Drop []FBill
	NextTotal float64
	Coins []FCoin
	OtherCoins float64
	Bills []FBill
	CoinTotal float64
	SubTotal float64
	OtherDrops float64
	Final float64
	Total float64
	CurrentTill float64
	PreviousTill float64
	Accountable float64
	Over float64
}
type LotterySlot struct {
	Slot float64
	Ticket float64
	Value float64
}

func checkNum(str string) string {
	reg := regexp.MustCompile("[0-9|.]+")
	newStr := reg.FindAllString(str, -1.0)

	result := "0"
	if len(newStr) > 0 {
		result = newStr[0]
	}
	return result
}
func strToFloat(str string) float64 {

	if str == "" {
		return 0
	}

	num, err := strconv.ParseFloat(str, 64)

	if err != nil {
		fmt.Println("strToFloat Error:", err)
	}

	return num
}
func floatToStr(float float64, prec int) string {

	if float < 0.01 && float > -0.01 {
		return "0"
	}

	str := strconv.FormatFloat(float, 'f', prec, 64)

	return str
}
func Round(val float64, roundOn float64, places int ) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}
func calculateDrop(drop float64, value float64, overBase float64) float64 {
	if drop <= 0 {
		return 0
	}
	if overBase > 0 {
		modulus := math.Mod(drop, value)

		result := drop - modulus

		if result <= overBase {
			return result
		}
		return overBase
	}
	return 0
}
func extraDrop(dp float64,db float64,b float64,v float64) (float64, float64, float64) {
	for dp >= v && b > 0 {
		dp -= v
		db += v
		b -= v
	}
	return dp, db, b
}
func till(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("public/layout.html", "public/till.html")
	
	coins := []Coin {
  	{"Pennies", "pennies", "pennyRoll"},
  	{"Nickels", "nickels", "nickelRoll"},
  	{"Dimes", "dimes", "dimeRoll"},
  	{"Quarter", "quarters", "quarterRoll"},
	}
	bills := []Bill {
  	{"Ones", "ones"},
  	{"Five", "fives"},
  	{"Tens", "tens"},
  	{"Twenties", "twenties"},
  	{"Fifties", "fifties"},
  	{"Hundreds", "hundreds"},
	}
	data := struct {
		Coins []Coin
		Bills []Bill
	}{coins, bills}
	tmpl.ExecuteTemplate(w, "layout", data)
}
func count(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()

	if err != nil {
		fmt.Fprintln(w, err.Error())
		return
	}

	//fmt.Fprintln(w, r.Form)

	coins := CoinTable {
		PenniesL: strToFloat(r.Form["pennies"][0]),
		PenniesR: strToFloat(r.Form["pennyRoll"][0]),
		NickelsL: strToFloat(r.Form["nickels"][0]),
		NickelsR: strToFloat(r.Form["nickelRoll"][0]),
		DimesL: strToFloat(r.Form["dimes"][0]),
		DimesR: strToFloat(r.Form["dimeRoll"][0]),
		QuartersL: strToFloat(r.Form["quarters"][0]),
		QuartersR: strToFloat(r.Form["quarterRoll"][0]),
		Other: strToFloat(r.Form["other"][0]),
	} // coins

	bills := BillTable {
		Ones: strToFloat(r.Form["ones"][0]),
		Fives: strToFloat(r.Form["fives"][0]),
		Tens: strToFloat(r.Form["tens"][0]),
		Twenties: strToFloat(r.Form["twenties"][0]),
		Fifties: strToFloat(r.Form["fifties"][0]),
		Hundreds: strToFloat(r.Form["hundreds"][0]),
	} // bills

	coinTotal := 0.0
	billTotal := 0.0

	coinR := reflect.ValueOf(coins)

	for i := 0; i < coinR.NumField(); i++ {
		val := coinR.Field(i).Float()
		coinTotal += val
	}

	billR := reflect.ValueOf(bills)

	for i := 0; i < billR.NumField(); i++ {
		val := billR.Field(i).Float()
		billTotal += val
	}

	subTotal := coinTotal + billTotal
	finalDrop := subTotal - 300.0

	otherDrops := strToFloat(r.Form["otherDrops"][0])

	total := finalDrop + otherDrops

	currentTill := strToFloat(r.Form["currentTill"][0])
	previousTill := strToFloat(r.Form["previousTill"][0])

	accountable := currentTill - previousTill
	over := total - accountable

	// Drop

	// Hundreds

	hundreds := bills.Hundreds
	dropHundreds := calculateDrop(finalDrop, 100.0, hundreds)
	hundreds -= dropHundreds

	// Fifties

	lessDrop := Round(finalDrop - dropHundreds,1,0)
	fifties := bills.Fifties
	dropFifties := calculateDrop(lessDrop, 50.0, fifties)
	fifties -= dropFifties

	// Twenties

	lessDrop -= dropFifties
	twenties := bills.Twenties
	over20 := twenties - 80
	dropTwenties := calculateDrop(lessDrop, 20.0, over20)
	twenties -= dropTwenties

	// Tens

	lessDrop -= dropTwenties
	tens := bills.Tens
	over10 := tens - 80
	dropTens := calculateDrop(lessDrop, 10.0, over10)
	tens -= dropTens

	// Fives

	lessDrop -= dropTens
	fives := bills.Fives
	over5 := fives - 60
	dropFives := calculateDrop(lessDrop, 5.0, over5)
	fives -= dropFives

	// Ones

	lessDrop -= dropFives
	ones := bills.Ones
	over1 := ones - 20
	dropOnes := calculateDrop(lessDrop, 1.0, over1)
	ones -= dropOnes

	// Extra TODO

	lessDrop -= dropOnes
	if lessDrop > 0 {
		lessDrop, dropTwenties, twenties = extraDrop(lessDrop, dropTwenties, twenties, 20.0)
		lessDrop, dropTens, tens = extraDrop(lessDrop, dropTens, tens, 10.0)
		lessDrop, dropFives, fives = extraDrop(lessDrop, dropFives, fives, 5.0)
		lessDrop, dropOnes, ones = extraDrop(lessDrop, dropOnes, ones, 1.0)
		fmt.Println(w, "Left Over:", lessDrop)
	}

	// Next Till

	nextTotal := coinTotal + ones + fives + tens + twenties + fifties + hundreds

	nextTill := Till {
		Drop: []FBill {
			{"Ones", dropOnes},
			{"Fives", dropFives},
			{"Tens", dropTens},
			{"Twenties", dropTwenties},
			{"Fifties", dropFifties},
			{"Hundreds", dropHundreds},
		},
		NextTotal: Round(nextTotal, 1, 2),
		Coins: []FCoin {
			{"Pennies", coins.PenniesL, coins.PenniesR},
			{"Nickels", coins.NickelsL, coins.NickelsR},
			{"Dimes", coins.DimesL, coins.DimesR},
			{"Quarters", coins.QuartersL, coins.QuartersR},
		},
		OtherCoins: coins.Other,
		Bills: []FBill {
			{"Ones", ones},
			{"Fives", fives},
			{"Tens", tens},
			{"Twenties", twenties},
			{"Fifties", fifties},
			{"Hundreds", hundreds},
		},
		CoinTotal: Round(coinTotal, 1, 2),
		SubTotal: Round(subTotal, 1, 2),
		OtherDrops: Round(otherDrops, 1, 0),
		Final: Round(finalDrop, 1, 0),
		Total: Round(total, 1, 2),
		CurrentTill: Round(currentTill, 1, 2),
		PreviousTill: Round(previousTill, 1, 2),
		Accountable: Round(accountable, 1, 2),
		Over: Round(over, 1, 2),
	}

	tmpl, _ := template.ParseFiles("public/layout.html", "public/count.html")
	errCount := tmpl.ExecuteTemplate(w, "layout", nextTill)
	
	if errCount != nil {
		fmt.Println(err)
	}
}
func lottery(w http.ResponseWriter, r *http.Request) {
	
	data := []string { "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12" }
	tmpl, _ := template.ParseFiles("public/layout.html", "public/lottery.html")
	err := tmpl.ExecuteTemplate(w, "layout", data)
	
	if err != nil {
		fmt.Println(err)
	}
}
func totals(w http.ResponseWriter, r *http.Request) {
	
	err := r.ParseForm()

	if err != nil {
		fmt.Fprintln(w, err.Error())
		return
	}
	
	slots := []LotterySlot {
		{ 1, strToFloat(r.Form["slot1"][0]), strToFloat(r.Form["value1"][0]) },
		{ 2, strToFloat(r.Form["slot2"][0]), strToFloat(r.Form["value2"][0]) },
		{ 3, strToFloat(r.Form["slot3"][0]), strToFloat(r.Form["value3"][0]) },
		{ 4, strToFloat(r.Form["slot4"][0]), strToFloat(r.Form["value4"][0]) },
		{ 5, strToFloat(r.Form["slot5"][0]), strToFloat(r.Form["value5"][0]) },
		{ 6, strToFloat(r.Form["slot6"][0]), strToFloat(r.Form["value6"][0]) },
		{ 7, strToFloat(r.Form["slot7"][0]), strToFloat(r.Form["value7"][0]) },
		{ 8, strToFloat(r.Form["slot8"][0]), strToFloat(r.Form["value8"][0]) },
		{ 9, strToFloat(r.Form["slot9"][0]), strToFloat(r.Form["value9"][0]) },
		{ 10, strToFloat(r.Form["slot10"][0]), strToFloat(r.Form["value10"][0]) },
		{ 11, strToFloat(r.Form["slot11"][0]), strToFloat(r.Form["value11"][0]) },
		{ 12, strToFloat(r.Form["slot12"][0]), strToFloat(r.Form["value12"][0]) },
	}
	
	total := 0.0
	
	for _, slot := range slots {
		total += slot.Value * ( slot.Ticket + 1 )
		//fmt.Fprintln(w, "Value", slot.Value, "Ticket", slot.Ticket, "Total", total)
	}
	
	data := struct {
		Slots []LotterySlot
		Total float64
	} { slots, total }
	tmpl, _ := template.ParseFiles("public/layout.html", "public/lotteryTotal.html")
	tmpl.ExecuteTemplate(w, "layout", data)
}
func about(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("public/layout.html", "public/about.html")
	tmpl.ExecuteTemplate(w, "layout", nil)
}
func serveResource(w http.ResponseWriter, req *http.Request) {
    path := "public" + req.URL.Path
    http.ServeFile(w, req, path)
}
func main() {
  http.HandleFunc("/", till)
  http.HandleFunc("/count", count)
  http.HandleFunc("/lottery", lottery)
  http.HandleFunc("/totals", totals)
  http.HandleFunc("/about", about)
  http.HandleFunc("/css/", serveResource)
  http.HandleFunc("/img/", serveResource)
  http.ListenAndServe(":3000", nil)
}