# go-video-organizer
Cli tool to organize videos files based on duration


```go
var folderStructure = []FolderInfo{
	{Name: "Micro", MaxDuration: 15},                       // 0 <= duration <= 15 seconds
	{Name: "Mini", MaxDuration: 60},                        // 15 < duration <= 60 seconds
	{Name: "Short", MaxDuration: 5 * 60},                   // 60 < duration <= 5*60 seconds (5 minutes)
	{Name: "Medium", MaxDuration: 15 * 60},                 // 560 < duration <= 1560 seconds (15 minutes)
	{Name: "Long", MaxDuration: 30 * 60},                   // 1560 < duration <= 3060 seconds (30 minutes)
	{Name: "Extended", MaxDuration: 60 * 60},               // 3060 < duration <= 6060 seconds (60 minutes)
	{Name: "Feature", MaxDuration: 120 * 60},               // 6060 < duration <= 12060 seconds (120 minutes)
	{Name: "Epic", MaxDuration: float64(^uint(0) >> 1)},    // > 120*60 seconds
}
```