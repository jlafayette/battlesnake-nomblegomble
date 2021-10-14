package wire

type GameState struct {
	Game  Game        `json:"game"`
	Turn  int         `json:"turn"`
	Board Board       `json:"board"`
	You   Battlesnake `json:"you"`
}

type Game struct {
	ID      string  `json:"id"`
	Ruleset Ruleset `json:"ruleset"`
	Timeout int32   `json:"timeout"`
}

type Ruleset struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Board struct {
	Height int           `json:"height"`
	Width  int           `json:"width"`
	Food   []Coord       `json:"food"`
	Snakes []Battlesnake `json:"snakes"`

	// Used in non-standard game modes
	Hazards []Coord `json:"hazards"`
}

func (b *Board) CapTo4Snakes(you Battlesnake) {
	// currently the tree search can only handle 4 snakes (including you),
	// so this function removes far away snakes until there are only 4
	if len(b.Snakes) <= 4 {
		return
	}
	// measure distance from your head to any segment of the other
	// snakes body (ignore tail would be a nice optimization)
	lengths := make(map[string]int, len(b.Snakes)-1)
	for _, snake := range b.Snakes {
		if snake.ID == you.ID {
			continue
		}
		head := you.Head
		minDist := 999
		for _, c := range snake.Body {
			minDist = min(minDist, head.distance(c))
		}
		lengths[snake.ID] = minDist
	}
	// Remove the snake with greatest distance until there are 4
	for {
		if len(b.Snakes) <= 4 {
			break
		}
		maxId := ""
		maxDistance := 0
		for k, v := range lengths {
			if v > maxDistance {
				maxId = k
				maxDistance = v
			}
		}
		indexToRemove := -1
		for i, snake := range b.Snakes {
			if snake.ID == maxId {
				indexToRemove = i
			}
		}
		b.Snakes = remove(b.Snakes, indexToRemove)
		delete(lengths, maxId)
	}
}

type Battlesnake struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Health  int32   `json:"health"`
	Body    []Coord `json:"body"`
	Head    Coord   `json:"head"`
	Length  int32   `json:"length"`
	Latency string  `json:"latency"`

	// Used in non-standard game modes
	Shout string `json:"shout"`
	Squad string `json:"squad"`
}

type Coord struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Response Structs

type BattlesnakeInfoResponse struct {
	APIVersion string `json:"apiversion"`
	Author     string `json:"author"`
	Color      string `json:"color"`
	Head       string `json:"head"`
	Tail       string `json:"tail"`
}

type BattlesnakeMoveResponse struct {
	Move  string `json:"move"`
	Shout string `json:"shout,omitempty"`
}

// Utils

func remove(slice []Battlesnake, s int) []Battlesnake {
	return append(slice[:s], slice[s+1:]...)
}

func (c1 Coord) distance(c2 Coord) int {
	return abs(c1.X-c2.X) + abs(c1.Y-c2.Y)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
