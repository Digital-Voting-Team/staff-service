package helpers

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"staff-service/resources"
	"strconv"
	"time"
)

func MergeErrors(validationErrors ...validation.Errors) validation.Errors {
	result := make(validation.Errors)
	for _, errs := range validationErrors {
		for key, err := range errs {
			result[key] = err
		}
	}
	return result
}

func IsInteger(value interface{}) error {
	if integer, ok := value.(*int64); ok {
		if *integer >= 0 {
			return nil
		}
	}

	if v, ok := value.(*string); ok {
		if integer, err := strconv.Atoi(*v); err == nil {
			if integer >= 0 {
				return nil
			}
			return errors.New("value is less or equal 0")
		}
		return errors.New("value is not a number")
	}

	return errors.New("unknown value type")
}

func IsDate(value interface{}) error {
	if _, ok := value.(*time.Time); ok {
		return nil
	}
	if _, ok := value.(**time.Time); ok {
		return nil
	}
	return errors.New("value is not an valid date")
}

func IsValidAccessLevel(value interface{}) error {
	if _, ok := value.(*resources.AccessLevel); ok {
		return nil
	}
	return errors.New("value is not an valid access level")
}

func IsValidWorkerStatus(value interface{}) error {
	if status, ok := value.(**resources.WorkerStatus); ok {
		if (*status).Validate(string(**status)) {
			return nil
		}
	}
	return errors.New("value is not an valid worker status")
}
