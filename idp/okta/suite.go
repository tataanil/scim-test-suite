package okta

import (
	"encoding/json"
	"github.com/di-wu/regen"
	"github.com/di-wu/scim-test-suite/util"
)

type TestSuite struct {
	util.Suite
	invalidID   func() string
	randomName  func() string
	randomEmail func() string
}

func (s *TestSuite) IsNumber(i interface{}) {
	s.IsType(json.Number("0"), i)
}

func (s *TestSuite) SetInvalidID(random func() string) {
	s.invalidID = random
}

func (s *TestSuite) InvalidID() string {
	if s.invalidID != nil {
		return s.invalidID()
	}

	gen, _ := regen.New(`\b[0-9a-f]{8}\b-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-\b[0-9a-f]{12}\b`)
	return gen.Generate()
}

func (s *TestSuite) SetRandomName(random func() string) {
	s.randomName = random
}

func (s *TestSuite) RandomName() string {
	if s.randomName != nil {
		return s.randomName()
	}

	gen, _ := regen.New(`^[a-zA-Z0-9]+`)
	return gen.Generate()
}

func (s *TestSuite) SetRandomEmail(random func() string) {
	s.randomEmail = random
}

func (s *TestSuite) RandomEmail() string {
	if s.randomEmail != nil {
		return s.randomEmail()
	}

	gen, _ := regen.New(`^[a-z0-9]+@[a-z0-9]+\.[a-z]{2,4}$`)
	return gen.Generate()
}
