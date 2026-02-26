package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/lornest/whoop-cli/internal/tui/style"
	"github.com/lornest/whoop-cli/internal/util"
	"github.com/lornest/whoop-cli/pkg/whoop"
)

func formatOutput(format, dataType string, data any) error {
	switch format {
	case "json":
		return formatJSON(data)
	case "text":
		return formatText(dataType, data)
	case "table":
		return formatTable(dataType, data)
	default:
		return fmt.Errorf("unknown format: %s (use table, json, or text)", format)
	}
}

func formatJSON(data any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}

func formatText(dataType string, data any) error {
	switch dataType {
	case "recovery":
		return formatRecoveryText(data.([]whoop.Recovery))
	case "sleep":
		return formatSleepText(data.([]whoop.Sleep))
	case "workouts":
		return formatWorkoutsText(data.([]whoop.Workout))
	case "cycles":
		return formatCyclesText(data.([]whoop.Cycle))
	case "profile":
		return formatProfileText(data.(*whoop.BodyMeasurement))
	default:
		return fmt.Errorf("unknown data type: %s", dataType)
	}
}

func formatTable(dataType string, data any) error {
	switch dataType {
	case "recovery":
		return formatRecoveryTable(data.([]whoop.Recovery))
	case "sleep":
		return formatSleepTable(data.([]whoop.Sleep))
	case "workouts":
		return formatWorkoutsTable(data.([]whoop.Workout))
	case "cycles":
		return formatCyclesTable(data.([]whoop.Cycle))
	case "profile":
		return formatProfileTable(data.(*whoop.BodyMeasurement))
	default:
		return fmt.Errorf("unknown data type: %s", dataType)
	}
}

// --- Recovery ---

func formatRecoveryTable(records []whoop.Recovery) error {
	if len(records) == 0 {
		fmt.Println("No recovery data available.")
		return nil
	}

	header := fmt.Sprintf("%-12s %-8s %-10s %-10s %-8s", "Date", "Score", "RHR", "HRV", "SpO2")
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(style.ColorDim)
	fmt.Println(headerStyle.Render(header))
	fmt.Println(headerStyle.Render(strings.Repeat("─", 50)))

	for _, r := range records {
		date := parseDate(r.CreatedAt)
		if r.Score == nil {
			fmt.Printf("%-12s %-8s %-10s %-10s %-8s\n", date, "--", "--", "--", "--")
			continue
		}
		s := r.Score
		scoreStr := util.FormatPercent(s.RecoveryScore)
		color := style.RecoveryColor(s.RecoveryScore)
		colorStyle := lipgloss.NewStyle().Foreground(color)

		fmt.Printf("%-12s %s %-10s %-10s %-8s\n",
			date,
			colorStyle.Render(fmt.Sprintf("%-8s", scoreStr)),
			fmt.Sprintf("%.0f bpm", s.RestingHeartRate),
			fmt.Sprintf("%.1f ms", s.HRVRmssdMilli),
			fmt.Sprintf("%.0f%%", s.SpO2Percentage),
		)
	}
	return nil
}

func formatRecoveryText(records []whoop.Recovery) error {
	for i, r := range records {
		if i > 0 {
			fmt.Println("---")
		}
		fmt.Printf("date: %s\n", parseDate(r.CreatedAt))
		if r.Score == nil {
			fmt.Println("score: pending")
			continue
		}
		s := r.Score
		fmt.Printf("recovery: %.0f%%\n", s.RecoveryScore)
		fmt.Printf("resting_heart_rate: %.0f bpm\n", s.RestingHeartRate)
		fmt.Printf("hrv: %.1f ms\n", s.HRVRmssdMilli)
		fmt.Printf("spo2: %.0f%%\n", s.SpO2Percentage)
		fmt.Printf("skin_temp: %.1f°C\n", s.SkinTempCelsius)
	}
	return nil
}

// --- Sleep ---

func formatSleepTable(records []whoop.Sleep) error {
	if len(records) == 0 {
		fmt.Println("No sleep data available.")
		return nil
	}

	header := fmt.Sprintf("%-12s %-10s %-10s %-10s %-10s %-10s %-10s %-10s %-8s", "Date", "In Bed", "Awake", "Light", "Deep", "REM", "Perf %", "Eff %", "Nap")
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(style.ColorDim)
	fmt.Println(headerStyle.Render(header))
	fmt.Println(headerStyle.Render(strings.Repeat("─", 90)))

	for _, s := range records {
		date := parseDate(s.Start)
		if s.Score == nil {
			fmt.Printf("%-12s %-10s %-10s %-10s %-10s %-10s %-10s %-10s %-8s\n", date, "--", "--", "--", "--", "--", "--", "--", "--")
			continue
		}
		sc := s.Score
		ss := sc.StageSummary
		inBed := util.MillisToDuration(ss.TotalInBedTimeMilli)
		awake := util.MillisToDuration(ss.TotalAwakeTimeMilli)
		light := util.MillisToDuration(ss.TotalLightSleepTimeMilli)
		deep := util.MillisToDuration(ss.TotalSlowWaveSleepTimeMilli)
		rem := util.MillisToDuration(ss.TotalREMSleepTimeMilli)
		perf := util.FormatPercent(sc.SleepPerformancePercentage)
		eff := util.FormatPercent(sc.SleepEfficiencyPercentage)
		nap := " "
		if s.Nap {
			nap = "yes"
		}
		fmt.Printf("%-12s %-10s %-10s %-10s %-10s %-10s %-10s %-10s %-8s\n", date, inBed, awake, light, deep, rem, perf, eff, nap)
	}
	return nil
}

func formatSleepText(records []whoop.Sleep) error {
	for i, s := range records {
		if i > 0 {
			fmt.Println("---")
		}
		fmt.Printf("date: %s\n", parseDate(s.Start))
		fmt.Printf("nap: %v\n", s.Nap)
		if s.Score == nil {
			fmt.Println("score: pending")
			continue
		}
		sc := s.Score
		fmt.Printf("performance: %.0f%%\n", sc.SleepPerformancePercentage)
		fmt.Printf("efficiency: %.0f%%\n", sc.SleepEfficiencyPercentage)
		fmt.Printf("respiratory_rate: %.1f/min\n", sc.RespiratoryRate)
		fmt.Printf("in_bed: %s\n", util.MillisToDuration(sc.StageSummary.TotalInBedTimeMilli))
		fmt.Printf("awake: %s\n", util.MillisToDuration(sc.StageSummary.TotalAwakeTimeMilli))
		fmt.Printf("light: %s\n", util.MillisToDuration(sc.StageSummary.TotalLightSleepTimeMilli))
		fmt.Printf("deep: %s\n", util.MillisToDuration(sc.StageSummary.TotalSlowWaveSleepTimeMilli))
		fmt.Printf("rem: %s\n", util.MillisToDuration(sc.StageSummary.TotalREMSleepTimeMilli))
	}
	return nil
}

// --- Workouts ---

var sportName = map[int]string{
	0: "Activity", 1: "Running", 44: "Weightlifting", 71: "Functional Fitness",
	63: "Cycling", 48: "Swimming", 52: "Yoga", 64: "Walking",
	-1: "Other",
}

func getSportNameCLI(id int) string {
	if name, ok := sportName[id]; ok {
		return name
	}
	return fmt.Sprintf("Sport %d", id)
}

func formatWorkoutsTable(records []whoop.Workout) error {
	if len(records) == 0 {
		fmt.Println("No workouts found.")
		return nil
	}

	header := fmt.Sprintf("%-14s %-20s %-10s %-10s %-10s %-10s", "Date", "Activity", "Strain", "Avg HR", "Max HR", "Duration")
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(style.ColorDim)
	fmt.Println(headerStyle.Render(header))
	fmt.Println(headerStyle.Render(strings.Repeat("─", 76)))

	for _, w := range records {
		date := parseDatetime(w.Start)
		sport := w.SportName
		if sport == "" {
			sport = getSportNameCLI(w.SportID)
		}

		strain, avgHR, maxHR, duration := "--", "--", "--", "--"
		if w.Score != nil {
			strain = util.FormatStrain(w.Score.Strain)
			avgHR = fmt.Sprintf("%d", w.Score.AverageHeartRate)
			maxHR = fmt.Sprintf("%d", w.Score.MaxHeartRate)
		}
		if w.End != "" {
			start, _ := time.Parse("2006-01-02T15:04:05.000Z", w.Start)
			end, _ := time.Parse("2006-01-02T15:04:05.000Z", w.End)
			duration = util.MillisToDuration(int(end.Sub(start).Milliseconds()))
		}

		fmt.Printf("%-14s %-20s %-10s %-10s %-10s %-10s\n", date, sport, strain, avgHR, maxHR, duration)
	}
	return nil
}

func formatWorkoutsText(records []whoop.Workout) error {
	for i, w := range records {
		if i > 0 {
			fmt.Println("---")
		}
		sport := w.SportName
		if sport == "" {
			sport = getSportNameCLI(w.SportID)
		}
		fmt.Printf("date: %s\n", parseDatetime(w.Start))
		fmt.Printf("activity: %s\n", sport)
		if w.Score != nil {
			fmt.Printf("strain: %s\n", util.FormatStrain(w.Score.Strain))
			fmt.Printf("avg_hr: %d bpm\n", w.Score.AverageHeartRate)
			fmt.Printf("max_hr: %d bpm\n", w.Score.MaxHeartRate)
			fmt.Printf("calories: %.0f kcal\n", util.KilojoulesToCalories(w.Score.Kilojoule))
			fmt.Printf("distance: %.2f mi\n", util.MetersToMiles(w.Score.DistanceMeter))
		}
	}
	return nil
}

// --- Cycles ---

func formatCyclesTable(records []whoop.Cycle) error {
	if len(records) == 0 {
		fmt.Println("No cycle data available.")
		return nil
	}

	header := fmt.Sprintf("%-12s %-10s %-12s %-10s %-10s", "Date", "Strain", "Kilojoule", "Avg HR", "Max HR")
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(style.ColorDim)
	fmt.Println(headerStyle.Render(header))
	fmt.Println(headerStyle.Render(strings.Repeat("─", 56)))

	for _, c := range records {
		date := parseDate(c.Start)
		if c.Score == nil {
			fmt.Printf("%-12s %-10s %-12s %-10s %-10s\n", date, "--", "--", "--", "--")
			continue
		}
		s := c.Score
		fmt.Printf("%-12s %-10s %-12.0f %-10d %-10d\n",
			date, util.FormatStrain(s.Strain), s.Kilojoule, s.AverageHeartRate, s.MaxHeartRate)
	}
	return nil
}

func formatCyclesText(records []whoop.Cycle) error {
	for i, c := range records {
		if i > 0 {
			fmt.Println("---")
		}
		fmt.Printf("date: %s\n", parseDate(c.Start))
		fmt.Printf("state: %s\n", c.ScoreState)
		if c.Score != nil {
			fmt.Printf("strain: %s\n", util.FormatStrain(c.Score.Strain))
			fmt.Printf("kilojoule: %.0f\n", c.Score.Kilojoule)
			fmt.Printf("avg_hr: %d bpm\n", c.Score.AverageHeartRate)
			fmt.Printf("max_hr: %d bpm\n", c.Score.MaxHeartRate)
		}
	}
	return nil
}

// --- Profile ---

func formatProfileTable(body *whoop.BodyMeasurement) error {
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(style.ColorDim)
	fmt.Println(headerStyle.Render("Body Measurements"))
	fmt.Println(headerStyle.Render(strings.Repeat("─", 40)))

	fmt.Printf("%-15s %s (%.2f m)\n", "Height:", util.MetersToFeetInches(body.HeightMeter), body.HeightMeter)
	fmt.Printf("%-15s %.0f lbs (%.1f kg)\n", "Weight:", util.KgToLbs(body.WeightKilogram), body.WeightKilogram)
	fmt.Printf("%-15s %d bpm\n", "Max HR:", body.MaxHeartRate)
	return nil
}

func formatProfileText(body *whoop.BodyMeasurement) error {
	fmt.Printf("height_m: %.4f\n", body.HeightMeter)
	fmt.Printf("height: %s\n", util.MetersToFeetInches(body.HeightMeter))
	fmt.Printf("weight_kg: %.1f\n", body.WeightKilogram)
	fmt.Printf("weight_lbs: %.0f\n", util.KgToLbs(body.WeightKilogram))
	fmt.Printf("max_hr: %d\n", body.MaxHeartRate)
	return nil
}

// --- Helpers ---

func parseDate(ts string) string {
	t, err := time.Parse("2006-01-02T15:04:05.000Z", ts)
	if err != nil {
		return ts
	}
	return t.Format("Jan 02")
}

func parseDatetime(ts string) string {
	t, err := time.Parse("2006-01-02T15:04:05.000Z", ts)
	if err != nil {
		return ts
	}
	return t.Format("Jan 02 15:04")
}
