package main

import "unicode"

type cliCommand struct {
    name        string
    description string
    callback    func(*Config, string) error
}

func cleanInput(text string) []string {
    var tokens = make([]string, 0)
    token := ""
    for _, c := range text {
        if c == ' ' {
            if len(token) > 0 {
                tokens = append(tokens, token);
                token = ""
            }
        } else {
            token += string(unicode.ToLower(c))
        }
    }
    if (len(token) > 0) {
        tokens = append(tokens, token);
    }
    return tokens
}
