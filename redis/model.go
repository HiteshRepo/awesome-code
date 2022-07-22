package main

type Model struct {
	Key1 string `json:"key1"`
	Key2 Key2   `json:"key2"`
}

type Key2 struct {
	InnerKey1 string `json:"inner-key-1"`
}