package models

import (
	"fmt"
)

const (
	Counter = "counter"
	Gauge   = "gauge"
)

// NOTE: Не усложняем пример, вводя иерархическую вложенность структур.
// Органичиваясь плоской моделью.
// Delta и Value объявлены через указатели,
// что бы отличать значение "0", от не заданного значения
// и соответственно не кодировать в структуру.
type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
	Hash  string   `json:"hash,omitempty"`
}

// Float64 — тип-обёртка для float64 с кастомной сериализацией
type Float64 float64

func (f Float64) MarshalJSON() ([]byte, error) {
	// %.17g — печатает до 17 значащих цифр, не добавляя лишнего
	return []byte(fmt.Sprintf("%.17g", f)), nil
}
