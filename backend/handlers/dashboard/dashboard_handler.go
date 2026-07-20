package dashboard

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
	"pondok-tani-backend/config"
	authmodels "pondok-tani-backend/models/auth"
	coremodels "pondok-tani-backend/models/core"
	financemodels "pondok-tani-backend/models/finance"
	maintenancemodels "pondok-tani-backend/models/maintenance"
	salesmodels "pondok-tani-backend/models/sales"

	"github.com/gin-gonic/gin"
)

type StatCard struct {
	Label      string  `json:"label"`
	Value      string  `json:"value"`
	Change     string  `json:"change"`
	IsPositive bool    `json:"isPositive"`
	Icon       string  `json:"icon"`
}

type MonthlyData struct {
	Month   string  `json:"month"`
	Revenue float64 `json:"revenue"`
	Costs   float64 `json:"costs"`
}

type PlantStatusData struct {
	Name  string `json:"name"`
	Value int64  `json:"value"`
	Color string `json:"color"`
}

type RevenueComparison struct {
	Month     string  `json:"month"`
	Estimated float64 `json:"estimated"`
	Actual    float64 `json:"actual"`
}

type SalesHistory struct {
	ID       string  `json:"id"`
	Date     string  `json:"date"`
	Grade    string  `json:"grade"`
	Quantity float64 `json:"quantity"`
	Revenue  string  `json:"revenue"`
	Buyer    string  `json:"buyer"`
}

type RecentActivity struct {
	Action string `json:"action"`
	Time   string `json:"time"`
	Type   string `json:"type"`
}

type UpcomingMaintenance struct {
	TaskName  string `json:"task_name"`
	DueDate   string `json:"due_date"`
	DaysDiff  int    `json:"days_diff"`
	IsOverdue bool   `json:"is_overdue"`
}

type HarvestStatus struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
	Color  string `json:"color"`
}

type ROIProgress struct {
	CurrentROI   float64 `json:"current_roi"`
	TargetROI    float64 `json:"target_roi"`
	MonthsToGoal int     `json:"months_to_goal"`
}

type DashboardResponse struct {
	Stats               []StatCard            `json:"stats"`
	MonthlyData         []MonthlyData         `json:"monthly_data"`
	PlantStatusData     []PlantStatusData     `json:"plant_status_data"`
	RevenueComparison   []RevenueComparison   `json:"revenue_comparison"`
	SalesHistory        []SalesHistory        `json:"sales_history"`
	RecentActivity      []RecentActivity      `json:"recent_activity"`
	UpcomingMaintenance []UpcomingMaintenance `json:"upcoming_maintenance"`
	HarvestStatus       []HarvestStatus       `json:"harvest_status"`
	ROIProgress         *ROIProgress          `json:"roi_progress"`
}

func getMonthsList() []time.Time {
	months := []time.Time{}
	now := time.Now()
	for i := 5; i >= 0; i-- {
		m := now.AddDate(0, -i, 0)
		months = append(months, m)
	}
	return months
}

func GetAdminDashboard(c *gin.Context) {
	db := config.DB

	var totalUsers int64
	db.Model(&authmodels.User{}).Where("role IN ?", []string{"investor", "mitra"}).Count(&totalUsers)

	var totalPlants int64
	db.Model(&coremodels.PlantBatch{}).Count(&totalPlants)

	var totalInvestment float64
	db.Model(&coremodels.Investment{}).Select("COALESCE(SUM(amount), 0)").Scan(&totalInvestment)

	var totalRevenue float64
	db.Model(&salesmodels.SalesOrder{}).Select("COALESCE(SUM(total_amount), 0)").Scan(&totalRevenue)

	var totalHarvestRequests int64
	db.Table("harvest_requests").Count(&totalHarvestRequests)

	var totalSales int64
	db.Model(&salesmodels.SalesOrder{}).Count(&totalSales)

	var totalOpCosts float64
	db.Model(&financemodels.OperationalCost{}).Select("COALESCE(SUM(amount), 0)").Scan(&totalOpCosts)
	
	roi := 0.0
	if totalOpCosts > 0 {
		roi = (totalRevenue - totalOpCosts) / totalOpCosts * 100
	}

	stats := []StatCard{
		{Label: "Total investasi", Value: formatMoney(totalInvestment), Change: "", IsPositive: true, Icon: "Banknote"},
		{Label: "Investor/Mitra aktif", Value: formatNumber(totalUsers), Change: "", IsPositive: true, Icon: "Users"},
		{Label: "Total batch tanaman", Value: formatNumber(totalPlants), Change: "", IsPositive: true, Icon: "Sprout"},
		{Label: "Permintaan panen", Value: formatNumber(totalHarvestRequests), Change: "", IsPositive: true, Icon: "Sprout"},
		{Label: "Total pesanan", Value: formatNumber(totalSales), Change: "", IsPositive: true, Icon: "TrendingUp"},
		{Label: "ROI rata-rata", Value: fmt.Sprintf("%.1f%%", roi), Change: "", IsPositive: roi >= 0, Icon: "TrendingUp"},
	}

	var plantStatuses []struct {
		Status string
		Count  int64
	}
	db.Model(&coremodels.PlantBatch{}).Select("status, count(*) as count").Group("status").Scan(&plantStatuses)

	statusData := []PlantStatusData{}
	colors := map[string]string{"active": "#3b82f6", "harvested": "#f59e0b", "dead": "#ef4444", "pending": "#6366f1"}
	for _, ps := range plantStatuses {
		color := colors[ps.Status]
		if color == "" {
			color = "#10b981"
		}
		statusData = append(statusData, PlantStatusData{
			Name:  ps.Status,
			Value: ps.Count,
			Color: color,
		})
	}

	// Monthly data
	var investments []coremodels.Investment
	db.Find(&investments)
	var costs []financemodels.OperationalCost
	db.Find(&costs)
	var sims []financemodels.RevenueSimulation
	db.Find(&sims)

	months := getMonthsList()
	monthly := []MonthlyData{}
	revenueComp := []RevenueComparison{}

	for _, m := range months {
		mName := m.Format("Jan")
		var sRev, sCost, sSim float64
		for _, inv := range investments {
			if inv.InvestmentDate.Month() == m.Month() && inv.InvestmentDate.Year() == m.Year() {
				sRev += inv.Amount
			}
		}
		for _, c := range costs {
			if c.CostDate.Month() == m.Month() && c.CostDate.Year() == m.Year() {
				sCost += c.Amount
			}
		}
		for _, s := range sims {
			if s.RevenueDate.Month() == m.Month() && s.RevenueDate.Year() == m.Year() {
				sSim += s.Amount
			}
		}
		monthly = append(monthly, MonthlyData{
			Month:   mName,
			Revenue: sRev,
			Costs:   sCost,
		})
		revenueComp = append(revenueComp, RevenueComparison{
			Month:     mName,
			Estimated: sSim,
			Actual:    sRev,
		})
	}

	var recentSales []salesmodels.SalesOrder
	db.Order("sales_date desc").Limit(5).Find(&recentSales)
	salesHistory := []SalesHistory{}
	for _, s := range recentSales {
		salesHistory = append(salesHistory, SalesHistory{
			ID:       s.OrderNumber,
			Date:     s.SalesDate.Format("2006-01-02"),
			Grade:    "-",
			Quantity: 0, 
			Revenue:  formatMoney(s.TotalAmount),
			Buyer:    s.Buyer,
		})
	}

	var recentMaintenances []maintenancemodels.MaintenanceActivity
	db.Preload("PlantBatch").Order("updated_at desc").Limit(5).Find(&recentMaintenances)
	recentActivities := []RecentActivity{}
	for _, m := range recentMaintenances {
		actionType := m.ActivityType
		batchCode := "N/A"
		if m.PlantBatch.BatchCode != "" {
			batchCode = m.PlantBatch.BatchCode
		}
		
		action := fmt.Sprintf("Validasi %s - Batch %s", actionType, batchCode)
		
		activityType := "info"
		if m.ValidationStatus == "pending" {
			activityType = "warning"
		} else if m.ValidationStatus == "approved" {
			activityType = "success"
		} else if m.ValidationStatus == "rejected" {
			activityType = "error"
		}

		recentActivities = append(recentActivities, RecentActivity{
			Action: action,
			Time:   m.UpdatedAt.Format("02 Jan 2006 15:04"),
			Type:   activityType,
		})
	}

	c.JSON(http.StatusOK, DashboardResponse{
		Stats:               stats,
		MonthlyData:         monthly,
		PlantStatusData:     statusData,
		RevenueComparison:   revenueComp,
		SalesHistory:        salesHistory,
		RecentActivity:      recentActivities,
		UpcomingMaintenance: []UpcomingMaintenance{},
		HarvestStatus:       []HarvestStatus{},
		ROIProgress:         nil,
	})
}

func GetUserDashboard(c *gin.Context) {
	userId, _ := c.Get("user_id")
	db := config.DB

	var myPlants int64
	db.Model(&coremodels.PlantBatch{}).Joins("JOIN investments ON investments.id = plant_batches.investment_id").Where("investments.user_id = ?", userId).Count(&myPlants)

	var myInvestments float64
	db.Model(&coremodels.Investment{}).Select("COALESCE(SUM(amount), 0)").Where("user_id = ?", userId).Scan(&myInvestments)

	var myRevenue float64
	db.Model(&salesmodels.SalesOrder{}).Select("COALESCE(SUM(total_amount), 0)").Where("buyer = (SELECT name FROM users WHERE id = ?)", userId).Scan(&myRevenue)

	roi := 0.0
	if myInvestments > 0 {
		roi = (myRevenue - myInvestments) / myInvestments * 100
	}

	stats := []StatCard{
		{Label: "Total investasi", Value: formatMoney(myInvestments), Change: "", IsPositive: true, Icon: "Banknote"},
		{Label: "Tanaman aktif", Value: formatNumber(myPlants), Change: "", IsPositive: true, Icon: "Sprout"},
		{Label: "Total pendapatan", Value: formatMoney(myRevenue), Change: "", IsPositive: true, Icon: "Banknote"},
		{Label: "ROI saat ini", Value: fmt.Sprintf("%.1f%%", roi), Change: "", IsPositive: roi >= 0, Icon: "TrendingUp"},
	}

	var plantStatuses []struct {
		Status string
		Count  int64
	}
	db.Model(&coremodels.PlantBatch{}).Joins("JOIN investments ON investments.id = plant_batches.investment_id").Where("investments.user_id = ?", userId).Select("plant_batches.status, count(*) as count").Group("plant_batches.status").Scan(&plantStatuses)

	statusData := []PlantStatusData{}
	colors := map[string]string{"active": "#3b82f6", "harvested": "#f59e0b", "dead": "#ef4444", "pending": "#6366f1"}
	for _, ps := range plantStatuses {
		color := colors[ps.Status]
		if color == "" {
			color = "#10b981"
		}
		statusData = append(statusData, PlantStatusData{
			Name:  ps.Status,
			Value: ps.Count,
			Color: color,
		})
	}

	var sales []salesmodels.SalesOrder
	db.Where("buyer = (SELECT name FROM users WHERE id = ?)", userId).Find(&sales)
	var costs []financemodels.OperationalCost
	db.Where("plant_batch_id IN (SELECT plant_batches.id FROM plant_batches JOIN investments ON investments.id = plant_batches.investment_id WHERE investments.user_id = ?)", userId).Find(&costs)
	var sims []financemodels.RevenueSimulation
	db.Where("plant_batch_id IN (SELECT plant_batches.id FROM plant_batches JOIN investments ON investments.id = plant_batches.investment_id WHERE investments.user_id = ?)", userId).Find(&sims)

	months := getMonthsList()
	monthly := []MonthlyData{}
	revenueComp := []RevenueComparison{}

	for _, m := range months {
		mName := m.Format("Jan")
		var sRev, sCost, sSim float64
		for _, s := range sales {
			if s.SalesDate.Month() == m.Month() && s.SalesDate.Year() == m.Year() {
				sRev += s.TotalAmount
			}
		}
		for _, c := range costs {
			if c.CostDate.Month() == m.Month() && c.CostDate.Year() == m.Year() {
				sCost += c.Amount
			}
		}
		for _, s := range sims {
			if s.RevenueDate.Month() == m.Month() && s.RevenueDate.Year() == m.Year() {
				sSim += s.Amount
			}
		}
		monthly = append(monthly, MonthlyData{
			Month:   mName,
			Revenue: sRev,
			Costs:   sCost,
		})
		revenueComp = append(revenueComp, RevenueComparison{
			Month:     mName,
			Estimated: sSim,
			Actual:    sRev,
		})
	}

	activities := []RecentActivity{}
	var recentBatches []coremodels.PlantBatch
	db.Joins("JOIN investments ON investments.id = plant_batches.investment_id").Where("investments.user_id = ?", userId).Order("plant_batches.created_at desc").Limit(3).Find(&recentBatches)
	for _, b := range recentBatches {
		activities = append(activities, RecentActivity{
			Action: fmt.Sprintf("Batch %s dibuat", b.BatchCode),
			Time:   b.CreatedAt.Format("02 Jan 2006"),
			Type:   "success",
		})
	}

	userSalesHistory := []SalesHistory{}
	for _, s := range sales {
		userSalesHistory = append(userSalesHistory, SalesHistory{
			ID:       s.OrderNumber,
			Date:     s.SalesDate.Format("2006-01-02"),
			Grade:    "-", 
			Quantity: 0, 
			Revenue:  formatMoney(s.TotalAmount),
			Buyer:    s.Buyer,
		})
	}

	var harvestRequests int64
	db.Table("harvest_requests").Joins("JOIN plant_batches ON plant_batches.id = harvest_requests.plant_batch_id").Joins("JOIN investments ON investments.id = plant_batches.investment_id").Where("investments.user_id = ?", userId).Count(&harvestRequests)

	var countDrying int64
	db.Model(&coremodels.PlantBatch{}).Joins("JOIN investments ON investments.id = plant_batches.investment_id").Where("investments.user_id = ? AND plant_batches.status = ?", userId, "drying").Count(&countDrying)
	
	var countSold int64
	db.Model(&coremodels.PlantBatch{}).Joins("JOIN investments ON investments.id = plant_batches.investment_id").Where("investments.user_id = ? AND plant_batches.status = ?", userId, "harvested").Count(&countSold)

	harvestStatus := []HarvestStatus{
		{Status: "Siap panen", Count: harvestRequests, Color: "#10b981"},
		{Status: "Pengeringan", Count: countDrying, Color: "#f59e0b"},
		{Status: "Terjual", Count: countSold, Color: "#3b82f6"},
	}

	roiProgress := &ROIProgress{
		CurrentROI:   roi,
		TargetROI:    40.0,
		MonthsToGoal: 8,
	}

	var upcomingSchedules []maintenancemodels.MaintenanceSchedule
	db.Preload("PlantBatch").Joins("JOIN plant_batches ON plant_batches.id = maintenance_schedules.plant_batch_id").
		Joins("JOIN investments ON investments.id = plant_batches.investment_id").
		Where("investments.user_id = ? AND maintenance_schedules.status != ?", userId, "completed").
		Order("maintenance_schedules.next_due_date asc").
		Limit(2).Find(&upcomingSchedules)

	upcomingMaintenance := []UpcomingMaintenance{}
	now := time.Now()
	for _, s := range upcomingSchedules {
		daysDiff := int(s.NextDueDate.Sub(now).Hours() / 24)
		isOverdue := daysDiff < 0
		
		taskName := s.ActivityType
		if s.PlantBatch.BatchCode != "" {
			taskName = fmt.Sprintf("%s - %s", s.ActivityType, s.PlantBatch.BatchCode)
		}
		
		upcomingMaintenance = append(upcomingMaintenance, UpcomingMaintenance{
			TaskName:  taskName,
			DueDate:   s.NextDueDate.Format("2006-01-02"),
			DaysDiff:  daysDiff,
			IsOverdue: isOverdue,
		})
	}

	c.JSON(http.StatusOK, DashboardResponse{
		Stats:               stats,
		MonthlyData:         monthly,
		PlantStatusData:     statusData,
		RevenueComparison:   revenueComp,
		SalesHistory:        userSalesHistory,
		RecentActivity:      activities,
		UpcomingMaintenance: upcomingMaintenance,
		HarvestStatus:       harvestStatus,
		ROIProgress:         roiProgress,
	})
}

func formatWithSeparator(n int64) string {
	in := strconv.FormatInt(n, 10)
	var out []byte
	if n < 0 {
		out = append(out, '-')
		in = in[1:]
	}
	
	for i, c := range in {
		if i > 0 && (len(in)-i)%3 == 0 {
			out = append(out, '.')
		}
		out = append(out, byte(c))
	}
	return string(out)
}

func formatMoney(val float64) string {
	return fmt.Sprintf("Rp %s", formatWithSeparator(int64(val)))
}

func formatNumber(val int64) string {
	return formatWithSeparator(val)
}
