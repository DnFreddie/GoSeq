package notes

import "DnFreddie/goseq/lib"

type DayilyNoteRetriver struct{}
func NewDRetriver()*DayilyNoteRetriver{
	retriver := DayilyNoteRetriver{}
	return  &retriver
}

func (d *DayilyNoteRetriver) GetNotes(p lib.Period) ([]DNote,error){
	
	notes,err := getNotes(p)
	return notes,err
}

