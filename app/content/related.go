// Package content provides intelligent related posts finding with tag-based scoring and randomization.
package content

import (
	"math"
	"math/rand/v2"
	"slices"
	"time"
)

// RelatedPostsConfig controls related posts algorithm behavior.
type RelatedPostsConfig struct {
	// MaxCandidates is how many top candidates to find before randomization (default: 10)
	MaxCandidates int
	// ResultLimit is how many posts to return (default: 6)
	ResultLimit int
	// TagMatchWeight is score points per matching tag (default: 10)
	TagMatchWeight int
	// RecencyDecayDays is how many days reduce score by 1 point (default: 30)
	RecencyDecayDays int
	// RandomBonusMax is maximum random bonus points for variety (default: 2)
	RandomBonusMax int
	// MinTagMatches is minimum tag overlap required (default: 1)
	MinTagMatches int
}

// DefaultRelatedConfig returns sensible defaults for related posts.
func DefaultRelatedConfig() RelatedPostsConfig {
	return RelatedPostsConfig{
		MaxCandidates:    10,
		ResultLimit:      6,
		TagMatchWeight:   10,
		RecencyDecayDays: 30,
		RandomBonusMax:   2,
		MinTagMatches:    1,
	}
}

// RelatedPostsFinder finds related posts using tag-based scoring.
type RelatedPostsFinder struct {
	tagIndex map[string][]*Node // tag -> posts with that tag
	config   RelatedPostsConfig
}

// NewRelatedPostsFinder creates a finder with tag index.
func NewRelatedPostsFinder(posts []*Node, config RelatedPostsConfig) *RelatedPostsFinder {
	finder := &RelatedPostsFinder{
		tagIndex: make(map[string][]*Node),
		config:   config,
	}
	finder.buildTagIndex(posts)
	return finder
}

// buildTagIndex creates inverted index: tag -> []posts for O(k) lookup.
func (rf *RelatedPostsFinder) buildTagIndex(posts []*Node) {
	for _, post := range posts {
		tags := rf.extractTags(post)
		for _, tag := range tags {
			rf.tagIndex[tag] = append(rf.tagIndex[tag], post)
		}
	}
}

// FindRelated returns related posts for the given page, with randomization for variety.
func (rf *RelatedPostsFinder) FindRelated(currentPage *Node) []*Node {
	currentTags := rf.extractTags(currentPage)
	if len(currentTags) == 0 {
		return []*Node{} // No tags = no related posts
	}

	// Build candidate pool with scores
	candidates := rf.findCandidates(currentPage, currentTags)
	if len(candidates) == 0 {
		return []*Node{}
	}

	// Sort by score descending
	slices.SortFunc(candidates, func(a, b *scoredPost) int {
		// Higher scores first
		return int(b.score - a.score)
	})

	// Take top MaxCandidates, then randomly select ResultLimit
	maxCandidates := min(len(candidates), rf.config.MaxCandidates)

	topCandidates := candidates[:maxCandidates]

	// Shuffle for randomness on each build
	rand.Shuffle(len(topCandidates), func(i, j int) {
		topCandidates[i], topCandidates[j] = topCandidates[j], topCandidates[i]
	})

	// Return up to ResultLimit posts
	resultLimit := min(len(topCandidates), rf.config.ResultLimit)

	results := make([]*Node, resultLimit)
	for i := range resultLimit {
		results[i] = topCandidates[i].post
	}

	return results
}

// scoredPost holds a post with its relevance score.
type scoredPost struct {
	post  *Node
	score float64
}

// findCandidates builds list of candidate posts with relevance scores.
func (rf *RelatedPostsFinder) findCandidates(currentPage *Node, currentTags []string) []*scoredPost {
	seen := make(map[string]bool)
	seen[currentPage.Permalink] = true // Exclude current page

	candidates := make([]*scoredPost, 0)

	// For each tag, get posts with that tag (O(k) where k = avg posts per tag)
	for _, tag := range currentTags {
		postsWithTag, exists := rf.tagIndex[tag]
		if !exists {
			continue
		}

		for _, candidate := range postsWithTag {
			if seen[candidate.Permalink] {
				continue // Already processed or is current page
			}
			seen[candidate.Permalink] = true

			score := rf.calculateScore(currentPage, candidate, currentTags)
			if score > 0 {
				candidates = append(candidates, &scoredPost{
					post:  candidate,
					score: score,
				})
			}
		}
	}

	return candidates
}

// calculateScore computes relevance score for a candidate post.
// Score = (MatchingTags × TagMatchWeight) - (DaysDifference ÷ RecencyDecayDays) + RandomBonus.
func (rf *RelatedPostsFinder) calculateScore(currentPage, candidate *Node, currentTags []string) float64 {
	candidateTags := rf.extractTags(candidate)

	// Count matching tags
	matchingTags := rf.countMatchingTags(currentTags, candidateTags)
	if matchingTags < rf.config.MinTagMatches {
		return 0 // Not enough overlap
	}

	// Base score from tag matches
	score := float64(matchingTags * rf.config.TagMatchWeight)

	// Recency penalty: older posts score slightly lower
	recencyPenalty := rf.calculateRecencyPenalty(currentPage, candidate)
	score -= recencyPenalty

	// Add random bonus for variety (0 to RandomBonusMax)
	randomBonus := rand.Float64() * float64(rf.config.RandomBonusMax)
	score += randomBonus

	return score
}

// countMatchingTags counts how many tags overlap between two posts.
func (rf *RelatedPostsFinder) countMatchingTags(tags1, tags2 []string) int {
	tagSet := make(map[string]bool)
	for _, tag := range tags1 {
		tagSet[tag] = true
	}

	matches := 0
	for _, tag := range tags2 {
		if tagSet[tag] {
			matches++
		}
	}
	return matches
}

// calculateRecencyPenalty returns penalty based on date difference.
func (rf *RelatedPostsFinder) calculateRecencyPenalty(currentPage, candidate *Node) float64 {
	currentDate := rf.extractDate(currentPage)
	candidateDate := rf.extractDate(candidate)

	if currentDate.IsZero() || candidateDate.IsZero() {
		return 0 // No penalty if dates missing
	}

	// Calculate day difference
	daysDiff := math.Abs(currentDate.Sub(candidateDate).Hours() / 24)

	// Penalty = daysDiff / RecencyDecayDays
	penalty := daysDiff / float64(rf.config.RecencyDecayDays)

	return penalty
}

// extractTags gets normalized tags from a node.
func (rf *RelatedPostsFinder) extractTags(node *Node) []string {
	if node.Config == nil {
		return []string{}
	}

	tagsInterface, ok := node.Config["tags"]
	if !ok {
		return []string{}
	}

	switch v := tagsInterface.(type) {
	case []string:
		return v
	case []any:
		tags := make([]string, 0, len(v))
		for _, t := range v {
			if str, ok := t.(string); ok {
				tags = append(tags, str)
			}
		}
		return tags
	default:
		return []string{}
	}
}

// extractDate gets the date from a node's config.
func (rf *RelatedPostsFinder) extractDate(node *Node) time.Time {
	if node.Config == nil {
		return time.Time{}
	}

	dateVal, ok := node.Config["date"]
	if !ok {
		return time.Time{}
	}

	switch v := dateVal.(type) {
	case time.Time:
		return v
	case string:
		// Try RFC3339
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			return t
		}
		// Try date-only format
		if t, err := time.Parse("2006-01-02", v); err == nil {
			return t
		}
		return time.Time{}
	default:
		return time.Time{}
	}
}
