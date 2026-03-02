package web

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"gorm.io/gorm"

	"ppl-study-planner/internal/model"
	"ppl-study-planner/internal/services"
)

type server struct {
	db *gorm.DB
}

type pageData struct {
	Title string
	Body  template.HTML
}

func Run(db *gorm.DB, host string, port int) error {
	s := &server{db: db}

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.dashboard)
	mux.HandleFunc("/study", s.study)
	mux.HandleFunc("/study/toggle", s.studyToggle)
	mux.HandleFunc("/budget", s.budget)
	mux.HandleFunc("/budget/update", s.budgetUpdate)
	mux.HandleFunc("/checklist", s.checklist)
	mux.HandleFunc("/checklist/toggle", s.checklistToggle)

	bindAddr := fmt.Sprintf("%s:%d", host, port)
	url := browserURL(host, port)
	fmt.Printf("Web UI starting on http://%s\n", bindAddr)
	fmt.Printf("Opening browser at %s\n", url)
	go func() {
		time.Sleep(350 * time.Millisecond)
		if err := openBrowser(url); err != nil {
			fmt.Printf("Browser open failed: %v\n", err)
		}
	}()

	httpServer := &http.Server{Addr: bindAddr, Handler: mux}
	return httpServer.ListenAndServe()
}

func (s *server) dashboard(w http.ResponseWriter, r *http.Request) {
	var tasks []model.DailyTask
	s.db.Find(&tasks)
	completed, total, percentage := services.CalculateProgress(tasks)
	body := template.HTML(fmt.Sprintf(`
<p><strong>Progress:</strong> %d/%d completed (%.1f%%)</p>
<ul>
  <li><a href="/study">Study tasks</a></li>
  <li><a href="/budget">Budget planner</a></li>
  <li><a href="/checklist">Checkride checklist</a></li>
</ul>
`, completed, total, percentage))
	renderPage(w, pageData{Title: "Dashboard", Body: body})
}

func (s *server) study(w http.ResponseWriter, r *http.Request) {
	var tasks []model.DailyTask
	s.db.Order("date asc, id asc").Find(&tasks)

	body := "<h3>Study Tasks</h3><table><tr><th>Date</th><th>Code</th><th>Description</th><th>Done</th><th></th></tr>"
	for _, t := range tasks {
		checked := ""
		if t.Completed {
			checked = "checked"
		}
		body += fmt.Sprintf(`<tr><td>%s</td><td>%s</td><td>%s</td><td><input type="checkbox" disabled %s></td><td>
<form method="POST" action="/study/toggle"><input type="hidden" name="id" value="%d"><button type="submit">Toggle</button></form>
</td></tr>`, t.Date.Format("2006-01-02"), template.HTMLEscapeString(extractCode(t.Title)), template.HTMLEscapeString(extractDesc(t.Title)), checked, t.ID)
	}
	body += "</table>"
	renderPage(w, pageData{Title: "Study", Body: template.HTML(body)})
}

func (s *server) studyToggle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/study", http.StatusSeeOther)
		return
	}
	id, _ := strconv.Atoi(r.FormValue("id"))
	if id > 0 {
		var task model.DailyTask
		if err := s.db.First(&task, id).Error; err == nil {
			task.Completed = !task.Completed
			s.db.Save(&task)
		}
	}
	http.Redirect(w, r, "/study", http.StatusSeeOther)
}

func (s *server) budget(w http.ResponseWriter, r *http.Request) {
	planeRate := readBudget(s.db, model.BudgetPlaneRate, 150)
	cfiRate := readBudget(s.db, model.BudgetCfiRate, 60)
	living := readBudget(s.db, model.BudgetLiving, 500)
	limit := readBudget(s.db, model.BudgetLimit, 10000)

	body := fmt.Sprintf(`
<h3>Budget</h3>
<form method="POST" action="/budget/update">
  <label>Plane rate: <input name="plane_rate" value="%.2f"></label><br>
  <label>CFI rate: <input name="cfi_rate" value="%.2f"></label><br>
  <label>Living cost: <input name="living" value="%.2f"></label><br>
  <label>Budget limit: <input name="budget_limit" value="%.2f"></label><br>
  <button type="submit">Save</button>
</form>
`, planeRate, cfiRate, living, limit)
	renderPage(w, pageData{Title: "Budget", Body: template.HTML(body)})
}

func (s *server) budgetUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/budget", http.StatusSeeOther)
		return
	}

	if v, err := strconv.ParseFloat(r.FormValue("plane_rate"), 64); err == nil {
		upsertBudget(s.db, model.BudgetPlaneRate, v)
	}
	if v, err := strconv.ParseFloat(r.FormValue("cfi_rate"), 64); err == nil {
		upsertBudget(s.db, model.BudgetCfiRate, v)
	}
	if v, err := strconv.ParseFloat(r.FormValue("living"), 64); err == nil {
		upsertBudget(s.db, model.BudgetLiving, v)
	}
	if v, err := strconv.ParseFloat(r.FormValue("budget_limit"), 64); err == nil {
		upsertBudget(s.db, model.BudgetLimit, v)
	}

	http.Redirect(w, r, "/budget", http.StatusSeeOther)
}

func (s *server) checklist(w http.ResponseWriter, r *http.Request) {
	var items []model.ChecklistItem
	s.db.Order("id asc").Find(&items)
	body := "<h3>Checkride Checklist</h3><table><tr><th>Category</th><th>Item</th><th>Done</th><th></th></tr>"
	for _, item := range items {
		checked := ""
		if item.Completed {
			checked = "checked"
		}
		body += fmt.Sprintf(`<tr><td>%s</td><td>%s</td><td><input type="checkbox" disabled %s></td><td>
<form method="POST" action="/checklist/toggle"><input type="hidden" name="id" value="%d"><button type="submit">Toggle</button></form>
</td></tr>`, template.HTMLEscapeString(string(item.Category)), template.HTMLEscapeString(item.Title), checked, item.ID)
	}
	body += "</table>"
	renderPage(w, pageData{Title: "Checklist", Body: template.HTML(body)})
}

func (s *server) checklistToggle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/checklist", http.StatusSeeOther)
		return
	}
	id, _ := strconv.Atoi(r.FormValue("id"))
	if id > 0 {
		var item model.ChecklistItem
		if err := s.db.First(&item, id).Error; err == nil {
			item.Completed = !item.Completed
			s.db.Save(&item)
		}
	}
	http.Redirect(w, r, "/checklist", http.StatusSeeOther)
}

func renderPage(w http.ResponseWriter, p pageData) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	layout := `<!doctype html>
<html><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1">
<title>{{.Title}} - openppl</title>
<style>body{font-family:ui-sans-serif,system-ui;padding:18px;max-width:1100px;margin:0 auto}nav a{margin-right:12px}table{border-collapse:collapse;width:100%}th,td{border:1px solid #ddd;padding:8px;text-align:left}button{padding:4px 8px}</style>
</head><body>
<h1>openppl web</h1>
<nav><a href="/">Dashboard</a><a href="/study">Study</a><a href="/budget">Budget</a><a href="/checklist">Checklist</a></nav>
<hr>
{{.Body}}
</body></html>`
	t := template.Must(template.New("layout").Parse(layout))
	_ = t.Execute(w, p)
}

func readBudget(db *gorm.DB, itemType model.BudgetItemType, fallback float64) float64 {
	var budget model.Budget
	if err := db.Where("item_type = ?", itemType).Order("id asc").Limit(1).Find(&budget).Error; err != nil {
		return fallback
	}
	if budget.ID == 0 {
		return fallback
	}
	return budget.Amount
}

func upsertBudget(db *gorm.DB, itemType model.BudgetItemType, amount float64) {
	var budget model.Budget
	if err := db.Where("item_type = ?", itemType).Order("id asc").Limit(1).Find(&budget).Error; err != nil {
		return
	}
	if budget.ID == 0 {
		_ = db.Create(&model.Budget{ItemType: itemType, Amount: amount}).Error
		return
	}
	budget.Amount = amount
	_ = db.Save(&budget).Error
}

func extractCode(title string) string {
	for i, c := range title {
		if c == ' ' {
			return title[:i]
		}
	}
	return title
}

func extractDesc(title string) string {
	for i, c := range title {
		if c == ' ' {
			if i+1 < len(title) {
				return title[i+1:]
			}
			break
		}
	}
	return title
}
