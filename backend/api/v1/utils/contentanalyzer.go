package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"
	"unicode"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// ContentAnalyzer represents the main structure for content analysis
type ContentAnalyzer struct {
	config        *ContentConfig
	logger        *zap.Logger
	prohibitedRe  *regexp.Regexp
	suspiciousRes []*regexp.Regexp
	urlRe         *regexp.Regexp
	emailRe       *regexp.Regexp
}

// ContentConfig holds all the configuration parameters for the content analyzer
type ContentConfig struct {
	ProhibitedWords     []string          `json:"prohibited_words"`
	SuspiciousPatterns  []string          `json:"suspicious_patterns"`
	ImageExtensions     []string          `json:"image_extensions"`
	MaxConsecutiveCaps  int               `json:"max_consecutive_caps"`
	MaxExclamationMarks int               `json:"max_exclamation_marks"`
	MaxLinks            int               `json:"max_links"`
	MaxEmails           int               `json:"max_emails"`
	SpamTriggers        map[string]int    `json:"spam_triggers"`
	ContentCategories   map[string][]word `json:"content_categories"`
}

// word represents a word and its weight in a specific category
type word struct {
	Text   string  `json:"text"`
	Weight float64 `json:"weight"`
}

// AnalysisResult represents the outcome of content analysis
type AnalysisResult struct {
	IsSafe             bool
	Message            string
	SpamScore          float64
	CategoryScores     map[string]float64
	SuspiciousPatterns []string
	LinkCount          int
	EmailCount         int
}

// NewContentAnalyzer creates a new instance of ContentAnalyzer
func NewContentAnalyzer(configPath string, logger *zap.Logger) (*ContentAnalyzer, error) {
	// Load configuration from file
	config, err := loadConfig(configPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load config")
	}

	// Compile regular expressions for prohibited words
	prohibitedRe := regexp.MustCompile(fmt.Sprintf(`(?i)\b(%s)\b`, strings.Join(config.ProhibitedWords, "|")))

	// Compile regular expressions for suspicious patterns
	var suspiciousRes []*regexp.Regexp
	for _, pattern := range config.SuspiciousPatterns {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid suspicious pattern: %s", pattern)
		}
		suspiciousRes = append(suspiciousRes, re)
	}

	// Compile regular expressions for URLs and email addresses
	urlRe := regexp.MustCompile(`https?://\S+`)
	emailRe := regexp.MustCompile(`\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`)

	return &ContentAnalyzer{
		config:        config,
		logger:        logger,
		prohibitedRe:  prohibitedRe,
		suspiciousRes: suspiciousRes,
		urlRe:         urlRe,
		emailRe:       emailRe,
	}, nil
}

// loadConfig reads and parses the JSON configuration file
func loadConfig(path string) (*ContentConfig, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config ContentConfig
	if err := json.Unmarshal(file, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// AnalyzeContent performs content analysis on the given text and attachments
func (ca *ContentAnalyzer) AnalyzeContent(ctx context.Context, text string, attachments []string) (AnalysisResult, error) {
	var (
		wg     sync.WaitGroup
		result AnalysisResult
		mu     sync.Mutex
	)

	result.CategoryScores = make(map[string]float64)

	// Analyze text content
	wg.Add(1)
	go func() {
		defer wg.Done()
		ca.analyzeText(text, &result, &mu)
	}()

	// Analyze attachments
	wg.Add(1)
	go func() {
		defer wg.Done()
		ca.analyzeAttachments(attachments, &result, &mu)
	}()

	// Wait for all analyses to complete
	wg.Wait()

	// Determine if the content is safe based on the analysis results
	result.IsSafe = result.SpamScore < 0.7 && len(result.SuspiciousPatterns) == 0

	// Generate a summary message
	result.Message = ca.generateSummaryMessage(&result)

	return result, nil
}

// analyzeText performs various checks on the text content
func (ca *ContentAnalyzer) analyzeText(text string, result *AnalysisResult, mu *sync.Mutex) {
	// Check for prohibited words
	if ca.prohibitedRe.MatchString(text) {
		mu.Lock()
		result.SpamScore += 0.5
		mu.Unlock()
	}

	// Check for suspicious patterns
	for _, re := range ca.suspiciousRes {
		if re.MatchString(text) {
			mu.Lock()
			result.SuspiciousPatterns = append(result.SuspiciousPatterns, re.String())
			mu.Unlock()
		}
	}

	// Check for excessive capitalization
	if ca.containsExcessiveCaps(text) {
		mu.Lock()
		result.SpamScore += 0.2
		mu.Unlock()
	}

	// Check for excessive exclamation marks
	if ca.containsExcessiveExclamations(text) {
		mu.Lock()
		result.SpamScore += 0.1
		mu.Unlock()
	}

	// Count links and emails
	links := ca.urlRe.FindAllString(text, -1)
	emails := ca.emailRe.FindAllString(text, -1)

	mu.Lock()
	result.LinkCount = len(links)
	result.EmailCount = len(emails)
	mu.Unlock()

	if len(links) > ca.config.MaxLinks {
		mu.Lock()
		result.SpamScore += 0.3
		mu.Unlock()
	}

	if len(emails) > ca.config.MaxEmails {
		mu.Lock()
		result.SpamScore += 0.3
		mu.Unlock()
	}

	// Check for spam triggers
	ca.checkSpamTriggers(text, result, mu)

	// Analyze content categories
	ca.analyzeContentCategories(text, result, mu)
}

// analyzeAttachments checks for suspicious attachments
func (ca *ContentAnalyzer) analyzeAttachments(attachments []string, result *AnalysisResult, mu *sync.Mutex) {
	for _, filename := range attachments {
		for _, ext := range ca.config.ImageExtensions {
			if strings.HasSuffix(strings.ToLower(filename), ext) {
				mu.Lock()
				result.SpamScore += 0.1
				mu.Unlock()
				break
			}
		}
	}
}

// containsExcessiveCaps checks for excessive use of capital letters
func (ca *ContentAnalyzer) containsExcessiveCaps(text string) bool {
	capsCount := 0
	for _, char := range text {
		if unicode.IsUpper(char) {
			capsCount++
			if capsCount > ca.config.MaxConsecutiveCaps {
				return true
			}
		} else {
			capsCount = 0
		}
	}
	return false
}

// containsExcessiveExclamations checks for excessive use of exclamation marks
func (ca *ContentAnalyzer) containsExcessiveExclamations(text string) bool {
	return strings.Count(text, "!") > ca.config.MaxExclamationMarks
}

// checkSpamTriggers looks for specific spam trigger words and phrases
func (ca *ContentAnalyzer) checkSpamTriggers(text string, result *AnalysisResult, mu *sync.Mutex) {
	lowercaseText := strings.ToLower(text)
	for trigger, weight := range ca.config.SpamTriggers {
		if strings.Contains(lowercaseText, strings.ToLower(trigger)) {
			mu.Lock()
			result.SpamScore += float64(weight) / 100
			mu.Unlock()
		}
	}
}

// analyzeContentCategories categorizes the content based on predefined categories and word weights
func (ca *ContentAnalyzer) analyzeContentCategories(text string, result *AnalysisResult, mu *sync.Mutex) {
	words := strings.Fields(strings.ToLower(text))
	for category, categoryWords := range ca.config.ContentCategories {
		score := 0.0
		for _, word := range categoryWords {
			if contains(words, word.Text) {
				score += word.Weight
			}
		}
		mu.Lock()
		result.CategoryScores[category] = score
		mu.Unlock()
	}
}

// generateSummaryMessage creates a summary message based on the analysis results
func (ca *ContentAnalyzer) generateSummaryMessage(result *AnalysisResult) string {
	if !result.IsSafe {
		if result.SpamScore >= 0.7 {
			return fmt.Sprintf("Content flagged as potential spam (score: %.2f)", result.SpamScore)
		}
		if len(result.SuspiciousPatterns) > 0 {
			return "Content contains suspicious patterns"
		}
	}
	return fmt.Sprintf("Content appears safe (spam score: %.2f)", result.SpamScore)
}

// contains checks if a slice of strings contains a specific string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
