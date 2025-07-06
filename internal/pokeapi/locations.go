package pokeapi

type LocationArea struct {
	Id        int
	Name      string
	GameIndex int
	// EncounterMethodRates
	// Location
	// Names
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}

type PokemonEncounter struct {
	Pokemon NamedAPIResource[Pokemon]
	// VersionDetails
}
