package util

func (suite *Suite) IsValidAttributeName(name string) bool {
	return suite.attrNameValidator([]byte(name)).Best() != nil
}
