package main

import (
    "fmt"
    "bufio"
    "os"
    "internal/pokecache"
    "math/rand"
)

func commandExit(conf *Config, arg string) error {
    fmt.Println("Closing the Pokedex... Goodbye!")
    os.Exit(0)
    return nil
}

func commandHelp(conf *Config, arg string) error {
    fmt.Println("Welcome to the Pokedex!\nUsage:")
    fmt.Println("")

    for _, v := range conf.Commands {
        fmt.Printf("%s: %s\n", v.name, v.description)
    }
    return nil
}

func commandMap(conf *Config, arg string) error {
    defUrl := "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"
    var url string
    if conf.Next == nil {
        url = defUrl
    } else {
        url = *conf.Next
    }
    areas, err := get[Areas](url, conf.C)
    if (err != nil) {
        return err
    }

    conf.Prev = areas.Previous
    conf.Next = areas.Next
    for _, r := range areas.Results {
        fmt.Println(r.Name)
        //fmt.Println(r.URL)
    }
    return nil
}

func commandMapb(conf *Config, arg string) error {
    defUrl := "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"
    var url string
    if conf.Prev== nil {
        url = defUrl
    } else {
        url = *conf.Prev
    }
    areas, err := get[Areas](url, conf.C)
    if (err != nil) {
        return err
    }
    conf.Prev = areas.Previous
    conf.Next = areas.Next
    for _, r := range areas.Results {
        fmt.Println(r.Name)
    }
    return nil
}

func commandExplore(conf *Config, areaName string) error {
    fmt.Println("Exploring " + areaName)
    url := "https://pokeapi.co/api/v2/location-area/" + areaName
    area, err := get[Area](url, conf.C)
    if (err != nil) {
        fmt.Println("Can't get area info")
        return err
    }
    if len(area.PokemonEncounters) == 0 {
        fmt.Println("No Pokemons found")
        return nil
    }
    fmt.Println("Found Pokemon:")
    for _, p := range area.PokemonEncounters {
        fmt.Println(" - " + p.Pokemon.Name)
    }
    return nil
}


func commandCatch(conf *Config, pokeName string) error {
    fmt.Println("Throwing a Pokeball at " + pokeName + "...")
    url := "https://pokeapi.co/api/v2/pokemon/" + pokeName
    poke, err := get[Pokemon](url, conf.C)
    if (err != nil) {
        fmt.Println("Can't get pokemon info")
        return err
    }
    if (rand.Intn(1000) > poke.BaseExperience) {
        fmt.Println(pokeName + " was caught!")
        conf.Pokemons[pokeName] = poke
    } else {
        fmt.Println(pokeName + " escaped!")
    }
    return nil
}

func commandInspect(conf *Config, pokeName string) error {
    poke, ok := conf.Pokemons[pokeName]
    if !ok {
        fmt.Println("You have not caught that pokemon")
        return nil
    }
    fmt.Println("Name: " + poke.Name)
    fmt.Printf("Height: %v\n", poke.Height)
    fmt.Printf("Weight: %v\n", poke.Weight)
    fmt.Println("Stats:")
    for _, s := range poke.Stats {
        fmt.Printf("  - %v: %v\n", s.Stat.Name, s.BaseStat)
    }
    fmt.Println("Types:")
    for _, t := range poke.Types {
        fmt.Printf("  - %v\n", t.Type.Name)
    }

    return nil
}

func commandPokedex(conf *Config, arg string) error {
    if len(conf.Pokemons) < 1 {
        fmt.Println("You have not caught any pokemons")
        return nil
    }
    fmt.Println("Your Pokedex:")
    for n, _ := range conf.Pokemons {
        fmt.Printf("  - %v\n", n)
    }
    return nil
}

func main() {
    commands := map[string]cliCommand{
        "exit": {
            name:        "exit",
            description: "Exit the Pokedex",
            callback:    commandExit,
        },
        "help": {
            name:        "help",
            description: "Displays a help message",
            callback:    commandHelp,
        },
        "map": {
            name:        "map",
            description: "Displays 20 map items",
            callback:    commandMap,
        },
        "mapb": {
            name:        "mapb",
            description: "Displays preb 20 map items",
            callback:    commandMapb,
        },
        "explore": {
            name:        "explore <area_name>",
            description: "Displays a list of all the Pokémon located in area",
            callback:    commandExplore,
        },
        "catch": {
            name:        "catch <pokemon_name>",
            description: "Tries to catch named Pokemon",
            callback:    commandCatch,
        },
        "inspect": {
            name:        "inspect <pokemon_name>",
            description: "Prints Pokemon info",
            callback:    commandInspect,
        },
        "pokedex": {
            name:        "pokedex",
            description: "Lists Pokemons you caught",
            callback:    commandPokedex,
        },
    }
    cache := pokecache.NewCache(300)
    conf := Config{nil, nil, &cache, commands, make(map[string]Pokemon)}

    scanner := bufio.NewScanner(os.Stdin)
    for true {
        fmt.Print("Pokedex > ")
        scanner.Scan()
        req := scanner.Text()
        tokens := cleanInput(req)
        if (len(tokens) > 0) {
            name := tokens[0]
            command, ok := commands[name]
            if !ok {
                fmt.Println("Unknown command")
                continue
            }
            arg := ""
            if(len(tokens) > 1) {
                arg = tokens[1]
            }
            command.callback(&conf, arg)
            //fmt.Printf("Your command was: %s\n", tokens[0])
        }
    }
}
