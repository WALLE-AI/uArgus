package agents

import (
	"fmt"
	"strings"
	"time"
)

// ---------------------------------------------------------------------------
// Summarize prompt builder — ported from TS _shared.ts::buildArticlePrompts
// ---------------------------------------------------------------------------

// SummarizePrompts holds the system and user prompts for a summarize request.
type SummarizePrompts struct {
	SystemPrompt string
	UserPrompt   string
}

// BuildSummarizePrompts constructs mode-aware, variant-aware prompts.
// Modes: "brief", "analysis", "translate", "" (default).
// Variants: "tech", "" (full/general).
func BuildSummarizePrompts(headlines []string, opts SummarizeOpts) SummarizePrompts {
	headlineText := formatNumberedHeadlines(headlines)
	intelSection := ""
	if opts.GeoContext != "" {
		intelSection = "\n\n" + opts.GeoContext
	}

	isTech := opts.Variant == "tech"
	dateContext := fmt.Sprintf("Current date: %s.", time.Now().UTC().Format("2006-01-02"))
	if !isTech {
		dateContext += " Provide geopolitical context appropriate for the current date."
	}

	langInstruction := ""
	if opts.Lang != "" && opts.Lang != "en" {
		langInstruction = fmt.Sprintf("\nIMPORTANT: Output the summary in %s language.", strings.ToUpper(opts.Lang))
	}

	switch opts.Mode {
	case "brief":
		return buildBriefPrompts(headlineText, intelSection, dateContext, langInstruction, isTech)
	case "analysis":
		return buildAnalysisPrompts(headlineText, intelSection, dateContext, langInstruction, isTech)
	case "translate":
		return buildTranslatePrompts(headlines, opts.Variant)
	default:
		return buildDefaultPrompts(headlineText, intelSection, dateContext, langInstruction, isTech)
	}
}

func buildBriefPrompts(headlineText, intelSection, dateContext, langInstruction string, isTech bool) SummarizePrompts {
	var systemPrompt string
	if isTech {
		systemPrompt = fmt.Sprintf(`%s

Summarize the single most important tech/startup headline in 2 concise sentences MAX (under 60 words total).
Rules:
- Each numbered headline below is a SEPARATE, UNRELATED story
- Pick the ONE most significant headline and summarize ONLY that story
- NEVER combine or merge facts, names, or details from different headlines
- Focus ONLY on technology, startups, AI, funding, product launches, or developer news
- IGNORE political news, trade policy, tariffs, government actions unless directly about tech regulation
- Lead with the company/product/technology name
- No bullet points, no meta-commentary, no elaboration beyond the core facts%s`, dateContext, langInstruction)
	} else {
		systemPrompt = fmt.Sprintf(`%s

Summarize the single most important headline in 2 concise sentences MAX (under 60 words total).
Rules:
- Each numbered headline below is a SEPARATE, UNRELATED story
- Pick the ONE most significant headline and summarize ONLY that story
- NEVER combine or merge people, places, or facts from different headlines into one sentence
- Lead with WHAT happened and WHERE - be specific
- NEVER start with "Breaking news", "Good evening", "Tonight", or TV-style openings
- Start directly with the subject of the chosen headline
- If intelligence context is provided, use it only if it relates to your chosen headline
- No bullet points, no meta-commentary, no elaboration beyond the core facts%s`, dateContext, langInstruction)
	}

	userPrompt := fmt.Sprintf("Each headline below is a separate story. Pick the most important ONE and summarize only that story:\n%s%s", headlineText, intelSection)
	return SummarizePrompts{SystemPrompt: systemPrompt, UserPrompt: userPrompt}
}

func buildAnalysisPrompts(headlineText, intelSection, dateContext, langInstruction string, isTech bool) SummarizePrompts {
	var systemPrompt string
	if isTech {
		systemPrompt = fmt.Sprintf(`%s

Analyze the most significant tech/startup development in 2 concise sentences MAX (under 60 words total).
Rules:
- Each numbered headline below is a SEPARATE, UNRELATED story
- Pick the ONE most significant story and analyze ONLY that
- NEVER combine facts from different headlines
- Focus ONLY on technology implications: funding trends, AI developments, market shifts, product strategy
- IGNORE political implications, trade wars, government unless directly about tech policy
- Lead with the insight, no filler or elaboration`, dateContext)
	} else {
		systemPrompt = fmt.Sprintf(`%s

Analyze the most significant development in 2 concise sentences MAX (under 60 words total). Be direct and specific.
Rules:
- Each numbered headline below is a SEPARATE, UNRELATED story
- Pick the ONE most significant story and analyze ONLY that
- NEVER combine or merge people, places, or facts from different headlines
- Lead with the insight - what's significant and why
- NEVER start with "Breaking news", "Tonight", "The key/dominant narrative is"
- Start with substance, no filler or elaboration
- If intelligence context is provided, use it only if it relates to your chosen headline`, dateContext)
	}

	var userPrompt string
	if isTech {
		userPrompt = fmt.Sprintf("Each headline is a separate story. What's the key tech trend?\n%s%s", headlineText, intelSection)
	} else {
		userPrompt = fmt.Sprintf("Each headline is a separate story. What's the key pattern or risk?\n%s%s", headlineText, intelSection)
	}
	_ = langInstruction // analysis mode does not append lang instruction in TS
	return SummarizePrompts{SystemPrompt: systemPrompt, UserPrompt: userPrompt}
}

func buildTranslatePrompts(headlines []string, targetLang string) SummarizePrompts {
	systemPrompt := fmt.Sprintf(`You are a professional news translator. Translate the following news headlines/summaries into %s.
Rules:
- Maintain the original tone and journalistic style.
- Do NOT add any conversational filler (e.g., "Here is the translation").
- Output ONLY the translated text.
- If the text is already in %s, return it as is.`, targetLang, targetLang)

	text := ""
	if len(headlines) > 0 {
		text = headlines[0]
	}
	userPrompt := fmt.Sprintf("Translate to %s:\n%s", targetLang, text)
	return SummarizePrompts{SystemPrompt: systemPrompt, UserPrompt: userPrompt}
}

func buildDefaultPrompts(headlineText, intelSection, dateContext, langInstruction string, isTech bool) SummarizePrompts {
	var systemPrompt string
	if isTech {
		systemPrompt = fmt.Sprintf("%s\n\nPick the most important tech headline and summarize it in 2 concise sentences (under 60 words). Each headline is a separate story - NEVER merge facts from different headlines. Focus on startups, AI, funding, products. Ignore politics unless directly about tech regulation.%s", dateContext, langInstruction)
	} else {
		systemPrompt = fmt.Sprintf("%s\n\nPick the most important headline and summarize it in 2 concise sentences (under 60 words). Each headline is a separate, unrelated story - NEVER merge people or facts from different headlines. Lead with substance. NEVER start with \"Breaking news\" or \"Tonight\".%s", dateContext, langInstruction)
	}

	userPrompt := fmt.Sprintf("Each headline is a separate story. Key takeaway from the most important one:\n%s%s", headlineText, intelSection)
	return SummarizePrompts{SystemPrompt: systemPrompt, UserPrompt: userPrompt}
}

// ---------------------------------------------------------------------------
// Classify prompt — refined categories aligned with TS _classifier.ts
// ---------------------------------------------------------------------------

// BuildClassifyPrompt returns the system prompt and user prompt for classification.
func BuildClassifyPrompt(text string) (systemPrompt, userPrompt string) {
	systemPrompt = `You are a news classifier. Given a headline, return a JSON object with:
- "categories": array of category strings from this set: "conflict", "protest", "disaster", "diplomatic", "economic", "terrorism", "cyber", "health", "environmental", "military", "crime", "infrastructure", "tech", "general"
- "confidence": a float between 0 and 1
Only return valid JSON, no explanation.`

	userPrompt = text
	return
}

// ---------------------------------------------------------------------------
// Sentiment prompt — LLM-based sentiment analysis
// ---------------------------------------------------------------------------

// BuildSentimentPrompt returns the system and user prompt for batch sentiment analysis.
func BuildSentimentPrompt(texts []string) (systemPrompt, userPrompt string) {
	systemPrompt = `You are a sentiment classifier for news headlines.
For each headline provided, classify the sentiment as "positive", "negative", or "neutral" with a confidence score between 0 and 1.
Return a JSON array where each element is: {"label": "positive"|"negative"|"neutral", "score": <float>}
The array must have exactly the same number of elements as input headlines, in the same order.
Only return valid JSON, no explanation.`

	var sb strings.Builder
	sb.WriteString("Classify the sentiment of each headline:\n")
	for i, t := range texts {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, t))
	}
	userPrompt = sb.String()
	return
}

// ---------------------------------------------------------------------------
// NER prompt — LLM-based named entity recognition
// ---------------------------------------------------------------------------

// BuildNERPrompt returns the system and user prompt for batch NER extraction.
func BuildNERPrompt(texts []string) (systemPrompt, userPrompt string) {
	systemPrompt = `You are a named entity recognition (NER) system for news headlines.
For each headline, extract all named entities with their type and confidence.
Entity types: PER (person), ORG (organization), LOC (location), MISC (miscellaneous).
Return a JSON array of arrays. Each inner array corresponds to one headline and contains objects: {"text": "<entity>", "type": "PER"|"ORG"|"LOC"|"MISC", "confidence": <float>}
The outer array must have exactly the same number of elements as input headlines, in the same order.
Only return valid JSON, no explanation.`

	var sb strings.Builder
	sb.WriteString("Extract named entities from each headline:\n")
	for i, t := range texts {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, t))
	}
	userPrompt = sb.String()
	return
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func formatNumberedHeadlines(headlines []string) string {
	var sb strings.Builder
	for i, h := range headlines {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, h))
	}
	return sb.String()
}
