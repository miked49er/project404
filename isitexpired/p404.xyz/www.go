package main

import (
  "html/template"
  "net/http"
  "strconv"
  "fmt"
  "time"
)

type Date struct {
  Month int
  Day int
  Year int
}

type DateString struct {
  Month string
  Day string
  Year string
}

type Brand struct {
  Brand string
  Pull string
  Expiration string
}

type DateCode struct {
  Code string
  Month int
}

func strToInt(str string) int {
  result, err := strconv.Atoi(str)

  if (err != nil) {
    fmt.Println("strToInt error", err)
  }

  return result
}

func intToStr(num int) string {
  return strconv.Itoa(num)
}

func getDateCodeArray() []DateCode {
  return []DateCode {
    {
      "A",
      1,
    },
    {
      "B",
      2,
    },
    {
      "C",
      3,
    },
    {
      "D",
      4,
    },
    {
      "E",
      5,
    },
    {
      "F",
      6,
    },
    {
      "G",
      7,
    },
    {
      "H",
      8,
    },
    {
      "J",
      9,
    },

    {
      "K",
      10,
    },
    {
      "L",
      11,
    },
    {
      "M",
      12,
    },
  }
}

func getMonthCode(month int) string {
  codes := getDateCodeArray()
  letter := ""

  for i := 0; i < len(codes); i++ {
    if month == codes[i].Month {
      letter = codes[i].Code
    }
  }
  return letter
}

func getYearCode(year int) string {
  return intToStr(year % 10)
}

// Letter for 5 Months ago
func getLetterCode(date Date) string {

  month5Ago := date.Month - 5
  currentCode := getMonthCode(month5Ago)

  return currentCode
}

func getDateCode(date Date) string {

  currentCode := getMonthCode(date.Month) + intToStr(date.Day) + getYearCode(date.Year)

  return currentCode
}

func expiration(w http.ResponseWriter, r *http.Request) {

  year, month, day := time.Now().Date()

  date := Date {
    int(month),
    day,
    year,
  }

  dateStr := DateString {
    month.String(),
    intToStr(day),
    intToStr(year),
  }

  letterCode := getLetterCode(date) + getYearCode(date.Year)
  dateCode := getMonthCode(date.Month) + intToStr(date.Day) + getYearCode(date.Year)
  snusVuse := getMonthCode(date.Month) + getYearCode(date.Year - 1)

  brands := []Brand {
    {
      "Grizzly | Kodiak",
      letterCode,
      "Last Day of " + dateStr.Month + " " + dateStr.Year,
    },
    {
      "Copenhagon | Skoal",
      dateStr.Month + " " + dateStr.Day + " " + dateStr.Year,
      dateStr.Month + " " + dateStr.Day + " " + dateStr.Year,
    },
    {
      "Timber Wolf | Longhorn",
      dateCode,
      dateStr.Month + " " + dateStr.Day + " " + dateStr.Year,
    },
    {
      "CAMEL SNUS | VUSE",
      snusVuse,
      "Last Day of " + dateStr.Month + " " + dateStr.Year,
    },
  }

  data := struct {
    Date string
    Brands []Brand
  } { dateStr.Month + " " + dateStr.Day + " " + dateStr.Year, brands }

  tmpl, _ := template.ParseFiles("public/layout.html", "public/templates/index.html")
  tmpl.ExecuteTemplate(w, "layout", data)
}

func dateCodes(w http.ResponseWriter, r *http.Request) {
  tmpl, _ := template.ParseFiles("public/layout.html", "public/templates/codeFormat.html")
  tmpl.ExecuteTemplate(w, "layout", nil)
}

func about(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("public/layout.html", "public/templates/about.html")
	tmpl.ExecuteTemplate(w, "layout", nil)
}

func serveResource(w http.ResponseWriter, req *http.Request) {
  path := "public" + req.URL.Path
  http.ServeFile(w, req, path)
}

func main() {
  http.HandleFunc("/", expiration)
  http.HandleFunc("/date-code", dateCodes)
  http.HandleFunc("/about", about)
  http.HandleFunc("/css/", serveResource)
  http.HandleFunc("/img/", serveResource)
  http.ListenAndServe(":3000", nil)
}
