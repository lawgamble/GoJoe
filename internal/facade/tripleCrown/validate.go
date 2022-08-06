package tripleCrown

import "goJoe/internal/service"

func validateRegistration(id string) bool {
	userData := service.GetUserData(id)
	if userData.Found == true {
		if userData.UserReg[0].Pcl == "Y" {
			if userData.UserReg[0].OculusName != "" {
				TCUserName = userData.UserReg[0].OculusName
				return true
			}
		}
	}
	return false
}
