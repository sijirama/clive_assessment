package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
	"worker/db"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
	"gorm.io/gorm"
)

var Store *gorm.DB
var err error
var llm *googleai.GoogleAI

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading the fucking .env file")
	}

	// init langchain with gemini
	ctx := context.Background()
	apiKey := os.Getenv("API_KEY")
	llm, err = googleai.New(ctx, googleai.WithAPIKey(apiKey), googleai.WithDefaultModel("gemini-2.0-flash"))
	if err != nil {
		log.Fatal(err)
	}

	Store, err = db.ConnectToDb()
	if err != nil {
		panic("Failed to connect to db")
	}

	//INFO: Cron stuff
	c := cron.New(cron.WithSeconds())
	cronSeconds := os.Getenv("INSIGHT_INTERVAL_SECONDS")

	// Convert to integer
	seconds, err := strconv.Atoi(cronSeconds)
	if err != nil {
		// Handle invalid number or missing env variable
		fmt.Printf("Invalid INSIGHT_INTERVAL_SECONDS value: %v, defaulting to hourly\n", err)
		c.AddFunc("@hourly", SummariseAndSave)
	} else {
		// Create cron schedule (e.g., "*/60 * * * * *" for every 60 seconds)
		cronSpec := fmt.Sprintf("*/%d * * * * *", seconds/3)
		_, err = c.AddFunc(cronSpec, SummariseAndSave)
		if err != nil {
			// Handle invalid cron spec
			fmt.Printf("Invalid cron schedule: %v, defaulting to hourly\n", err)
			c.AddFunc("@hourly", SummariseAndSave)
		}
	}

	c.Start()

	router := gin.Default()

	router.GET("/insights/latest", LatestReport)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.Run() // listen and serve on 0.0.0.0:8080
}

func LatestReport(c *gin.Context) {
	var report db.InsightReport
	if err := Store.Where("report_date IS NOT NULL").Order("report_date DESC").First(&report).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{
				"error": "No insight report found",
			})
			return
		}
		log.Printf("Failed to read latest insight report: %v", err)
		c.JSON(500, gin.H{
			"error": "Internal server error",
		})
		return
	}

	c.JSON(200, gin.H{
		"report": report,
	})
}

func SummariseAndSave() {
	var invoices []db.Invoice
	if err := Store.Find(&invoices).Error; err != nil {
		log.Fatalf("Failed to read invoices: %v", err)
	}

	var totalSpend float64
	vendorSpend := make(map[string]float64)
	var overdueCount int
	var anomalies []map[string]interface{}

	now := time.Now()

	for _, inv := range invoices {
		totalSpend += inv.Amount
		vendorSpend[inv.Vendor] += inv.Amount
		if !inv.IsPaid && inv.DueDate.Before(now) {
			overdueCount++
		}
		// Simple anomaly detection: flag invoices > 2x average
		if len(invoices) > 0 && inv.Amount > (totalSpend/float64(len(invoices)))*2 {
			anomalies = append(anomalies, map[string]interface{}{
				"invoice_id": inv.ID,
				"amount":     inv.Amount,
				"vendor":     inv.Vendor,
				"reason":     "Amount significantly above average",
			})
		}
	}
	fmt.Println(anomalies)

	// Find largest vendor
	var largestVendor string
	var maxSpend float64
	for vendor, spend := range vendorSpend {
		if spend > maxSpend {
			maxSpend = spend
			largestVendor = vendor
		}
	}

	// Prepare LLM input
	metrics := map[string]interface{}{
		"total_spend":    totalSpend,
		"largest_vendor": largestVendor,
		"overdue_count":  overdueCount,
		"anomalies":      anomalies,
	}
	metricsJSON, err := json.Marshal(metrics)
	if err != nil {
		log.Fatalf("Failed to marshal metrics: %v", err)
	}

	recommendation, err := GenerateSummary(context.Background(), string(metricsJSON))

	if err != nil {
		log.Fatalf("Failed to generate LLM summary: %v", err)
	}

	// Save report to insight_reports
	report := db.InsightReport{
		ID:                       uuid.New().String(),
		TotalSpend:               totalSpend,
		LargestVendor:            largestVendor,
		OverdueCount:             overdueCount,
		Anomalies:                db.JSONBMap{"anomalies": anomalies},
		CostSavingRecommendation: recommendation,
		ReportDate:               now,
	}
	if err := Store.Create(&report).Error; err != nil {
		log.Fatalf("Failed to save insight report: %v", err)
	}

	log.Println("Insight report generated and saved successfully")
}

func GenerateSummary(ctx context.Context, input string) (string, error) {
	prompt := fmt.Sprintf("Based on the data: %s, consider negotiating bulk discounts with the largest vendor.", input)
	return llms.GenerateFromSinglePrompt(ctx, llm, prompt)
}
