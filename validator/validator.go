package validator

type Validator interface {
	Validate() (bool, string)
	FieldName() string
}

func ValidateAll(validators []Validator) (bool, map[string]string) {
	errors := map[string]string{}

	for _, eachValidator := range validators {
		valid, err := eachValidator.Validate()
		if !valid {
			errors[eachValidator.FieldName()] = err
		}
	}

	if len(errors) > 0 {
		return false, errors
	}

	return true, nil
}
