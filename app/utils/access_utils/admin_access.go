package access_utils

import (
	"app/models"
	"app/utils"
	"fmt"
	"net/http"
)

// Permit request only for Admin users
func IsCrtUserAdmin(crtUser *models.AuthUser, respWr http.ResponseWriter) bool {
	fmt.Println("[utils] Current User is:", crtUser.Role)
	if crtUser.Role != models.UserRole_Admin {
		fmt.Println("[utils] Unauthorized request")
		utils.SendErrorResponse(respWr, 401, "Unauthorized")
		return false
	}
	return true
}
