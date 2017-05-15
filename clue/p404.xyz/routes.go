package main

import (
  "html/template"
  "net/http"
)

type Item struct {
  Id string
  Name string
  Selection int
}

type List struct {
  Title string
  Items []Item
}

func main() {
  tmpl := template.Must(template.ParseFiles("public/index.html"))
  lists := []List{
    {"Guests", []Item{
        {"mickeymouse", "Mickey Mouse", 0},
        {"minniemouse", "Minnie Mouse", 0},
        {"donaldduck", "Donald Duck", 0},
        {"daisyduck", "Daisy Duck", 0},
        {"goofy", "Goofy", 0},
        {"pluto", "Pluto", 0},
    }},
  {"Ghosts", []Item{
    {"bride", "Bride", 0},
    {"mariner", "Mariner", 0},
    {"operasinger", "Opera Singer", 0},
    {"skeleton", "Skeleton", 0},
    {"prisoner", "Prisoner", 0},
    {"traveler", "Traveler", 0},
    }},
  {"Rooms", []Item{
    {"foyer", "Foyer", 0},
    {"portraitchamber", "Portrait Chamber", 0},
    {"seanceroom", "Seance Room", 0},
    {"graveyard", "Graveyard", 0},
    {"ballroom", "Ballroom", 0},
    {"conservatory", "Conservatory", 0},
    {"mausoleum", "Mausoleum", 0},
    {"library", "Library", 0},
    {"attic", "Attic", 0},
    }},
  }

  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    tmpl.Execute(w, struct{ Lists []List }{lists})
  })

  http.ListenAndServe(":3000", nil)
}
