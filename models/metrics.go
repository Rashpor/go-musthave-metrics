package models

import (
	"fmt"
)

const (
	Counter = "counter"
	Gauge   = "gauge"
)

type Float64 float64

// NOTE: Не усложняем пример, вводя иерархическую вложенность структур.
// Органичиваясь плоской моделью.
// Delta и Value объявлены через указатели,
// что бы отличать значение "0", от не заданного значения
// и соответственно не кодировать в структуру.
type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *Float64 `json:"value,omitempty"`
	Hash  string   `json:"hash,omitempty"`
}

/*
func (f Float64) MarshalJSON() ([]byte, error) {
	// Используем 'g', 17, 64 для полной точности
	return []byte(fmt.Sprintf("%.17g", f)), nil
}*/

func (f *Float64) MarshalJSON() ([]byte, error) {
	if f == nil {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("%.17g", *f)), nil
}
