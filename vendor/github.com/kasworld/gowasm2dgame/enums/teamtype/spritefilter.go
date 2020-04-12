package teamtype

type iv struct {
	Index int
	Value int
}

var SpriteFilter = [TeamType_Count]struct {
	Name string
	IV   []iv
}{
	Red:    {"red", []iv{{0, 255}}},
	Blue:   {"blue", []iv{{1, 255}}},
	Green:  {"green", []iv{{2, 255}}},
	RRed:   {"rred", []iv{{0, 0}}},
	RBlue:  {"rblue", []iv{{1, 0}}},
	RGreen: {"rgreen", []iv{{2, 0}}},
}
