package pokeapi

const BaseURL = "https://pokeapi.co/api/v2"
const (
	LocationAreaURL = BaseURL + "/location-area"
	PokemonURL      = BaseURL + "/pokemon"
)

type NamedAPIResourceList[T any] struct {
	Count    int    // Total items at endpoint
	Next     string // Next URL in the sequence
	Previous string // Previous URL in the sequence
	Results  []NamedAPIResource[T]
}

type NamedAPIResource[_ any] struct {
	Name string
	URL  string
}

type Pokemon struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	IsDefault      bool   `json:"is_default"`
	Order          int    `json:"order"`
	Weight         int    `json:"weight"`
	// Abilities []PokemonAbility
	// Forms []NamedAPIResource[PokemonForm]
	// GameIndices []VersionGameIndex `json:"game_indices"`
	// HeldItems []PokemonHeldItem `json:"held_items"`
	LocationAreaEncounters string `json:"location_area_encounters"`
	// Moves []PokemonMove
	// PastTypes []PokemonTypePast
	// PastAbilities []PokemonAbiliityPast
	// Sprites PokemonSprites
	// Cries PokemonCries
	// Species NamedAPIResource[PokemonSpecies]
	Stats []PokemonStat
	// Types []PokemonType
}

type PokemonStat struct {
	Stat      NamedAPIResource[Stat]
	Effort    int
	Base_stat int
}

type Stat struct {
	// Id
	// Name
	// GameIndex
	// IsBattleOnly
	// AffectingMoves
	// AffectingNatures
	// Characteristics
	// MoveDamageClass
	// Names
}
