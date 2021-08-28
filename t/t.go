package t

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

func (c Coord) OnEdge(w, h int) bool {
	return c.X == 0 || c.Y == 0 || c.X == w-1 || c.Y == h-1
}

func (c Coord) Moved(dir string) Coord {
	switch dir {
	case "up":
		return Coord{c.X, c.Y + 1}
	case "down":
		return Coord{c.X, c.Y - 1}
	case "left":
		return Coord{c.X - 1, c.Y}
	case "right":
		return Coord{c.X + 1, c.Y}
	}
	panic("invalid move")
	// return Coord{head.X, head.Y} // invalid move - shouldn't happen
}

func (c Coord) OutOfBounds(width, height int) bool {
	if c.X < 0 {
		return true
	} else if c.X > width-1 {
		return true
	}
	if c.Y < 0 {
		return true
	} else if c.Y > height-1 {
		return true
	}
	return false
}

func (p1 Coord) Distance(p2 Coord) int {
	// In a plane with p1 at (x1, y1) and p2 at (x2, y2), it is |x1 - x2| + |y1 - y2|
	return abs(p1.X-p2.X) + abs(p1.Y-p2.Y)
}

func (a Coord) SamePos(b Coord) bool {
	return a.X == b.X && a.Y == b.Y
}

// Abs returns the absolute value of x.
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
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
