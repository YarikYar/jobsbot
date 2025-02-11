package config

import (
    "log"
    "os"

    "gopkg.in/yaml.v3"
)

func LoadConfig(filePath string) map[string]interface{} {
    data, err := os.ReadFile(filePath)
    if err != nil {
        log.Fatalf("Ошибка чтения файла конфигурации: %v", err)
    }

    var result map[string]interface{}
    err = yaml.Unmarshal(data, &result)
    if err != nil {
        log.Fatalf("Ошибка парсинга YAML: %v", err)
    }

    return result
}

