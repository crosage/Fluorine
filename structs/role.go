package structs

const (
	// Permissions
	ViewSalary       = 2
	CheckIn          = 3
	ApplyLeave       = 4
	ReClock          = 5
	EditInfo         = 6
	ApproveLeave     = 7
	ManageSalary     = 8
	ManagePermission = 1
)

type Role struct {
	ID          int
	Name        string
	Permissions []Permission
}

var roles = []Role{
	{
		ID:   1,
		Name: "普通员工",
		Permissions: []Permission{
			{ID: ViewSalary, Name: "查看薪水", Description: "查看自己薪水"},
			{ID: CheckIn, Name: "签到打卡", Description: "签到"},
			{ID: ApplyLeave, Name: "申请请假", Description: "申请请假"},
		},
	},
	{
		ID:   2,
		Name: "部门组长",
		Permissions: []Permission{
			{ID: ReClock, Name: "补打卡", Description: "为员工补签"},
			{ID: EditInfo, Name: "编辑信息", Description: "编辑员工信息"},
			{ID: ApproveLeave, Name: "审批请假", Description: "审批员工请假"},
		},
	},
	{
		ID:   3,
		Name: "人事管理",
		Permissions: []Permission{
			{ID: EditInfo, Name: "编辑信息", Description: "编辑员工信息"},
			{ID: ManageSalary, Name: "管理薪水", Description: "管理员工薪资"},
		},
	},
	{
		ID:   9,
		Name: "总经理",
		Permissions: []Permission{
			{ID: ManagePermission, Name: "管理权限", Description: "管理和赋予其他用户权限"},
		},
	},
}

func GetRoleByName(name string) *Role {
	for _, role := range roles {
		if role.Name == name {
			return &role
		}
	}
	return nil
}
