package models

type Card struct {
	id         int
	name       string
	modifierId int
}

type Class struct {
	id     int
	name   string
	health uint
}

type Modifier struct {
	id     int
	damage int
	buff   bool
	debuff bool
}
