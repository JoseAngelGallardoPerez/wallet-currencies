package feed

type Feed string

func (f Feed) String() string {
	return string(f)
}

const (
	ECB = Feed("ECB")
)

var knownFeeds = []Feed{ECB}

//GetKnownFeeds returns a list of known feeds
func GetKnownFeeds() []Feed {
	result := make([]Feed, len(knownFeeds))
	for i, knownFeed := range knownFeeds {
		result[i] = knownFeed
	}
	return result
}

//IsKnown verifies if given feed is known
func IsKnown(feed Feed) bool {
	for _, knownFeed := range knownFeeds {
		if feed == knownFeed {
			return true
		}
	}
	return false
}

//IsKnownString verifies if given feed string is known
func IsKnownString(feed string) bool {
	return IsKnown(Feed(feed))
}
