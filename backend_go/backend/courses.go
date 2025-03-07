package backend

import "gradgrind/backend/base"

func init() {
	commandMap["GET_COURSES"] = getCourses
}

// TODO: This is probably not what is needed. Rather I will want the data
// for the courses tables, or rather a single course table?
// I could pass a class, teacher (or subject?) and get the data to be
// displayed?

/* What is needed? For a class, teacher (or subject?) a list of courses.
The table has fields:
    subject, groups, teachers, lesson lengths, properties(?), rooms
The side panel ("form") allows editing of the selected course and has
the additional fields:
    block-name, block-subject, constraints
This information should be fetchable separately, but if the back-end has
a concept of "selected", this will need to be fetched with the list data
for the selected course.
There might also be summary information about the number of lessons for
the teacher or groups.
*/

/* A subgroup can "belong" to more than one supergroup, which complicates
things a bit. It is understandable in that the subgroup basically – as far as
the timetabling is concerned – just functions as a set of "resources"
(teachers, groups, etc.). The supercourses provide only the lessons (and a
name/tag).
Editing possibilities:
    - add course to block (there is a problem with existing lessons in the
    course, if the block is new – has no lessons – they could be transferred,
    otherwise need confirmation that they should be discarded)
    - remove course from block (if block now has no subcourses it can be
    discarded?), option to keep or discard fully freed block?
*/

func getCourses(cmd map[string]any, outmap map[string]any) bool {
	outmap["Courses"] = DB.Courses
	outmap["SuperCourses"] = DB.SuperCourses
	outmap["SubCourses"] = DB.SubCourses
	return true
}

/* The handler for the courses gui ...

1) Full initialization of the view.
 - decide on classes, teachers or subjects
 - set the selection list
 - collect the data for the initial selection

*/

func CourseViewTypes() []string {
	vals := []string{`<>Class>`, `<>Teacher>`, `<>Subject>`}
	res := make([]string, len(vals))
	for i, s := range vals {
		res[i], _ = base.I18N(s)
	}
	return res
}

type CoursesState struct {
	// Maintains the state of the current courses view.
	// Note that this is dependent on the current database and needs
	// to be renewed when a new data set is loaded.
	tableType int
	tableShow string
}

func coursesInit(courseState *CoursesState) {
	if courseState == nil {
		//TODO ... Where is the state saved?
		courseState = &CoursesState{}
	}
	cmd := map[string]any{}
	switch courseState.tableType {
	case 0: // show a class's courses
		// Set the classes
		clist := make([]string, len(DB.Classes))
		for i, c := range DB.Classes {
			clist[i] = c.Name
		}
		cmd["GUI"] = "COMBOBOX_SET_ITEMS"
		cmd["OBJECT"] = "table_show"
		cmd["DATA"] = clist
		// I could do this as a function call, like
		//   send_cmd("COMBOBOX_SET_ITEMS", "table_show", clist)
		// if all gui commands can have the same structure.

		// Then I would need the courses involving the currently selected
		// class to set up the table rows. The printing module may help with
		// some ideas here ...
		// The rows could be sent as a list of string-lists.

		// Finally, one of the courses (by default the first?) would be
		// selected (leading to its display in the editor panel).

	case 1: // show a teacher's courses

	case 2: // show the courses in a particular subject

	default:
		base.Report(`<Bug>coursesInit: courseState.tableType = %d>`,
			courseState.tableType)
	}
}

func gui(cmd string, object string, data any) {
	payload := map[string]any{
		"GUI":    cmd,
		"OBJECT": object,
		"DATA":   data,
	}
	SenderChannel <- payload
}
