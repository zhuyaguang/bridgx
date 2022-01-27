package request

import (
	"github.com/galaxy-future/BridgX/internal/model"
	"github.com/galaxy-future/BridgX/internal/service"
)

type TagRequest struct {
	ClusterName string            `json:"cluster_name" binding:"required,min=1"`
	Tags        map[string]string `json:"tags" binding:"required"`
}

type GetTagsRequest struct {
	ClusterName string `json:"cluster_name" binding:"required,min=1" form:"cluster_name"`
	TagKey      string `json:"tag_key" form:"tag_key"`
	PageNumber  int    `json:"page_number" form:"page_number"`
	PageSize    int    `json:"page_size" form:"page_size"`
}

type SetExpectInstanceCountRequest struct {
	ClusterName string `json:"cluster_name"`
	ExpectCount int    `json:"expect_count"`
}

type ExpandClusterRequest struct {
	TaskName    string `json:"task_name"`
	ClusterName string `json:"cluster_name" binding:"required"`
	Count       int    `json:"count" binding:"required,min=1,max=10000"`
}

type ShrinkClusterRequest struct {
	TaskName    string   `json:"task_name"`
	ClusterName string   `json:"cluster_name" binding:"required"`
	IPs         []string `json:"ips"`
	Count       int      `json:"count" binding:"required,min=1,max=10000"`
}

type ShrinkAllInstancesRequest struct {
	TaskName    string `json:"task_name"`
	ClusterName string `json:"cluster_name" binding:"required"`
}

type CreateVpcRequest struct {
	Provider  string `json:"provider"`
	RegionId  string `json:"region_id"`
	VpcName   string `json:"vpc_name"`
	CidrBlock string `json:"cidr_block"`
	Ak        string `json:"ak"`
}

func (c *CreateVpcRequest) Check() bool {
	return c.Provider != "" && c.RegionId != "" && c.VpcName != "" && c.Ak != ""
}

type DescribeVpcRequest struct {
	Provider   string `form:"provider" binding:"required"`
	RegionId   string `form:"region_id" binding:"required"`
	VpcName    string `form:"vpc_name"`
	AccountKey string `form:"account_key" binding:"required"`
}

type CreateSwitchRequest struct {
	SwitchName string `json:"switch_name"`
	RegionId   string `json:"region_id"`
	VpcId      string `json:"vpc_id"`
	CidrBlock  string `json:"cidr_block"`
	GatewayIp  string `json:"gateway_ip"`
	ZoneId     string `json:"zone_id"`
}

func (c *CreateSwitchRequest) Check() bool {
	return c.SwitchName != "" && c.VpcId != "" && c.CidrBlock != "" && c.ZoneId != ""
}

type CreateSecurityGroupRequest struct {
	VpcId             string `json:"vpc_id"`
	RegionId          string `json:"region_id"`
	SecurityGroupName string `json:"security_group_name"`
	SecurityGroupType string `json:"security_group_type"`
}

func (c *CreateSecurityGroupRequest) Check() bool {
	return c.SecurityGroupName != "" && c.RegionId != "" && c.VpcId != ""
}

type AddSecurityGroupRuleRequest struct {
	VpcId           string              `json:"vpc_id"`
	RegionId        string              `json:"region_id"`
	SecurityGroupId string              `json:"security_group_id"`
	Rules           []service.GroupRule `json:"rules"`
}

func (c *AddSecurityGroupRuleRequest) Check() bool {
	return c.VpcId != "" && c.RegionId != "" && c.SecurityGroupId != "" &&
		len(c.Rules) > 0
}

type CreateSecurityGroupWithRuleRequest struct {
	VpcId             string              `json:"vpc_id"`
	RegionId          string              `json:"region_id"`
	SecurityGroupName string              `json:"security_group_name"`
	SecurityGroupType string              `json:"security_group_type"`
	Rules             []service.GroupRule `json:"rules"`
}

func (c *CreateSecurityGroupWithRuleRequest) Check() bool {
	return c.SecurityGroupName != "" && c.RegionId != "" && c.VpcId != ""
}

type CreateNetworkRequest struct {
	Provider          string `json:"provider" binding:"required,mustIn=cloud"`
	RegionId          string `json:"region_id" binding:"required"`
	CidrBlock         string `json:"cidr_block" binding:"required"`
	VpcName           string `json:"vpc_name" binding:"required"`
	ZoneId            string `json:"zone_id" binding:"required"`
	SwitchCidrBlock   string `json:"switch_cidr_block" binding:"required"`
	SwitchName        string `json:"switch_name" binding:"required"`
	SecurityGroupName string `json:"security_group_name" binding:"required"`
	SecurityGroupType string `json:"security_group_type"`
	Ak                string `json:"ak" binding:"required"`
}

type LoginRequest struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

type CreateCloudAccountRequest struct {
	AccountName   string `json:"account_name"`
	Provider      string `json:"provider" binding:"required,mustIn=cloud"`
	AccountKey    string `json:"account_key" binding:"required"`
	AccountSecret string `json:"account_secret" binding:"required"`
}

type EditCloudAccountRequest struct {
	AccountId   string `json:"account_id"`
	AccountName string `json:"account_name"`
	Provider    string `json:"provider"`
}

type EditOrgRequest struct {
	OrgId   int64  `json:"org_id" binding:"required"`
	OrgName string `json:"org_name"`
}

type CreateUserRequest struct {
	UserName string  `json:"username"`
	Password string  `json:"password"`
	RoleIds  []int64 `json:"role_ids"`
}

type ModifyAdminPasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type ModifyUsernameRequest struct {
	UserId      string `json:"user_id"`
	NewUsername string `json:"new_username"`
}

type UserStatusRequest struct {
	UserNames []string `json:"usernames"`
	Action    string   `json:"action"`
}

type ModifyUserRequest struct {
	UserId     int64   `json:"user_id" binding:"required"`
	Username   string  `json:"username" binding:"required"`
	UserStatus string  `json:"user_status" binding:"required"`
	RoleIds    []int64 `json:"role_ids"`
}

type CreateOrgRequest struct {
	OrgName  string `json:"org_name"`
	UserName string `json:"username"`
	Password string `json:"password"`
}

type ListClusterByTagsRequest struct {
	Tags       map[string]string `json:"tags" binding:"required"`
	PageNumber int               `json:"page_number"`
	PageSize   int               `json:"page_size"`
}

type SyncInstanceExpireTimeRequest struct {
	ClusterName string `json:"cluster_name" binding:"required"`
}

type CustomPublicCloudClusterRequest struct {
	ClusterName  string                        `json:"name" binding:"required"`
	ClusterDesc  string                        `json:"desc"`
	Provider     string                        `json:"provider" binding:"required,mustIn=cloud"`
	AccountKey   string                        `json:"account_key"`
	InstanceList []model.CustomClusterInstance `json:"instance_list" binding:"required,min=1"`
}

type CustomPrivateCloudClusterRequest struct {
	ClusterName  string                        `json:"name" binding:"required"`
	ClusterDesc  string                        `json:"desc"`
	InstanceList []model.CustomClusterInstance `json:"instance_list" binding:"required,min=1"`
}

type CheckInstanceConnectableRequest struct {
	InstanceList []model.CustomClusterInstance `json:"instance_list"`
}

/*role start*/

type CreateRoleRequest struct {
	Name    string  `json:"name" binding:"required"`   //角色名称
	Code    string  `json:"code" binding:"required"`   //角色编码
	Status  *int8   `json:"status" binding:"required"` //状态  0:禁用 1:启用
	Sort    int     `json:"sort" binding:"required"`   //排序 值越小越靠前
	MenuIds []int64 `json:"menu_ids"`                  //菜单ID列表
}

type UpdateRoleRequest struct {
	Id      int64   `json:"id" binding:"required"`   //角色ID
	Name    string  `json:"name" binding:"required"` //角色名称
	Code    string  `json:"code" binding:"required"` //角色编码
	Sort    int     `json:"sort" binding:"required"` //排序 值越小越靠前
	MenuIds []int64 `json:"menu_ids"`                //菜单ID列表
}

type UpdateRoleStatusRequest struct {
	Id     []int64 `json:"ids" binding:"required"`    //角色ID集合
	Status *int8   `json:"status" binding:"required"` //状态  0:禁用 1:启用
}

type RoleListRequest struct {
	Name       string `form:"name"`   //角色名称
	Status     *int8  `form:"status"` //状态  0:禁用 1:启用
	PageNumber int    `form:"page_number" binding:"required"`
	PageSize   int    `form:"page_size" binding:"required"`
}

/*role end*/

/*menu start*/

type CreateMenuRequest struct {
	ParentId      *int64  `json:"parent_id" binding:"required"`       //父节点ID
	Name          string  `json:"name" binding:"required"`            //菜单名称
	Icon          string  `json:"icon"`                               //图标
	Type          *int8   `json:"type" binding:"required"`            //菜单类型 0:目录 1:菜单 2:按钮
	Path          string  `json:"path"`                               //路径
	Component     string  `json:"component"`                          //组件
	Permission    string  `json:"permission"`                         //权限编码
	Visible       *int8   `json:"visible" binding:"required"`         //是否展示  0:否 1:是
	OuterLinkFlag *int8   `json:"outer_link_flag" binding:"required"` //外链标识  0:否 1:是
	Sort          *int    `json:"sort" binding:"required"`            //排序 值越小越靠前
	ApiIds        []int64 `json:"api_ids"`                            //接口ID列表
}

type UpdateMenuRequest struct {
	Id            int64   `json:"id" binding:"required"`              //角色ID
	ParentId      *int64  `json:"parent_id" binding:"required"`       //父节点ID
	Name          string  `json:"name" binding:"required"`            //菜单名称
	Icon          string  `json:"icon"`                               //图标
	Type          *int8   `json:"type" binding:"required"`            //菜单类型 0:目录 1:菜单 2:按钮
	Path          string  `json:"path"`                               //路径
	Component     string  `json:"component"`                          //组件
	Permission    string  `json:"permission"`                         //权限编码
	Visible       *int8   `json:"visible" binding:"required"`         //是否展示  0:否 1:是
	OuterLinkFlag *int8   `json:"outer_link_flag" binding:"required"` //外链标识  0:否 1:是
	Sort          *int    `json:"sort" binding:"required"`            //排序 值越小越靠前
	ApiIds        []int64 `json:"api_ids"`                            //接口ID列表
}

type MenuListRequest struct {
	Name       string `form:"name"`    //菜单名称
	Visible    *int8  `form:"visible"` //是否展示  0:否 1:是
	PageNumber int    `form:"page_number" binding:"required"`
	PageSize   int    `form:"page_size" binding:"required"`
}

/*menu end*/

/*api start*/

type CreateApiRequest struct {
	Name   string `json:"name" binding:"required"`   //接口名称
	Path   string `json:"path" binding:"required"`   //地址
	Method string `json:"method" binding:"required"` //请求方法
	Status *int8  `json:"status" binding:"required"` //状态  0:禁用 1:启用
}

type UpdateApiRequest struct {
	Id     int64  `json:"id" binding:"required"`     //apiID
	Name   string `json:"name" binding:"required"`   //接口名称
	Path   string `json:"path" binding:"required"`   //地址
	Method string `json:"method" binding:"required"` //请求方法
}

type UpdateApiStatusRequest struct {
	Id     []int64 `json:"ids" binding:"required"`    //API ID集合
	Status *int8   `json:"status" binding:"required"` //状态  0:禁用 1:启用
}

type ApiListRequest struct {
	Name       string `json:"name"`   //接口名称
	Path       string `json:"path"`   //地址
	Method     string `json:"method"` //请求方法
	Status     *int8  `json:"status"` //状态  0:禁用 1:启用
	PageNumber int    `form:"page_number" binding:"required"`
	PageSize   int    `form:"page_size" binding:"required"`
}

/*api end*/
