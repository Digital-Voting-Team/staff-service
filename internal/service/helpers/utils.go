package helpers

import "net/http"

func GetIdsForGivenUser(r *http.Request, userId int64) (int64, int64, int64, error) {
	resultStaff, err := StaffQ(r).FilterByUserID(userId).Get()
	if err != nil {
		return 0, 0, 0, err
	}
	resultPerson, err := PersonsQ(r).FilterByID(resultStaff.PersonID).Get()
	if err != nil {
		return 0, 0, 0, err
	}
	return resultStaff.ID, resultPerson.ID, resultPerson.AddressID, nil
}
