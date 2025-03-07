package w365tt

import (
	"fmt"
	"gradgrind/backend/base"
	"strconv"
	"strings"
)

const roomtagerror = `<Error>Room Tag (Shortcut) defined twice: %s>`

func (db *DbTopLevel) readRooms(newdb *base.DbTopLevel) {
	db.RealRooms = map[Ref]string{}
	db.RoomTags = map[string]Ref{}
	db.RoomChoiceNames = map[string]Ref{}
	for _, e := range db.Rooms {
		// Perform some checks and add to the RoomTags map.
		_, nok := db.RoomTags[e.Tag]
		if nok {
			base.Report(
				roomtagerror,
				e.Tag)
			continue
		}
		db.RoomTags[e.Tag] = e.Id
		// Copy to base db.
		tsl := db.handleZeroAfternoons(e.NotAvailable, 1)
		r := newdb.NewRoom(e.Id)
		r.Tag = e.Tag
		r.Name = e.Name
		r.NotAvailable = tsl
		db.RealRooms[e.Id] = e.Tag
	}
}

// In the case of RoomGroups, cater for empty Tags (Shortcuts).
func (db *DbTopLevel) readRoomGroups(newdb *base.DbTopLevel) {
	db.RoomGroupMap = map[Ref]bool{}
	for _, e := range db.RoomGroups {
		if e.Tag != "" {
			_, nok := db.RoomTags[e.Tag]
			if nok {
				base.Report(
					roomtagerror,
					e.Tag)
				continue
			}
			db.RoomTags[e.Tag] = e.Id
		}
		// Copy to base db.
		r := newdb.NewRoomGroup(e.Id)
		r.Tag = e.Tag
		r.Name = e.Name
		r.Rooms = e.Rooms
		db.RoomGroupMap[e.Id] = true
	}
}

// Call this after all room types have been "read".
func (db *DbTopLevel) checkRoomGroups(newdb *base.DbTopLevel) {
	for _, e := range newdb.RoomGroups {
		// Collect the Ids and Tags of the component rooms.
		taglist := []string{}
		reflist := []Ref{}
		for _, rref := range e.Rooms {
			rtag, ok := db.RealRooms[rref]
			if ok {
				reflist = append(reflist, rref)
				taglist = append(taglist, rtag)
				continue

			}
			base.Report(
				`<Error>Invalid Room in RoomGroup %s:\n  -- %s>`,
				e.Tag, rref)
		}
		if e.Tag == "" {
			// Make a new Tag
			var tag string
			i := 0
			for {
				i++
				tag = "{" + strconv.Itoa(i) + "}"
				_, nok := db.RoomTags[tag]
				if !nok {
					break
				}
			}
			e.Tag = tag
			db.RoomTags[tag] = e.Id
			// Also extend the name
			if e.Name == "" {
				e.Name = strings.Join(taglist, ",")
			} else {
				e.Name = strings.Join(taglist, ",") + ":: " + e.Name
			}
		} else if e.Name == "" {
			e.Name = strings.Join(taglist, ",")
		}
		e.Rooms = reflist
	}
}

func (db *DbTopLevel) makeRoomChoiceGroup(
	newdb *base.DbTopLevel,
	rooms []Ref,
) (Ref, string) {
	erlist := []string{} // Error messages
	// Collect the Ids and Tags of the component rooms.
	taglist := []string{}
	reflist := []Ref{}
	for _, rref := range rooms {
		rtag, ok := db.RealRooms[rref]
		if ok {
			reflist = append(reflist, rref)
			taglist = append(taglist, rtag)
			continue
		}
		erlist = append(erlist,
			fmt.Sprintf(
				"  ++ Invalid Room in new RoomChoiceGroup:\n  %s\n", rref))
	}
	name := strings.Join(taglist, ",")
	// Reuse existing Element when the rooms match.
	id, ok := db.RoomChoiceNames[name]
	if !ok {
		// Make a new Tag
		var tag string
		i := 0
		for {
			i++
			tag = "[" + strconv.Itoa(i) + "]"
			_, nok := db.RoomTags[tag]
			if !nok {
				break
			}
		}
		// Add new Element
		r := newdb.NewRoomChoiceGroup("")
		id = r.Id
		r.Tag = tag
		r.Name = name
		r.Rooms = reflist
		db.RoomTags[tag] = id
		db.RoomChoiceNames[name] = id
	}
	return id, strings.Join(erlist, "")
}
