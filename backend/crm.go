package main

import (
	"time"

	"github.com/lib/pq"
)

// ========== 客户相关业务函数 ==========

// getCustomers 获取客户列表
func getCustomers(page, limit int, search string) ([]*CustomerResponse, int64) {
	var customers []Customer
	var total int64

	query := DB.Model(&Customer{})
	if search != "" {
		query = query.Where("name LIKE ? OR contact_name LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	query.Count(&total)
	query.Offset((page - 1) * limit).Limit(limit).Find(&customers)

	responses := make([]*CustomerResponse, len(customers))
	for i, customer := range customers {
		responses[i] = CustomerToResponse(&customer)
	}

	return responses, total
}

// getCustomer 获取单个客户
func getCustomer(id uint64) *CustomerResponse {
	var customer Customer
	DB.First(&customer, id)
	return CustomerToResponse(&customer)
}

// createCustomer 创建客户
func createCustomer(req CustomerRequest) *CustomerResponse {
	customer := &Customer{
		Name:         req.Name,
		ContactName:  req.ContactName,
		Phones:       pq.StringArray(req.Phones),
		Wechats:      pq.StringArray(req.Wechats),
		Province:     req.Province,
		City:         req.City,
		District:     req.District,
		Address:      req.Address,
		Products:     pq.StringArray(req.Products),
		Category:     req.Category,
		Tags:         pq.StringArray(req.Tags),
		State:        req.State,
		Level:        req.Level,
		Source:       req.Source,
		ImportSource: req.ImportSource,
		Remark:       req.Remark,
		SallerName:   req.SallerName,
		Sellers:      pq.Int64Array(req.Sellers),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	DB.Create(customer)
	return CustomerToResponse(customer)
}

// updateCustomer 更新客户
func updateCustomer(id uint64, req CustomerRequest) *CustomerResponse {
	var customer Customer
	DB.First(&customer, id)

	customer.Name = req.Name
	customer.ContactName = req.ContactName
	customer.Phones = pq.StringArray(req.Phones)
	customer.Wechats = pq.StringArray(req.Wechats)
	customer.Province = req.Province
	customer.City = req.City
	customer.District = req.District
	customer.Address = req.Address
	customer.Products = pq.StringArray(req.Products)
	customer.Category = req.Category
	customer.Tags = pq.StringArray(req.Tags)
	customer.State = req.State
	customer.Level = req.Level
	customer.Source = req.Source
	customer.ImportSource = req.ImportSource
	customer.Remark = req.Remark
	customer.SallerName = req.SallerName
	customer.Sellers = pq.Int64Array(req.Sellers)
	customer.UpdatedAt = time.Now()

	DB.Save(&customer)
	return CustomerToResponse(&customer)
}

// deleteCustomer 删除客户
func deleteCustomer(id uint64) {
	DB.Delete(&Customer{}, id)
}

// ========== 待办事项相关业务函数 ==========

// getTodos 获取待办事项列表
func getTodos(customerID uint64, page, pageSize int) ([]TodoResponse, int64) {
	var todos []Todo
	var total int64

	query := DB.Model(&Todo{}).Preload("Customer").Preload("Creator").Preload("Executor")
	if customerID > 0 {
		query = query.Where("customer_id = ?", customerID)
	}

	query.Count(&total)
	query.Offset((page - 1) * pageSize).Limit(pageSize).Order("planned_time DESC").Find(&todos)

	responses := make([]TodoResponse, len(todos))
	for i, todo := range todos {
		responses[i] = TodoResponse{
			Todo:         todo,
			CreatorName:  todo.Creator.Name,
			ExecutorName: todo.Executor.Name,
			CustomerName: todo.Customer.Name,
			IsOverdue:    todo.IsOverdue(),
			DaysLeft:     todo.GetDaysLeft(),
		}
	}

	return responses, total
}

// createTodo 创建待办事项
func createTodo(req TodoCreateRequest) *TodoResponse {
	todo := &Todo{
		CustomerID:     req.CustomerID,
		CreatorID:      1, // TODO: 从上下文获取当前用户ID
		ExecutorID:     req.ExecutorID,
		Title:          req.Title,
		Content:        req.Content,
		PlannedTime:    req.PlannedTime,
		IsReminder:     req.IsReminder,
		ReminderType:   req.ReminderType,
		ReminderUserID: req.ReminderUserID,
		ReminderTime:   req.ReminderTime,
		Priority:       req.Priority,
		Tags:           req.Tags,
	}

	DB.Create(todo)
	DB.Preload("Customer").Preload("Creator").Preload("Executor").First(todo, todo.ID)

	return &TodoResponse{
		Todo:         *todo,
		CreatorName:  todo.Creator.Name,
		ExecutorName: todo.Executor.Name,
		CustomerName: todo.Customer.Name,
		IsOverdue:    todo.IsOverdue(),
		DaysLeft:     todo.GetDaysLeft(),
	}
}

// updateTodo 更新待办事项
func updateTodo(id uint64, req TodoUpdateRequest) *TodoResponse {
	var todo Todo
	DB.Preload("Customer").Preload("Creator").Preload("Executor").First(&todo, id)

	if req.Title != nil {
		todo.Title = *req.Title
	}
	if req.Content != nil {
		todo.Content = *req.Content
	}
	if req.Status != nil {
		todo.Status = *req.Status
	}
	if req.PlannedTime != nil {
		todo.PlannedTime = *req.PlannedTime
	}
	if req.Priority != nil {
		todo.Priority = *req.Priority
	}

	DB.Save(&todo)

	return &TodoResponse{
		Todo:         todo,
		CreatorName:  todo.Creator.Name,
		ExecutorName: todo.Executor.Name,
		CustomerName: todo.Customer.Name,
		IsOverdue:    todo.IsOverdue(),
		DaysLeft:     todo.GetDaysLeft(),
	}
}

// ========== 跟进记录相关业务函数 ==========

// getActivities 获取跟进记录列表
func getActivities(customerID uint64, page, pageSize int) ([]ActivityResponse, int64) {
	var activities []Activity
	var total int64

	query := DB.Model(&Activity{}).Preload("Customer").Preload("User")
	if customerID > 0 {
		query = query.Where("customer_id = ?", customerID)
	}

	query.Count(&total)
	query.Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&activities)

	responses := make([]ActivityResponse, len(activities))
	for i, activity := range activities {
		responses[i] = ActivityResponse{
			Activity:     activity,
			UserName:     activity.User.Name,
			CustomerName: activity.Customer.Name,
			TimeAgo:      activity.GetTimeAgo(),
		}

		// 从Data中提取字段
		if activity.Data != nil {
			if content, ok := activity.Data["content"].(string); ok {
				responses[i].Content = content
			}
			if result, ok := activity.Data["result"].(string); ok {
				responses[i].Result = result
			}
			if amount, ok := activity.Data["amount"].(float64); ok {
				responses[i].Amount = amount
			}
		}
	}

	return responses, total
}

// createActivity 创建跟进记录
func createActivity(req ActivityCreateRequest) *ActivityResponse {
	activity := &Activity{
		CustomerID:     req.CustomerID,
		UserID:         1, // TODO: 从上下文获取当前用户ID
		Kind:           req.Kind,
		Title:          req.Title,
		Remark:         req.Remark,
		Duration:       req.Duration,
		Location:       req.Location,
		NextFollowTime: req.NextFollowTime,
		Attachments:    req.Attachments,
	}

	// 设置Data字段
	data := make(JSONB)
	data["content"] = req.Content
	data["result"] = req.Result
	data["amount"] = req.Amount
	data["cost"] = req.Cost
	data["feedback"] = req.Feedback
	data["satisfaction"] = req.Satisfaction
	activity.Data = data

	DB.Create(activity)
	DB.Preload("Customer").Preload("User").First(activity, activity.ID)

	return &ActivityResponse{
		Activity:     *activity,
		UserName:     activity.User.Name,
		CustomerName: activity.Customer.Name,
		Content:      req.Content,
		Result:       req.Result,
		Amount:       req.Amount,
		Cost:         req.Cost,
		Feedback:     req.Feedback,
		Satisfaction: req.Satisfaction,
		TimeAgo:      activity.GetTimeAgo(),
	}
}

// ========== 用户相关业务函数 ==========

// getUsers 获取用户列表
func getUsers() []UserResponse {
	var users []User
	DB.Find(&users)

	responses := make([]UserResponse, len(users))
	for i, user := range users {
		responses[i] = UserResponse{
			ID:         user.ID,
			Name:       user.Name,
			Department: user.Department,
			Position:   user.Position,
			Email:      user.Email,
			Phone:      user.Phone,
			Status:     user.Status,
			AvatarURL:  user.AvatarURL,
		}
	}

	return responses
}

// getUserDetail 获取用户详情（智能判断员工/客户身份）
func getUserDetail(id uint64) *UserDetailResponse {
	var user User
	if err := DB.First(&user, id).Error; err != nil {
		return nil
	}

	// 检查是否为员工（ID存在于customers表的sellers字段中）
	var customerCount int64
	DB.Table("customers").Where("? = ANY(sellers)", id).Count(&customerCount)
	isEmployee := customerCount > 0

	var displayInfo string
	if isEmployee {
		// 员工：显示部门+职位
		displayInfo = user.Department
		if user.Position != "" {
			if displayInfo != "" {
				displayInfo += user.Position
			} else {
				displayInfo = user.Position
			}
		}
	} else {
		// 客户：显示所在公司名称（从customers表查找）
		var customer Customer
		if err := DB.Where("name = ?", user.Name).First(&customer).Error; err == nil {
			displayInfo = customer.Name
		}
	}

	// 计算今日跟进数量（待办+活动记录）
	today := time.Now().Format("2006-01-02")
	var todayTodos int64
	var todayActivities int64

	DB.Model(&Todo{}).
		Where("executor_id = ? AND DATE(planned_time) = ? AND status != 'completed'", id, today).
		Count(&todayTodos)

	DB.Model(&Activity{}).
		Where("user_id = ? AND DATE(created_at) = ?", id, today).
		Count(&todayActivities)

	return &UserDetailResponse{
		ID:           user.ID,
		Name:         user.Name,
		DisplayInfo:  displayInfo,
		IsEmployee:   isEmployee,
		TodayRevenue: 0, // 暂时默认为0
		TodayFollows: int(todayTodos + todayActivities),
		AvatarURL:    user.AvatarURL,
	}
}

// searchDashboardData 仪表板数据搜索
func searchDashboardData(req DashboardSearchRequest) ([]DashboardSearchResponse, int64) {
	var todos []Todo
	var total int64

	query := DB.Model(&Todo{}).
		Preload("Customer").
		Where("executor_id = ? AND status != 'completed' AND is_deleted = false", req.UserID)

	// 时间筛选条件：时间和客户状态维度
	switch req.TimeFilter {
	case "今日待跟进":
		today := time.Now().Format("2006-01-02")
		query = query.Where("DATE(planned_time) = ?", today)
	case "近期待跟进":
		tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
		query = query.Where("DATE(planned_time) >= ?", tomorrow)
	case "从未联系":
		query = query.Joins("JOIN customers ON customers.id = todos.customer_id").
			Where("customers.last_called IS NULL")
	case "从未下单":
		query = query.Joins("JOIN customers ON customers.id = todos.customer_id").
			Where("customers.last_order_date IS NULL")
	case "公海":
		query = query.Joins("JOIN customers ON customers.id = todos.customer_id").
			Where("customers.sellers IS NULL OR array_length(customers.sellers, 1) = 0")
	case "不用跟进":
		query = query.Joins("JOIN customers ON customers.id = todos.customer_id").
			Where("customers.remark LIKE '%不用跟进%'")
	case "黑名单":
		query = query.Joins("JOIN customers ON customers.id = todos.customer_id").
			Where("'黑名单' = ANY(customers.tags)")
	}

	// 状态筛选条件：待办状态维度
	switch req.StatusFilter {
	case "全部":
		// 不添加额外条件
	case "待办":
		today := time.Now().Format("2006-01-02")
		query = query.Where("DATE(planned_time) = ?", today)
	case "定期":
		query = query.Joins("JOIN activities ON activities.todo_id = todos.id").
			Where("activities.is_regular = true")
	case "已发样":
		query = query.Where("title LIKE '%已发样%'")
	case "已发货":
		query = query.Where("title LIKE '%已发货%'")
	case "半年未下单":
		sixMonthsAgo := time.Now().AddDate(0, -6, 0).Format("2006-01-02")
		query = query.Joins("JOIN customers ON customers.id = todos.customer_id").
			Where("customers.last_order_date IS NULL OR customers.last_order_date < ?", sixMonthsAgo)
	case "一直未下单":
		query = query.Joins("JOIN customers ON customers.id = todos.customer_id").
			Where("customers.last_order_date IS NULL")
	}

	query.Count(&total)
	query.Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).
		Order("planned_time ASC").Find(&todos)

	// 按客户分组聚合待办内容
	customerTodos := make(map[uint64][]Todo)
	for _, todo := range todos {
		customerTodos[todo.CustomerID] = append(customerTodos[todo.CustomerID], todo)
	}

	var responses []DashboardSearchResponse
	for customerID, todoList := range customerTodos {
		if len(todoList) == 0 {
			continue
		}

		customer := todoList[0].Customer
		var contents []string
		var plannedTime string

		for _, todo := range todoList {
			if todo.Content != "" {
				contents = append(contents, todo.Content)
			}
			if plannedTime == "" {
				plannedTime = todo.PlannedTime.Format("2006-01-02 15:04")
			}
		}

		response := DashboardSearchResponse{
			CustomerID:   customerID,
			ContactName:  customer.ContactName,
			CustomerName: customer.Name,
			Tags:         []string(customer.Tags),
			TodoContents: joinStrings(contents, "，"),
			TodoCount:    len(todoList),
			PlannedTime:  plannedTime,
		}

		if customer.LastCalled != nil {
			response.LastCallTime = customer.LastCalled.Format("2006-01-02")
		}
		if customer.LastOrderDate != nil {
			response.LastOrderTime = customer.LastOrderDate.Format("2006-01-02")
		}

		responses = append(responses, response)
	}

	return responses, total
}

// searchCustomers 客户搜索
func searchCustomers(keyword, systemTagsStr string) []CustomerSearchResponse {
	var customers []Customer

	query := DB.Model(&Customer{})

	// 关键词搜索（客户名称或联系人）
	if keyword != "" {
		query = query.Where("name LIKE ? OR contact_name LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 系统标签搜索（多选OR拼接）
	if systemTagsStr != "" {
		// 解析逗号分隔的标签ID
		tagIDs := parseCommaSeparatedInt64(systemTagsStr)
		if len(tagIDs) > 0 {
			// 构建OR条件：任一标签匹配即可
			var orConditions []string
			var args []interface{}
			for _, tagID := range tagIDs {
				orConditions = append(orConditions, "? = ANY(system_tags)")
				args = append(args, tagID)
			}
			if len(orConditions) > 0 {
				orQuery := "(" + joinStrings(orConditions, " OR ") + ")"
				query = query.Where(orQuery, args...)
			}
		}
	}

	// 查询所有匹配的客户
	query.Order("updated_at DESC").Find(&customers)

	// 转换为响应格式
	responses := make([]CustomerSearchResponse, len(customers))
	for i, customer := range customers {
		response := CustomerSearchResponse{
			ID:          uint64(customer.ID),
			Name:        customer.Name,
			ContactName: customer.ContactName,
			Category:    customer.Category,
			Tags:        []string(customer.Tags),
			SystemTags:  []int64(customer.SystemTags),
			Province:    customer.Province,
			City:        customer.City,
			State:       customer.State,
			Level:       customer.Level,
		}

		// 设置主要电话
		if len(customer.Phones) > 0 {
			response.Phone = customer.Phones[0]
		}

		responses[i] = response
	}

	return responses
}

// ========== 标签相关业务函数 ==========

// getTags 获取标签列表
func getTags() []TagResponse {
	var tags []Tag
	DB.Preload("Dimension").Find(&tags)

	responses := make([]TagResponse, len(tags))
	for i, tag := range tags {
		responses[i] = TagResponse{
			ID:            tag.ID,
			DimensionID:   tag.DimensionID,
			DimensionName: tag.Dimension.Name,
			Name:          tag.Name,
			Color:         tag.Color,
			Description:   tag.Description,
			SortOrder:     tag.SortOrder,
		}
	}

	return responses
}

// createTag 创建标签
func createTag(req TagCreateRequest) *TagResponse {
	tag := &Tag{
		DimensionID: req.DimensionID,
		Name:        req.Name,
		Color:       req.Color,
		Description: req.Description,
		SortOrder:   req.SortOrder,
	}

	DB.Create(tag)
	DB.Preload("Dimension").First(tag, tag.ID)

	return &TagResponse{
		ID:            tag.ID,
		DimensionID:   tag.DimensionID,
		DimensionName: tag.Dimension.Name,
		Name:          tag.Name,
		Color:         tag.Color,
		Description:   tag.Description,
		SortOrder:     tag.SortOrder,
	}
}

// ========== 提醒相关业务函数 ==========

// getReminders 获取提醒列表
func getReminders(userID uint64, page, pageSize int) ([]ReminderResponse, int64) {
	var reminders []Reminder
	var total int64

	query := DB.Model(&Reminder{}).Preload("Todo").Preload("User")
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}

	query.Count(&total)
	query.Offset((page - 1) * pageSize).Limit(pageSize).Order("schedule_time DESC").Find(&reminders)

	responses := make([]ReminderResponse, len(reminders))
	for i, reminder := range reminders {
		responses[i] = ReminderResponse{
			Reminder: reminder,
			UserName: reminder.User.Name,
		}
	}

	return responses, total
}

// createReminder 创建提醒
func createReminder(req ReminderCreateRequest) *ReminderResponse {
	reminder := &Reminder{
		TodoID:       req.TodoID,
		UserID:       req.UserID,
		Type:         req.Type,
		Title:        req.Title,
		Content:      req.Content,
		Frequency:    req.Frequency,
		ScheduleTime: req.ScheduleTime,
		MaxRetries:   req.MaxRetries,
	}

	DB.Create(reminder)
	DB.Preload("Todo").Preload("User").First(reminder, reminder.ID)

	return &ReminderResponse{
		Reminder: *reminder,
		UserName: reminder.User.Name,
	}
}
