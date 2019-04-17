package model

import (
	"encoding/json"
	"fmt"
)

type signal struct{}

var liveEntriesLit []string = []string{
	"STR",
	"DEX",
	"CON",
	"INT",
	"WIS",
	"CHA",
	"Acrobatics",
	"Aegis Level",
	"Aegis|Level",
	"Alchemist Level",
	"Alchemist|Level",
	"Appraise",
	"Arcanist Level",
	"Arcanist|Level",
	"Archivist Level",
	"Archivist|Level",
	"Ardent Level",
	"Ardent|Level",
	"Artificer Level",
	"Artificer|Level",
	"Athletics",
	"Autohypnosis",
	"Balance",
	"Barbarian Level",
	"Barbarian|Level",
	"Bard Level",
	"Bard|Level",
	"Battle Dancer Level",
	"Battle Dancer|Level",
	"Beguiler Level",
	"Beguiler|Level",
	"Binder Level",
	"Binder|Level",
	"Bloodrager Level",
	"Bloodrager|Level",
	"Bluff",
	"Brawler Level",
	"Brawler|Level",
	"Cavalier Level",
	"Cavalier|Level",
	"Cleric Level",
	"Cleric|Level",
	"Climb",
	"Computers", 
	"Concentration",
	"Control Shape",
	"Craft (alchemy)",
	"Craft (armour)", 
	"Craft (armourer)", 
	"Craft (baskets)", 
	"Craft (blacksmithing)", 
	"Craft (bookbinding)", 
	"Craft (books)",
	"Craft (bowmaking)",
	"Craft (bows)", 
	"Craft (brewery)", 
	"Craft (calligraphy)",
	"Craft (carpentry)", 
	"Craft (carver)", 
	"Craft (chemist)", 
	"Craft (cloth)", 
	"Craft (clothing)",
	"Craft (cooking)",
	"Craft (founder)",
	"Craft (glass)", 
	"Craft (glassblowing)",
	"Craft (jeweler)",
	"Craft (jewelry)",
	"Craft (leather)", 
	"Craft (leatherworking)",
	"Craft (locks)", 
	"Craft (masonry)",
	"Craft (mechanic)", 
	"Craft (painting)",
	"Craft (paintings)", 
	"Craft (pottery)",
	"Craft (sculpting)", 
	"Craft (sculptures)",
	"Craft (ships)",
	"Craft (shoes)",
	"Craft (stone masonry)",
	"Craft (stonemasonry)", 
	"Craft (tailor)",
	"Craft (traps)", 
	"Craft (weapons)",
	"Craft (weaving)",
	"Craft (winemaking)",
	"Crusader Level",
	"Crusader|Level",
	"Cryptic Level",
	"Cryptic|Level",
	"Culture",
	"Death Master Level",
	"Death Master|Level",
	"Decipher Script",
	"Diplomacy",
	"Disable Device",
	"Disarm Traps",
	"Disguise",
	"Divine Mind Level",
	"Divine Mind|Level",
	"Dragon Shaman Level",
	"Dragon Shaman|Level",
	"Dragonfire Adept Level",
	"Dragonfire Adept|Level",
	"Dread Level",
	"Dread Necromancer Level",
	"Dread Necromancer|Level",
	"Dread|Level",
	"Druid Level",
	"Druid|Level",
	"Duskblade Level",
	"Duskblade|Level",
	"Eidolon Level",
	"Eidolon|Level",
	"Engineering",
	"Escape Artist",
	"Factotum Level",
	"Factotum|Level",
	"Favoured Soul Level",
	"Favoured Soul|Level",
	"Fighter Level",
	"Fighter|Level",
	"Finesse",
	"Fly",
	"Forgery",
	"Gather Information",
	"Gunslinger Level",
	"Gunslinger|Level",
	"Handle Animal",
	"Harbinger Level",
	"Harbinger|Level",
	"Heal",
	"Hexblade Level",
	"Hexblade|Level",
	"Hide",
	"High Guard Level",
	"High Guard|Level",
	"Hunter Level",
	"Hunter|Level",
	"Iaijutsu Focus",
	"Imperial Man-at-arms Level",
	"Imperial Man-at-arms|Level",
	"Incarnate Level",
	"Incarnate|Level",
	"Influence",
	"Inquisitor Level",
	"Inquisitor|Level",
	"Intimidate",
	"Investigator Level",
	"Investigator|Level",
	"Jester Level",
	"Jester|Level",
	"Jump",
	"Khalid Asad Level",
	"Khalid Asad|Level",
	"Kineticist Level",
	"Kineticist|Level",
	"Knowledge (aeronautics)",
	"Knowledge (arcana)",
	"Knowledge (dungeoneering)",
	"Knowledge (engineering)",
	"Knowledge (geography)",
	"Knowledge (history)",
	"Knowledge (local)",
	"Knowledge (nature)",
	"Knowledge (nobility)",
	"Knowledge (planes)",
	"Knowledge (psionics)",
	"Knowledge (religion)",
	"Life Science",
	"Linguistics",
	"Listen",
	"Locate Traps",
	"Lurk Level",
	"Lurk|Level",
	"Magus Level",
	"Magus|Level",
	"Marksman Level",
	"Marksman|Level",
	"Martial Adept Level",
	"Martial Adept|Level",
	"Martial Lore",
	"Medicine", 
	"Medium Level",
	"Medium|Level",
	"Mesmerist Level",
	"Mesmerist|Level",
	"Monk Level",
	"Monk|Level",
	"Mountebank Level",
	"Mountebank|Level",
	"Move Silently",
	"Mystic Level",
	"Mystic|Level",
	"Mysticism",
	"Mythic Tier",
	"Nature",
	"Ninja Level",
	"Ninja|Level",
	"Occultist Level",
	"Occultist|Level",
	"Open Lock",
	"Oracle Level",
	"Oracle|Level",
	"Paladin Level",
	"Paladin|Level",
	"Panther Warrior Level",
	"Panther Warrior|Level",
	"Perception",
	"Perform (act)",
	"Perform (comedy)",
	"Perform (dance)",
	"Perform (keyboard)",
	"Perform (oratory)",
	"Perform (percussion)",
	"Perform (sing)",
	"Perform (string)",
	"Perform (wind)",
	"Performance",
	"Physical Science", 
	"Piloting", 
	"Prepared Level",
	"Prepared|Level",
	"Prestige Level",
	"Prestige|Level",
	"Priest Level",
	"Priest|Level",
	"Profession (architect)",
	"Profession (baker)",
	"Profession (barrister)", 
	"Profession (brewer)",
	"Profession (butcher)", 
	"Profession (clerk)",
	"Profession (cook)",
	"Profession (courtesan)",
	"Profession (driver)",
	"Profession (engineer)", 
	"Profession (farmer)", 
	"Profession (fisherman)",
	"Profession (gambler)",
	"Profession (gardener)",
	"Profession (herbalist)", 
	"Profession (innkeeper)", 
	"Profession (librarian)",
	"Profession (merchant)", 
	"Profession (midwife)", 
	"Profession (miller)",
	"Profession (miner)", 
	"Profession (porter)",
	"Profession (sailor)", 
	"Profession (scribe)", 
	"Profession (shepherd)",
	"Profession (soldier)", 
	"Profession (stable master)",
	"Profession (tanner)", 
	"Profession (trapper)", 
	"Profession (woodcutter)",
	"Profession",
	"Psicraft",
	"Psion Level",
	"Psion|Level",
	"Psychic Level",
	"Psychic Warrior Level",
	"Psychic Warrior|Level",
	"Psychic|Level",
	"Ranger Level",
	"Ranger|Level",
	"Religion",
	"Ride",
	"Rogue Level",
	"Rogue|Level",
	"Samurai Level",
	"Samurai|Level",
	"Savant Level",
	"Savant|Level",
	"Scout Level",
	"Scout|Level",
	"Scry",
	"Search",
	"Sense Motive",
	"Sha'ir Level",
	"Sha'ir|Level",
	"Shadowcaster Level",
	"Shadowcaster|Level",
	"Shaman Level",
	"Shaman|Level",
	"Shugenja Level",
	"Shugenja|Level",
	"Skald Level",
	"Skald|Level",
	"Slayer Level",
	"Slayer|Level",
	"Sleight of Hand",
	"Society",
	"Sorcerer Level",
	"Sorcerer|Level",
	"Soulborn Level",
	"Soulborn|Level",
	"Soulknife Level",
	"Soulknife|Level",
	"Spellcaster Level",
	"Spellcaster|Level",
	"Spellcraft",
	"Spellthief Level",
	"Spellthief|Level",
	"Spirit Shaman Level",
	"Spirit Shaman|Level",
	"Spiritualist Level",
	"Spiritualist|Level",
	"Spontaneous Level",
	"Spontaneous|Level",
	"Spot",
	"Stalker Level",
	"Stalker|Level",
	"Stealth",
	"Summoner Level",
	"Summoner|Level",
	"Survival",
	"Swashbuckler Level",
	"Swashbuckler|Level",
	"Swim",
	"Swordsage Level",
	"Swordsage|Level",
	"Tactician Level",
	"Tactician|Level",
	"Totemist Level",
	"Totemist|Level",
	"Track",
	"Truenamer Level",
	"Truenamer|Level",
	"Truespeak",
	"Tumble",
	"Urban Druid Level",
	"Urban Druid|Level",
	"Use Magic Device",
	"Use Rope",
	"Vigilante Level",
	"Vigilante|Level",
	"Vitalist Level",
	"Vitalist|Level",
	"Warblade Level",
	"Warblade|Level",
	"Warder Level",
	"Warder|Level",
	"Warlock Level",
	"Warlock|Level",
	"Warmage Level",
	"Warmage|Level",
	"Warpriest Level",
	"Warpriest|Level",
	"Wild Empathy",
	"Wilder Level",
	"Wilder|Level",
	"Witch Hunter Level",
	"Witch Hunter|Level",
	"Witch Level",
	"Witch|Level",
	"Wizard Level",
	"Wizard|Level",
	"Wu Jen Level",
	"Wu Jen|Level",
	"Wyrdcaster Level",
	"Wyrdcaster|Level",
	"Zealot Level",
	"Zealot|Level",
}

func liveEntries() []*StackedEntry {
	// scan for classes, skills etc
	strings := make(map[string]signal, 512)
	for _, game := range []string{"pathfinder", "dnd35"} {
		gameData := ReadGameData(game)
		if gameData != nil {
			// All skills
			// string += gameData skills | { ^ displayName }
			for _, skill := range gameData.Skills {
				strings[skill.SkillName()] = signal{}
			}

			// All classes
			// strings += gameData class | { ^ displayName }
			for _, class := range gameData.Classes {
				strings[class.Name+" Level"] = signal{}
			}
		}
	}
	// bring in the manual list as well, just to make sure
	for _, str := range liveEntriesLit {
		strings[str] = signal{}
	}

	entries := make([]*Entry, 0, len(strings))
	for str, _ := range strings {
		entry := Entry{str, str}
		entries = append(entries, &entry)
	}
	stacked := stackEntries(entries)
	fmt.Println("Found", len(entries), "entries to translate")
	return stacked
}

// GetLiveTranslations gives the translations needed by the Composer app
func GetLiveTranslations() []byte {
	entries := liveEntries()

	// translations := make([]*StackedTranslation, 0, len(entries)*len(Languages))

	var liveTranslations LiveTranslations
	liveTranslations.Languages = make([]LiveTranslationsLanguage, 0, len(Languages))
	for _, language := range Languages {
		languageTranslations := make([]LiveTranslationEntry, 0, len(entries))
		for _, entry := range entries {
			translations := entry.GetTranslations(language)
			selected := PickPreferredTranslation(entry.RankTranslations(translations, false))
			if selected != nil {
				for i, part := range entry.Entries {
					languageTranslations = append(languageTranslations, LiveTranslationEntry{
						Original:    part.Original,
						Translation: selected.Parts[i].Translation,
					})
				}
			}
		}
		fmt.Println(" -", language, "-", len(languageTranslations), "translations")

		if len(languageTranslations) > 0 {
			liveTranslations.Languages = append(liveTranslations.Languages, LiveTranslationsLanguage{
				Name:         LanguagePaths[language],
				Translations: languageTranslations,
			})
		}
	}
	// fmt.Println("Exporting:", liveTranslations)
	return liveTranslations.export()
}

type LiveTranslations struct {
	Languages []LiveTranslationsLanguage `json:"languages"`
}

type LiveTranslationsLanguage struct {
	Name         string `json:"name"`
	Translations []LiveTranslationEntry `json:"translations"`
}

type LiveTranslationEntry struct {
	Original    string `json:"original"`
	Translation string `json:"translation"`
}

func (liveTranslations LiveTranslations) export() []byte {
	data, err := json.Marshal(liveTranslations)
	if err != nil {
		return nil
	}
	return data
}

func GetMasterInjectionEntries() []*StackedEntry {
	// entries := make()

	pathfinder := ReadGameData("pathfinder")
	// dnd35 := ReadGameData("dnd35")
/*
	pathfinderCharInfoPages := []string{
        "Pathfinder/Core/Character Info.ai", 
        "Pathfinder/Core/Animal Companion.ai",
        "Pathfinder/Core/Barbarian/Barbarian - Character Info.ai",
        "Pathfinder/Core/Ranger/Ranger - Character Info.ai",
        "Pathfinder/GM/NPC.ai",
        "Pathfinder/GM/NPC Group.ai",
        "Pathfinder/Archetypes/Druid/World Walker - Character Info.ai",
    }*/
	
	pathfinderSkills := make(map[string]signal, len(pathfinder.Skills))
	// string += gameData skills | { ^ displayName }
	for _, skill := range pathfinder.Skills {
		pathfinderSkills[skill.SkillName()] = signal{}
	}

	/*
	var pathfinderSkills = pathfinder skills ^ skillName
	*/

	return nil
}