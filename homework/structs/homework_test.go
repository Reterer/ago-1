package main

import (
	"math"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

type Option func(*GamePerson)

func WithName(name string) func(*GamePerson) {
	return func(person *GamePerson) {
		l := copy(person.name[:], []byte(name))
		// Решил сохранять длину строки, так как так проще и быстрее вытаскивать строку.
		// К тому же было свободное место в структуре
		person.dataB = setData(person.dataB, dataB(l), dataBNameLenShift, dataBNameLen)
	}
}

func WithCoordinates(x, y, z int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.x = int32(x)
		person.y = int32(y)
		person.z = int32(z)
	}
}

func WithGold(gold int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.gold = uint32(gold)
	}
}

func WithMana(mana int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.dataA = setData(person.dataA, dataA(mana), dataAManaShift, dataAMana)
	}
}

func WithHealth(health int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.dataA = setData(person.dataA, dataA(health), dataAHealthShift, dataAHealth)
	}
}

func WithRespect(respect int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.dataA = setData(person.dataA, dataA(respect), dataARespectShift, dataARespect)
	}
}

func WithStrength(strength int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.dataA = setData(person.dataA, dataA(strength), dataAStrengthShift, dataAStrength)
	}
}

func WithExperience(experience int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.dataA = setData(person.dataA, dataA(experience), dataAExpirienceShift, dataAExpirience)
	}
}

func WithLevel(level int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.dataB = setData(person.dataB, dataB(level), dataBLevelShift, dataBLevel)
	}
}

func WithHouse() func(*GamePerson) {
	return func(person *GamePerson) {
		person.dataB = setData(person.dataB, 1, dataBHouseShift, dataBHouse)
	}
}

func WithGun() func(*GamePerson) {
	return func(person *GamePerson) {
		person.dataB = setData(person.dataB, 1, dataBGunShift, dataBGun)
	}
}

func WithFamily() func(*GamePerson) {
	return func(person *GamePerson) {
		person.dataB = setData(person.dataB, 1, dataBFamilyShift, dataBFamily)
	}
}

func WithType(personType int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.dataB = setData(person.dataB, dataB(personType), dataBTypeShift, dataBHouse)
	}
}

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)

type GamePerson struct {
	name  [42]byte
	dataB dataB // hosuse | gun | family | type | level | nameLen
	dataA dataA // mana | health | respect | strength | expirience
	gold  uint32

	x int32
	y int32
	z int32
}

type dataB uint16

const (
	dataBHouse   dataB = 0b_00000000_00000001
	dataBGun     dataB = 0b_00000000_00000010
	dataBFamily  dataB = 0b_00000000_00000100
	dataBNameLen dataB = 0b_00000011_11110000
	dataBLevel   dataB = 0b_00111100_00000000
	dataBType    dataB = 0b_11000000_00000000

	dataBHouseShift   dataB = 0
	dataBGunShift     dataB = 1
	dataBFamilyShift  dataB = 2
	dataBNameLenShift dataB = 4
	dataBLevelShift   dataB = 10
	dataBTypeShift    dataB = 14
)

type dataA uint32

const (
	dataAMana       dataA = 0b_00000000_00000000_00000011_11111111
	dataAHealth     dataA = 0b_00000000_00001111_11111100_00000000
	dataARespect    dataA = 0b_00000000_11110000_00000000_00000000
	dataAStrength   dataA = 0b_00001111_00000000_00000000_00000000
	dataAExpirience dataA = 0b_11110000_00000000_00000000_00000000

	dataAManaShift       dataA = 0
	dataAHealthShift     dataA = 10
	dataARespectShift    dataA = 20
	dataAStrengthShift   dataA = 24
	dataAExpirienceShift dataA = 28
)

func setData[T dataA | dataB](data T, value T, shift T, mask T) T {
	data |= (value << shift) & mask
	return data
}

func getData[T dataA | dataB](data T, shift T, mask T) T {
	return (data & mask) >> shift
}

func NewGamePerson(options ...Option) GamePerson {
	var gp GamePerson

	for _, opt := range options {
		opt(&gp)
	}

	return gp
}

func (p *GamePerson) Name() string {
	l := getData(p.dataB, dataBNameLenShift, dataBNameLen)
	return string(p.name[:l])
}

func (p *GamePerson) X() int {
	return int(p.x)
}

func (p *GamePerson) Y() int {
	return int(p.y)
}

func (p *GamePerson) Z() int {
	return int(p.z)
}

func (p *GamePerson) Gold() int {
	return int(p.gold)
}

func (p *GamePerson) Mana() int {
	return int(getData(p.dataA, dataAManaShift, dataAMana))
}

func (p *GamePerson) Health() int {
	return int(getData(p.dataA, dataAHealthShift, dataAHealth))
}

func (p *GamePerson) Respect() int {
	return int(getData(p.dataA, dataARespectShift, dataARespect))
}

func (p *GamePerson) Strength() int {
	return int(getData(p.dataA, dataAStrengthShift, dataAStrength))
}

func (p *GamePerson) Experience() int {
	return int(getData(p.dataA, dataAExpirienceShift, dataAExpirience))
}

func (p *GamePerson) Level() int {
	return int(getData(p.dataB, dataBLevelShift, dataBLevel))
}

func (p *GamePerson) HasHouse() bool {
	return getData(p.dataB, dataBHouseShift, dataBHouse) == 1
}

func (p *GamePerson) HasGun() bool {
	return getData(p.dataB, dataBGunShift, dataBGun) == 1
}

func (p *GamePerson) HasFamilty() bool {
	return getData(p.dataB, dataBFamilyShift, dataBFamily) == 1
}

func (p *GamePerson) Type() int {
	return int(getData(p.dataB, dataBTypeShift, dataBType))
}

func TestGamePerson(t *testing.T) {
	assert.LessOrEqual(t, unsafe.Sizeof(GamePerson{}), uintptr(64))

	const x, y, z = math.MinInt32, math.MaxInt32, 0
	const name = "aaaaaaaaaaaaa_bbbbbbbbbbbbb_cccccccccccccc"
	const personType = BuilderGamePersonType
	const gold = math.MaxInt32
	const mana = 1000
	const health = 1000
	const respect = 10
	const strength = 10
	const experience = 10
	const level = 10

	options := []Option{
		WithName(name),
		WithCoordinates(x, y, z),
		WithGold(gold),
		WithMana(mana),
		WithHealth(health),
		WithRespect(respect),
		WithStrength(strength),
		WithExperience(experience),
		WithLevel(level),
		WithHouse(),
		WithFamily(),
		WithType(personType),
	}

	person := NewGamePerson(options...)
	assert.Equal(t, name, person.Name())
	assert.Equal(t, x, person.X())
	assert.Equal(t, y, person.Y())
	assert.Equal(t, z, person.Z())
	assert.Equal(t, gold, person.Gold())
	assert.Equal(t, mana, person.Mana())
	assert.Equal(t, health, person.Health())
	assert.Equal(t, respect, person.Respect())
	assert.Equal(t, strength, person.Strength())
	assert.Equal(t, experience, person.Experience())
	assert.Equal(t, level, person.Level())
	assert.True(t, person.HasHouse())
	assert.True(t, person.HasFamilty())
	assert.False(t, person.HasGun())
	assert.Equal(t, personType, person.Type())
}
