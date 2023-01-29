package repositories

import (
	"github.com/abr-ooo/hodhod/connections"
	"github.com/abr-ooo/hodhod/inputs"
	"github.com/abr-ooo/hodhod/models"
	// log "github.com/sirupsen/logrus"
)

func ExistProjectAlertGroup(input inputs.ProjectAlertGroupBody) bool {
	var projectAlertGroup models.ProjectAlertGroup

	if result := connections.DB.Where(models.ProjectAlertGroup{
		Service:   input.Service,
		ProjectID: input.ProjectID,
	}).First(&projectAlertGroup); result.Error != nil {
		return false
	}
	return true
}

func ExistNameProjectAlertGroup(name string, alertGroupID int) bool {
	var projectAlertGroup models.ProjectAlertGroup

	if result := connections.DB.Where(models.ProjectAlertGroup{
		Name:         name,
		AlertGroupID: alertGroupID,
	}).First(&projectAlertGroup); result.Error != nil {
		return false
	}
	return true
}

func GetProjectsAlertGroup(alertGroupID int) (*[]models.ProjectAlertGroup, error) {
	var projectsAlertGroup []models.ProjectAlertGroup

	if result := connections.DB.Find(&projectsAlertGroup, models.ProjectAlertGroup{AlertGroupID: alertGroupID}); result.RowsAffected < 1 {
		return nil, result.Error
	}

	return &projectsAlertGroup, nil
}

func GetProjectAlertGroup(projectAlertGroupID int) (*models.ProjectAlertGroup, error) {
	var projectAlertGroup models.ProjectAlertGroup

	if result := connections.DB.First(&projectAlertGroup, projectAlertGroupID); result.Error != nil {
		return nil, result.Error
	}

	return &projectAlertGroup, nil
}

func CreateProjectAlertGroup(userID int, alertGroupID int, input inputs.ProjectAlertGroupBody) (*models.ProjectAlertGroup, error) {

	projectAlertGroup := models.ProjectAlertGroup{
		Name:         input.Name,
		UserID:       userID,
		ProjectID:    input.ProjectID,
		Service:      input.Service,
		AlertGroupID: alertGroupID,
	}

	if result := connections.DB.Create(&projectAlertGroup); result.Error != nil {
		return nil, result.Error
	}

	return &projectAlertGroup, nil
}

func DeleteProjectAlertGroup(projectAlertGroupID int) error {

	if result := connections.DB.Where("id = ?", projectAlertGroupID).Delete(&models.ProjectAlertGroup{}); result.Error != nil {
		return result.Error
	}

	if result := connections.DB.Where("project_alert_group_id = ?", projectAlertGroupID).Delete(&models.ContactProjectAlertGroup{}); result.Error != nil {
		return result.Error
	}

	return nil
}

// func GetProjectUserID(projectAlertGroupID int) (*int, error) {

// 	userID := 0

// 	_ = connections.DB.Table("alert_groups").Select("alert_groups.user_id").
// 		Joins("join project_alert_groups on alert_groups.id = project_alert_groups.alert_group_id and project_alert_groups.deleted_at is NULL and alert_groups.deleted_at is NULL and project_alert_groups.id = " + strconv.Itoa(projectAlertGroupID)).
// 		Scan(&userID)

// 	fmt.Println(userID)

// 	if userID == 0 {
// 		return nil, errors.New("there isn't any alert group")
// 	}

// 	// query := "select t1.user_id from alert_groups t1 inner join project_alert_groups t2 on t1.id = t2.alert_group_id where t2.id=$1"

// 	// result := connections.DB.Raw(query, projectAlertGroupID)
// 	// if result.Error != nil {
// 	// 	return nil, result.Error
// 	// }
// 	// result.Scan(&userID)

// 	return &userID, nil
// }
