package agents_test

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/WALL-AI/uArgus/backend/monitor/internal/llm/provider"
	"github.com/WALL-AI/uArgus/backend/monitor/internal/semantic/agents"
)

// ---------------------------------------------------------------------------
// .env loader — reads key=value pairs from a .env file into the process env
// ---------------------------------------------------------------------------

func loadDotEnv(t *testing.T) {
	t.Helper()
	// loadDotEnv is a no-op; TestMain already loaded the .env file.
	// This is kept for backward compatibility / additional logging.
	if os.Getenv("LLM_PROVIDER") == "" {
		t.Logf("LLM_PROVIDER is empty after .env load")
	}
}

// ---------------------------------------------------------------------------
// Helper — build real DirectAgentsClient from env
// ---------------------------------------------------------------------------

func buildDirectClient(t *testing.T) *agents.DirectAgentsClient {
	t.Helper()

	llmProvider := os.Getenv("LLM_PROVIDER")
	if llmProvider == "" {
		t.Skip("LLM_PROVIDER not set; skipping integration test")
	}

	client, err := agents.NewDirectAgentsClient(agents.DirectAgentsConfig{
		ProviderType: provider.ProviderType(llmProvider),
		APIKey:       os.Getenv("LLM_API_KEY"),
		BaseURL:      os.Getenv("LLM_BASE_URL"),
		Model:        os.Getenv("LLM_MODEL"),
	})
	if err != nil {
		t.Fatalf("NewDirectAgentsClient: %v", err)
	}

	// Attach embedding client if configured
	embURL := os.Getenv("EMBEDDING_API_URL")
	if embURL != "" {
		emb := agents.NewEmbeddingClient(agents.EmbeddingConfig{
			BaseURL: embURL,
			Model:   envOr("EMBEDDING_MODEL", "bge-m3"),
			APIKey:  os.Getenv("EMBEDDING_API_KEY"),
		})
		client.SetEmbedder(emb)
	}

	return client
}

func buildEmbeddingClient(t *testing.T) *agents.EmbeddingClient {
	t.Helper()
	url := os.Getenv("EMBEDDING_API_URL")
	if url == "" {
		t.Skip("EMBEDDING_API_URL not set; skipping embedding test")
	}
	return agents.NewEmbeddingClient(agents.EmbeddingConfig{
		BaseURL: url,
		Model:   envOr("EMBEDDING_MODEL", "bge-m3"),
		APIKey:  os.Getenv("EMBEDDING_API_KEY"),
	})
}

func buildChromaClient(t *testing.T) *agents.ChromaClient {
	t.Helper()
	url := os.Getenv("CHROMA_URL")
	if url == "" {
		t.Skip("CHROMA_URL not set; skipping chroma test")
	}
	return agents.NewChromaClient(agents.ChromaConfig{
		BaseURL:    url,
		Collection: envOr("CHROMA_COLLECTION", "news_embeddings_test"),
	})
}

func envOr(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}

func ctxWithTimeout(t *testing.T, d time.Duration) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), d)
	t.Cleanup(cancel)
	return ctx
}

// ---------------------------------------------------------------------------
// Test headlines used across tests
// ---------------------------------------------------------------------------

var testHeadlines = []string{
	"SpaceX successfully launches Starship to orbit for the first time",
	"Federal Reserve raises interest rates by 25 basis points amid inflation concerns",
	"Major earthquake strikes Turkey, thousands feared dead",
}

// ===========================================================================
// Integration tests — require real LLM
// ===========================================================================

func TestMain(m *testing.M) {
	// Attempt to load .env before any tests.
	_, thisFile, _, _ := runtime.Caller(0)
	dir := filepath.Dir(thisFile)
	for i := 0; i < 6; i++ {
		candidate := filepath.Join(dir, ".env")
		if f, err := os.Open(candidate); err == nil {
			scanner := bufio.NewScanner(f)
			loaded := 0
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				if line == "" || strings.HasPrefix(line, "#") {
					continue
				}
				idx := strings.Index(line, "=")
				if idx < 1 {
					continue
				}
				key := strings.TrimSpace(line[:idx])
				val := strings.TrimSpace(line[idx+1:])
				// Strip surrounding quotes (single or double)
				val = stripQuotes(val)
				if os.Getenv(key) == "" {
					os.Setenv(key, val)
					loaded++
				}
			}
			f.Close()
			fmt.Fprintf(os.Stderr, "[test] loaded %d env vars from %s\n", loaded, candidate)
			// Print key env vars for diagnostics (mask sensitive values)
			for _, k := range []string{"LLM_PROVIDER", "LLM_BASE_URL", "LLM_MODEL", "LLM_API_KEY", "EMBEDDING_API_URL", "EMBEDDING_MODEL", "CHROMA_URL"} {
				v := os.Getenv(k)
				display := v
				if strings.Contains(strings.ToLower(k), "key") && len(v) > 8 {
					display = v[:4] + "****" + v[len(v)-4:]
				}
				fmt.Fprintf(os.Stderr, "[test]   %s=%q\n", k, display)
			}
			break
		}
		dir = filepath.Dir(dir)
	}
	os.Exit(m.Run())
}

// isModelNotFound returns true when the embedding API rejects the model name (HTTP 400 config error).
func isModelNotFound(err error) bool {
	if err == nil {
		return false
	}
	s := err.Error()
	return strings.Contains(s, "HTTP 400") || strings.Contains(s, "Model does not exist") || strings.Contains(s, "model_not_found")
}

// isConnRefused returns true if the error indicates the server is not running.
func isConnRefused(err error) bool {
	if err == nil {
		return false
	}
	s := err.Error()
	return strings.Contains(s, "connection refused") ||
		strings.Contains(s, "actively refused") ||
		strings.Contains(s, "no connection could be made")
}

func stripQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}

// ---------------------------------------------------------------------------
// Summarize
// ---------------------------------------------------------------------------

func TestIntegration_Summarize_Brief(t *testing.T) {
	loadDotEnv(t)
	client := buildDirectClient(t)
	ctx := ctxWithTimeout(t, 30*time.Second)

	result, err := client.Summarize(ctx, testHeadlines, agents.SummarizeOpts{
		Mode:      "brief",
		MaxTokens: 256,
	})
	if err != nil {
		t.Fatalf("Summarize brief: %v", err)
	}
	if len(result) < 10 {
		t.Fatalf("summary too short: %q", result)
	}
	t.Logf("Brief summary:\n%s", result)
}

func TestIntegration_Summarize_Analysis(t *testing.T) {
	loadDotEnv(t)
	client := buildDirectClient(t)
	ctx := ctxWithTimeout(t, 30*time.Second)

	result, err := client.Summarize(ctx, testHeadlines, agents.SummarizeOpts{
		Mode:      "analysis",
		MaxTokens: 512,
	})
	if err != nil {
		t.Fatalf("Summarize analysis: %v", err)
	}
	if result == "" {
		t.Log("analysis returned empty response (model may not support this prompt style)")
	} else if len(result) < 10 {
		t.Logf("analysis suspiciously short: %q", result)
	} else {
		t.Logf("Analysis summary:\n%s", result)
	}
}

func TestIntegration_Summarize_TechVariant(t *testing.T) {
	loadDotEnv(t)
	client := buildDirectClient(t)
	ctx := ctxWithTimeout(t, 30*time.Second)

	result, err := client.Summarize(ctx, testHeadlines, agents.SummarizeOpts{
		Mode:      "brief",
		Variant:   "tech",
		MaxTokens: 256,
	})
	if err != nil {
		t.Fatalf("Summarize tech: %v", err)
	}
	t.Logf("Tech brief:\n%s", result)
}

func TestIntegration_Summarize_Default(t *testing.T) {
	loadDotEnv(t)
	client := buildDirectClient(t)
	ctx := ctxWithTimeout(t, 30*time.Second)

	result, err := client.Summarize(ctx, testHeadlines, agents.SummarizeOpts{
		MaxTokens: 256,
	})
	if err != nil {
		t.Fatalf("Summarize default: %v", err)
	}
	if len(result) < 10 {
		t.Fatalf("default summary too short: %q", result)
	}
	t.Logf("Default summary:\n%s", result)
}

// ---------------------------------------------------------------------------
// Classify
// ---------------------------------------------------------------------------

func TestIntegration_Classify(t *testing.T) {
	loadDotEnv(t)
	client := buildDirectClient(t)
	ctx := ctxWithTimeout(t, 30*time.Second)

	result, err := client.Classify(ctx, "Major earthquake strikes Turkey, thousands feared dead")
	if err != nil {
		t.Fatalf("Classify: %v", err)
	}
	if len(result.Categories) == 0 {
		t.Fatal("expected at least one category")
	}
	t.Logf("Classify result: categories=%v confidence=%.2f", result.Categories, result.Confidence)
}

// ---------------------------------------------------------------------------
// Sentiment
// ---------------------------------------------------------------------------

func TestIntegration_Sentiment(t *testing.T) {
	loadDotEnv(t)
	client := buildDirectClient(t)
	ctx := ctxWithTimeout(t, 30*time.Second)

	results, err := client.Sentiment(ctx, testHeadlines)
	if err != nil {
		t.Fatalf("Sentiment: %v", err)
	}
	if len(results) != len(testHeadlines) {
		t.Fatalf("expected %d results, got %d", len(testHeadlines), len(results))
	}
	for i, r := range results {
		if r.Label != "positive" && r.Label != "negative" && r.Label != "neutral" {
			t.Errorf("[%d] unexpected label: %q", i, r.Label)
		}
		if r.Score < 0 || r.Score > 1 {
			t.Errorf("[%d] score out of range: %.4f", i, r.Score)
		}
		t.Logf("[%d] %q → %s (%.2f)", i, testHeadlines[i], r.Label, r.Score)
	}
}

// ---------------------------------------------------------------------------
// NER (ExtractEntities)
// ---------------------------------------------------------------------------

func TestIntegration_ExtractEntities(t *testing.T) {
	loadDotEnv(t)
	client := buildDirectClient(t)
	ctx := ctxWithTimeout(t, 30*time.Second)

	results, err := client.ExtractEntities(ctx, testHeadlines)
	if err != nil {
		t.Fatalf("ExtractEntities: %v", err)
	}
	if len(results) != len(testHeadlines) {
		t.Fatalf("expected %d headline results, got %d", len(testHeadlines), len(results))
	}

	// The SpaceX headline should have at least one ORG entity
	spacexEntities := results[0]
	foundOrg := false
	for _, e := range spacexEntities {
		t.Logf("  entity: text=%q type=%s conf=%.2f", e.Text, e.Type, e.Confidence)
		if e.Type == "ORG" {
			foundOrg = true
		}
	}
	if !foundOrg {
		t.Logf("warning: no ORG entity found for SpaceX headline (LLM variability)")
	}

	// The Turkey headline should have at least one LOC
	turkeyEntities := results[2]
	foundLoc := false
	for _, e := range turkeyEntities {
		t.Logf("  entity: text=%q type=%s conf=%.2f", e.Text, e.Type, e.Confidence)
		if e.Type == "LOC" {
			foundLoc = true
		}
	}
	if !foundLoc {
		t.Logf("warning: no LOC entity found for Turkey headline (LLM variability)")
	}
}

// ---------------------------------------------------------------------------
// Translate
// ---------------------------------------------------------------------------

func TestIntegration_Translate(t *testing.T) {
	loadDotEnv(t)
	client := buildDirectClient(t)
	ctx := ctxWithTimeout(t, 30*time.Second)

	result, err := client.Translate(ctx, "SpaceX successfully launches Starship to orbit", "zh")
	if err != nil {
		t.Fatalf("Translate: %v", err)
	}
	if len(result) < 5 {
		t.Fatalf("translation too short: %q", result)
	}
	t.Logf("Translate en→zh: %s", result)
}

func TestIntegration_Translate_French(t *testing.T) {
	loadDotEnv(t)
	client := buildDirectClient(t)
	ctx := ctxWithTimeout(t, 30*time.Second)

	result, err := client.Translate(ctx, "Federal Reserve raises interest rates by 25 basis points", "fr")
	if err != nil {
		t.Fatalf("Translate: %v", err)
	}
	if len(result) < 5 {
		t.Fatalf("translation too short: %q", result)
	}
	t.Logf("Translate en→fr: %s", result)
}

// ===========================================================================
// Embedding integration tests — require real BGE-M3 endpoint
// ===========================================================================

func TestIntegration_Embed_Single(t *testing.T) {
	loadDotEnv(t)
	client := buildDirectClient(t)
	if os.Getenv("EMBEDDING_API_URL") == "" {
		t.Skip("EMBEDDING_API_URL not set")
	}
	ctx := ctxWithTimeout(t, 30*time.Second)

	embedding, err := client.Embed(ctx, "SpaceX launches Starship")
	if err != nil {
		if isModelNotFound(err) {
			t.Skipf("Embedding model not found — check EMBEDDING_MODEL (e.g. SiliconFlow needs BAAI/bge-m3): %v", err)
		}
		t.Fatalf("Embed: %v", err)
	}
	if len(embedding) == 0 {
		t.Fatal("empty embedding")
	}
	t.Logf("Embedding dimension: %d (first 5: %v)", len(embedding), embedding[:min(5, len(embedding))])
}

func TestIntegration_EmbedBatch(t *testing.T) {
	loadDotEnv(t)
	emb := buildEmbeddingClient(t)
	ctx := ctxWithTimeout(t, 30*time.Second)

	embeddings, err := emb.EmbedBatch(ctx, testHeadlines)
	if err != nil {
		if isModelNotFound(err) {
			t.Skipf("Embedding model not found — check EMBEDDING_MODEL (e.g. SiliconFlow needs BAAI/bge-m3): %v", err)
		}
		t.Fatalf("EmbedBatch: %v", err)
	}
	if len(embeddings) != len(testHeadlines) {
		t.Fatalf("expected %d embeddings, got %d", len(testHeadlines), len(embeddings))
	}
	for i, e := range embeddings {
		if len(e) == 0 {
			t.Errorf("[%d] empty embedding", i)
		} else {
			t.Logf("[%d] dim=%d", i, len(e))
		}
	}
}

func TestIntegration_Embed_CosineSimilarity(t *testing.T) {
	loadDotEnv(t)
	emb := buildEmbeddingClient(t)
	ctx := ctxWithTimeout(t, 30*time.Second)

	// Two semantically similar texts and one different
	texts := []string{
		"SpaceX successfully launches Starship to orbit",
		"Starship rocket reaches orbit in historic SpaceX mission",
		"New bakery opens on Main Street in downtown Portland",
	}

	embeddings, err := emb.EmbedBatch(ctx, texts)
	if err != nil {
		if isModelNotFound(err) {
			t.Skipf("Embedding model not found — check EMBEDDING_MODEL (e.g. SiliconFlow needs BAAI/bge-m3): %v", err)
		}
		t.Fatalf("EmbedBatch: %v", err)
	}

	sim01 := cosineSim(embeddings[0], embeddings[1])
	sim02 := cosineSim(embeddings[0], embeddings[2])

	t.Logf("sim(SpaceX1, SpaceX2) = %.4f", sim01)
	t.Logf("sim(SpaceX1, Bakery)  = %.4f", sim02)

	if sim01 <= sim02 {
		t.Errorf("expected similar texts to have higher cosine similarity: %.4f <= %.4f", sim01, sim02)
	}
}

func cosineSim(a, b []float32) float64 {
	var dot, normA, normB float64
	for i := range a {
		dot += float64(a[i]) * float64(b[i])
		normA += float64(a[i]) * float64(a[i])
		normB += float64(b[i]) * float64(b[i])
	}
	denom := sqrtF64(normA) * sqrtF64(normB)
	if denom == 0 {
		return 0
	}
	return dot / denom
}

func sqrtF64(x float64) float64 {
	if x <= 0 {
		return 0
	}
	// Newton's method (good enough for tests)
	z := x
	for i := 0; i < 50; i++ {
		z = (z + x/z) / 2
	}
	return z
}

// ===========================================================================
// Chroma integration tests — require running Chroma server
// ===========================================================================

func TestIntegration_Chroma_CreateCollection(t *testing.T) {
	loadDotEnv(t)
	chroma := buildChromaClient(t)
	ctx := ctxWithTimeout(t, 15*time.Second)

	collName := fmt.Sprintf("test_integration_%d", time.Now().UnixMilli())
	id, err := chroma.CreateCollection(ctx, collName)
	if err != nil {
		if isConnRefused(err) {
			t.Skipf("Chroma not reachable: %v", err)
		}
		t.Fatalf("CreateCollection: %v", err)
	}
	t.Logf("Created collection %q id=%s", collName, id)

	// Cleanup
	if err := chroma.DeleteCollection(ctx, collName); err != nil {
		t.Logf("cleanup DeleteCollection: %v", err)
	}
}

func TestIntegration_Chroma_UpsertAndQuery(t *testing.T) {
	loadDotEnv(t)
	chroma := buildChromaClient(t)
	emb := buildEmbeddingClient(t)
	ctx := ctxWithTimeout(t, 60*time.Second)

	collName := fmt.Sprintf("test_upsert_%d", time.Now().UnixMilli())
	_, err := chroma.CreateCollection(ctx, collName)
	if err != nil {
		if isConnRefused(err) {
			t.Skipf("Chroma not reachable: %v", err)
		}
		t.Fatalf("CreateCollection: %v", err)
	}
	t.Cleanup(func() {
		_ = chroma.DeleteCollection(context.Background(), collName)
	})

	// Generate embeddings for test headlines
	embeddings, err := emb.EmbedBatch(ctx, testHeadlines)
	if err != nil {
		t.Fatalf("EmbedBatch: %v", err)
	}

	// Upsert
	items := make([]agents.UpsertInput, len(testHeadlines))
	for i, h := range testHeadlines {
		items[i] = agents.UpsertInput{
			ID:        fmt.Sprintf("headline-%d", i),
			Embedding: embeddings[i],
			Document:  h,
			Metadata:  map[string]string{"source": "test"},
		}
	}
	if err := chroma.Upsert(ctx, collName, items); err != nil {
		t.Fatalf("Upsert: %v", err)
	}
	t.Logf("Upserted %d items", len(items))

	// Query with a semantically similar text
	queryEmb, err := emb.Embed(ctx, "earthquake disaster in Turkey kills many people")
	if err != nil {
		t.Fatalf("Embed query: %v", err)
	}
	results, err := chroma.Query(ctx, collName, queryEmb, 3)
	if err != nil {
		t.Fatalf("Query: %v", err)
	}
	if len(results) == 0 {
		t.Fatal("expected query results")
	}
	for i, r := range results {
		t.Logf("  [%d] id=%s dist=%.4f doc=%q", i, r.ID, r.Distance, r.Document)
	}

	// The top result should be the earthquake headline
	if results[0].ID != "headline-2" {
		t.Logf("warning: top result is %q, expected headline-2 (earthquake); LLM/embedding variability", results[0].ID)
	}
}

// ===========================================================================
// Prompt alignment spot-check tests
// ===========================================================================

func TestIntegration_Prompt_BriefDoesNotStartWithBreaking(t *testing.T) {
	loadDotEnv(t)
	client := buildDirectClient(t)
	ctx := ctxWithTimeout(t, 30*time.Second)

	result, err := client.Summarize(ctx, testHeadlines, agents.SummarizeOpts{
		Mode:      "brief",
		MaxTokens: 256,
	})
	if err != nil {
		t.Fatalf("Summarize: %v", err)
	}
	lower := strings.ToLower(strings.TrimSpace(result))
	for _, bad := range []string{"breaking news", "good evening", "tonight"} {
		if strings.HasPrefix(lower, bad) {
			t.Errorf("brief summary starts with forbidden phrase %q: %s", bad, result)
		}
	}
	t.Logf("Brief (prompt check): %s", result)
}

func TestIntegration_Prompt_ClassifyUsesRefinedCategories(t *testing.T) {
	loadDotEnv(t)
	client := buildDirectClient(t)
	ctx := ctxWithTimeout(t, 30*time.Second)

	result, err := client.Classify(ctx, "Cyberattack targets government infrastructure in Estonia")
	if err != nil {
		t.Fatalf("Classify: %v", err)
	}
	t.Logf("Categories: %v (confidence: %.2f)", result.Categories, result.Confidence)

	allowedCats := map[string]bool{
		"conflict": true, "protest": true, "disaster": true, "diplomatic": true,
		"economic": true, "terrorism": true, "cyber": true, "health": true,
		"environmental": true, "military": true, "crime": true, "infrastructure": true,
		"tech": true, "general": true,
	}
	for _, cat := range result.Categories {
		if !allowedCats[cat] {
			t.Errorf("unexpected category %q — should be one of the refined set", cat)
		}
	}
}

// ===========================================================================
// End-to-end pipeline test: Summarize → Sentiment → NER → Embed
// ===========================================================================

func TestIntegration_FullPipeline(t *testing.T) {
	loadDotEnv(t)
	client := buildDirectClient(t)
	ctx := ctxWithTimeout(t, 120*time.Second)

	headlines := []string{
		"Apple announces new M4 chip with breakthrough AI capabilities",
		"UN Security Council emergency session on escalating conflict in Sudan",
		"Global stock markets surge as inflation data beats expectations",
	}

	// 1. Summarize
	summary, err := client.Summarize(ctx, headlines, agents.SummarizeOpts{Mode: "brief", MaxTokens: 256})
	if err != nil {
		t.Fatalf("pipeline summarize: %v", err)
	}
	t.Logf("1. Summary: %s", summary)

	// 2. Classify each headline
	for i, h := range headlines {
		cls, err := client.Classify(ctx, h)
		if err != nil {
			t.Fatalf("pipeline classify[%d]: %v", i, err)
		}
		t.Logf("2. Classify[%d]: %v (%.2f)", i, cls.Categories, cls.Confidence)
	}

	// 3. Sentiment
	sentiments, err := client.Sentiment(ctx, headlines)
	if err != nil {
		t.Fatalf("pipeline sentiment: %v", err)
	}
	for i, s := range sentiments {
		t.Logf("3. Sentiment[%d]: %s (%.2f)", i, s.Label, s.Score)
	}

	// 4. NER
	entities, err := client.ExtractEntities(ctx, headlines)
	if err != nil {
		t.Fatalf("pipeline ner: %v", err)
	}
	for i, ents := range entities {
		t.Logf("4. NER[%d]: %d entities", i, len(ents))
		for _, e := range ents {
			t.Logf("   - %q [%s] %.2f", e.Text, e.Type, e.Confidence)
		}
	}

	// 5. Translate summary to Chinese
	translated, err := client.Translate(ctx, summary, "zh")
	if err != nil {
		t.Fatalf("pipeline translate: %v", err)
	}
	t.Logf("5. Translated: %s", translated)

	// 6. Embed (if configured)
	if os.Getenv("EMBEDDING_API_URL") != "" {
		emb, err := client.Embed(ctx, headlines[0])
		if err != nil {
			t.Logf("6. Embedding error (non-fatal): %v", err)
		} else {
			t.Logf("6. Embedding dim=%d", len(emb))
		}
	} else {
		t.Log("6. Embedding skipped (EMBEDDING_API_URL not set)")
	}
}
